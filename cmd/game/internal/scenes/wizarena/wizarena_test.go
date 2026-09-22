package wizarena

import (
	"testing"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
)

func testSceneInit() *WizArena {
	w := NewWizArena(nil)
	w.wizard = shared.NewWizard(shared.NewPlayer(&protocol.PlayerData{}, 0))

	pCount := 12
	x := 5
	y := 5
	w.wizard.X = 20
	w.wizard.Y = 20
	for range pCount {
		projectile := w.wizard.ShootProjectile(
			w.projectileImageCache[shared.Fireball],
			float64(x),
			float64(y),
			shared.Fireball,
		)
		projectile.X = float64(x + 1)
		projectile.Y = float64(y + 1)
		x += 10
		y += 10
		w.projectiles = append(w.projectiles, projectile)
	}

	w.camera = shared.NewCamera(
		-(w.wizard.X+shared.HalfTile)+shared.ScreenWidth/2.0,
		-(w.wizard.Y+shared.HalfTile)+shared.ScreenHeight/2.0,
	)
	return w
}

func BenchmarkDraw(b *testing.B) {
	w := testSceneInit()
	screen := ebiten.NewImage(int(shared.ScreenWidth), int(shared.ScreenHeight))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Draw(screen)
	}
}

func BenchmarkUpdate(b *testing.B) {
	w := testSceneInit()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Update([]protocol.Message{})
	}
}
