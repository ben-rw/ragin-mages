package spritesheet

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	TileWidth     int
	TileHeight    int
}

func NewSpriteSheet(w, h, tw, th int) *SpriteSheet {
	return &SpriteSheet{
		w, h, tw, th,
	}
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.TileWidth
	y := (index / s.WidthInTiles) * s.TileHeight

	return image.Rect(x, y, x+s.TileWidth, y+s.TileHeight)
}

func LoadFrames(s *SpriteSheet, img *ebiten.Image) []*ebiten.Image {
	frames := make([]*ebiten.Image, s.WidthInTiles*s.HeightInTiles)
	for i := range s.WidthInTiles * s.HeightInTiles {
		frame := img.SubImage(
			s.Rect(i),
		).(*ebiten.Image)
		frames[i] = frame
	}
	return frames
}
