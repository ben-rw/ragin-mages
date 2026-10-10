package wizServer

import "math/rand"

type Wizard struct {
	Username        string
	X               float64
	Y               float64
	Dx              float64
	Dy              float64
	AttackCooldown  float64
	ProjectileSpeed float64
	ProjectileScale float64
	Knockback       float64
	Health          float64
	Score           int
}

// ordered to space out the spawns of the first 4 players
var PlayerSpawns = map[int]struct{ X, Y float64 }{
	0: {X: 9, Y: 20},
	1: {X: 90, Y: 39},
	2: {X: 20, Y: 50},
	3: {X: 79, Y: 50},
	4: {X: 79, Y: 9},
	5: {X: 9, Y: 39},
	6: {X: 90, Y: 20},
	7: {X: 20, Y: 9},
}

// gives players random, non-overlapping spawns
func (w *WizardServer) AssignPlayerSpawns() {
	pLen := len(w.wizards)
	playerSpawnIndex := make([]int, pLen)
	for i := range pLen {
		playerSpawnIndex[i] = i
	}

	for _, player := range w.wizards {
		i := rand.Intn(len(playerSpawnIndex))
		s := playerSpawnIndex
		spawnIndex := s[i]
		s[i] = s[len(s)-1]
		playerSpawnIndex = s[:len(s)-1]

		player.X = PlayerSpawns[spawnIndex].X * TileSize
		player.Y = PlayerSpawns[spawnIndex].Y * TileSize
	}
}
