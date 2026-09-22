package wizarena

import "fmt"

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
