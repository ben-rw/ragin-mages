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
	FireballWidth         = 63
	FireballHeight        = 48
	FireballWidthInTiles  = 6
	FireballHeightInTiles = 1
	FireballAnimSpeed     = 1
)

var projectileImgPaths = map[ProjectileType]string{
	Fireball: FireballPath,
}

type Projectile struct {
	*Sprite
	Speed       float64
	Size        float64
	Knockback   float64
	Rotation    float64
	CenterX     float64
	CenterY     float64
	TicksToLive int
}

func LoadProjectile(projectileType ProjectileType) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, projectileImgPaths[projectileType])
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SpawnProjectile(img *ebiten.Image, speed, size, knockback, playerX, playerY, cursorX, cursorY float64, ticksToLive int) *Projectile {
	vX := cursorX - playerX
	vY := cursorY - playerY
	vlen := math.Hypot(vX, vY)
	if vlen == 0 {
		return &Projectile{}
	}
	normX := vX / vlen
	normY := vY / vlen
	rotation := math.Atan2(vY, vX)

	s := spritesheet.NewSpriteSheet(FireballWidthInTiles, FireballHeightInTiles, FireballWidth, FireballHeight)
	anim := animations.NewAnimation(0, 5, 1, FireballAnimSpeed)

	return &Projectile{
		&Sprite{
			Img:         img,
			X:           playerX + HalfTile + normX*HalfTile,
			Y:           playerY + HalfTile + normY*HalfTile,
			Dx:          normX * speed,
			Dy:          normY * speed,
			SpriteSheet: s,
			Animations: map[EntityState]*animations.Animation{
				FireballFly: anim,
			},
			ActiveAnimation: anim,
		},
		speed,
		size,
		knockback,
		rotation,
		-FireballWidth / 2,
		-FireballHeight / 2,
		ticksToLive,
	}
}

func (p *Projectile) Update() {
	p.X += p.Dx
	p.Y += p.Dy
	p.TicksToLive -= 1
}
