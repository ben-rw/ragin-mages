package wizServer

import (
	"math"
	"time"
)

const (
	EnemyRespawnTimer         = 20
	EnemySpawnCount           = 16
	RoundStartEnemySpawnTimer = 5
)

const EnemyMoveSpeed = 1.5

type Enemy struct {
	*Body
	MoveSpeed     float64
	Health        float64
	ID            uint16
	FollowsPlayer bool
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

func (w *WizardServer) getEnemyID() uint16 {
	w.enemyID++
	return w.enemyID
}

func (w *WizardServer) enemyRespawn() {
	if w.enemyRespawnTimer.IsZero() || time.Until(w.enemyRespawnTimer) <= 0 {
		w.enemyRespawnTimer = time.Now().Add(EnemyRespawnTimer * time.Second)
		w.spawnEnemies()
	}
}

func (w *WizardServer) spawnEnemies() {
	// i := 0
	// for range EnemySpawnCount {
	// 	enemy := w.newEnemy(
	// 		EnemySpawns[i].X*TileSize,
	// 		EnemySpawns[i].Y*TileSize,
	// 	)
	// 	w.Mu.Lock()
	// 	w.enemies[enemy.ID] = enemy
	// 	w.Mu.Unlock()
	// 	i++
	// }
}

func (w *WizardServer) newEnemy(x, y float64) *Enemy {
	return &Enemy{
		Body: &Body{
			X:      x,
			Y:      y,
			Noclip: true,
		},
		ID:            w.getEnemyID(),
		FollowsPlayer: true,
		MoveSpeed:     0,
	}
}

func (w *WizardServer) updateEnemies() {
	for _, enemy := range w.enemies {
		enemy.Dx = 0
		enemy.Dy = 0

		if enemy.MoveSpeed < EnemyMoveSpeed {
			enemy.MoveSpeed += EnemyMoveSpeed / 150
		}

		if enemy.FollowsPlayer {
			// make enemy follow nearest living player
			var dx, dy, dist float64
			var minDx, minDy float64
			var minDist = 10000.0 // arbitrary big number
			for _, wizard := range w.wizards {
				if _, ok := w.deadWizards[wizard.Username]; ok {
					continue
				}
				dx = wizard.X - enemy.X
				dy = wizard.Y - enemy.Y
				dist = math.Hypot(dx, dy)
				if dist < minDist {
					minDist = dist
					minDx = dx
					minDy = dy
				}
			}

			closeEnough := 8.0 //8px

			if minDist > closeEnough {
				normX := minDx / minDist
				normY := minDy / minDist

				speed := enemy.MoveSpeed

				if enemy.MoveSpeed > minDist {
					speed = minDist
				}

				enemy.Dx = normX * speed
				enemy.Dy = normY * speed
			}
		}

		// apply knockback when hit
		if enemy.KnockbackDx != 0 || enemy.KnockbackDy != 0 {
			enemy.Dx += enemy.KnockbackDx
			enemy.Dy += enemy.KnockbackDy
			enemy.KnockbackDx = 0
			enemy.KnockbackDy = 0
		}
	}
}

func (w *WizardServer) checkEnemyCollisions() {
	for _, enemy := range w.enemies {
		// turn on noclip if enemy is knocked outside of the map
		if enemy.X < 0 ||
			enemy.Y < 0 ||
			enemy.X > float64(w.tilemapJSON.Layers[0].Width)*16.0-TileSize || // padding to prevent enemies getting stuck outside the map
			enemy.Y > float64(w.tilemapJSON.Layers[0].Height)*16.0-TileSize {
			enemy.Noclip = true
		} else if enemy.Noclip {
			// lets enemies noclip until they're fully out of the hole
			if !CheckCollisionHorizontal(enemy.Body, w.holes) &&
				!CheckCollisionVertical(enemy.Body, w.holes) {
				enemy.Noclip = false
			}
		}

		// turn noclip back on if enemy gets knocked into hole
		if CheckPointInRect(
			int(enemy.X+HalfTile),
			int(enemy.Y+14), // between 13 and 15, which are the y coords used in CheckCollisionHorizontal/Vertical
			w.holes,
		) {
			enemy.Noclip = true
		}

		enemy.X += enemy.Dx
		if !enemy.Noclip {
			CheckCollisionHorizontal(enemy.Body, w.holes)
		}
		enemy.Y += enemy.Dy
		if !enemy.Noclip {
			CheckCollisionVertical(enemy.Body, w.holes)
		}
	}
}
