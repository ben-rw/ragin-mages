package shared

var ImgRootPath string = "assets/images/ninja_adventure/Actor/Character/"
var ImgPathEnd string = "/SpriteSheet.png"

var PlayerSpriteIndex map[int]string = map[int]string{
	0: ImgRootPath + "OldMan3" + ImgPathEnd,
	1: ImgRootPath + "Monk" + ImgPathEnd,
	2: ImgRootPath + "Samurai" + ImgPathEnd,
	3: ImgRootPath + "Sultan" + ImgPathEnd,
	4: ImgRootPath + "Master" + ImgPathEnd,
	5: ImgRootPath + "NinjaMageBlack" + ImgPathEnd,
	6: ImgRootPath + "MaskRaccoon" + ImgPathEnd,
	7: ImgRootPath + "MaskFrog" + ImgPathEnd,
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
