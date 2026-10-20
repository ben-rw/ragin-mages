package wizarena

import (
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizArena) startNewRound() {
	w.round += 1
	for _, wizard := range w.wizards {
		wizard.Alpha = 1.0
		wizard.Combat.Dead = false
		wizard.Combat.Fell = false
		w.camera.DoneMoving = false
		wizard.Combat.SetHealth(PlayerHealth)
		wizard.Combat.SetIFrames(0)
		wizard.KnockbackDx = 0
		wizard.KnockbackDy = 0
	}

	clear(w.openedChests)
	w.enemies = make(map[uint16]*Enemy, 32)
	w.projectiles = make(map[uint16]*Projectile, 16)

	w.backgroundAlpha = 0

	w.phase = wizServer.Playing
}

func (w *WizArena) updateRoundStart() {
	w.scoreboardUpdated = false

	for _, wizard := range w.wizards {
		wizard.Dx = 0
		wizard.Dy = 0

		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		wizard.ActiveAnimation.Update()

		wizard.NameTag.X = wizard.X + wizServer.TileSize/2
		wizard.NameTag.Y = wizard.Y + wizServer.TileSize + 2
	}

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
}

func (w *WizArena) updateRoundOver() {
	if !w.scoreboardUpdated {
		w.NewScoreboard()
		w.scoreboardUpdated = true
	}

	if w.wizard.Combat.Dead {
		if !w.camera.DoneMoving {
			// keep following player to position knockback puts them
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

			tolerance := 0.1
			if w.camera.X-w.wizard.DeathX < tolerance &&
				w.camera.Y-w.wizard.DeathY < tolerance {
				w.camera.DoneMoving = true
			}
		}
	}

	for _, wizard := range w.wizards {
		// add a fade effect when players die
		if wizard.Combat.Dead {
			if !wizard.Combat.Fell {
				if wizard.Y > -wizServer.HalfTile {
					wizard.Y -= 3.0
				}
			}
		}
		if wizard.Alpha > 0.0 {
			wizard.Alpha -= .02 //fade speed
		} else {
			wizard.Alpha = 0.0
		}

		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		wizard.ActiveAnimation.Update()

		wizard.NameTag.X = wizard.X + wizServer.TileSize/2
		wizard.NameTag.Y = wizard.Y + wizServer.TileSize + 2
	}

	for _, enemy := range w.enemies {
		if enemy.Combat.Dead {
			enemy.DieAnim = true
		}
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()

		enemy.X += enemy.Dx / 2 // slow motion woahhhh
		enemy.Y += enemy.Dy / 2

		if enemy.Alpha > 0.0 {
			enemy.Alpha -= .01 //fade speed
		}

	}

	for _, projectile := range w.projectiles {
		projectile.ActiveAnimation.Update()

		projectile.X += projectile.Dx / 2 // slow motion woahhhh
		projectile.Y += projectile.Dy / 2

		if projectile.Alpha > 0.0 {
			projectile.Alpha -= .01 //fade speed
		}
	}

	for _, wizard := range w.wizards {
		wizard.Reflect.ActiveAnimation = wizard.Reflect.GetActiveAnimation()

		if wizard.Reflect.ActiveAnimation != nil {
			wizard.Reflect.ActiveAnimation.Update()

			wizard.Reflect.X = wizard.X + wizServer.HalfTile
			wizard.Reflect.Y = wizard.Y + wizServer.HalfTile
		}
	}

	w.UpdateTraps()
}

func (w *WizArena) gameStart() {
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

	w.backgroundAlpha = 255
}

func (w *WizArena) gameOver() {
	if time.Until(w.phaseTimer) <= 2*time.Second {
		_ = sound.FadeOut(w.audioPlayer)
	}

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.Conn.WriteMsg(
			protocol.SceneChange,
			&protocol.SceneChangeData{
				SceneType: protocol.LobbyScene,
			},
		)
	}
}
