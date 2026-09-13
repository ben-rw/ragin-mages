package wizarena

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"log"
)

func (w *WizArena) Update(messages []protocol.Message) error {
	for _, message := range messages {
		switch message.Type {
		case protocol.JoinResponse:
			err := w.HandleJoinResponse(message)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, player := range w.Players {
				if _, ok := w.wizards[player.Data.Name]; !ok {
					wizard := NewWizard(player)
					wizard.JustJoined = false
					w.wizards[wizard.Data.Name] = wizard
				}
			}

			w.wizard = w.wizards[w.Player.Data.Name]

			playerUpdateData := protocol.PlayerUpdateData{
				PlayerData: w.Player.Data,
			}

			w.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

		case protocol.PlayerUpdate:
			err := w.HandlePlayerUpdate(message)
			if err != nil {
				log.Println(err)
				continue
			}

		default:
		}
	}

	w.wizard.Dx = 0
	w.wizard.Dy = 0

	if w.wizard.Combat.Dead {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			w.camera.X -= shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			w.camera.X += shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			w.camera.Y += shared.FreeCamSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			w.camera.Y -= shared.FreeCamSpeed
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			w.wizard.Dx += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			w.wizard.Dx += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			w.wizard.Dy += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			w.wizard.Dy += 1
		}

		//normalize diagonal movement
		if w.wizard.Dx != 0 && w.wizard.Dy != 0 {
			w.wizard.Dx *= 0.7071
			w.wizard.Dy *= 0.7071
		}
	}

	for _, collider := range w.colliders {
		if collider.Overlaps(image.Rect(
			int(w.wizard.X),
			int(w.wizard.Y),
			int(w.wizard.X)+16,
			int(w.wizard.Y)+16,
		)) {
			if w.wizard.Dy > 0.0 {
				w.wizard.Y = float64(collider.Min.Y) - shared.TileSize
			} else if w.wizard.Dy < 0.0 {
				w.wizard.Y = float64(collider.Max.Y)
			}
		}
	}

	for _, wizard := range w.wizards {
		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		wizard.ActiveAnimation.Update()
	}

	for _, enemy := range w.enemies {
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer {
			tolerance := 1.0
			if enemy.X <= w.wizard.X-tolerance {
				enemy.Dx = 1
			}
			if enemy.X >= w.wizard.X+tolerance {
				enemy.Dx = -1
			}
			if enemy.Y <= w.wizard.Y-tolerance {
				enemy.Dy = 1
			}
			if enemy.Y >= w.wizard.Y+tolerance {
				enemy.Dy = -1
			}

			//normalize diagonal movement
			if enemy.Dx != 0 && enemy.Dy != 0 {
				enemy.Dx *= 0.7071
				enemy.Dy *= 0.7071
			}

			enemy.X += enemy.Dx * enemy.Combat.MoveSpeed()
			shared.CheckCollisionHorizontal(enemy.Sprite, w.colliders)
			enemy.Y += enemy.Dy * enemy.Combat.MoveSpeed()
			shared.CheckCollisionVertical(enemy.Sprite, w.colliders)
		}
	}

	for _, enemy := range w.enemies {
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()
	}

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	cX, cY := ebiten.CursorPosition()
	cX -= int(w.camera.X)
	cY -= int(w.camera.Y)
	w.wizard.Combat.Update()

	if clicked && !w.wizard.Combat.Dead {
		projectile := shared.SpawnProjectile(
			w.projectileCache[shared.Fireball],
			w.wizard.Combat.ProjectileSpeed(),
			w.wizard.Combat.ProjectileSize(),
			w.wizard.Combat.Knockback(),
			w.wizard.Combat.AttackRange(),
			w.wizard.X,
			w.wizard.Y,
			float64(cX),
			float64(cY),
		)
		w.projectiles = append(w.projectiles, projectile)
		log.Println(w.projectiles)
	}

	deadProjectiles := make(map[int]struct{})
	for i, projectile := range w.projectiles {
		projectile.Update()
		projectile.ActiveAnimation.Update()
		if projectile.TicksToLive < 1 {
			deadProjectiles[i] = struct{}{}
		}
	}

	if len(deadProjectiles) > 0 {
		newProjectiles := make([]*shared.Projectile, 0)
		for i, projectile := range w.projectiles {
			if _, ok := deadProjectiles[i]; !ok {
				newProjectiles = append(newProjectiles, projectile)
			}
		}
		w.projectiles = newProjectiles
	}

	wizardRect := image.Rect(
		int(w.wizard.X),
		int(w.wizard.Y),
		int(w.wizard.X)+shared.TileSize,
		int(w.wizard.Y)+shared.TileSize,
	)

	deadEnemies := make(map[int]struct{})
	for i, enemy := range w.enemies {
		enemy.Combat.Update()
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+shared.TileSize,
			int(enemy.Y)+shared.TileSize,
		)

		if rect.Overlaps(wizardRect) {
			if enemy.Combat.Attack() {
				w.wizard.Combat.Damage(enemy.Combat.AttackPower())
				log.Printf("wiz hp: %v\n", w.wizard.Combat.Health())

				// player pushed away by enemy
				// find vector, divide by vector length, add to player's velocity
				vX := w.wizard.X - enemy.X
				vY := w.wizard.Y - enemy.Y
				vlen := math.Sqrt(math.Pow(vX, 2) + math.Pow(vY, 2))
				normX := vX / vlen
				normY := vY / vlen

				w.wizard.Dx = normX * shared.TileSize * enemy.Combat.Knockback()
				w.wizard.Dy = normY * shared.TileSize * enemy.Combat.Knockback()

				if w.wizard.Combat.Health() <= 0 {
					w.wizard.Combat.Dead = true
					// setting Dying to true plays death anim
					w.wizard.Dying = true
				}
			}
		}

		if cX > rect.Min.X &&
			cX < rect.Max.X &&
			cY > rect.Min.Y &&
			cY < rect.Max.Y {
			if clicked {
				enemy.Combat.Damage(w.wizard.Combat.AttackPower())

				if enemy.Combat.Health() <= 0 {
					deadEnemies[i] = struct{}{}
					// player who last hit the enemy gets a stat boost
					w.wizard.Combat.RandomBoost(shared.KillEnemyBoost)
				}
			}
		}
	}
	if len(deadEnemies) > 0 {
		newEnemies := make([]*shared.Enemy, 0)
		for i, enemy := range w.enemies {
			if _, ok := deadEnemies[i]; !ok {
				newEnemies = append(newEnemies, enemy)
			}
		}
		w.enemies = newEnemies
	}

	w.wizard.X += w.wizard.Dx * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.X = w.wizard.X + shared.TileSize/2
	shared.CheckCollisionHorizontal(w.wizard.Sprite, w.colliders)

	w.wizard.Y += w.wizard.Dy * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.Y = w.wizard.Y + shared.TileSize + 2
	shared.CheckCollisionVertical(w.wizard.Sprite, w.colliders)

	if !w.wizard.Combat.Dead {
		w.camera.FollowTarget(w.wizard.X, w.wizard.Y)
	}
	w.camera.Constrain(
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)

	if !w.audioPlayer.IsPlaying() {
		w.audioPlayer.SetVolume(0.2)
		w.audioPlayer.SetBufferSize(300)
		w.audioPlayer.Play()
	}

	//TODO: when 1 player is left, start a new round
	// while preserving stat boosts.
	// at the end of the third round, announce the winner
	// and reload into lobby.

	return nil
}

