package wizarena

import (
	"image"
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) CheckEnemyHitWizard() {
	wizardRect := image.Rect(
		int(w.wizard.X),
		int(w.wizard.Y),
		int(w.wizard.X)+wizServer.TileSize,
		int(w.wizard.Y)+wizServer.TileSize,
	)

	for _, enemy := range w.enemies {
		enemy.Combat.Update()

		if enemy.Combat.Health() <= 0 {
			enemy.Dx = 0
			enemy.Dy = 0

			if enemy.DieAnim {
				// add a fade effect when enemies die
				if enemy.Y > -wizServer.HalfTile {
					enemy.Y -= 3.0
				}
				if enemy.Alpha > 0.0 {
					enemy.Alpha -= .02 //fade speed
				} else {
					enemy.Alpha = 0.0
				}
			}

			if enemy.DieAnim == false {
				delete(w.enemies, enemy.ID)
				continue
			}
		}

		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+wizServer.TileSize,
			int(enemy.Y)+wizServer.TileSize,
		)

		if rect.Overlaps(wizardRect) && w.wizard.Combat.IFrames() == 0 && enemy.DieAnim == false {
			if enemy.Combat.Attack() {
				w.wizard.Combat.Damage(enemy.Combat.AttackPower())
				// w.wizard.Combat.LastHitBy = "" // prevent a player from getting credit for an enemy's kill
				if w.wizard.Combat.Health() <= 0 && w.wizard.Combat.LastHitBy != "" {
					pScale, pSpeed, kb, ac := StealStats(w.wizards[w.wizard.Combat.LastHitBy].Combat, w.wizard.Combat)
					w.WriteGivePoints(w.wizard.Combat.LastHitBy, int(KillPlayerBoost), pScale, pSpeed, kb, ac, false)
				}

				w.WriteWizardStatUpdate()

				// player pushed away by enemy
				// find vector, divide by vector length, add to player's velocity
				vX := w.wizard.X - enemy.X
				vY := w.wizard.Y - enemy.Y
				vlen := math.Sqrt(vX*vX + vY*vY)
				normX := vX / vlen
				normY := vY / vlen

				w.wizard.Dx = normX * wizServer.TileSize * enemy.Combat.Knockback()
				w.wizard.Dy = normY * wizServer.TileSize * enemy.Combat.Knockback()
				w.wizard.Noclip = true
				w.WriteWizardMovementUpdate()

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
