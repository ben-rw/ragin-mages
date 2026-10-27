package lobby

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"log"
)

func ValidateMsg(msg *protocol.Message, l *LobbyServer) error {
	var data any
	switch msg.Type {
	case protocol.PlayerUpdate:
		l.Broadcast(msg)
		return nil
	case protocol.FadeMusic:
		l.Broadcast(msg)
		return nil
	default:
	}

	log.Printf("couldn't validate msg data: %v\n", data)
	return nil
}
