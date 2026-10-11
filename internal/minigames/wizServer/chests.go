package wizServer

import (
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

const ChestRespawnTimer = 30

// spawn chests on a timer
func (w *WizardServer) chestRespawn() {
	if w.chestRespawnTimer.IsZero() || time.Until(w.chestRespawnTimer) <= 0 {
		w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)
	}
}

func (w *WizardServer) BroadcastMsg(mt protocol.MessageType, msgData any) {
	msg, err := protocol.MarshalToMessage(mt, msgData)
	if err != nil {
		log.Println(err)
		return
	}
	w.Broadcast(msg)
}
