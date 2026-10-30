package wizServer

import (
	"context"
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

func (w *WizardServer) Update(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	const tickRate = 60
	ticker := time.NewTicker(time.Second / tickRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return
		case player := <-w.playerJoined:
			w.PlayerConnectedMidGame(player)
		case <-ticker.C:
			switch w.phase {
			case Playing:
				w.updatePlaying()
			case RoundOver:
				w.updateRoundOver()
			case RoundStart:
				w.updateRoundStart()
			case GameOver:
				w.gameOver()
			case GameStart:
				w.gameStart()
			case WaitingForClients:
				w.waitForClients()
			case SyncScores:
				w.syncScores()
			default:
				log.Printf("wizServer: unrecognized round phase: %v\n", w.phase)
			}
		}
	}

}

func (w *WizardServer) PlayerConnectedMidGame(player *room.Player) {
	if _, ok := w.wizards[player.Username]; !ok {
		// create new wizard
		w.wizards[player.Username] = NewWizard(player.Username)

		// send other clients info about the new wizard
		playerUpdateMsg := &protocol.PlayerUpdateData{
			PlayerData: player.PlayerToPlayerData(),
		}

		pumsg, err := protocol.MarshalToMessage(protocol.PlayerUpdate, playerUpdateMsg)
		if err != nil {
			log.Printf("couldn't connect player %v mid-game, failed to marshal player update message: %v\n", player.Username, err)
		}
		w.BroadcastExcludeSender(pumsg, player.Username)
	}

	// send the connecting wizard's client the list of the room's players and identify the connecting wizard to themself
	playerDataList, err := w.GetPlayerDataList()
	if err != nil {
		log.Printf("couldn't connect player %v mid-game, failed to get player data list: %v\n", player.Username, err)
	}

	joinResponseMessage := &protocol.JoinResponseData{
		PlayerData: player.PlayerToPlayerData(),
		PlayerList: playerDataList,
	}

	w.BroadcastMsg(protocol.JoinResponse, joinResponseMessage)
}

func (w *WizardServer) updatePlaying() {
	w.updateEnemies()
	w.checkEnemyCollisions()

	w.updateMovement()
	w.updateSyncTimer()
	w.enemyRespawn()
	w.chestRespawn()
	w.updateRound()
}

func (w *WizardServer) updateSyncTimer() {
	if w.syncTimer.IsZero() || time.Until(w.syncTimer) <= 0 {
		w.BroadcastMsg(protocol.RoundTimerUpdate, &protocol.RoundTimerUpdateData{
			TimeRemaining: int(time.Until(w.phaseTimer).Milliseconds()),
		})
		w.syncTimer = time.Now().Add(SyncTime * time.Second)
	}
}

func (w *WizardServer) updateMovement() {
	w.Mu.RLock()
	wizardsMovementData := []*protocol.WizardMovementUpdateData{}
	for _, wizard := range w.wizards {
		wmd := &protocol.WizardMovementUpdateData{
			X:        wizard.X,
			Y:        wizard.Y,
			Dx:       wizard.Dx,
			Dy:       wizard.Dy,
			Username: wizard.Username,
		}
		wizardsMovementData = append(wizardsMovementData, wmd)
	}

	enemiesMovementData := []*protocol.EnemyMovementUpdateData{}
	for _, enemy := range w.enemies {
		emd := &protocol.EnemyMovementUpdateData{
			X:  enemy.X,
			Y:  enemy.Y,
			Dx: enemy.Dx,
			Dy: enemy.Dy,
			ID: enemy.ID,
		}
		enemiesMovementData = append(enemiesMovementData, emd)
	}
	w.Mu.RUnlock()

	w.BroadcastMsg(protocol.WizArenaMovement, &protocol.WizArenaMovementData{
		WizardsMovementData: wizardsMovementData,
		EnemiesMovementData: enemiesMovementData,
	})
}
