package room

import (
	"context"
	"errors"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/coder/websocket/wsjson"
	"log"
	"sync"
)

const (
	MaxPlayers = 8
)

type Room struct {
	PlayerSpriteIndex *[]int
	Mu                *sync.RWMutex
	ID                string
	Players           []*Player
	Scene             protocol.SceneType
}

// check for username collisions, set player.Room, add player to Room
func (r *Room) AppendPlayer(player *Player) error {
	r.Mu.Lock()
	player.Name = nameCollisionSolver(player.Name, r.Players)
	player.Room = r
	r.Players = append(r.Players, player)
	r.Mu.Unlock()

	return nil
}

func (r *Room) RemovePlayer(currentPlayer *Player) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	var i int
	var match bool
	for index, player := range r.Players {
		if player == currentPlayer {
			i = index
			match = true
			break
		}
	}
	if !match {
		return errors.New("player not found")
	}

	r.Players = append(r.Players[:i], r.Players[i+1:]...)

	return nil
}

func (r *Room) GetPlayerList() ([]*Player, error) {
	r.Mu.RLock()
	playerList := make([]*Player, len(r.Players))
	_ = copy(playerList, r.Players)
	r.Mu.RUnlock()
	return playerList, nil
}

func (r *Room) Broadcast(msg *protocol.Message) {
	r.Mu.RLock()
	playerList := make([]*Player, len(r.Players))
	_ = copy(playerList, r.Players)
	r.Mu.RUnlock()

	for _, player := range playerList {
		if player.Client != nil {
			err := wsjson.Write(context.Background(), player.Client.Conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
