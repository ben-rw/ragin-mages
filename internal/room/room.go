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
	Scene             *protocol.Scene
}

// check for username collisions, set player.Room, add player to Room
func (r *Room) AppendPlayer(player *Player) error {
	r.Mu.Lock()
	player.Username = nameCollisionSolver(player.Username, r.Players)
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

func (r *Room) GetPlayerDataList() ([]*protocol.PlayerData, error) {
	playerList, err := r.GetPlayerList()
	if err != nil {
		return nil, err
	}

	dataList := make([]*protocol.PlayerData, len(playerList))
	for i, player := range playerList {
		dataList[i] = player.PlayerToPlayerData()
	}

	return dataList, nil
}

func (r *Room) Broadcast(msg *protocol.Message) {
	r.Mu.RLock()
	playerList := make([]*Player, len(r.Players))
	_ = copy(playerList, r.Players)
	r.Mu.RUnlock()

	for _, player := range playerList {
		player.Mu.RLock()
		client := player.Client
		player.Mu.RUnlock()

		if client != nil {
			err := wsjson.Write(context.Background(), client.Conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (r *Room) BroadcastExcludeSender(msg *protocol.Message, senderName string) {
	r.Mu.RLock()
	playerList := make([]*Player, len(r.Players))
	_ = copy(playerList, r.Players)
	r.Mu.RUnlock()

	for _, player := range playerList {
		if player.Username == senderName {
			continue
		}

		player.Mu.RLock()
		client := player.Client
		player.Mu.RUnlock()

		if client != nil {
			err := wsjson.Write(context.Background(), client.Conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (r *Room) BroadcastMsg(mt protocol.MessageType, msgData any) {
	msg, err := protocol.MarshalToMessage(mt, msgData)
	if err != nil {
		log.Printf("couldn't broadcast message of type %v: %v\n", mt, msgData)
		return
	}
	r.Broadcast(msg)
}
