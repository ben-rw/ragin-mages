package wizards

import (
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

func ValidateMsg(msg *protocol.Message) (*protocol.Message, error) {
	switch msg.Type {
	case protocol.EnemyMovementUpdate:
	case protocol.EnemyStatUpdate:
	case protocol.NewProjectile:
	case protocol.PlayerUpdate:
	case protocol.ProjectileHitEnemy:
	case protocol.ProjectileHitWizard:
	case protocol.ProjectileReflected:
	case protocol.WizArenaMovementState:
	case protocol.WizardMovementUpdate:
	case protocol.WizardStatUpdate:
	default:
		log.Printf("unexpected protocol.MessageType: %#v", msg.Type)
	}

	return &protocol.Message{Type: protocol.Unset}, nil
}

type WizardServer struct {
}

func NewWizardServer(room *room.Room) {

}

func (w *WizardServer) Update() {

	// switch w.phase {
	// case Playing:
	//
	//	w.updatePlaying()
	//
	// case RoundOver:
	//
	//	w.updateRoundOver()
	//
	// case RoundStart:
	//
	//	w.updateRoundStart()
	//
	// case GameOver:
	//
	//	w.gameOver()
	//
	// case GameStart:
	//
	//		w.gameStart()
	//	}
}

func (w *WizardServer) updatePlaying() {
	panic("unimplemented")
}

func (w *WizardServer) updateRoundOver() {
	panic("unimplemented")
}

func (w *WizardServer) updateRoundStart() {
	panic("unimplemented")
}

func (w *WizardServer) gameOver() {
	panic("unimplemented")
}

func (w *WizardServer) gameStart() {
	panic("unimplemented")
}
