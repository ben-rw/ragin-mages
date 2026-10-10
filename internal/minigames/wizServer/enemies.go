package wizServer

import "time"

const (
	EnemyRespawnTimer         = 20
	EnemySpawnCount           = 16
	RoundStartEnemySpawnTimer = 5
)

type Enemy struct {
	X, Y, Dx, Dy float64
	Health       int
	ID           uint16
}

var EnemySpawns = map[int]struct{ X, Y float64 }{
	// separated into quadrants
	// top left, top right, bottom left, bottom right
	0: {X: 0, Y: 9},
	1: {X: 9, Y: 0},
	2: {X: 33, Y: 14},
	3: {X: 33, Y: 25},

	4: {X: 99, Y: 9},
	5: {X: 90, Y: 0},
	6: {X: 66, Y: 14},
	7: {X: 66, Y: 25},

	8:  {X: 0, Y: 50},
	9:  {X: 9, Y: 59},
	10: {X: 33, Y: 45},
	11: {X: 33, Y: 34},

	12: {X: 99, Y: 50},
	13: {X: 90, Y: 59},
	14: {X: 66, Y: 45},
	15: {X: 66, Y: 34},
}

var enemyID uint16 = 0

func getEnemyID() uint16 {
	enemyID++
	return enemyID
}

func (w *WizardServer) enemyRespawn() {
	if w.enemyRespawnTimer.IsZero() || time.Until(w.enemyRespawnTimer) <= 0 {
		w.enemyRespawnTimer = time.Now().Add(EnemyRespawnTimer * time.Second)
		w.spawnEnemies()
	}
}

func (w *WizardServer) spawnEnemies() {
	i := 0
	for range EnemySpawnCount {
		enemy := newEnemy(
			EnemySpawns[i].X*TileSize,
			EnemySpawns[i].Y*TileSize,
		)
		w.enemies[enemy.ID] = enemy
		i++
	}
}

func newEnemy(x, y float64) *Enemy {
	return &Enemy{
		X:  x,
		Y:  y,
		ID: getEnemyID(),
	}
}
