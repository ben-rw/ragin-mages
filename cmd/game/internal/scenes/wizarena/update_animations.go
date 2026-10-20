package wizarena

import (
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) UpdateAnimations() {
	for _, wizard := range w.wizards {
		if wizard.Combat.Attacking() {
			wizard.AttackAnim = true //starts attack animation
		}

		if wizard.Combat.Health() <= 0 && !wizard.Combat.Dead {
			wizard.Combat.Dead = true
			wizard.DeathX = wizard.X
			wizard.DeathY = wizard.Y
			wizard.DieAnim = true //starts death animation
		}

		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		wizard.ActiveAnimation.Update()

		wizard.Reflect.ActiveAnimation = wizard.Reflect.GetActiveAnimation()

		if wizard.Reflect.ActiveAnimation != nil {
			wizard.Reflect.ActiveAnimation.Update()

			wizard.Reflect.X = wizard.X + wizServer.HalfTile
			wizard.Reflect.Y = wizard.Y + wizServer.HalfTile
		}
	}

	// enemy animations
	for _, enemy := range w.enemies {
		if enemy.Combat.Attacking() {
			enemy.AttackAnim = true
		}
		if enemy.Combat.Dead {
			enemy.DieAnim = true
		}
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()
	}

	// projectile animations
	for _, projectile := range w.projectiles {
		projectile.ActiveAnimation.Update()
	}
}
