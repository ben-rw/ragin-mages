package room

import (
	"context"
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/coder/websocket/wsjson"
)

func (r *Room) SendAllPlayerDataToNewScene() {
	playerList, err := r.GetPlayerList()
	if err != nil {
		log.Printf("couldn't get player list: %v\n", err)
		return
	}

	playerDataList, err := r.GetPlayerDataList()
	if err != nil {
		log.Printf("couldn't get player data list: %v\n", err)
		return
	}

	for _, player := range playerList {
		scenePlayerData := &protocol.JoinResponseData{
			PlayerData: player.PlayerToPlayerData(),
			PlayerList: playerDataList,
		}

		if player.Client != nil {
			msg, err := protocol.MarshalToMessage(protocol.JoinResponse, scenePlayerData)
			if err != nil {
				log.Printf("couldn't marshal msg: %v\n", err)
				continue
			}

			err = wsjson.Write(context.Background(), player.Client.Conn, msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
