package room

import (
	"github.com/coder/websocket"
)

type Player struct {
	Room        *Room
	Client      *Client
	ID          string
	Name        string
	Score       int
	SpriteIndex int
	X, Y        float64
	Host        bool
}

type Client struct {
	*websocket.Conn
	Player *Player
}
