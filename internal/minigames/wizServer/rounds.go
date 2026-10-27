package wizServer

import (
	"log"
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
	SyncScores
)

const (
	Rounds             = 3   // 3 rounds
	RoundTime          = 150 // 2.5 min per round
	RoundOverTime      = 4   // 4 seconds displaying scoreboard
	RoundStartTime     = 1   // 1 seconds displaying round number
	GameOverTime       = 4   // 8 (4 from RoundOver + 4 additional) seconds at the end of the game
	GameStartTime      = 1
	WaitForClientsTime = 5 // 5 seconds max waiting for clients to connect
	SyncTime           = 5 // sync round timer every 5 seconds
	SyncScoresTime     = 2
)

const WinRoundBoost = 10.0

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizardServer) startNewRound() {
	w.Mu.Lock()
	w.enemies = make(map[uint16]*Enemy, 32)
	w.deadWizards = make(map[string]struct{}, len(w.Players))
	w.Mu.Unlock()

	w.phaseTimer = time.Now().Add(RoundTime * time.Second)
	w.syncTimer = time.Now().Add(SyncTime * time.Second)
	w.enemyRespawnTimer = time.Now().Add(RoundStartEnemySpawnTimer * time.Second)
	w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)

	w.phase = Playing
	w.BroadcastMsg(protocol.RoundTimerUpdate, &protocol.RoundTimerUpdateData{
		TimeRemaining: int(time.Until(w.phaseTimer).Milliseconds()),
	})
}

func (w *WizardServer) updateRound() {
	if w.checkRoundOver() || w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.determineAndAwardWinner()

		w.round += 1

		w.phase = SyncScores
		w.phaseTimer = time.Now().Add(SyncScoresTime * time.Second)
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
	}
}

func (w *WizardServer) syncScores() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.phase = RoundOver
		w.phaseTimer = time.Now().Add(RoundOverTime * time.Second)
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
	}
}

// single player plays until killed or out of time
// last wizard alive wins the round in a multiplayer game
func (w *WizardServer) checkRoundOver() bool {
	if len(w.wizards) == 1 && len(w.deadWizards) == 1 {
		return true
	} else if len(w.wizards) > 1 && len(w.deadWizards) == len(w.wizards)-1 {
		return true
	} else {
		return false
	}
}

func (w *WizardServer) updateRoundStart() {
	if !w.spawnsAssigned {
		w.AssignPlayerSpawns()
		w.updateMovement()
		w.spawnsAssigned = true
	}

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.phase = StartNewRound
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.spawnsAssigned = false
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
			w.BroadcastMsg(protocol.RoundTimerUpdate, &protocol.RoundTimerUpdateData{
				TimeRemaining: int(time.Until(w.phaseTimer).Milliseconds()),
			})
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
	if !w.spawnsAssigned {
		w.AssignPlayerSpawns()
		w.updateMovement()
		w.spawnsAssigned = true
	}

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.phase = StartNewRound
		w.BroadcastMsg(protocol.PhaseUpdate, &protocol.PhaseUpdateData{
			Phase: int(w.phase),
		})
		w.spawnsAssigned = false
		w.startNewRound()
	}
}

func (w *WizardServer) gameOver() {
}

func (w *WizardServer) waitForClients() {
	playerList, err := w.GetPlayerList()
	if err != nil {
		log.Printf("error getting player list: %v\n", err)
	}
	if int(w.clientLoadedCount.Load()) >= (len(playerList)) {
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

func (w *WizardServer) determineAndAwardWinner() {
	aliveAtRoundEnd := make([]string, 0, len(w.wizards)-len(w.deadWizards))
	for _, wizard := range w.wizards {
		if _, ok := w.deadWizards[wizard.Username]; !ok {
			aliveAtRoundEnd = append(aliveAtRoundEnd, wizard.Username)
		}
	}

	if len(aliveAtRoundEnd) == 0 {
		// if the last two wizards die on the same tick, nobody is awarded points
		return
	} else if len(aliveAtRoundEnd) == 1 {
		w.BroadcastMsg(protocol.GivePoints, &protocol.GivePointsData{
			RecipientName: aliveAtRoundEnd[0],
			Random:        WinRoundBoost,
		})
	} else {
		// undecided - for now, nobody gets points if multiple wizards are alive at round end
		return
	}
}
