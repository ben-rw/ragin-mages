package shared

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EntityState int

const (
	Idle EntityState = iota
	Down
	Up
	Left
	Right
	Join
	Attack
	Dying
	FireballFly
)

type Sprite struct {
	Img             *ebiten.Image
	X, Y, Dx, Dy    float64
	SpriteSheet     *spritesheet.SpriteSheet
	Animations      map[EntityState]*animations.Animation
	ActiveAnimation *animations.Animation
	JustJoined      bool
	Dying           bool
}

func (s *Sprite) GetActiveAnimation() *animations.Animation {
	if s.JustJoined {
		s.ActiveAnimation = s.Animations[Join]
		if s.ActiveAnimation.Over == true {
			s.JustJoined = false
		} else {
			return s.ActiveAnimation
		}
	}
	if s.Dying {
		s.ActiveAnimation = s.Animations[Dying]
		if s.ActiveAnimation.Over == true {
			s.Dying = false
		} else {
			return s.ActiveAnimation
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return s.Animations[Attack]
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		return s.Animations[Attack]
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return s.Animations[Attack]
	}
	if s.Dx > 0 {
		return s.Animations[Right]
	}
	if s.Dx < 0 {
		return s.Animations[Left]
	}
	if s.Dy > 0 {
		return s.Animations[Down]
	}
	if s.Dy < 0 {
		return s.Animations[Up]
	}
	return s.Animations[Idle]
}
