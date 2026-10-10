package wizarena

import (
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (w *WizArena) CheckProjectiles() {
	// check for mouse input
	leftClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	rightClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	cX, cY := ebiten.CursorPosition()
	cX -= int(w.camera.X)
	cY -= int(w.camera.Y)

	// spawn new fireballs
	if leftClicked && !w.wizard.Combat.Dead && w.wizard.Combat.Attack() {
		projectile := w.wizard.NewProjectile(
			w.projectileImageCache[Fireball],
			float64(cX),
			float64(cY),
			Fireball,
		)
		w.WriteNewProjectile(cX, cY)
		w.projectiles[projectile.ID] = projectile

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
		w.WriteReflecting(cX, cY)
	}

	// check fireball collisions
	for _, projectile := range w.projectiles {
		projectile.Update()
		if projectile.Caster == w.wizard {
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

				if CheckCollisionCircle(
					projectile.X+projectile.HitboxOffsetX,
					projectile.Y+projectile.HitboxOffsetY,
					projectile.ScaledRadius,
					wizard.X+wizServer.HalfTile,
					wizard.Y+wizServer.HalfTile,
					wizardRadius,
				) {
					if wizard.Combat.Reflecting() {
						wizard.ReflectProjectile(projectile, float64(cX), float64(cY))
						w.WriteProjectileReflected(wizard.Data.Username, projectile.ID)
						wizard.Combat.ResetReflectCooldown()
					} else {
						wizard.Combat.Damage(projectile.Damage)
						projectile.AlreadyHit[wizard] = struct{}{}
						StealStats(projectile.Caster.Combat, wizard.Combat)

						wizard.Dx = projectile.NormX * wizServer.TileSize * projectile.Knockback
						wizard.Dy = projectile.NormY * wizServer.TileSize * projectile.Knockback
						wizard.Noclip = true

						w.WriteProjectileHitWizard(wizard, wizard.Dx, wizard.Dy)

						if wizard.Combat.Health() <= 0 {
							// player who last hit the player gets a stat boost
							projectile.Caster.Combat.RandomBoost(KillPlayerBoost, StandardMult)
							w.WriteWizardStatUpdate()
						}

						w.dynamicText = w.NewDynamicTextMap()
					}

				}
			}
			for _, enemy := range w.enemies {
				if _, ok := projectile.AlreadyHit[enemy]; ok {
					continue
				}

				if CheckCollisionCircle(
					projectile.X+projectile.HitboxOffsetX,
					projectile.Y+projectile.HitboxOffsetY,
					projectile.ScaledRadius,
					enemy.X+wizServer.HalfTile,
					enemy.Y+wizServer.HalfTile,
					enemy.HurtboxRadius,
				) {
					enemy.Combat.Damage(projectile.Damage)
					projectile.AlreadyHit[enemy] = struct{}{}

					enemy.Dx = projectile.NormX * wizServer.TileSize * projectile.Knockback
					enemy.Dy = projectile.NormY * wizServer.TileSize * projectile.Knockback

					if enemy.Combat.Health() <= 0 {
						delete(w.enemies, enemy.ID)
						// player who last hit the enemy gets a stat boost
						projectile.Caster.Combat.RandomBoost(KillEnemyBoost, StandardMult)
						w.WriteWizardStatUpdate()
						w.dynamicText = w.NewDynamicTextMap()
					}
				}
			}
			for _, otherProjectile := range w.projectiles {
				if projectile == otherProjectile {
					continue
				}

				if CheckCollisionCircle(
					projectile.X+projectile.HitboxOffsetX,
					projectile.Y+projectile.HitboxOffsetY,
					projectile.ScaledRadius,
					otherProjectile.X+otherProjectile.HitboxOffsetX,
					otherProjectile.Y+otherProjectile.HitboxOffsetY,
					otherProjectile.ScaledRadius,
				) {
					if projectile.Scale > otherProjectile.Scale {
						delete(w.projectiles, otherProjectile.ID)
					} else if projectile.Scale < otherProjectile.Scale {
						delete(w.projectiles, projectile.ID)
					} else {
						delete(w.projectiles, projectile.ID)
						delete(w.projectiles, otherProjectile.ID)
					}
				}
			}
		}

		if projectile.TicksToLive < 1 {
			delete(w.projectiles, projectile.ID)
		}
	}
}
