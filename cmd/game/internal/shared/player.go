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

// walks player from one place to another at speed proportionate
// to the initial distance between the player's location and the
// destination with no input from the player
func (p *Player) ScriptedWalk(destX, destY float64) {
	distX := destX - p.X
	stepX := distX / 100000
	distY := destY - p.Y
	stepY := distY / 100000
	for destX != p.X || destY != p.Y {
		if destX-p.X > stepX {
			p.X = destX
		} else {
			p.X += stepX
		}
		if destY-p.Y > stepY {
			p.Y = destY
		} else {
			p.Y += stepY
		}
		log.Println(p.X, p.Y)
	}
}

func NewPlayer(data *protocol.PlayerData, joinOrder int) *Player {
	imgPath := PlayerSpriteIndex[data.SpriteIndex]
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, imgPath)
	if err != nil {
		log.Fatal(err)
	}

	startPosition := StartingPositions[joinOrder]

	return &Player{
		Sprite: &Sprite{
			Img:         playerImg,
			X:           startPosition.X,
			Y:           startPosition.Y,
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
				Attack: animations.NewAnimation(16, 16, 0, 20.0),
				Dying:  animations.NewAnimation(25, 25, 0, 10.0),
			},
			JustJoined: true,
		},
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
