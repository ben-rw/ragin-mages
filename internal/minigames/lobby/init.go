package lobby

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

type LobbyServer struct {
	*room.Room
}

func Init(r *room.Room) {
	r.Scene.Type = protocol.LobbyScene
	if r.Scene.Cancel != nil {
		r.Scene.Cancel()
		r.Scene.Cancel = nil
	}

	l := &LobbyServer{
		Room: r,
	}
	r.Scene.State = l
}
