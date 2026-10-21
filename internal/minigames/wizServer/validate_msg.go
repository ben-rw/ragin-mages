package wizServer

import (
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func ValidateMsg(msg *protocol.Message, w *WizardServer) error {
	switch msg.Type {
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
	case protocol.ProjectileDestroyed:
		w.handleProjectileDestroyed(msg)
	case protocol.WizardMovementUpdate:
		w.handleWizardMovementUpdate(msg)
	case protocol.WizardStatUpdate:
		w.handleWizardStatUpdate(msg)
	case protocol.ChestUpdate:
		w.handleChestUpdate(msg)
	case protocol.ClientLoaded:
		w.handleClientLoaded(msg)
	case protocol.Reflecting:
		w.handleReflecting(msg)
	case protocol.GivePoints:
		w.handleGivePoints(msg)
	default:
		log.Printf("unexpected protocol.MessageType: %#v\n", msg.Type)
	}

	return nil
}
