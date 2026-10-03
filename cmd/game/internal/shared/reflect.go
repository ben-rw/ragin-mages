package shared

import (
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
)

const (
	ReflectPath          = "assets/images/warped_shooting_fx/hits-1/spritesheet.png"
	ReflectWidth         = 32
	ReflectHeight        = 32
	ReflectWidthInTiles  = 5
	ReflectHeightInTiles = 1
	ReflectAnimSpeed     = 1
)

type ReflectType int

const (
	Circle ReflectType = iota
)

var ReflectSpriteSheet = spritesheet.NewSpriteSheet(ReflectWidthInTiles, ReflectHeightInTiles, ReflectWidth, ReflectHeight)

type Reflect struct {
	*Sprite
	Active bool
}

func NewReflect(x, y float64) *Reflect {
	ReflectAnimations := map[EntityState]*animations.Animation{
		ReflectCircle: animations.NewAnimation(0, 4, 1, ReflectAnimSpeed),
	}
	return &Reflect{
		Sprite: &Sprite{
			X:               x + HalfTile,
			Y:               y + HalfTile,
			SpriteSheet:     ReflectSpriteSheet,
			Animations:      ReflectAnimations,
			ActiveAnimation: nil,
		},
		Active: false,
	}
}

func (r *Reflect) GetActiveAnimation() *animations.Animation {
	if !r.Active {
		return nil
	}

	anim := r.Animations[ReflectCircle]
	if anim.Over {
		anim.Over = false
		r.Active = false
	}
	return anim
}

func (w *WizardPlayer) ReflectProjectile(p *Projectile, cursorX, cursorY float64) {
	vX := cursorX - w.X
	vY := cursorY - w.Y
	vlen := math.Hypot(vX, vY)
	if vlen == 0 {
		p.Despawn()
		return
	}
	normX := vX / vlen
	normY := vY / vlen
	rotation := math.Atan2(vY, vX)

	scaledOffsetX := (FBInnerBallX - FBCenterX) * w.Combat.ProjectileScale()
	scaledOffsetY := (FBInnerBallY - FBCenterY) * w.Combat.ProjectileScale()

	rotatedOffsetX := (scaledOffsetX * normX) - (scaledOffsetY * normY)
	rotatedOffsetY := (scaledOffsetX * normY) + (scaledOffsetY * normX)

	p.ReflectCount++

	boostPercent := .5 / float64(p.ReflectCount)

	p.Caster = w
	// increase stats of reflected projectile
	p.Speed *= (1.0 + boostPercent)
	p.Knockback *= (1.0 + boostPercent)
	p.Scale *= (1.0 + boostPercent)
	p.TicksToLive = w.Combat.attackCooldown
	p.Dx = normX * p.Speed
	p.Dy = normY * p.Speed
	p.NormX = normX
	p.NormY = normY
	p.ScaledRadius = FireballHitboxRadius * p.Scale
	p.HitboxOffsetX = rotatedOffsetX
	p.HitboxOffsetY = rotatedOffsetY
	p.Rotation = rotation
	p.AlreadyHit = make(map[any]struct{})
}
