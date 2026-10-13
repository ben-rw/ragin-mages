package wizServer

import (
	"context"
	"image"
	"log"
	"sync/atomic"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

type WizardServer struct {
	*room.Room
	wizards           map[string]*Wizard
	deadWizards       map[string]struct{}
	enemies           map[uint16]*Enemy
	holes             []image.Rectangle
	tilemapJSON       *TilemapJSON
	enemyRespawnTimer time.Time
	phaseTimer        time.Time
	syncTimer         time.Time
	chestRespawnTimer time.Time
	clientLoadedCount atomic.Int32
	phase             RoundPhase
	round             uint8
	enemyID           uint16
	spawnsAssigned    bool
}

func Init(r *room.Room) *WizardServer {
	r.Scene.Type = protocol.WizardsScene
	if r.Scene.Cancel != nil {
		r.Scene.Cancel()
		r.Scene.Cancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	players, err := r.GetPlayerList()
	if err != nil {
		log.Printf("error getting player list: %v\n", err)
	}

	wizards := make(map[string]*Wizard, len(players))
	deadWizards := make(map[string]struct{}, len(players))

	for _, player := range players {
		wizards[player.Username] = NewWizard(player.Username)
	}

	tilemap, err := NewTilemapJSON(TilemapPath)
	if err != nil {
		log.Printf("server couldn't load tilemap: %v\n", err)
	}

	w := &WizardServer{
		Room:              r,
		wizards:           wizards,
		deadWizards:       deadWizards,
		enemies:           make(map[uint16]*Enemy, 32), // 16 enemies spawn at at time, doubled for headroom
		phaseTimer:        time.Now().Add(GameStartTime * time.Second),
		enemyRespawnTimer: time.Now().Add(RoundStartEnemySpawnTimer * time.Second), // wait 5 seconds before spawning first group of enemies
		tilemapJSON:       tilemap,
		holes:             make([]image.Rectangle, 0, 2092), // 2092 holes on map
		phase:             WaitingForClients,
		round:             1,
		enemyID:           0,
		spawnsAssigned:    false,
	}
	r.Scene.State = w

	w.TileBounds()

	go w.Update(ctx, cancel)

	return w
}
