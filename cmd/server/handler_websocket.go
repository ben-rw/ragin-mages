package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ben-rw/ragin-mages/internal/minigames"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// get roomID + username, upgrade to websocket, register client to hub, check for reconnection,
// send username to client, start read loop, remove client when read loop returns
func (cfg *config) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	_, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		http.Error(w, "no active room with that room code", http.StatusBadRequest)
		return
	}

	currentRoom := cfg.RoomReg.ActiveRooms[roomID]

	userCookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "missing user session data", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	username := userCookie.Value

	idCookie, err := r.Cookie("id")
	if err != nil {
		http.Error(w, "missing user session data", http.StatusUnauthorized)
		log.Println(err)
		return
	}
	id := idCookie.Value

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "couldn't upgrade connection to websocket", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	client := room.Client{
		Conn: conn,
	}

	playerDataList, err := getPlayerDataList(currentRoom)
	if err != nil {
		log.Println(err)
	}

	var host bool
	if len(playerDataList) == 0 {
		host = true
	}

	existingPlayer, exists, err := checkForReconnect(currentRoom, id)
	if err != nil {
		log.Println(err)
	}

	var player *room.Player
	if exists {
		existingPlayer.Client = &client
		client.Player = existingPlayer
		player = existingPlayer
	} else {
		player = &room.Player{
			ID:          id,
			Name:        username,
			Score:       0,
			Host:        host,
			SpriteIndex: currentRoom.AssignPlayerSprite(),
			Room:        currentRoom,
			Client:      &client,
		}
		err = currentRoom.AppendPlayer(player)
		if err != nil {
			log.Println(err)
		}
		client.Player = player
	}

	playerData := playerToPlayerData(player)

	if !exists {
		playerDataList = append(playerDataList, playerData)
	}

	data := protocol.JoinResponseData{
		PlayerData: playerData,
		PlayerList: playerDataList,
	}

	msg, err := protocol.MarshalToMessage(protocol.JoinResponse, data)
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		log.Println(err)
		return
	}

	readLoop(currentRoom, &client)
	player.Client = nil
}

// read incoming messages, close dead connections, drop messages if buffer fills up
func readLoop(r *room.Room, c *room.Client) {
	defer c.Close(websocket.StatusInternalError, "connection closed unexpectedly")

	msg := protocol.Message{}

	for {
		err := wsjson.Read(context.Background(), c.Conn, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				c.Close(websocket.StatusNormalClosure, "connection closed normally")
				return
			} else {
				log.Printf("server: wsjson read error: %v\n", err)
				log.Println("server: terminating read loop")
				return
			}
		}

		updateMsg, err := minigames.ValidateMessage(&msg, c.Player.Room.Scene)
		if updateMsg.Type != protocol.Unset {
			log.Printf("msg: %v, updateMsg: %v\n", string(msg.Data), string(updateMsg.Data))
		}
		if err != nil {
			log.Printf("couldn't validate message: %v\n", err)
			continue
		} else if updateMsg.Type == protocol.Unset {
			log.Println("message not broadcasted: empty message")
			continue
		} else {
			r.Broadcast(updateMsg)
		}

		if updateMsg.Type == protocol.SceneChange {
			playerList, err := r.GetPlayerList()
			if err != nil {
				log.Printf("couldn't get player list: %v\n", err)
				continue
			}

			playerDataList, err := getPlayerDataList(r)
			if err != nil {
				log.Printf("couldn't get player data list: %v\n", err)
				continue
			}

			for _, player := range playerList {

				scenePlayerData := &protocol.JoinResponseData{
					PlayerData: playerToPlayerData(player),
					PlayerList: playerDataList,
				}

				msg, err := protocol.MarshalToMessage(protocol.JoinResponse, scenePlayerData)
				if err != nil {
					log.Printf("couldn't marshal msg: %v\n", err)
					continue
				}

				if player.Client != nil {
					err := wsjson.Write(context.Background(), player.Client.Conn, msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func getPlayerDataList(r *room.Room) ([]*protocol.PlayerData, error) {
	playerList, err := r.GetPlayerList()
	if err != nil {
		return nil, err
	}

	dataList := make([]*protocol.PlayerData, len(playerList))
	for i, player := range playerList {
		dataList[i] = playerToPlayerData(player)
	}

	return dataList, nil
}

func playerToPlayerData(player *room.Player) *protocol.PlayerData {
	return &protocol.PlayerData{
		Name:        player.Name,
		Score:       player.Score,
		Host:        player.Host,
		SpriteIndex: player.SpriteIndex,
	}
}

func checkForReconnect(r *room.Room, id string) (*room.Player, bool, error) {
	playerList, err := r.GetPlayerList()
	if err != nil {
		return nil, false, err
	}

	for _, player := range playerList {
		if player.ID == id {
			return player, true, nil
		}
	}

	return nil, false, nil
}
