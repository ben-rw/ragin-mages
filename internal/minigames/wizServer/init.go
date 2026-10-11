package wizServer

import (
	"context"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

type WizardServer struct {
	*room.Room
	wizards           map[string]*Wizard
	enemies           map[uint16]*Enemy
	enemyRespawnTimer time.Time
	phaseTimer        time.Time
	syncTimer         time.Time
	chestRespawnTimer time.Time
	phase             RoundPhase
	round             uint8
	clientLoadedCount uint8
}

func Init(r *room.Room) *WizardServer {
	r.Scene.Type = protocol.WizardsScene
	if r.Scene.Cancel != nil {
		r.Scene.Cancel()
		r.Scene.Cancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	wizards := make(map[string]*Wizard, len(r.Players))
	for _, player := range r.Players {
		wizards[player.Username] = &Wizard{Username: player.Username}
	}

	w := &WizardServer{
		Room:              r,
		wizards:           wizards,
		enemies:           make(map[uint16]*Enemy, 32), // 16 enemies spawn at at time, doubled for headroom
		phaseTimer:        time.Now().Add(GameStartTime * time.Second),
		enemyRespawnTimer: time.Now().Add(RoundStartEnemySpawnTimer * time.Second), // wait 5 seconds before spawning first group of enemies
		phase:             WaitingForClients,
		round:             1,
		clientLoadedCount: 0,
	}
	r.Scene.State = w

	go w.Update(ctx, cancel)

	return w
}
