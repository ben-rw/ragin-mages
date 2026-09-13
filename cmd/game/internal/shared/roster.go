package shared

import (
	"errors"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"log"
)

type Roster struct {
	Players map[string]*Player
	Player  *Player
}

func (r *Roster) HandleJoinResponse(message protocol.Message) error {
	joinResponseData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}

	data, ok := joinResponseData.(*protocol.JoinResponseData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	for _, playerData := range data.PlayerList {
		if _, ok := r.Players[playerData.Name]; ok {
			r.Players[playerData.Name].Data = playerData
		} else {
			player := NewPlayer(playerData, len(r.Players))
			r.Players[playerData.Name] = player
		}
	}

	r.Player = r.Players[data.PlayerData.Name]
	return nil
}

func (r *Roster) HandlePlayerUpdate(message protocol.Message) error {
	playerUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := playerUpdateData.(*protocol.PlayerUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if _, ok := r.Players[data.PlayerData.Name]; ok {
		r.Players[data.PlayerData.Name].Data = data.PlayerData
	} else {
		player := NewPlayer(data.PlayerData, len(r.Players))
		r.Players[data.PlayerData.Name] = player
	}

	log.Printf("updated: %v", *r.Players[data.PlayerData.Name].Data)

	return nil
}
