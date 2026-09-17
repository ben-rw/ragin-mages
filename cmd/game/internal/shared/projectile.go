package shared

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

type ProjectileType int

const (
	Fireball ProjectileType = iota
)

const (
	FireballPath          = "assets/images/warped_shooting_fx/charged/spritesheet.png"
	FireballWidth         = 63.0
	FireballHeight        = 48.0
	FireballWidthInTiles  = 6
	FireballHeightInTiles = 1
	FireballAnimSpeed     = 1
	FireballHitboxRadius  = 12.0 // radius of ball at the tip of projectile
	FBInnerBallX          = 47.0 // 47x24 px from the sprites's top left corner to center of blast ball
	FBInnerBallY          = 24.0
	FBCenterX             = FireballWidth * .5
	FBCenterY             = FireballHeight * .5
)

var projectileImgPaths = map[ProjectileType]string{
	Fireball: FireballPath,
}

type Projectile struct {
	*Sprite
	Caster        *WizardPlayer
	Damage        float64
	Speed         float64
	Knockback     float64
	NormX         float64
	NormY         float64
	Radius        float64
	ScaledRadius  float64
	Scale         float64
	HitboxOffsetX float64
	HitboxOffsetY float64
	Rotation      float64
	CenterX       float64
	CenterY       float64
	TicksToLive   int
	AlreadyHit    map[any]struct{}
}

func LoadProjectile(projectileType ProjectileType) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, projectileImgPaths[projectileType])
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (w *WizardPlayer) ShootProjectile(img *ebiten.Image, cursorX, cursorY float64) *Projectile {
	vX := cursorX - w.X
	vY := cursorY - w.Y
	vlen := math.Hypot(vX, vY)
	if vlen == 0 {
		return &Projectile{}
	}
	normX := vX / vlen
	normY := vY / vlen
	rotation := math.Atan2(vY, vX)

	scaledOffsetX := (FBInnerBallX - FBCenterX) * w.Combat.ProjectileScale()
	scaledOffsetY := (FBInnerBallY - FBCenterY) * w.Combat.ProjectileScale()

	rotatedOffsetX := (scaledOffsetX * normX) - (scaledOffsetY * normY)
	rotatedOffsetY := (scaledOffsetX * normY) + (scaledOffsetY * normX)

	s := spritesheet.NewSpriteSheet(FireballWidthInTiles, FireballHeightInTiles, FireballWidth, FireballHeight)

	var ProjectileAnimations = map[EntityState]*animations.Animation{
		FireballFly: animations.NewAnimation(0, 5, 1, FireballAnimSpeed),
	}

	return &Projectile{
		Sprite: NewSprite(
			img,
			w.X+HalfTile+normX*HalfTile,
			w.Y+HalfTile+normY*HalfTile,
			normX*w.Combat.projectileSpeed,
			normY*w.Combat.projectileSpeed,
			s,
			ProjectileAnimations,
			ProjectileAnimations[FireballFly],
			FireballFly,
			false,
			1.0,
		),
		Caster:        w,
		Damage:        w.Combat.AttackPower(),
		Speed:         w.Combat.ProjectileSpeed(),
		Knockback:     w.Combat.Knockback(),
		NormX:         normX,
		NormY:         normY,
		Radius:        FireballHitboxRadius,
		ScaledRadius:  FireballHitboxRadius * w.Combat.ProjectileScale(),
		Scale:         w.Combat.ProjectileScale(),
		HitboxOffsetX: rotatedOffsetX,
		HitboxOffsetY: rotatedOffsetY,
		Rotation:      rotation,
		CenterX:       FireballWidth / 2,
		CenterY:       FireballHeight / 2,
		TicksToLive:   w.Combat.AttackCooldown(),
		AlreadyHit:    make(map[any]struct{}),
	}
}

func (p *Projectile) Update() {
	p.X += p.Dx
	p.Y += p.Dy
	p.TicksToLive -= 1
}
