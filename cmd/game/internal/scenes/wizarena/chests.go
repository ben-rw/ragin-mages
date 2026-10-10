package wizarena

import (
	"image"

	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) CheckChestCollisions() {
	for _, chest := range w.chests {
		if chest.Overlaps(image.Rect(
			int(w.wizard.X),
			int(w.wizard.Y),
			int(w.wizard.X)+16,
			int(w.wizard.Y)+16,
		)) {
			if _, ok := w.openedChests[wizServer.OpenedChest{X: chest.Min.X, Y: chest.Min.Y}]; !ok {
				w.wizard.Combat.RandomBoost(ChestBoost, StandardMult)
				w.openedChests[wizServer.OpenedChest{X: chest.Min.X, Y: chest.Min.Y}] = struct{}{}
				w.WriteWizardStatUpdate()
				w.dynamicText = w.NewDynamicTextMap()
			}
		}
	}
}
