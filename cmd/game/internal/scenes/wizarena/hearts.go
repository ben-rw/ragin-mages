package wizarena

import "github.com/ben-rw/ragin-mages/cmd/game/internal/shared"

var HealthHeartLocations = map[int]func() (float64, float64){
	0: shared.Heart1TopLeft,
	1: shared.Heart2TopLeft,
	2: shared.Heart3TopLeft,
}
