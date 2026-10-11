package minigames

import (
	"errors"
	"math/rand"

	"github.com/ben-rw/ragin-mages/internal/minigames/lobby"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

// lobby, random omitted as they are not games
var minigameMap = map[int]protocol.SceneType{
	0: protocol.WizardsScene,
}

func ValidateMessage(message *protocol.Message, r *room.Room) error {
	if message.Type == protocol.SceneChange {
		sceneChangeData, err := message.UnmarshalMessageData()
		if err != nil {
			return err
		}
		scdata, ok := sceneChangeData.(*protocol.SceneChangeData)
		if !ok {
			err = errors.New("couldn't validate message: failed type assertion")
			return err
		}

		if scdata.SceneType == protocol.RandomScene {
			scdata.SceneType = minigameMap[rand.Intn(len(minigameMap))]
		}

		switch scdata.SceneType {
		case protocol.LobbyScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.LobbyScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return err
			}

			lobby.Init(r)
			r.Broadcast(msg)
			r.SendAllPlayerDataToNewScene()

			return nil

		case protocol.MemoryScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.MemoryScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return err
			}

			r.Broadcast(msg)
			r.SendAllPlayerDataToNewScene()

			return nil

		case protocol.WizardsScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.WizardsScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return err
			}

			wizServer.Init(r)
			r.Broadcast(msg)
			r.SendAllPlayerDataToNewScene()

			return nil
		}
	} else {
		switch r.Scene.Type {
		case protocol.WizardsScene:
			err := wizServer.ValidateMsg(message, r.Scene.State.(*wizServer.WizardServer))
			if err != nil {
				return err
			}
			return nil

		case protocol.LobbyScene:
			err := lobby.ValidateMsg(message, r.Scene.State.(*lobby.LobbyServer))
			if err != nil {
				return err
			}
			return nil

		case protocol.MemoryScene:
			err := lobby.ValidateMsg(message, r.Scene.State.(*lobby.LobbyServer))
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}
