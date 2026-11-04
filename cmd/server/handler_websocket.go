package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ben-rw/ragin-mages/internal/minigames"
	"github.com/ben-rw/ragin-mages/internal/minigames/lobby"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// get roomID + username, upgrade to websocket, register client to hub, check for reconnection,
// send username to client, start read loop, remove client when read loop returns
func (cfg *config) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	currentRoom, ok := cfg.RoomReg.Get(roomID)
	if !ok {
		http.Error(w, "no active room with that room code", http.StatusBadRequest)
		return
	}

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

	hostCookie, err := r.Cookie("host")
	if err != nil {
		http.Error(w, "missing user session data", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	var host bool
	if hostCookie.Value == "true" {
		host = true
	} else {
		host = false
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "couldn't upgrade connection to websocket", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	client := room.Client{
		Conn: conn,
	}

	playerDataList, err := currentRoom.GetPlayerDataList()
	if err != nil {
		log.Println(err)
	}

	existingPlayer, exists, err := checkForReconnect(currentRoom, id)
	if err != nil {
		log.Println(err)
	}

	var player *room.Player
	if exists {
		existingPlayer.Mu.Lock()
		client.Player = existingPlayer
		existingPlayer.Client = &client
		player = existingPlayer
		player.ClientInitDone = false
		player.ClientStateUpdated = false
		existingPlayer.Mu.Unlock()
	} else {
		player = &room.Player{
			ID:                 id,
			Username:           username,
			Mu:                 &sync.RWMutex{},
			Score:              0,
			Host:               host,
			SpriteIndex:        currentRoom.AssignPlayerSprite(),
			Room:               currentRoom,
			Client:             &room.Client{},
			ClientInitDone:     false,
			ClientStateUpdated: false,
		}
		client.Player = player
		player.Client = &client

		err = currentRoom.AppendPlayer(player)
		if err != nil {
			log.Println(err)
		}
	}

	// send scene change message to get new client in the active scene
	currentRoom.Mu.RLock()
	scdata := &protocol.SceneChangeData{
		SceneType: currentRoom.Scene.Type,
	}
	currentRoom.Mu.RUnlock()

	scmsg, err := protocol.MarshalToMessage(protocol.SceneChange, scdata)
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	err = wsjson.Write(ctx, conn, scmsg)
	if err != nil {
		log.Println(err)
		return
	}

	playerData := player.PlayerToPlayerData()
	if !exists {
		playerDataList = append(playerDataList, playerData)
	}

	// sync scene state
	currentRoom.Mu.Lock()
	sceneState := currentRoom.Scene.State
	currentRoom.Mu.Unlock()

	if wizServer, ok := sceneState.(*wizServer.WizardServer); ok {
		wizServer.PlayerJoined(player, playerDataList)
	} else if lobbyServer, ok := sceneState.(*lobby.LobbyServer); ok {
		lobbyServer.PlayerJoined(player, playerDataList, playerData, exists)
	}

	// start read loop
	readLoop(currentRoom, &client)

	// nil obsolete clients
	player.Mu.Lock()
	if player.Client == &client {
		player.Client = nil
	}
	player.Mu.Unlock()

	// TODO: remove disconnected players
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

		err = minigames.ValidateMessage(&msg, r, c)
		// if updateMsg.Type != protocol.Unset {
		// 	log.Printf("msg: %v, updateMsg: %v\n", string(msg.Data), string(updateMsg.Data))
		// }
		if err != nil {
			log.Printf("couldn't validate message: %v\n", err)
			continue
		}
		// else if updateMsg.Type == protocol.Unset {
		// 	log.Println("message not broadcasted: empty message")
		// 	continue
		// } else {
		// 	r.Broadcast(updateMsg)
		// }
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