func (w *WizArena) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range w.tilemapJSON.Layers {
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := i % layer.Width
			y := i / layer.Width

			x *= shared.TileSize
			y *= shared.TileSize

			tile := shared.Tile{}

			if id&int(shared.FlagFlippedHorizontally) != 0 {
				tile.Flips.HorizontalFlip = true
			}
			if id&int(shared.FlagFlippedVertically) != 0 {
				tile.Flips.VerticalFlip = true
			}
			if id&int(shared.FlagFlippedDiagonally) != 0 {
				tile.Flips.DiagonalFlip = true
			}

			id &= ^(int(shared.FlagFlippedHorizontally) |
				int(shared.FlagFlippedVertically) |
				int(shared.FlagFlippedDiagonally) |
				int(shared.FlagRotatedHexagonal120))

			switch {
			case tile.Flips.HorizontalFlip && tile.Flips.VerticalFlip:
				opts.GeoM.Scale(-1, -1)
				x += 16
				y += 16
			case tile.Flips.DiagonalFlip && tile.Flips.HorizontalFlip:
				opts.GeoM.Translate(-shared.TileSize/2, -shared.TileSize/2)
				opts.GeoM.Rotate(math.Pi / 2)
				opts.GeoM.Translate(shared.TileSize/2, shared.TileSize/2)
			case tile.Flips.HorizontalFlip:
				opts.GeoM.Scale(-1, 1)
				x += 16
			case tile.Flips.DiagonalFlip && tile.Flips.VerticalFlip:
				opts.GeoM.Translate(-shared.TileSize/2, -shared.TileSize/2)
				opts.GeoM.Rotate(3 * math.Pi / 2)
				opts.GeoM.Translate(shared.TileSize/2, shared.TileSize/2)
			case tile.Flips.VerticalFlip:
				opts.GeoM.Scale(1, -1)
				y += 16
			default:
			}

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(w.camera.X, w.camera.Y)

			screen.DrawImage(
				w.tileCache[id].Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	for _, collider := range w.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(w.camera.X),
			float32(collider.Min.Y)+float32(w.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dx()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
		opts.GeoM.Reset()
	}

	for _, wizard := range w.wizards {
		opts.GeoM.Translate(wizard.X, wizard.Y)

		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		wizard.ActiveAnimation = wizard.GetActiveAnimation()
		screen.DrawImage(
			wizard.Img.SubImage(
				wizard.SpriteSheet.Rect(wizard.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, enemy := range w.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)

		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		screen.DrawImage(
			enemy.Img.SubImage(
				enemy.SpriteSheet.Rect(enemy.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, projectile := range w.projectiles {
		opts.GeoM.Translate(projectile.CenterX, projectile.CenterY)

		opts.GeoM.Scale(0.5, 0.5)
		opts.GeoM.Rotate(projectile.Rotation)

		opts.GeoM.Translate(projectile.X, projectile.Y)
		opts.GeoM.Translate(w.camera.X, w.camera.Y)

		screen.DrawImage(
			projectile.Img.SubImage(
				projectile.SpriteSheet.Rect(projectile.ActiveAnimation.Frame()),
			).(*ebiten.Image),
			&opts,
		)
	}

	opts.GeoM.Reset()

	for i := range w.wizard.Combat.Health() {
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

		text.Draw(screen, wizard.Data.Name, wizard.NameTag.Face, &textOpts)

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

	stat1text := fmt.Sprintf("Fireball Size: %v", w.wizard.Combat.ProjectileSize())
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
}
