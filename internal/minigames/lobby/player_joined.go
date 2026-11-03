package lobby

import (
	"context"
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
	"github.com/coder/websocket/wsjson"
)

func (l *LobbyServer) PlayerJoined(player *room.Player, playerDataList []*protocol.PlayerData, playerData *protocol.PlayerData, exists bool) {
	// send new player the list of players and identify the new player to theirself
	data := protocol.LobbyStateSnapshotData{
		PlayerData: playerData,
		PlayerList: playerDataList,
	}

	msg, err := protocol.MarshalToMessage(protocol.LobbyStateSnapshot, data)
	if err != nil {
		log.Println(err)
	}

	err = wsjson.Write(context.Background(), player.Client.Conn, msg)
	if err != nil {
		log.Println(err)
		return
	}
}
