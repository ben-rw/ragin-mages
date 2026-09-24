package shared

import "image/color"

var BackgroundColor = color.RGBA{0, 105, 150, 1}

const (
	defaultScreenWidth  = 640
	defaultScreenHeight = 360
)

var (
	ScreenWidth  float64 = defaultScreenWidth
	ScreenHeight float64 = defaultScreenHeight
)

func Center() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.5
}

func CenterRight() (float64, float64) {
	return ScreenWidth * 0.9, ScreenHeight * 0.5
}

func BottomCenter() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.9
}

func TopCenter() (float64, float64) {
	return ScreenWidth * 0.5, ScreenHeight * 0.05
}

func BottomRight() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.9
}

func BottomRightFurther() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.96
}

func TopRight() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.1
}

func TopRightFurther() (float64, float64) {
	return ScreenWidth * 0.99, ScreenHeight * 0.01
}

func TopCenterRight() (float64, float64) {
	return ScreenWidth * 0.81, ScreenHeight * 0.07
}

func TopLeftFurther() (float64, float64) {
	return ScreenWidth * 0.01, ScreenHeight * 0.01
}

func Stat3BottomLeft() (float64, float64) {
	return ScreenWidth * 0.01, ScreenHeight * 0.96
}

func Stat2BottomLeft() (float64, float64) {
	return ScreenWidth * 0.01, ScreenHeight * 0.91
}

func Stat1BottomLeft() (float64, float64) {
	return ScreenWidth * 0.01, ScreenHeight * 0.86
}

// duplicate of TopLeftFurther
func Heart1TopLeft() (float64, float64) {
	return ScreenWidth * 0.01, ScreenHeight * 0.01
}

func Heart2TopLeft() (float64, float64) {
	return ScreenWidth * 0.034, ScreenHeight * 0.01
}

func Heart3TopLeft() (float64, float64) {
	return ScreenWidth * 0.060, ScreenHeight * 0.01
}
