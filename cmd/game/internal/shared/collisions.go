package shared

import "image"

const (
	collisionBoxLeft   = 7
	collisionBoxRight  = 9
	collisionBoxTop    = 13
	collisionBoxBottom = 15
)

func CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) bool {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			// puts hitbox close to feet
			int(sprite.X)+collisionBoxLeft,
			int(sprite.Y)+collisionBoxTop,
			int(sprite.X)+collisionBoxRight,
			int(sprite.Y)+collisionBoxBottom,
		)) {
			if sprite.Noclip == false {
				if sprite.Dx > 0.0 {
					sprite.X = float64(collider.Min.X) - collisionBoxRight
				} else if sprite.Dx < 0.0 {
					sprite.X = float64(collider.Max.X) - collisionBoxLeft
				}
			}
			return true
		}
	}
	return false
}

func CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) bool {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			// puts hitbox close to feet
			int(sprite.X)+collisionBoxLeft,
			int(sprite.Y)+collisionBoxTop,
			int(sprite.X)+collisionBoxRight,
			int(sprite.Y)+collisionBoxBottom,
		)) {
			if sprite.Noclip == false {
				if sprite.Dy > 0.0 {
					sprite.Y = float64(collider.Min.Y) - collisionBoxBottom
				} else if sprite.Dy < 0.0 {
					sprite.Y = float64(collider.Max.Y) - collisionBoxTop
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
