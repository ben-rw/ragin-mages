package shared

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	*Sprite
	Combat        *EnemyCombat
	enemyType     EnemyType
	FollowsPlayer bool
	HurtboxRadius float64
}

type EnemyType int

const (
	Skeleton EnemyType = iota
)

func NewEnemyImageCache(enemyTypes []EnemyType) (map[EnemyType]*ebiten.Image, error) {
	imgs := make(map[EnemyType]*ebiten.Image, len(enemyTypes))
	for _, enemyType := range enemyTypes {
		enemyImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, EnemySpriteIndex[enemyType])
		if err != nil {
			return nil, err
		}
		imgs[enemyType] = enemyImg

	}
	return imgs, nil
}

func NewEnemy(img *ebiten.Image, enemyType EnemyType, followsPlayer bool, x, y float64) *Enemy {

	return &Enemy{
		NewSprite(
			img,
			x,
			y,
			0,
			0,
			spritesheet.NewSpriteSheet(4, 7, TileSize, TileSize),
			map[EntityState]*animations.Animation{
				WalkDown:       animations.NewAnimation(4, 12, 4, 20.0),
				WalkUp:         animations.NewAnimation(5, 13, 4, 20.0),
				WalkLeft:       animations.NewAnimation(6, 14, 4, 20.0),
				WalkRight:      animations.NewAnimation(7, 15, 4, 20.0),
				Idle:           animations.NewAnimation(0, 16, 16, 20.0),
				AttackingDown:  animations.NewAnimation(16, 16, 1, 20.0),
				AttackingUp:    animations.NewAnimation(17, 17, 1, 20.0),
				AttackingLeft:  animations.NewAnimation(18, 18, 1, 20.0),
				AttackingRight: animations.NewAnimation(19, 19, 1, 20.0),
				Die:            animations.NewAnimation(25, 25, 0, 10.0),
			},
			&animations.Animation{},
			true,
			false,
			0.0,
		),
		NewEnemyCombat(EnemyAttackCooldown, EnemyHealth, EnemyAttackPower, 0, 0, 0, EnemyKnockBack),
		Skeleton,
		followsPlayer,
		DefaultHurtboxRadius,
	}
}
