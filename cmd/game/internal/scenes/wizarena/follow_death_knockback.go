package wizarena

import "math"

func (w *WizArena) FollowDeathKnockback() {
	preCameraX := w.camera.X
	preCameraY := w.camera.Y
	w.camera.FollowTarget(
		w.wizard.DeathX,
		w.wizard.DeathY,
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)
	w.camera.Constrain(
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)

	tolerance := .1
	if math.Abs(w.camera.X-preCameraX) < tolerance &&
		math.Abs(w.camera.Y-preCameraY) < tolerance {
		w.camera.DoneMoving = true
	}
}
