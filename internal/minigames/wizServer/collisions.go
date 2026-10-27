package wizServer

import (
	"image"
)

const (
	collisionBoxLeft   = 7
	collisionBoxRight  = 9
	collisionBoxTop    = 13
	collisionBoxBottom = 15
)

func CheckCollisionHorizontal(body *Body, colliders []image.Rectangle) bool {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			// puts hitbox close to feet
			int(body.X)+collisionBoxLeft,
			int(body.Y)+collisionBoxTop,
			int(body.X)+collisionBoxRight,
			int(body.Y)+collisionBoxBottom,
		)) {
			if body.Noclip == false {
				if body.Dx > 0.0 {
					body.X = float64(collider.Min.X) - collisionBoxRight
				} else if body.Dx < 0.0 {
					body.X = float64(collider.Max.X) - collisionBoxLeft
				}
			}
			return true
		}
	}
	return false
}

func CheckCollisionVertical(body *Body, colliders []image.Rectangle) bool {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			// puts hitbox close to feet
			int(body.X)+collisionBoxLeft,
			int(body.Y)+collisionBoxTop,
			int(body.X)+collisionBoxRight,
			int(body.Y)+collisionBoxBottom,
		)) {
			if body.Noclip == false {
				if body.Dy > 0.0 {
					body.Y = float64(collider.Min.Y) - collisionBoxBottom
				} else if body.Dy < 0.0 {
					body.Y = float64(collider.Max.Y) - collisionBoxTop
				}
			}
			return true
		}
	}
	return false
}

func CheckCollisionCircle(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x1 - x2
	dy := y1 - y2

	rSum := r1 + r2

	return (dx*dx + dy*dy) < (rSum * rSum)
}

func CheckPointInRect(x, y int, rect []image.Rectangle) bool {
	point := image.Point{x, y}
	for _, rect := range rect {
		if point.In(rect) {
			return true
		}
	}
	return false
}
