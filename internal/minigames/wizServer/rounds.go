package wizServer

import (
	"time"
)

type RoundPhase int

const (
	Playing RoundPhase = iota
	RoundOver
	RoundStart
	GameOver
	GameStart
)

const (
	Rounds         = 3   // 3 rounds
	RoundTime      = 150 // 2.5 min per round
	RoundOverTime  = 3   // 3 seconds displaying scoreboard
	RoundStartTime = 1   // 1 seconds displaying round number
	GameOverTime   = 5   // 8 (3 from intermission + 5 additional) seconds at the end of the game
	GameStartTime  = 1
)

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizardServer) startNewRound() {
	clear(w.openedChests)
	w.enemies = make(map[uint16]*Enemy, 32)

	w.phaseTimer = time.Now().Add(RoundTime * time.Second)
	w.enemyRespawnTimer = time.Now().Add(RoundStartEnemySpawnTimer * time.Second)
	w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)

	w.phase = Playing
}

func (w *WizardServer) updateRound() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.round += 1

		w.phase = RoundOver
		w.phaseTimer = time.Now().Add(RoundOverTime * time.Second)
		// TODO:Send scores
	}
}

func (w *WizardServer) updateRoundStart() {
	w.AssignPlayerSpawns()

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.startNewRound()
	}
}

func (w *WizardServer) updateRoundOver() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		if w.round > Rounds {
			w.phase = GameOver
			w.phaseTimer = time.Now().Add(GameOverTime * time.Second)
			w.gameOver()
			return
		}

		w.phase = RoundStart
		w.phaseTimer = time.Now().Add(RoundStartTime * time.Second)
	}
}

func (w *WizardServer) gameStart() {
	w.AssignPlayerSpawns()

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.startNewRound()
	}
}

func (w *WizardServer) gameOver() {
}
