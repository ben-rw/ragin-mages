package wizServer

import (
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

const ChestRespawnTimer = 30

// spawn chests on a timer
func (w *WizardServer) chestRespawn() {
	if w.chestRespawnTimer.IsZero() || time.Until(w.chestRespawnTimer) <= 0 {
		// send 0, 0 to indicate openChests should cleared
		w.BroadcastMsg(protocol.ChestUpdate, &protocol.ChestUpdateData{
			X: 0,
			Y: 0,
		})
		w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)
	}
}
