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
	Join
	AttackingDown
	AttackingUp
	AttackingLeft
	AttackingRight
	Die
	FireballFly
)

type Sprite struct {
	Img             *ebiten.Image
	SpriteSheet     *spritesheet.SpriteSheet
	Animations      map[EntityState]*animations.Animation
	ActiveAnimation *animations.Animation
	X, Y, Dx, Dy    float64
	AttackDirection EntityState
	Alpha           float32
	Noclip          bool
	JoinAnim        bool
	DieAnim         bool
	AttackAnim      bool
}

func NewSprite(img *ebiten.Image, x, y, dx, dy float64, spritesheet *spritesheet.SpriteSheet, animations map[EntityState]*animations.Animation, activeAnimation *animations.Animation, noclip, joinAnim bool, alpha float32) *Sprite {
	return &Sprite{
		Img:             img,
		X:               x,
		Y:               y,
		Dx:              dx,
		Dy:              dy,
		SpriteSheet:     spritesheet,
		Animations:      animations,
		ActiveAnimation: activeAnimation,
		Noclip:          noclip,
		JoinAnim:        joinAnim,
		Alpha:           alpha,
	}
}

func (s *Sprite) GetActiveAnimation() *animations.Animation {
	var anim *animations.Animation
	if s.JoinAnim {
		anim = s.Animations[Join]
		if anim.Over {
			anim.Over = false
			s.JoinAnim = false
			return s.Animations[Idle]
		} else {
			return anim
		}
	}
	if s.DieAnim {
		anim = s.Animations[Die]
		if s.ActiveAnimation.Over {
			anim.Over = false
			s.DieAnim = false
			return &animations.Animation{}
		} else {
			return anim
		}
	}
	if s.AttackAnim {
		switch s.AttackDirection {
		case AttackingDown:
			anim = s.Animations[AttackingDown]
		case AttackingUp:
			anim = s.Animations[AttackingUp]
		case AttackingRight:
			anim = s.Animations[AttackingRight]
		case AttackingLeft:
			anim = s.Animations[AttackingLeft]
		}

		if anim.Over {
			anim.Over = false
			s.AttackAnim = false
		} else {
			return anim
		}
	}
	switch {
	case s.Dx > 0:
		return s.Animations[WalkRight]
	case s.Dx < 0:
		return s.Animations[WalkLeft]
	case s.Dy > 0:
		return s.Animations[WalkDown]
	case s.Dy < 0:
		return s.Animations[WalkUp]
	default:
		return s.Animations[Idle]
	}
}
