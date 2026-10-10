package wizServer

import (
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizardServer) ValidateMsg(msg *protocol.Message) (*protocol.Message, error) {
	switch msg.Type {
	case protocol.EnemyMovementUpdate:
		w.handleEnemyMovementUpdate(msg)
	case protocol.EnemyStatUpdate:
		w.handleEnemyStatUpdate(msg)
	case protocol.NewProjectile:
		w.handleNewProjectile(msg)
	case protocol.ProjectileHitEnemy:
		w.handleProjectileHitEnemy(msg)
	case protocol.ProjectileHitWizard:
		w.handleProjectileHitWizard(msg)
	case protocol.ProjectileReflected:
		w.handleProjectileReflected(msg)
	case protocol.WizArenaMovement:
		w.handleWizArenaMovementUpdate(msg)
	case protocol.WizardMovementUpdate:
		w.handleWizArenaMovementUpdate(msg)
	case protocol.WizardStatUpdate:
		w.handleWizardStatUpdate(msg)
	case protocol.ChestUpdate:
		w.handleChestUpdate(msg)
	default:
		log.Printf("unexpected protocol.MessageType: %#v", msg.Type)
	}

	return &protocol.Message{Type: protocol.Unset}, nil
}
