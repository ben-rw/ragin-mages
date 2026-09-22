package wizarena

import (
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
)

func (w *WizArena) SpawnEnemies() {
	for i := range shared.EnemySpawnCount {
		w.enemies = append(w.enemies, shared.NewEnemy(w.enemyImageCache[shared.Skeleton], shared.Skeleton, true, shared.EnemySpawns[i].X*shared.TileSize, shared.EnemySpawns[i].Y*shared.TileSize))
	}
}

func (w *WizArena) EnemyRespawn() {
	if w.enemyRespawnTimer.IsZero() || time.Until(w.enemyRespawnTimer) <= 0 {
		w.enemyRespawnTimer = time.Now().Add(shared.EnemyRespawnTimer * time.Second)
		w.SpawnEnemies()
	}
}
