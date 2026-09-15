package shared

import "image"

func CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - TileSize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - TileSize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func CheckCollisionCircle(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x1 - x2
	dy := y1 - y2

	rSum := r1 + r2

	return (dx*dx + dy*dy) < (rSum * rSum)
}
