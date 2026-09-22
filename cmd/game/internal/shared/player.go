package shared

import (
	"log"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	nameTagSize = 4
)

type Player struct {
	*Sprite
	Data    *protocol.PlayerData
	NameTag *NameTag
}

type NameTag struct {
	Face          *text.GoTextFace
	X, Y          float64
	LayoutOptions text.LayoutOptions
}

var PlayerSpriteSheet = spritesheet.NewSpriteSheet(4, 7, TileSize, TileSize)

func NewPlayer(data *protocol.PlayerData, joinOrder int) *Player {
	imgPath := PlayerSpriteIndex[data.SpriteIndex]
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	startPosition := StartingPositions[joinOrder]

	return &Player{
		Sprite: NewSprite(
			playerImg,
			startPosition.X,
			startPosition.Y,
			0,
			0,
			PlayerSpriteSheet,
			map[EntityState]*animations.Animation{
				WalkDown:       animations.NewAnimation(4, 12, 4, 20.0),
				WalkUp:         animations.NewAnimation(5, 13, 4, 20.0),
				WalkLeft:       animations.NewAnimation(6, 14, 4, 20.0),
				WalkRight:      animations.NewAnimation(7, 15, 4, 20.0),
				Idle:           animations.NewAnimation(0, 16, 16, 20.0),
				Join:           animations.NewAnimation(26, 27, 1, 60.0),
				AttackingDown:  animations.NewAnimation(16, 16, 1, 20.0),
				AttackingUp:    animations.NewAnimation(17, 17, 1, 20.0),
				AttackingLeft:  animations.NewAnimation(18, 18, 1, 20.0),
				AttackingRight: animations.NewAnimation(19, 19, 1, 20.0),
				Die:            animations.NewAnimation(25, 25, 1, 80.0),
			},
			&animations.Animation{},
			false,
			true,
			1.0,
		),
		Data: &protocol.PlayerData{
			Name:        data.Name,
			Score:       data.Score,
			Host:        data.Host,
			SpriteIndex: data.SpriteIndex,
		},
		NameTag: &NameTag{
			Face: &text.GoTextFace{
				Source: FontSrc,
				Size:   nameTagSize,
			},
			X: startPosition.X + TileSize/2,
			Y: startPosition.Y + TileSize + 2,
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign: 1,
			},
		},
	}
}
