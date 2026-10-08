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
