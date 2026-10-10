package wizarena

import (
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) UpdatePlayerMovementAndCheckCollisions() {
	w.wizard.X += w.wizard.Dx * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.X = w.wizard.X + wizServer.TileSize/2
	if !w.wizard.Noclip {
		CheckCollisionHorizontal(w.wizard.Sprite, w.colliders)
	}

	w.wizard.Y += w.wizard.Dy * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.Y = w.wizard.Y + wizServer.TileSize + 2
	if !w.wizard.Noclip {
		CheckCollisionVertical(w.wizard.Sprite, w.colliders)
	}

	// send w.wizard location to server
	if !w.wizard.Combat.Dead {
		w.WriteWizardMovementUpdate()
	}

	w.CheckChestCollisions()

	if w.trapsUp && w.wizard.Combat.IFrames() == 0 {
		if CheckPointInRect(
			int(w.wizard.X+wizServer.HalfTile),
			int(w.wizard.Y+14), // puts hitbox close to feet
			w.traps,
		) {
			w.wizard.Combat.Damage(1.0)
		}
	}
	w.UpdateTraps()

	// check if wizard fell
	w.wizard.Combat.Fell = CheckPointInRect(
		int(w.wizard.X+wizServer.HalfTile),
		int(w.wizard.Y+14), // puts hitbox close to feet
		w.holes,
	)
	if w.wizard.X < 0 ||
		w.wizard.Y < 0 ||
		w.wizard.X > float64(w.tilemapJSON.Layers[0].Width)*16.0 ||
		w.wizard.Y > float64(w.tilemapJSON.Layers[0].Height)*16.0 {
		w.wizard.Combat.Fell = true
	}
	if w.wizard.Combat.Fell {
		w.wizard.Combat.SetHealth(0)
	}
}
