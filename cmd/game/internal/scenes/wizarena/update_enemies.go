package wizarena

import (
	"math"

	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

func (w *WizArena) UpdateEnemies() {
	for _, enemy := range w.enemies {
		enemy.Combat.Update()
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer {
			// make enemy follow nearest living player
			var dx, dy, dist float64
			var minDx, minDy float64
			var minDist = 10000.0 // arbitrary big number
			for _, w := range w.wizards {
				if w.Combat.Dead {
					continue
				}
				dx = w.X - enemy.X
				dy = w.Y - enemy.Y
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

				speed := enemy.Combat.MoveSpeed()

				if enemy.Combat.MoveSpeed() > minDist {
					speed = minDist
				}

				enemy.Dx = normX * speed
				enemy.Dy = normY * speed
			}
		}
	}
}

func (w *WizArena) CheckEnemyCollisions() {
	for _, enemy := range w.enemies {
		// turn on noclip if enemy is knocked outside of the map
		if enemy.X < 0 ||
			enemy.Y < 0 ||
			enemy.X > float64(w.tilemapJSON.Layers[0].Width)*16.0-wizServer.TileSize || // padding to prevent enemies getting stuck outside the map
			enemy.Y > float64(w.tilemapJSON.Layers[0].Height)*16.0-wizServer.TileSize {
			enemy.Noclip = true
		} else if enemy.Noclip {
			// lets enemies noclip until they're fully out of the hole
			if !CheckCollisionHorizontal(enemy.Sprite, w.holes) &&
				!CheckCollisionVertical(enemy.Sprite, w.holes) {
				enemy.Noclip = false
			}
		}

		// turn noclip back on if enemy gets knocked into hole
		if CheckPointInRect(
			int(enemy.X+wizServer.HalfTile),
			int(enemy.Y+14), // between 13 and 15, which are the y coords used in CheckCollisionHorizontal/Vertical
			w.holes,
		) {
			enemy.Noclip = true
		}

		// fade in and ramp movement speed after spawn
		if enemy.Alpha < 1.0 {
			enemy.Alpha += 0.01
		}
		if enemy.Combat.MoveSpeed() < EnemyMoveSpeed {
			enemy.Combat.SetMoveSpeed(enemy.Combat.MoveSpeed() + EnemyMoveSpeed/150)
		}

		enemy.X += enemy.Dx
		if !enemy.Noclip {
			CheckCollisionHorizontal(enemy.Sprite, w.holes)
		}
		enemy.Y += enemy.Dy
		if !enemy.Noclip {
			CheckCollisionVertical(enemy.Sprite, w.holes)
		}
	}
}
