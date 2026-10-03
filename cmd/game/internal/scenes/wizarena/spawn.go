package wizarena

import (
	"math/rand"
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

// gives players random, non-overlapping spawns
func (w *WizArena) AssignPlayerSpawns() {
	pLen := len(w.Players)
	playerSpawnIndex := make([]int, pLen)
	for i := 0; i < pLen; i++ {
		playerSpawnIndex[i] = i
	}

	for _, player := range w.Players {
		i := rand.Intn(len(playerSpawnIndex))
		s := playerSpawnIndex
		spawnIndex := s[i]
		s[i] = s[len(s)-1]
		playerSpawnIndex = s[:len(s)-1]

		player.X = shared.PlayerSpawns[spawnIndex].X * shared.TileSize
		player.Y = shared.PlayerSpawns[spawnIndex].Y * shared.TileSize
	}
}
