package memory

import (
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
)

// starts bottom center, to bottom left, then top left
var PlayerPositions = map[int]struct{ X, Y float64 }{
	0: {X: shared.ScreenWidth * .5, Y: shared.ScreenHeight * .9},
	1: {X: shared.ScreenWidth * .25, Y: shared.ScreenHeight * .9},
	2: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .9},
	3: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .74},
	4: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .58},
	5: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .42},
	6: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .26},
	7: {X: shared.ScreenWidth * .1, Y: shared.ScreenHeight * .1},
}

// walks player from one place to another at speed proportionate
// to the initial distance between the player's location and the
// destination with no input from the player
func (m *Memory) ScriptedWalk(destX, destY float64) {
	distX := destX - m.Player.X
	stepX := distX / 100000
	distY := destY - m.Player.Y
	stepY := distY / 100000
	for destX != m.Player.X || destY != m.Player.Y {
		if destX-m.Player.X > stepX {
			m.Player.X = destX
		} else {
			m.Player.X += stepX
		}
		if destY-m.Player.Y > stepY {
			m.Player.Y = destY
		} else {
			m.Player.Y += stepY
		}
	}
}

// put player at the back of the line after thier turn, remove furthest positions
// based on number of eliminated players, making the line shorter
func (m *Memory) getInLine(totalPlayers, eliminatedPlayers int) {
	lineIndex := totalPlayers - eliminatedPlayers
	m.ScriptedWalk(PlayerPositions[lineIndex].X, PlayerPositions[lineIndex].Y)
}
