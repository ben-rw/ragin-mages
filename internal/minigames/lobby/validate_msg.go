package lobby

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"log"
)

func ValidateMsg(msg *protocol.Message, l *LobbyServer) error {
	var data any
	if msg.Type == protocol.PlayerUpdate {
		l.Broadcast(msg)
		return nil
	}

	log.Printf("couldn't validate msg data: %v", data)
	return nil
}
