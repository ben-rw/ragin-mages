package wizarena

import (
	"fmt"
	"math"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Stat int

const (
	ProjectileScale = iota
	ProjectileSpeed
	Knockback
)

func (w *WizArena) NewStatText() map[Stat]string {
	return map[Stat]string{
		ProjectileScale: fmt.Sprintf("Fireball Size: %v", w.wizard.Combat.ProjectileScale()),
		ProjectileSpeed: fmt.Sprintf("Fireball Speed: %v", w.wizard.Combat.ProjectileSpeed()),
		Knockback:       fmt.Sprintf("F.B. Knockback: %v", w.wizard.Combat.Knockback()),
	}
}

var (
	// explicity initialized as text.Face to avoid interface boxing
	fontFace8  text.Face = &text.GoTextFace{Source: shared.FontSrc, Size: 8}
	fontFace24 text.Face = &text.GoTextFace{Source: shared.FontSrc, Size: 24}
)

type StaticText int

const (
	instructionText1 = iota
	instructionText2
	deadText
)

var (
	instructionText1Img = NewTextLabel(
		"You need more POWER!",
		fontFace8,
	)
	instructionText2Img = NewTextLabel(
		"Get boosts from chests, skeletons, and your friends!",
		fontFace8,
	)
	deadTextImg = NewTextLabel(
		"YOU DIED",
		fontFace24,
	)
)

var opts text.DrawOptions
var StaticTextMap = map[StaticText]*ebiten.Image{
	instructionText1: instructionText1Img,
	instructionText2: instructionText2Img,
	deadText:         deadTextImg,
}

func NewTextLabel(str string, face text.Face) *ebiten.Image {
	w, h := text.Measure(str, face, 1.0)

	imgW := max(1, int(math.Ceil(w)))
	imgH := max(1, int(math.Ceil(h)))
	img := ebiten.NewImage(imgW, imgH)

	text.Draw(img, str, face, nil)
	return img
}
