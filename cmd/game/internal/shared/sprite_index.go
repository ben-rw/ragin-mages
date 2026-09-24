package shared

var imgRootPath string = "assets/images/ninja_adventure/Actor/Character/"
var imgPathEnd string = "/SpriteSheet.png"

var PlayerSpriteIndex map[int]string = map[int]string{
	0: imgRootPath + "OldMan3" + imgPathEnd,
	1: imgRootPath + "Monk" + imgPathEnd,
	2: imgRootPath + "Samurai" + imgPathEnd,
	3: imgRootPath + "Sultan" + imgPathEnd,
	4: imgRootPath + "Master" + imgPathEnd,
	5: imgRootPath + "NinjaMageBlack" + imgPathEnd,
	6: imgRootPath + "MaskRaccoon" + imgPathEnd,
	7: imgRootPath + "MaskFrog" + imgPathEnd,
}

var StartingPositions = map[int]struct{ X, Y float64 }{
	0: {X: 18, Y: 20},
	1: {X: 18, Y: 56},
	2: {X: 18, Y: 92},
	3: {X: 18, Y: 128},
	4: {X: 54, Y: 20},
	5: {X: 54, Y: 56},
	6: {X: 54, Y: 92},
	7: {X: 54, Y: 128},
}

var EnemySpriteIndex map[EnemyType]string = map[EnemyType]string{
	Skeleton: imgRootPath + "Skeleton" + imgPathEnd,
}
