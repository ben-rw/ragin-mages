package wizarena

import (
	"image"
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
)

func (w *WizArena) CheckEnemyHitWizard() {
	wizardRect := image.Rect(
		int(w.wizard.X),
		int(w.wizard.Y),
		int(w.wizard.X)+shared.TileSize,
		int(w.wizard.Y)+shared.TileSize,
	)

	for _, enemy := range w.enemies {
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+shared.TileSize,
			int(enemy.Y)+shared.TileSize,
		)

		if rect.Overlaps(wizardRect) && w.wizard.Combat.IFrames() == 0 {
			if enemy.Combat.Attack() {
				w.wizard.Combat.Damage(enemy.Combat.AttackPower())

				// player pushed away by enemy
				// find vector, divide by vector length, add to player's velocity
				vX := w.wizard.X - enemy.X
				vY := w.wizard.Y - enemy.Y
				vlen := math.Sqrt(vX*vX + vY*vY)
				normX := vX / vlen
				normY := vY / vlen

				w.wizard.Dx = normX * shared.TileSize * enemy.Combat.Knockback()
				w.wizard.Dy = normY * shared.TileSize * enemy.Combat.Knockback()
				w.wizard.Noclip = true

				if math.Abs(normX) > math.Abs(normY) {
					if normX > 0 {
						enemy.AttackDirection = shared.AttackingRight
					} else {
						enemy.AttackDirection = shared.AttackingLeft
					}
				} else {
					if normY > 0 {
						enemy.AttackDirection = shared.AttackingDown
					} else {
						enemy.AttackDirection = shared.AttackingUp
					}
				}
			}
		}
	}
}
