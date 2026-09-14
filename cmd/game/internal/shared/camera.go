package shared

import "math"

const FreeCamSpeed = 10

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) FollowTarget(targetX, targetY float64) {
	//center of the player
	destX := -(targetX + TileSize/2.0) + ScreenWidth/2.0
	destY := -(targetY + TileSize/2.0) + ScreenHeight/2.0

	//move towards the center of the player
	lerp := 0.1
	c.X += (destX - c.X) * lerp
	c.Y += (destY - c.Y) * lerp
}

func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, ScreenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, ScreenHeight-tilemapHeightPixels)
}
