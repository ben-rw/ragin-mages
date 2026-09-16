package shared

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

const EnemyRespawnTimer = 30

var EnemySpawns = map[int]struct{ X, Y float64 }{
	0: {X: 0, Y: 9},
	1: {X: 9, Y: 0},
	2: {X: 33, Y: 14},
	3: {X: 33, Y: 25},
	// 4: {X: 79, Y: 9},
	// 5: {X: 9, Y: 39},
	// 6: {X: 90, Y: 20},
	// 7: {X: 20, Y: 9},
}
