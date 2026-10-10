package wizarena

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	*shared.Sprite
	Combat        *EnemyCombat
	Type          EnemyType
	FollowsPlayer bool
	HurtboxRadius float64
	ID            uint16
}

type EnemyType int

const (
	Skeleton EnemyType = iota
)

var EnemySpriteSheet = spritesheet.NewSpriteSheet(4, 7, wizServer.TileSize, wizServer.TileSize)

var EnemySpriteIndex map[EnemyType]string = map[EnemyType]string{
	Skeleton: shared.ImgRootPath + "Skeleton" + shared.ImgPathEnd,
}

func NewEnemyImageCache(enemyTypes []EnemyType) (map[EnemyType]*ebiten.Image, error) {
	imgs := make(map[EnemyType]*ebiten.Image, len(enemyTypes))
	for _, enemyType := range enemyTypes {
		enemyImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, EnemySpriteIndex[enemyType])
		if err != nil {
			return nil, err
		}
		imgs[enemyType] = enemyImg

	}
	return imgs, nil
}

func NewEnemy(img *ebiten.Image, enemyType EnemyType, followsPlayer bool, x, y float64, id uint16) *Enemy {
	return &Enemy{
		Sprite: shared.NewSprite(
			img,
			x,
			y,
			0,
			0,
			EnemySpriteSheet,
			map[shared.EntityState]*animations.Animation{
				shared.WalkDown:       animations.NewAnimation(4, 12, 4, 20.0),
				shared.WalkUp:         animations.NewAnimation(5, 13, 4, 20.0),
				shared.WalkLeft:       animations.NewAnimation(6, 14, 4, 20.0),
				shared.WalkRight:      animations.NewAnimation(7, 15, 4, 20.0),
				shared.Idle:           animations.NewAnimation(0, 16, 16, 20.0),
				shared.AttackingDown:  animations.NewAnimation(16, 16, 1, 20.0),
				shared.AttackingUp:    animations.NewAnimation(17, 17, 1, 20.0),
				shared.AttackingLeft:  animations.NewAnimation(18, 18, 1, 20.0),
				shared.AttackingRight: animations.NewAnimation(19, 19, 1, 20.0),
				shared.Die:            animations.NewAnimation(25, 25, 0, 10.0),
			},
			&animations.Animation{},
			true,
			false,
			0.0,
		),
		Combat:        NewEnemyCombat(EnemyAttackCooldown, EnemyHealth, EnemyAttackPower, 0, 0, 0, EnemyKnockBack),
		Type:          Skeleton,
		FollowsPlayer: followsPlayer,
		HurtboxRadius: HurtboxRadius,
		ID:            id,
	}
}

func (w *WizArena) SpawnEnemies(serverEnemies map[uint16]*wizServer.Enemy) {
	i := 0
	for k := range serverEnemies {
		enemy := NewEnemy(
			w.enemyImageCache[Skeleton],
			Skeleton,
			true,
			serverEnemies[k].X,
			serverEnemies[k].Y,
			serverEnemies[k].ID,
		)
		w.enemies[enemy.ID] = enemy
		i++
	}
}
