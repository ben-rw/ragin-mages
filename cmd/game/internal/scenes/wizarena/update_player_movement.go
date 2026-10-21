package wizarena

import (
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) UpdatePlayerMovementAndCheckCollisions() {
	w.wizard.X += w.wizard.Dx * w.wizard.Combat.MoveSpeed()
	if !w.wizard.Noclip {
		wizServer.CheckCollisionHorizontal(w.wizard.Body, w.colliders)
	}

	w.wizard.Y += w.wizard.Dy * w.wizard.Combat.MoveSpeed()
	if !w.wizard.Noclip {
		wizServer.CheckCollisionVertical(w.wizard.Body, w.colliders)
	}

	for _, wizard := range w.wizards {
		wizard.NameTag.X = wizard.X + wizServer.TileSize/2
		wizard.NameTag.Y = wizard.Y + wizServer.TileSize + 2
	}

	// send w.wizard location to server
	if !w.wizard.Combat.Dead {
		w.WriteWizardMovementUpdate()
	}

	w.CheckChestCollisions()

	if w.trapsUp && w.wizard.Combat.IFrames() == 0 {
		if wizServer.CheckPointInRect(
			int(w.wizard.X+wizServer.HalfTile),
			int(w.wizard.Y+14), // puts hitbox close to feet
			w.traps,
		) {
			w.wizard.Combat.Damage(TrapDamage)
			if w.wizard.Combat.Health() <= 0 && w.wizard.Combat.LastHitBy != "" {
				pScale, pSpeed, kb, ac := StealStats(w.wizards[w.wizard.Combat.LastHitBy].Combat, w.wizard.Combat)
				w.WriteGivePoints(w.wizard.Combat.LastHitBy, int(KillPlayerBoost), pScale, pSpeed, kb, ac)
			}
			w.WriteWizardStatUpdate()
		}
	}
	w.UpdateTraps()

}

func (w *WizArena) checkIfFell(wizard *WizardPlayer) bool {
	// check if wizard fell in a hole
	fell := wizServer.CheckPointInRect(
		int(wizard.X+wizServer.HalfTile+wizard.KnockbackDx),
		int(wizard.Y+14+wizard.KnockbackDy), // puts hitbox close to feet
		w.holes,
	)
	if wizard.X+wizard.KnockbackDx < 0 ||
		wizard.Y+wizard.KnockbackDy < 0 ||
		wizard.X+wizard.KnockbackDx > float64(w.tilemapJSON.Layers[0].Width)*16.0 ||
		wizard.Y+wizard.KnockbackDy > float64(w.tilemapJSON.Layers[0].Height)*16.0 {
		fell = true
	}

	return fell
}
