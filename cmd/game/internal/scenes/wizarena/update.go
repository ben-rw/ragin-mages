package wizarena

import (
	"image"
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

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
					wizard := shared.NewWizard(player)
					wizard.JoinAnim = false
					w.wizards[wizard.Data.Name] = wizard
				}
			}

			// dummy test player
			w.wizards["CPU"] = shared.NewWizard(shared.NewPlayer(&protocol.PlayerData{Name: "CPU"}, 7))
			w.wizards["CPU"].X = 100
			w.wizards["CPU"].Y = 100

			w.wizard = w.wizards[w.Player.Data.Name]
			w.statText = w.NewStatText()

			w.camera = shared.NewCamera(
				-(w.wizard.X+shared.HalfTile)+shared.ScreenWidth/2.0,
				-(w.wizard.Y+shared.HalfTile)+shared.ScreenHeight/2.0,
			)

			// playerUpdateData := protocol.PlayerUpdateData{
			// 	PlayerData: w.wizard.Data,
			// }
			//
			// w.Conn.WriteMsg(protocol.PlayerUpdate, playerUpdateData)

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

	if w.wizard.DieAnim {
		// add a fade effect when players die
		if !w.wizard.Combat.Fell {
			if w.wizard.Y > -shared.HalfTile {
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

		w.camera.Constrain(
			float64(w.tilemapJSON.Layers[0].Width)*16.0,
			float64(w.tilemapJSON.Layers[0].Height)*16.0,
		)

	} else {
		// add velocity to wizard based on player input
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			w.wizard.Dy += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			w.wizard.Dy += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			w.wizard.Dx += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			w.wizard.Dx += -1
		}

		//normalize diagonal movement
		if w.wizard.Dx != 0 && w.wizard.Dy != 0 {
			w.wizard.Dx *= 0.70710678
			w.wizard.Dy *= 0.70710678
		}

		//define normal camera behavior
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
	for _, enemy := range w.enemies {
		enemy.Combat.Update()
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer {
			// make enemy follow nearest living player
			var dx, dy, dist float64
			var minDx, minDy float64
			var minDist = 10000.0 // arbitrary big number
			for _, w := range w.wizards {
				if w.Combat.Dead {
					continue
				}
				dx = w.X - enemy.X
				dy = w.Y - enemy.Y
				dist = math.Hypot(dx, dy)
				if dist < minDist {
					minDist = dist
					minDx = dx
					minDy = dy
				}
			}

			closeEnough := 8.0 //8px

			if minDist > closeEnough {
				normX := minDx / minDist
				normY := minDy / minDist

				speed := enemy.Combat.MoveSpeed()

				if enemy.Combat.MoveSpeed() > minDist {
					speed = minDist
				}

				enemy.Dx = normX * speed
				enemy.Dy = normY * speed

			}
		}
	}

	// check for enemies colliding with/damaging player
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

	deadEnemies := make(map[int]struct{})
	deadProjectiles := make(map[int]struct{})

	leftClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	rightClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	cX, cY := ebiten.CursorPosition()
	cX -= int(w.camera.X)
	cY -= int(w.camera.Y)

	// spawn new fireballs
	if leftClicked && !w.wizard.Combat.Dead && w.wizard.Combat.Attack() {
		projectile := w.wizard.ShootProjectile(
			w.projectileImageCache[shared.Fireball],
			float64(cX),
			float64(cY),
			shared.Fireball,
		)
		w.projectiles = append(w.projectiles, projectile)

		// figure out attack direction for animation
		vX := float64(cX) - w.wizard.X
		vY := float64(cY) - w.wizard.Y
		vlen := math.Hypot(vX, vY)
		if vlen == 0 {
			projectile.Despawn() // prevents division by zero
		}
		normX := vX / vlen
		normY := vY / vlen

		if math.Abs(normX) > math.Abs(normY) {
			if normX > 0 {
				w.wizard.AttackDirection = shared.AttackingRight
			} else {
				w.wizard.AttackDirection = shared.AttackingLeft
			}
		} else {
			if normY > 0 {
				w.wizard.AttackDirection = shared.AttackingDown
			} else {
				w.wizard.AttackDirection = shared.AttackingUp
			}
		}
	}

	// spawn new reflects
	if rightClicked && !w.wizard.Combat.Dead && w.wizard.Combat.Reflect() {
		w.wizard.ReflectAnim = true
		w.wizard.Reflect.Active = true
	}

	// check fireball collisions
	for i, projectile := range w.projectiles {
		projectile.Update()
		log.Println(projectile.Caster.Data.Name)
		for _, wizard := range w.wizards {
			if wizard == projectile.Caster {
				continue
			}
			if wizard.Combat.IFrames() > 0 {
				continue
			}
			if _, ok := projectile.AlreadyHit[wizard]; ok {
				continue
			}

			var wizardRadius = w.wizard.HurtboxRadius
			if wizard.Combat.Reflecting() {
				wizardRadius *= 2
			}

			if shared.CheckCollisionCircle(
				projectile.X+projectile.HitboxOffsetX,
				projectile.Y+projectile.HitboxOffsetY,
				projectile.ScaledRadius,
				wizard.X+shared.HalfTile,
				wizard.Y+shared.HalfTile,
				wizardRadius,
			) {
				if wizard.Combat.Reflecting() {
					wizard.ReflectProjectile(projectile, float64(cX), float64(cY))
					wizard.Combat.ResetReflectCooldown()
				} else if wizard.Data.Name == "CPU" {
					wizard.ReflectProjectile(projectile, 200, 100)
				} else {
					wizard.Combat.Damage(projectile.Damage)
					projectile.AlreadyHit[wizard] = struct{}{}

					w.wizard.Dx = projectile.NormX * shared.TileSize * projectile.Knockback
					w.wizard.Dy = projectile.NormY * shared.TileSize * projectile.Knockback
					w.wizard.Noclip = true

					if wizard.Combat.Health() <= 0 {
						// player who last hit the player gets a stat boost
						projectile.Caster.Combat.RandomBoost(shared.KillPlayerBoost)
						w.statText = w.NewStatText()
					}
				}

			}
		}
		for i, enemy := range w.enemies {
			if _, ok := projectile.AlreadyHit[enemy]; ok {
				continue
			}

			if shared.CheckCollisionCircle(
				projectile.X+projectile.HitboxOffsetX,
				projectile.Y+projectile.HitboxOffsetY,
				projectile.ScaledRadius,
				enemy.X+shared.HalfTile,
				enemy.Y+shared.HalfTile,
				enemy.HurtboxRadius,
			) {
				enemy.Combat.Damage(projectile.Damage)
				projectile.AlreadyHit[enemy] = struct{}{}

				enemy.Dx = projectile.NormX * shared.TileSize * projectile.Knockback
				enemy.Dy = projectile.NormY * shared.TileSize * projectile.Knockback

				if enemy.Combat.Health() <= 0 {
					deadEnemies[i] = struct{}{}
					// player who last hit the enemy gets a stat boost
					projectile.Caster.Combat.RandomBoost(shared.KillEnemyBoost)
					w.statText = w.NewStatText()
				}
			}
		}
		for j, otherProjectile := range w.projectiles {
			if projectile == otherProjectile {
				continue
			}

			if shared.CheckCollisionCircle(
				projectile.X+projectile.HitboxOffsetX,
				projectile.Y+projectile.HitboxOffsetY,
				projectile.ScaledRadius,
				otherProjectile.X+otherProjectile.HitboxOffsetX,
				otherProjectile.Y+otherProjectile.HitboxOffsetY,
				otherProjectile.ScaledRadius,
			) {
				if projectile.Scale > otherProjectile.Scale {
					deadProjectiles[j] = struct{}{}
				} else if projectile.Scale < otherProjectile.Scale {
					deadProjectiles[i] = struct{}{}
				} else {
					deadProjectiles[j] = struct{}{}
					deadProjectiles[i] = struct{}{}
				}
			}
		}

		if projectile.TicksToLive < 1 {
			deadProjectiles[i] = struct{}{}
		}
	}

	// despawn dead entities
	if len(deadProjectiles) > 0 {
		newProjectiles := make([]*shared.Projectile, 0)
		for i, projectile := range w.projectiles {
			if _, ok := deadProjectiles[i]; !ok {
				newProjectiles = append(newProjectiles, projectile)
			}
		}
		w.projectiles = newProjectiles
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

	// check enemy collisions
	for _, enemy := range w.enemies {
		// lets enemies noclip until they're fully out of the hole
		if enemy.Noclip {
			if !shared.CheckCollisionHorizontal(enemy.Sprite, w.holes) &&
				!shared.CheckCollisionVertical(enemy.Sprite, w.holes) {
				enemy.Noclip = false
			}
		}

		// turn noclip back on if enemy gets knocked into hole
		if shared.CheckPointInRect(
			int(enemy.X+shared.HalfTile),
			int(enemy.Y+shared.HalfTile),
			w.holes,
		) {
			enemy.Noclip = true
		}

		// fade in and ramp movement speed after spawn
		if enemy.Alpha < 1.0 {
			enemy.Alpha += 0.01
		}
		if enemy.Combat.MoveSpeed() < shared.EnemyMoveSpeed {
			enemy.Combat.SetMoveSpeed(enemy.Combat.MoveSpeed() + shared.EnemyMoveSpeed/150)
		}

		enemy.X += enemy.Dx
		if !enemy.Noclip {
			shared.CheckCollisionHorizontal(enemy.Sprite, w.holes)
		}
		enemy.Y += enemy.Dy
		if !enemy.Noclip {
			shared.CheckCollisionVertical(enemy.Sprite, w.holes)
		}
	}

	// check player collisions
	w.wizard.X += w.wizard.Dx * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.X = w.wizard.X + shared.TileSize/2
	if !w.wizard.Noclip {
		shared.CheckCollisionHorizontal(w.wizard.Sprite, w.colliders)
	}

	w.wizard.Y += w.wizard.Dy * w.wizard.Combat.MoveSpeed()
	w.wizard.NameTag.Y = w.wizard.Y + shared.TileSize + 2
	if !w.wizard.Noclip {
		shared.CheckCollisionVertical(w.wizard.Sprite, w.colliders)
	}

	w.CheckChestCollisions()

	if w.trapsUp && w.wizard.Combat.IFrames() == 0 {
		if shared.CheckPointInRect(
			int(w.wizard.X+shared.HalfTile),
			int(w.wizard.Y+14), // puts hitbox close to feet
			w.traps,
		) {
			w.wizard.Combat.Damage(1.0)
		}
	}
	w.UpdateTraps()

	// check if wizard fell
	w.wizard.Combat.Fell = shared.CheckPointInRect(
		int(w.wizard.X+shared.HalfTile),
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

	// determine animations
	if w.wizard.Combat.Attacking() {
		w.wizard.AttackAnim = true //starts attack animation
	}

	if w.wizard.Combat.Health() <= 0 && !w.wizard.Combat.Dead {
		w.wizard.Combat.Dead = true
		w.wizard.DieAnim = true //starts death animation
	}

	w.wizard.ActiveAnimation = w.wizard.GetActiveAnimation()
	w.wizard.ActiveAnimation.Update()

	for _, wizard := range w.wizards {
		log.Println(w.wizard.Data.Name, w.wizard.Combat.Reflecting())
		wizard.Reflect.ActiveAnimation = wizard.Reflect.GetActiveAnimation()

		if wizard.Reflect.ActiveAnimation != nil {
			wizard.Reflect.ActiveAnimation.Update()

			wizard.Reflect.X = wizard.X + shared.HalfTile
			wizard.Reflect.Y = wizard.Y + shared.HalfTile
		}
	}

	for _, enemy := range w.enemies {
		if enemy.Combat.Attacking() {
			enemy.AttackAnim = true
		}
		if enemy.Combat.Dead {
			enemy.DieAnim = true
		}
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()
	}

	for _, projectile := range w.projectiles {
		projectile.ActiveAnimation.Update()
	}

	// spawn timers
	w.EnemyRespawn()
	w.ChestRespawn()

	// toggle hitbox indicators
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		w.debug = !w.debug
	}

	// play background music
	if !w.audioPlayer.IsPlaying() {
		w.audioPlayer.SetVolume(0.2)
		w.audioPlayer.SetBufferSize(8192)
		w.audioPlayer.Play()
	}

	//TODO: when 1 player is left, start a new round
	// while preserving stat boosts.
	// at the end of the third round, announce the winner
	// and reload into lobby.

	return nil
}
