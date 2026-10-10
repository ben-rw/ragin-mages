package wizarena

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (w *WizArena) NewScoreboard() {
	wizSlice := make([]*WizardPlayer, 0, 8)
	for _, wizard := range w.wizards {
		wizSlice = append(wizSlice, wizard)
	}

	if len(wizSlice) > 1 {
		slices.SortFunc(wizSlice, func(a, b *WizardPlayer) int {
			return a.Combat.GetTotalPower() - b.Combat.GetTotalPower()
		})
		slices.Reverse(wizSlice)
	}

	w.scoreboard = "POWER RANKINGS\n\n"
	for i, wizard := range wizSlice {
		if i == 0 && w.round <= wizServer.Rounds {
			// round winner gets stats boosted
			w.scoreboard += fmt.Sprintf("%v - %v + %v\n", wizard.Data.Username, wizard.Combat.GetTotalPower(), WinRoundBoost)
			wizard.Combat.RandomBoost(WinRoundBoost, StandardMult)
			w.WriteWizardStatUpdate()
			w.dynamicText = w.NewDynamicTextMap()
		} else {
			w.scoreboard += fmt.Sprintf("%v - %v\n", wizard.Data.Username, wizard.Combat.GetTotalPower())
		}
	}

	if w.round > wizServer.Rounds {
		w.scoreboard += fmt.Sprintf("\n%v WINS!", wizSlice[0].Data.Username)
	}
}

type DynamicText int

const (
	ProjectileScale = iota
	ProjectileSpeed
	Knockback
	AttackCooldown
	RoundTimer
)

func (w *WizArena) NewDynamicTextMap() map[DynamicText]string {
	return map[DynamicText]string{
		// display stats to player starting at 1 and scaling by number of boosts applied to each stat
		// to let player see how much they've increased each stat, rather than internal stat numbers
		ProjectileSpeed: fmt.Sprintf("Fireball Speed: %v", w.wizard.Combat.ProjectileSpeed()*(1/PSpeedMult)-(1/PSpeedMult*BaseProjectileSpeed-1)),
		ProjectileScale: fmt.Sprintf("Fireball Size: %v", w.wizard.Combat.ProjectileScale()*(1/PScaleMult)-(1/PScaleMult*BaseProjectileScale-1)),
		Knockback:       fmt.Sprintf("F.B. Knockback: %v", w.wizard.Combat.Knockback()*(1/KnockbackMult)-(1/KnockbackMult*BaseKnockback-1)),
		AttackCooldown:  fmt.Sprintf("Attack Cooldown: %.2fs", float64(w.wizard.Combat.AttackCooldown())/60.0),
	}
}

// can't use these methods for stats now, maybe after upcoming refactor
func (w *WizArena) UpdateProjectileScaleText() {
	w.dynamicText[ProjectileSpeed] = fmt.Sprintf("Fireball Speed: %v", w.wizard.Combat.ProjectileSpeed()*2-5)
}

func (w *WizArena) UpdateProjectileSpeedText() {
	w.dynamicText[ProjectileScale] = fmt.Sprintf("Fireball Size: %v", w.wizard.Combat.ProjectileScale()*4-1)
}

func (w *WizArena) UpdateKnockbackText() {
	w.dynamicText[Knockback] = fmt.Sprintf("F.B. Knockback: %v", w.wizard.Combat.Knockback()*2-5)
}

func (w *WizArena) UpdateAttackCDText() {
	w.dynamicText[AttackCooldown] = fmt.Sprintf("Attack Cooldown: %.2fs", float64(w.wizard.Combat.AttackCooldown())/60.0)
}

func (w *WizArena) UpdateRoundTimerText() {
	remaining := w.phaseTimer.Sub(time.Now())
	currentSecond := int(remaining.Seconds())
	if currentSecond != w.lastFormattedSecond {
		w.roundTimerString = fmt.Sprintf(
			"%d:%02d",
			int(remaining/time.Minute),
			int(remaining/time.Second)%60,
		)
		w.lastFormattedSecond = currentSecond
	}
}

var (
	// explicity initialized as text.Face to avoid interface boxing
	fontFace8  text.Face = &text.GoTextFace{Source: shared.FontSrc, Size: 8}
	fontFace16 text.Face = &text.GoTextFace{Source: shared.FontSrc, Size: 16}
	fontFace24 text.Face = &text.GoTextFace{Source: shared.FontSrc, Size: 24}
)

type StaticText int

const (
	instructionText1 = iota
	instructionText2
	deadText
	round1Text
	round2Text
	round3Text
)

var (
	instructionText1Img = NewTextLabel(
		"You need more POWER!",
		fontFace8,
	)
	instructionText2Img = NewTextLabel(
		"Get POWER from chests, skeletons, and your friends!",
		fontFace8,
	)
	deadTextImg = NewTextLabel(
		"YOU DIED",
		fontFace24,
	)
	Round1TextImg = NewTextLabel(
		"ROUND 1",
		fontFace24,
	)
	Round2TextImg = NewTextLabel(
		"ROUND 2",
		fontFace24,
	)
	Round3TextImg = NewTextLabel(
		"FINAL ROUND",
		fontFace24,
	)
)

var opts text.DrawOptions
var StaticTextMap = map[StaticText]*ebiten.Image{
	instructionText1: instructionText1Img,
	instructionText2: instructionText2Img,
	deadText:         deadTextImg,
	round1Text:       Round1TextImg,
	round2Text:       Round2TextImg,
	round3Text:       Round3TextImg,
}

func NewTextLabel(str string, face text.Face) *ebiten.Image {
	w, h := text.Measure(str, face, 1.0)

	imgW := max(1, int(math.Ceil(w)))
	imgH := max(1, int(math.Ceil(h)))
	img := ebiten.NewImage(imgW, imgH)

	text.Draw(img, str, face, nil)
	return img
}
