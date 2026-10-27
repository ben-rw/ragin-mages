package wizarena

import (
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizArena) startNewRound() {
	// wizard state is reset in HandlePhaseUpdate
	w.round += 1
	w.camera.DoneMoving = false

	clear(w.openedChests)
	w.enemies = make(map[uint16]*Enemy, 32)
	w.projectiles = make(map[uint16]*Projectile, 16)

	w.backgroundAlpha = 0

	w.phase = wizServer.Playing
}

func (w *WizArena) updateRoundStart() {
	for _, wizard := range w.wizards {
		wizard.Dx = 0
		wizard.Dy = 0

		w.updateRoundEndAnimations()
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

func (w *WizArena) syncScores() {
	w.updateRoundEndAnimations()

	if !w.scoresUpdated {
		w.waitForScores()
	} else if w.scoresUpdated && !w.scoreboardUpdated {
		w.NewScoreboard()
		w.scoreboardUpdated = true
	}
}

func (w *WizArena) updateRoundOver() {
	w.scoreboardUpdated = false
	w.scoresUpdated = false

	w.updateRoundEndAnimations()
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

func (w *WizArena) waitForScores() {
	if int(w.scoresUpdatedCount.Load()) >= (len(w.wizards) - 1) { // excludes this client's wizard
		w.scoresUpdated = true
	}
}

func (w *WizArena) updateRoundEndAnimations() {
	if w.wizard.Combat.Dead {
		if !w.camera.DoneMoving {
			// keep following player to position knockback puts them
			w.FollowDeathKnockback()
		}
	}

	for _, wizard := range w.wizards {
		wizard.Dx = 0
		wizard.Dy = 0

		// add a fade effect when players die
		if wizard.Combat.Dead {
			if !wizard.Combat.Fell {
				if wizard.Y > -wizServer.HalfTile {
					wizard.Y -= 3
				}
			}
			if wizard.Alpha > 0.0 {
				wizard.Alpha -= .02 //fade speed
			} else {
				wizard.Alpha = 0.0
			}
		} else {
			wizard.Combat.Update()
			if wizard.Combat.IFrames() > 0 {
				wizard.IFrameFlicker()
			}
		}

		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		wizard.ActiveAnimation.Update()

		wizard.NameTag.X = wizard.X + wizServer.TileSize/2
		wizard.NameTag.Y = wizard.Y + wizServer.TileSize + 2
	}

	for _, enemy := range w.enemies {
		if enemy.Combat.Health() <= 0 {
			enemy.Combat.Dead = true
			enemy.DieAnim = true

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
			}
		}
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()

		enemy.X += enemy.Dx / 2 // slow motion woahhhh
		enemy.Y += enemy.Dy / 2

		if enemy.Alpha > 0.0 && enemy.DieAnim == false {
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

	// sort wizards, enemies, and projectiles for consistent draw ordering
	w.SortEntities()
}
