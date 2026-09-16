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

func (c *Camera) FollowTarget(targetX, targetY, tilemapWidthPixels, tilemapHeightPixels float64) {
	//center of the player
	destX := -(targetX + HalfTile) + ScreenWidth/2.0
	destY := -(targetY + HalfTile) + ScreenHeight/2.0

	minX := ScreenWidth - tilemapWidthPixels
	minY := ScreenHeight - tilemapHeightPixels

	//move towards the center of the player
	lerp := 0.1
	if destX <= 0.0 && destX >= minX {
		c.X += (destX - c.X) * lerp
	}
	if destY <= 0.0 && destY >= minY {
		c.Y += (destY - c.Y) * lerp
	}
}

func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, ScreenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, ScreenHeight-tilemapHeightPixels)
}
