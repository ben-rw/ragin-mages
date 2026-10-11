package wizServer

import (
	"context"
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizardServer) Update(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	const tickRate = 60
	ticker := time.NewTicker(time.Second / tickRate)

	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return
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
			default:
				log.Printf("wizServer: unrecognized round phase: %v", w.phase)
			}
		}
	}

}

func (w *WizardServer) updatePlaying() {
	w.updateMovement()
	w.updateSyncTimer()

	w.enemyRespawn()
	w.chestRespawn()

	w.updateRound()
}

func (w *WizardServer) updateSyncTimer() {
	if w.syncTimer.IsZero() || time.Until(w.syncTimer) <= 0 {
		w.BroadcastMsg(protocol.RoundTimerUpdate, &protocol.RoundTimerUpdateData{
			TimeRemaining: int(time.Until(w.phaseTimer) * time.Millisecond),
		})
		w.syncTimer = time.Now().Add(SyncTime * time.Second)
	}
}

func (w *WizardServer) updateMovement() {
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

	w.BroadcastMsg(protocol.WizArenaMovement, &protocol.WizArenaMovementData{
		WizardsMovementData: wizardsMovementData,
		EnemiesMovementData: enemiesMovementData,
	})
}
