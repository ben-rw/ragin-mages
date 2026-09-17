package shared

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/animations"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type EntityState int

const (
	Idle EntityState = iota
	WalkDown
	WalkUp
	WalkLeft
	WalkRight
	JustJoined
	AttackingDown
	AttackingUp
	AttackingLeft
	AttackingRight
	Dying
	FireballFly
)

type Sprite struct {
	Img             *ebiten.Image
	X, Y, Dx, Dy    float64
	SpriteSheet     *spritesheet.SpriteSheet
	Animations      map[EntityState]*animations.Animation
	ActiveAnimation *animations.Animation
	State           EntityState
	Noclip          bool
	Alpha           float32
}

func NewSprite(img *ebiten.Image, x, y, dx, dy float64, spritesheet *spritesheet.SpriteSheet, animations map[EntityState]*animations.Animation, activeAnimation *animations.Animation, state EntityState, noclip bool, alpha float32) *Sprite {
	return &Sprite{
		Img:             img,
		X:               x,
		Y:               y,
		Dx:              dx,
		Dy:              dy,
		SpriteSheet:     spritesheet,
		Animations:      animations,
		ActiveAnimation: activeAnimation,
		State:           state,
		Noclip:          noclip,
		Alpha:           alpha,
	}
}

func (s *Sprite) SetActiveAnimation(state EntityState) {
	if s.State == state {
		return
	}
	s.State = state
	switch s.State {
	case JustJoined:
		s.ActiveAnimation = s.Animations[JustJoined]
	case Dying:
		s.ActiveAnimation = s.Animations[Dying]
	case AttackingDown:
		s.ActiveAnimation = s.Animations[AttackingDown]
	case WalkRight:
		s.ActiveAnimation = s.Animations[WalkRight]
	case WalkLeft:
		s.ActiveAnimation = s.Animations[WalkLeft]
	case WalkDown:
		s.ActiveAnimation = s.Animations[WalkDown]
	case WalkUp:
		s.ActiveAnimation = s.Animations[WalkUp]
	default:
		s.ActiveAnimation = s.Animations[Idle]
	}
}

func (s *Sprite) AnimationOver() {
	switch s.State {
	case JustJoined, Dying, AttackingDown, AttackingUp, AttackingLeft, AttackingRight:
		if s.ActiveAnimation.Over {
			s.State = Idle
		}
	}
}
