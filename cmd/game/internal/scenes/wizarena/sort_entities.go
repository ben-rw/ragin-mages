package wizarena

import (
	"cmp"
	"slices"
)

func (w *WizArena) SortEntities() {
	w.sortedWizards = make([]*WizardPlayer, 0, 8)
	w.sortedEnemies = make([]*Enemy, 0, 32)
	w.sortedProjectiles = make([]*Projectile, 0, 16)

	for _, wizard := range w.wizards {
		if !wizard.Combat.Dead {
			w.sortedWizards = append(w.sortedWizards, wizard)
		}
	}
	if len(w.sortedWizards) > 1 {
		slices.SortFunc(w.sortedWizards, func(a, b *WizardPlayer) int {
			return cmp.Compare(a.Data.SpriteIndex, b.Data.SpriteIndex)
		})
	}

	for _, enemy := range w.enemies {
		w.sortedEnemies = append(w.sortedEnemies, enemy)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedEnemies, func(a, b *Enemy) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}

	for _, enemy := range w.enemies {
		w.sortedEnemies = append(w.sortedEnemies, enemy)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedEnemies, func(a, b *Enemy) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}

	for _, projectile := range w.projectiles {
		w.sortedProjectiles = append(w.sortedProjectiles, projectile)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedProjectiles, func(a, b *Projectile) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}
}
