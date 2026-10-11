package wizServer

import (
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

type RoundPhase int

const (
	Playing RoundPhase = iota
	RoundOver
	RoundStart
	GameOver
	GameStart
	WaitingForClients
	StartNewRound
)

const (
	Rounds             = 3   // 3 rounds
	RoundTime          = 150 // 2.5 min per round
	RoundOverTime      = 3   // 3 seconds displaying scoreboard
	RoundStartTime     = 1   // 1 seconds displaying round number
	GameOverTime       = 5   // 8 (3 from intermission + 5 additional) seconds at the end of the game
	GameStartTime      = 1
	WaitForClientsTime = 5 // 5 seconds max waiting for clients to connect
	SyncTime           = 5 // sync round timer every 5 seconds
)

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizardServer) startNewRound() {
	w.enemies = make(map[uint16]*Enemy, 32)

	w.phaseTimer = time.Now().Add(RoundTime * time.Second)
	w.syncTimer = time.Now().Add(SyncTime * time.Second)
	w.enemyRespawnTimer = time.Now().Add(RoundStartEnemySpawnTimer * time.Second)
	w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)

	w.phase = Playing
	w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
		Phase: int(w.phase),
	})
}

func (w *WizardServer) updateRound() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.round += 1

		w.phase = RoundOver
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.phaseTimer = time.Now().Add(RoundOverTime * time.Second)
	}
}

func (w *WizardServer) updateRoundStart() {
	w.AssignPlayerSpawns()
	w.updateMovement()

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.phase = StartNewRound
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.startNewRound()
	}
}

func (w *WizardServer) updateRoundOver() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		if w.round > Rounds {
			w.phase = GameOver
			w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
				Phase: int(w.phase),
			})
			w.phaseTimer = time.Now().Add(GameOverTime * time.Second)
			return
		}

		w.phase = RoundStart
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.phaseTimer = time.Now().Add(RoundStartTime * time.Second)
	}
}

func (w *WizardServer) gameStart() {
	w.AssignPlayerSpawns()
	w.updateMovement()

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.startNewRound()
	}
}

func (w *WizardServer) gameOver() {
}

func (w *WizardServer) waitForClients() {
	if w.clientLoadedCount == uint8(len(w.Players)) {
		w.phaseTimer = time.Now()
	}

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.phase = GameStart
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.phaseTimer = time.Now().Add(GameStartTime * time.Second)
	}
}
