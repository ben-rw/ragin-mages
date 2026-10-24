package wizarena

import (
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
	}

	// spawn new reflects
	if rightClicked && !w.wizard.Combat.Dead && w.wizard.Combat.Reflect() {
		w.wizard.ReflectAnim = true
		w.wizard.Reflect.Active = true
		w.wizard.Reflect.cursorX = float64(cX)
		w.wizard.Reflect.cursorY = float64(cY)
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

				var wizardRadius = wizard.HurtboxRadius
				if wizard.Combat.Reflecting() {
					wizardRadius *= 2
				}

				if wizServer.CheckCollisionCircle(
					projectile.X+projectile.HitboxOffsetX,
					projectile.Y+projectile.HitboxOffsetY,
					projectile.ScaledRadius,
					wizard.X+wizServer.HalfTile,
					wizard.Y+wizServer.HalfTile,
					wizardRadius,
				) {
					if wizard.Combat.Reflecting() {
						w.WriteProjectileReflected(w.wizard.Data.Username, wizard.Data.Username, wizard.Reflect.cursorX, wizard.Reflect.cursorY, projectile.ID)
						wizard.ReflectProjectile(projectile)
						wizard.Combat.ResetReflectCooldown()
					} else {
						wizard.Combat.Damage(projectile.Damage)
						projectile.AlreadyHit[wizard] = struct{}{}

						wizard.KnockbackDx = projectile.NormX * wizServer.TileSize * projectile.Knockback
						wizard.KnockbackDy = projectile.NormY * wizServer.TileSize * projectile.Knockback
						wizard.Noclip = true

						if wizard.Combat.Health() <= 0 {
							// player who last hit the player gets a flat stat boost + steals stats from victim
							StealStats(projectile.Caster.Combat, wizard.Combat)
							w.wizard.Combat.RandomBoost(KillPlayerBoost, StandardMult)
						}

						w.WriteProjectileHitWizard(wizard, wizard.KnockbackDx, wizard.KnockbackDy)

						w.dynamicText = w.NewDynamicTextMap()
					}

				}
			}
			for _, enemy := range w.enemies {
				if _, ok := projectile.AlreadyHit[enemy]; ok {
					continue
				}

				if wizServer.CheckCollisionCircle(
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

					w.WriteProjectileHitEnemy(enemy, enemy.Dx, enemy.Dy)

					if enemy.Combat.Health() <= 0 {
						enemy.DieAnim = true
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

				if wizServer.CheckCollisionCircle(
					projectile.X+projectile.HitboxOffsetX,
					projectile.Y+projectile.HitboxOffsetY,
					projectile.ScaledRadius,
					otherProjectile.X+otherProjectile.HitboxOffsetX,
					otherProjectile.Y+otherProjectile.HitboxOffsetY,
					otherProjectile.ScaledRadius,
				) {
					if projectile.Scale > otherProjectile.Scale {
						delete(w.projectiles, otherProjectile.ID)
						w.WriteProjectileDestroyed(otherProjectile.Caster.Data.Username, otherProjectile.ID)
					} else if projectile.Scale < otherProjectile.Scale {
						delete(w.projectiles, projectile.ID)
						w.WriteProjectileDestroyed(projectile.Caster.Data.Username, projectile.ID)
					} else {
						delete(w.projectiles, projectile.ID)
						delete(w.projectiles, otherProjectile.ID)
						w.WriteProjectileDestroyed(projectile.Caster.Data.Username, projectile.ID)
						w.WriteProjectileDestroyed(otherProjectile.Caster.Data.Username, otherProjectile.ID)
					}
				}
			}
		}

		if projectile.TicksToLive < 1 {
			delete(w.projectiles, projectile.ID)
		}
	}
}
