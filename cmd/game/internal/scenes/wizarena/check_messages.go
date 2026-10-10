package wizarena

import (
	"log"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizArena) CheckMessages(messages []*protocol.Message) {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := w.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, player := range w.Players {
				if _, ok := w.wizards[player.Data.Username]; !ok {
					wizard := NewWizard(player)
					wizard.JoinAnim = false
					w.wizards[wizard.Data.Username] = wizard
				}
			}

			w.wizard = w.wizards[w.Player.Data.Username]
			w.dynamicText = w.NewDynamicTextMap()

			w.camera = shared.NewCamera(
				-(w.wizard.X+wizServer.HalfTile)+shared.ScreenWidth/2.0,
				-(w.wizard.Y+wizServer.HalfTile)+shared.ScreenHeight/2.0,
			)

		case protocol.PlayerUpdate:
			err := w.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		case protocol.WizardMovementUpdate:
			err := w.HandleWizardMovementUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.EnemyMovementUpdate:
			err := w.HandleEnemyMovementUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.EnemyStatUpdate:
			err := w.HandleEnemyStatUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.WizArenaMovement:
			err := w.HandleWizArenaMovementUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.WizardStatUpdate:
			err := w.HandleWizardStatUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.NewProjectile:
			err := w.HandleNewProjectile(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.ProjectileReflected:
			err := w.HandleProjectileReflected(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.ProjectileHitWizard:
			err := w.HandleProjectileHitWizard(message)
			if err != nil {
				log.Println(err)
				continue
			}
		case protocol.ProjectileHitEnemy:
			err := w.HandleProjectileHitEnemy(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}
}
