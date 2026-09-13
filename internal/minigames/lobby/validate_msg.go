package lobby

import (
	// "errors"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"log"
)

func ValidateMsg(msg *protocol.Message) (*protocol.Message, error) {
	var data any
	if msg.Type == protocol.PlayerUpdate {
		// playerUpdateData, err := msg.UnmarshalMessageData()
		// if err != nil {
		// 	return &protocol.Message{}, err
		// }
		// data, ok := playerUpdateData.(*protocol.PlayerUpdateData)
		// if !ok {
		// 	log.Println(errors.New("couldn't validate message: failed type assertion"))
		// 	return &protocol.Message{}, nil
		// }
		// if data.PlayerData.Name != "" {
		// log.Printf("validated msg data: %v, scene: %v", data, scene)
		return msg, nil
		// }
	}

	log.Printf("couldn't validate msg data: %v", data)
	return &protocol.Message{Type: protocol.Unset}, nil
}
