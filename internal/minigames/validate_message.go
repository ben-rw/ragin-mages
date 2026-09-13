package minigames

import (
	"errors"
	"github.com/ben-rw/ragin-mages/internal/minigames/lobby"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"math/rand"
)

// lobby, random omitted as they are not games
var minigameMap = map[int]protocol.SceneType{
	0: protocol.WizardsScene,
}

func ValidateMessage(message *protocol.Message, scene protocol.SceneType) (*protocol.Message, error) {
	if message.Type == protocol.SceneChange {
		sceneChangeData, err := message.UnmarshalMessageData()
		if err != nil {
			return &protocol.Message{Type: protocol.Unset}, err
		}
		scdata, ok := sceneChangeData.(*protocol.SceneChangeData)
		if !ok {
			err = errors.New("couldn't validate message: failed type assertion")
			return &protocol.Message{Type: protocol.Unset}, err
		}

		switch scdata.SceneType {
		case protocol.RandomScene:
			data := protocol.SceneChangeData{
				SceneType: minigameMap[rand.Intn(len(minigameMap))],
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.LobbyScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.LobbyScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.MemoryScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.MemoryScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.WizardsScene:
			data := protocol.SceneChangeData{
				SceneType: protocol.WizardsScene,
			}
			msg, err := protocol.MarshalToMessage(protocol.SceneChange, data)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		}
	} else {
		switch scene {
		case protocol.LobbyScene:
			msg, err := lobby.ValidateMsg(message)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.MemoryScene:
			msg, err := lobby.ValidateMsg(message)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		case protocol.WizardsScene:
			msg, err := lobby.ValidateMsg(message)
			if err != nil {
				return &protocol.Message{Type: protocol.Unset}, err
			}
			return msg, nil
		}
	}

	return &protocol.Message{Type: protocol.Unset}, nil
}
