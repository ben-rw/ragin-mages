package shared

import (
	"log"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	*Sprite
	Combat        *EnemyCombat
	enemyType     EnemyType
	FollowsPlayer bool
}

type EnemyType string

const (
	Skeleton EnemyType = "skeleton"
)

func NewEnemy(enemyType EnemyType, followsPlayer bool, x, y float64) *Enemy {
	imgPath := EnemySpriteIndex[0]
	enemyImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	return &Enemy{
		&Sprite{
			Img:         enemyImg,
			X:           x,
			Y:           y,
			Dx:          0,
			Dy:          0,
			SpriteSheet: spritesheet.NewSpriteSheet(4, 7, TileSize, TileSize),
			Animations: map[EntityState]*animations.Animation{
				Up:     animations.NewAnimation(5, 13, 4, 20.0),
				Down:   animations.NewAnimation(4, 12, 4, 20.0),
				Left:   animations.NewAnimation(6, 14, 4, 20.0),
				Right:  animations.NewAnimation(7, 15, 4, 20.0),
				Idle:   animations.NewAnimation(0, 16, 16, 20.0),
				Join:   animations.NewAnimation(26, 27, 1, 60.0),
				Attack: animations.NewAnimation(16, 16, 0, 20),
			},
			JustJoined: false,
		},
		NewEnemyCombat(EnemyHealth, EnemyAttackPower, EnemyAttackCooldown, EnemyMoveSpeed, 0, 0, EnemyKnockBack),
		Skeleton,
		followsPlayer,
	}
}
