package lobby

import "github.com/ben-rw/ragin-mages/internal/room"

type LobbyServer struct {
	*room.Room
}

func Init(r *room.Room) any {
	r.Scene.Cancel = nil

	return &LobbyServer{}
}
