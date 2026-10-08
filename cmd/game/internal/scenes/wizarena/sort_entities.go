package wizarena

import (
	"cmp"
	"slices"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
)

func (w *WizArena) SortEntities() {
	w.sortedWizards = make([]*shared.WizardPlayer, 0, 8)
	w.sortedEnemies = make([]*shared.Enemy, 0, 32)
	w.sortedProjectiles = make([]*shared.Projectile, 0, 16)

	for _, wizard := range w.wizards {
		if !wizard.Combat.Dead {
			w.sortedWizards = append(w.sortedWizards, wizard)
		}
	}
	if len(w.sortedWizards) > 1 {
		slices.SortFunc(w.sortedWizards, func(a, b *shared.WizardPlayer) int {
			return cmp.Compare(a.Data.SpriteIndex, b.Data.SpriteIndex)
		})
	}

	for _, enemy := range w.enemies {
		w.sortedEnemies = append(w.sortedEnemies, enemy)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedEnemies, func(a, b *shared.Enemy) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}

	for _, enemy := range w.enemies {
		w.sortedEnemies = append(w.sortedEnemies, enemy)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedEnemies, func(a, b *shared.Enemy) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}

	for _, projectile := range w.projectiles {
		w.sortedProjectiles = append(w.sortedProjectiles, projectile)
	}
	if len(w.sortedEnemies) > 1 {
		slices.SortFunc(w.sortedProjectiles, func(a, b *shared.Projectile) int {
			return cmp.Compare(a.ID, b.ID)
		})
	}
}
