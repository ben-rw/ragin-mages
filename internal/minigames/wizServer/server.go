package wizServer

import (
	"context"
	"time"

	"github.com/ben-rw/ragin-mages/internal/room"
)

type WizardServer struct {
	*room.Room
	wizards           map[string]*Wizard
	enemies           map[uint16]*Enemy
	enemyRespawnTimer time.Time
	phaseTimer        time.Time
	chestRespawnTimer time.Time
	openedChests      map[OpenedChest]struct{}
	phase             RoundPhase
	round             uint8
}

func StartWizardServer(r *room.Room) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wizards := make(map[string]*Wizard, len(r.Players))
	for _, wizard := range r.Players {
		wizards[wizard.Name] = &Wizard{}
	}

	w := &WizardServer{
		Room:              r,
		wizards:           wizards,
		enemies:           make(map[uint16]*Enemy, 32), // 16 enemies spawn at at time, doubled for headroom
		phaseTimer:        time.Now().Add(GameStartTime * time.Second),
		enemyRespawnTimer: time.Now().Add(RoundStartEnemySpawnTimer * time.Second), // wait 5 seconds before spawning first group of enemies
		phase:             GameStart,
	}

	w.Update(ctx)
}

func (w *WizardServer) Update(ctx context.Context) {

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
	default:
	}

}

func (w *WizardServer) updatePlaying() {
	w.enemyRespawn()
	w.chestRespawn()
	w.updateRound()
}
