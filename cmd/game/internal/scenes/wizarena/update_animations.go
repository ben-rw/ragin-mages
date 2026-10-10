package wizarena

import (
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) UpdateAnimations() {
	if w.wizard.Combat.Attacking() {
		w.wizard.AttackAnim = true //starts attack animation
	}

	if w.wizard.Combat.Health() <= 0 && !w.wizard.Combat.Dead {
		w.wizard.Combat.Dead = true
		w.wizard.DieAnim = true //starts death animation
	}

	w.wizard.ActiveAnimation = w.wizard.GetActiveAnimation()
	w.wizard.ActiveAnimation.Update()

	for _, wizard := range w.wizards {
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
