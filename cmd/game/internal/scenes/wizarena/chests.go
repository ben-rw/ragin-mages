package wizarena

import (
	"image"
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
)

const ChestRespawnTimer = 30

type OpenedChest struct {
	X, Y int
}

func (w *WizArena) CheckChestCollisions() {
	for _, chest := range w.chests {
		if chest.Overlaps(image.Rect(
			int(w.wizard.X),
			int(w.wizard.Y),
			int(w.wizard.X)+16,
			int(w.wizard.Y)+16,
		)) {
			if _, ok := w.openedChests[OpenedChest{chest.Min.X, chest.Min.Y}]; !ok {
				w.wizard.Combat.RandomBoost(shared.ChestBoost)
				w.openedChests[OpenedChest{chest.Min.X, chest.Min.Y}] = struct{}{}
				w.statText = w.NewStatText()
			}
		}
	}
}

// spawn chests on a timer
func (w *WizArena) ChestRespawn() {
	if w.chestRespawnTimer.IsZero() || time.Until(w.chestRespawnTimer) <= 0 {
		w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)
		w.openedChests = map[OpenedChest]struct{}{}
	}
}
