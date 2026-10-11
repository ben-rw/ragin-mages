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
	phase             RoundPhase
	round             uint8
	clientLoadedCount uint8
}

func Init(r *room.Room) *WizardServer {
	ctx, cancel := context.WithCancel(context.Background())

	wizards := make(map[string]*Wizard, len(r.Players))
	for _, wizard := range r.Players {
		wizards[wizard.Username] = &Wizard{}
	}

	w := &WizardServer{
		Room:              r,
		wizards:           wizards,
		enemies:           make(map[uint16]*Enemy, 32), // 16 enemies spawn at at time, doubled for headroom
		phaseTimer:        time.Now().Add(GameStartTime * time.Second),
		enemyRespawnTimer: time.Now().Add(RoundStartEnemySpawnTimer * time.Second), // wait 5 seconds before spawning first group of enemies
		phase:             GameStart,
		round:             1,
		clientLoadedCount: 0,
	}

	go w.Update(ctx, cancel)

	return w
}

func (w *WizardServer) Update(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

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
	}

}

func (w *WizardServer) updatePlaying() {
	w.enemyRespawn()
	w.chestRespawn()
	w.updateRound()
}
