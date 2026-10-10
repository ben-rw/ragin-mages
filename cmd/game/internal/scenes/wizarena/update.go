package wizarena

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (w *WizArena) Update(messages []*protocol.Message) error {
	w.CheckMessages(messages)

	switch w.phase {
	case wizServer.Playing:
		w.updatePlaying()
	case wizServer.RoundOver:
		w.updateRoundOver()
	case wizServer.RoundStart:
		w.updateRoundStart()
	case wizServer.GameOver:
		w.gameOver()
	case wizServer.GameStart:
		w.gameStart()
	}

	return nil
}

func (w *WizArena) updatePlaying() {
	w.wizard.Dx = 0
	w.wizard.Dy = 0

	if w.wizard.DieAnim {
		// add a fade effect when players die
		if !w.wizard.Combat.Fell {
			if w.wizard.Y > -wizServer.HalfTile {
				w.wizard.Y -= 3.0
			}
		}
		if w.wizard.Alpha > 0.0 {
			w.wizard.Alpha -= .02 //fade speed
		} else {
			w.wizard.Alpha = 0.0
		}

	} else if w.wizard.Combat.Dead {
		// define camera behavior for dead players
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			w.camera.X -= shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			w.camera.X += shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			w.camera.Y += shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			w.camera.Y -= shared.FreeCamSpeed
		}

		w.camera.Constrain(
			float64(w.tilemapJSON.Layers[0].Width)*16.0,
			float64(w.tilemapJSON.Layers[0].Height)*16.0,
		)

	} else {
		// add velocity to wizard based on player input
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			w.wizard.Dy += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			w.wizard.Dy += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			w.wizard.Dx += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			w.wizard.Dx += -1
		}

		// normalize diagonal movement
		if w.wizard.Dx != 0 && w.wizard.Dy != 0 {
			w.wizard.Dx *= 0.70710678
			w.wizard.Dy *= 0.70710678
		}

		// define normal camera behavior
		w.camera.FollowTarget(
			w.wizard.X,
			w.wizard.Y,
			float64(w.tilemapJSON.Layers[0].Width)*16.0,
			float64(w.tilemapJSON.Layers[0].Height)*16.0,
		)
		w.camera.Constrain(
			float64(w.tilemapJSON.Layers[0].Width)*16.0,
			float64(w.tilemapJSON.Layers[0].Height)*16.0,
		)

		w.wizard.Combat.Update()
		if w.wizard.Combat.IFrames() > 0 {
			w.wizard.IFrameFlicker()
		}
		w.wizard.Noclip = false

	}

	// define enemy behavior
	w.UpdateEnemies()

	// check for enemies damaging player
	w.CheckEnemyHitWizard()

	// spawn projectiles and check projectile collisions
	w.CheckProjectiles()

	// apply knockback when hit
	if w.wizard.KnockbackDx != 0 || w.wizard.KnockbackDy != 0 {
		w.wizard.Dx += w.wizard.KnockbackDx
		w.wizard.Dy += w.wizard.KnockbackDy
		w.wizard.KnockbackDx = 0
		w.wizard.KnockbackDy = 0
	}

	// check enemy collisions
	w.CheckEnemyCollisions()

	// move player, check collisions with walls, holes, chests, and traps
	w.UpdatePlayerMovementAndCheckCollisions()

	// determine animations
	w.UpdateAnimations()

	// sort wizards, enemies, and projectiles for consistent draw ordering
	w.SortEntities()

	// toggle hitbox indicators
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		w.debug = !w.debug
	}

	// toggle FPS/TPS display
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		w.displayFPS = !w.displayFPS
	}

	// play background music
	if !w.audioPlayer.IsPlaying() {
		w.audioPlayer.SetVolume(0.2)
		w.audioPlayer.SetBufferSize(8192)
		w.audioPlayer.Play()
	}

	w.updateRound()

	w.UpdateRoundTimerText()
}
