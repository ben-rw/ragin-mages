package wizarena

import (
	"fmt"
	"image/color"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (w *WizArena) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range w.tilemapJSON.Layers {
		if layer.Name == "traps_up" && !w.trapsUp {
			continue
		}
		if layer.Name == "traps_down" && w.trapsUp {
			continue
		}
		for _, tile := range w.tiles[layer.Name] {
			if layer.Name == "chests" {
				if _, ok := w.openedChests[OpenedChest{tile.X, tile.Y}]; ok {
					continue
				}
			}

			if tile.Flips.HorizontalFlip ||
				tile.Flips.VerticalFlip ||
				tile.Flips.DiagonalFlip {
				shared.FixRotatedTile(tile, &opts)
			}

			opts.GeoM.Translate(float64(tile.X), float64(tile.Y))
			opts.GeoM.Translate(w.camera.X, w.camera.Y)

			screen.DrawImage(
				tile.Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	for _, wizard := range w.wizards {
		opts.GeoM.Translate(wizard.X, wizard.Y)

		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		if wizard.Alpha < 1.0 {
			opts.ColorScale.ScaleAlpha(wizard.Alpha)
		}

		screen.DrawImage(
			w.wizardFrameCache[w.wizard.Data.SpriteIndex][w.wizard.ActiveAnimation.Frame()],
			&opts,
		)

		opts.ColorScale.Reset()
		opts.GeoM.Reset()
	}

	for _, enemy := range w.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)

		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		if enemy.Alpha < 1.0 {
			opts.ColorScale.ScaleAlpha(enemy.Alpha)
		}

		screen.DrawImage(
			w.enemyFrameCache[enemy.Type][enemy.ActiveAnimation.Frame()],
			&opts,
		)

		opts.ColorScale.Reset()
		opts.GeoM.Reset()
	}

	for _, projectile := range w.projectiles {

		opts.GeoM.Translate(-projectile.CenterX, -projectile.CenterY)

		opts.GeoM.Scale(projectile.Scale, projectile.Scale)
		opts.GeoM.Rotate(projectile.Rotation)

		opts.GeoM.Translate(projectile.X, projectile.Y)
		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		screen.DrawImage(
			w.projectileFrameCache[projectile.Type][projectile.ActiveAnimation.Frame()],
			&opts,
		)

		opts.GeoM.Reset()
	}

	for i := range int(w.wizard.Combat.Health()) {
		opts.GeoM.Translate(HealthHeartLocations[i]())
		screen.DrawImage(w.heartImage, &opts)

		opts.GeoM.Reset()
	}

	for _, wizard := range w.wizards {
		textOpts := text.DrawOptions{
			LayoutOptions: wizard.NameTag.LayoutOptions,
		}
		textOpts.GeoM.Translate(wizard.NameTag.X, wizard.NameTag.Y)

		textOpts.GeoM.Translate(w.camera.X, w.camera.Y)

		if wizard.Alpha < 1.0 {
			textOpts.ColorScale.ScaleAlpha(wizard.Alpha)
		}

		text.Draw(screen, wizard.Data.Name, wizard.NameTag.Face, &textOpts)

		textOpts.ColorScale.Reset()
		textOpts.GeoM.Reset()
	}

	instructionText1 := "You need more POWER!"
	textOpts := text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	}
	textOpts.GeoM.Translate(shared.TopRightFurther())
	text.Draw(screen, instructionText1, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	instructionText2 := "Get boosts from chests, skeletons, and your friends!"
	textOpts = text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignEnd,
		},
	}
	textOpts.GeoM.Translate(shared.BottomRightFurther())
	text.Draw(screen, instructionText2, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	stat1text := fmt.Sprintf("Fireball Size: %v", w.wizard.Combat.ProjectileScale())
	textOpts = text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	}
	textOpts.GeoM.Translate(shared.Stat1BottomLeft())
	text.Draw(screen, stat1text, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	stat2text := fmt.Sprintf("Fireball Speed: %v", w.wizard.Combat.ProjectileSpeed())
	textOpts = text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	}
	textOpts.GeoM.Translate(shared.Stat2BottomLeft())
	text.Draw(screen, stat2text, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	stat3text := fmt.Sprintf("F.B. Knockback: %v", w.wizard.Combat.Knockback())
	textOpts = text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignStart,
		},
	}
	textOpts.GeoM.Translate(shared.Stat3BottomLeft())
	text.Draw(screen, stat3text, &text.GoTextFace{Source: shared.FontSrc, Size: 8}, &textOpts)

	textOpts.GeoM.Reset()

	if w.wizard.Combat.Dead {
		deadText := "YOU DIED"
		textOpts = text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: text.AlignCenter,
			},
		}
		textOpts.GeoM.Translate(shared.Center())
		text.Draw(screen, deadText, &text.GoTextFace{Source: shared.FontSrc, Size: 24}, &textOpts)
	}

	textOpts.GeoM.Reset()

	if w.debug {
		for _, projectile := range w.projectiles {

			m := ebiten.GeoM{}
			m.Translate(-shared.FBCenterX, -shared.FBCenterY)
			m.Scale(projectile.Scale, projectile.Scale)
			m.Rotate(projectile.Rotation)
			m.Translate(projectile.X, projectile.Y)

			hx, hy := m.Apply(shared.FBInnerBallX, shared.FBInnerBallY)
			shared.CheckCollisionCircle(
				hx,
				hy,
				projectile.ScaledRadius,
				w.wizard.X+shared.HalfTile,
				w.wizard.Y+shared.HalfTile,
				w.wizard.HurtboxRadius,
			)

			vector.StrokeCircle(
				screen,
				float32(hx+w.camera.X),
				float32(hy+w.camera.Y),
				float32(projectile.ScaledRadius),
				1.0,
				color.RGBA{255, 0, 0, 255},
				false,
			)
		}

		for _, enemies := range w.enemies {
			vector.StrokeCircle(
				screen,
				float32(enemies.X+shared.HalfTile+w.camera.X),
				float32(enemies.Y+shared.HalfTile+w.camera.Y),
				float32(enemies.HurtboxRadius),
				1.0,
				color.RGBA{255, 0, 0, 255},
				false,
			)
		}

		// for _, collider := range w.colliders {
		// 	vector.StrokeRect(
		// 		screen,
		// 		float32(collider.Min.X)+float32(w.camera.X),
		// 		float32(collider.Min.Y)+float32(w.camera.Y),
		// 		float32(collider.Dx()),
		// 		float32(collider.Dy()),
		// 		1.0,
		// 		color.RGBA{255, 0, 0, 255},
		// 		false,
		// 	)
		// 	opts.GeoM.Reset()
		// }

		for _, hole := range w.holes {
			vector.StrokeRect(
				screen,
				float32(hole.Min.X)+float32(w.camera.X),
				float32(hole.Min.Y)+float32(w.camera.Y),
				float32(hole.Dx()),
				float32(hole.Dy()),
				1.0,
				color.RGBA{0, 0, 255, 255},
				false,
			)
		}

		for _, spike := range w.traps {
			vector.StrokeRect(
				screen,
				float32(spike.Min.X)+float32(w.camera.X),
				float32(spike.Min.Y)+float32(w.camera.Y),
				float32(spike.Dx()),
				float32(spike.Dy()),
				1.0,
				color.RGBA{0, 255, 0, 255},
				false,
			)
		}
	}
}
