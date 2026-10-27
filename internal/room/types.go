package room

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/coder/websocket"
)

type Player struct {
	Room        *Room
	Client      *Client
	ID          string
	Username    string
	Score       int
	SpriteIndex int
	X, Y        float64
	Host        bool
}

type Client struct {
	*websocket.Conn
	Player *Player
}

func (p *Player) PlayerToPlayerData() *protocol.PlayerData {
	return &protocol.PlayerData{
		Username:    p.Username,
		Score:       p.Score,
		Host:        p.Host,
		SpriteIndex: p.SpriteIndex,
	}
}
