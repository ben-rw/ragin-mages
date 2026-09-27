package shared

import (
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
)

const (
	ReflectPath          = "assets/images/warped_shooting_fx/hits-1/spritesheet.png"
	ReflectWidth         = 160
	ReflectHeight        = 32
	ReflectWidthInTiles  = 5
	ReflectHeightInTiles = 1
	ReflectAnimSpeed     = 1
)

var ReflectAnimations = map[EntityState]*animations.Animation{
	ReflectCircle: animations.NewAnimation(0, 4, 1, ReflectAnimSpeed),
}

var ReflectSpriteSheet = spritesheet.NewSpriteSheet(ReflectWidthInTiles, ReflectHeightInTiles, ReflectWidth, ReflectHeight)

type Reflect struct {
	*Sprite
}

func NewReflect(x, y float64) *Reflect {
	return &Reflect{
		Sprite: &Sprite{
			X:               x + HalfTile,
			Y:               y + HalfTile,
			SpriteSheet:     ReflectSpriteSheet,
			Animations:      ReflectAnimations,
			ActiveAnimation: nil,
		},
	}
}

func (r *Reflect) CheckReflectAnimation() {
	if r.ActiveAnimation.Over {
		r.ActiveAnimation.Over = false
		r.ActiveAnimation = nil
	}
}

func (w *WizardPlayer) ReflectProjectile(p *Projectile, cursorX, cursorY float64) {
	vX := cursorX - w.X
	vY := cursorY - w.Y
	vlen := math.Hypot(vX, vY)
	if vlen == 0 {
		p.Despawn()
	}
	normX := vX / vlen
	normY := vY / vlen
	rotation := math.Atan2(vY, vX)

	scaledOffsetX := (FBInnerBallX - FBCenterX) * w.Combat.ProjectileScale()
	scaledOffsetY := (FBInnerBallY - FBCenterY) * w.Combat.ProjectileScale()

	rotatedOffsetX := (scaledOffsetX * normX) - (scaledOffsetY * normY)
	rotatedOffsetY := (scaledOffsetX * normY) + (scaledOffsetY * normX)

	p.Caster = w
	// increase stats of reflected projectile
	p.Speed *= 1.25
	p.Knockback *= 1.10
	p.Scale *= 1.25
	p.TicksToLive += DefaultAttackCooldown
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
