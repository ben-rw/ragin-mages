package wizarena

import (
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

// rounds - round start, screen black, fades in over 3 seconds - ticking down the alpha each update loop
// if only 1 player alive or time runs out, round ends

// need a round countdown screen with a small delay

type RoundPhase int

const (
	Playing RoundPhase = iota
	RoundOver
	RoundStart
	GameOver
	GameStart
)

const (
	Rounds         = 3   // 3 rounds
	RoundTime      = 150 // 2.5 min per round
	RoundOverTime  = 3   // 3 seconds displaying scoreboard
	RoundStartTime = 1   // 1 seconds displaying round number
	GameOverTime   = 5   // 8 (3 from intermission + 5 additional) seconds at the end of the game
	GameStartTime  = 1
)

// reset enemies, projectiles, chests, timers, put players in a random spawn location
func (w *WizArena) startNewRound() {
	for _, wizard := range w.wizards {
		wizard.Combat.Dead = false
		wizard.Combat.Fell = false
		wizard.Combat.SetHealth(shared.PlayerHealth)
		wizard.Combat.SetIFrames(0)
		wizard.Alpha = 1.0
	}

	clear(w.openedChests)
	w.enemies = make([]*shared.Enemy, 0, 32)
	w.projectiles = make([]*shared.Projectile, 0, 16)

	w.phaseTimer = time.Now().Add(RoundTime * time.Second)
	w.enemyRespawnTimer = time.Now().Add(shared.RoundStartEnemySpawnTimer * time.Second)
	w.chestRespawnTimer = time.Now().Add(ChestRespawnTimer * time.Second)

	w.phase = Playing

	w.backgroundAlpha = 0
}

func (w *WizArena) updateRound() {
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 || w.wizard.Combat.Dead {
		w.round += 1

		w.phase = RoundOver
		w.phaseTimer = time.Now().Add(RoundOverTime * time.Second)
	}
}

func (w *WizArena) updateRoundStart() {
	w.AssignPlayerSpawns()

	w.wizard.Alpha = 1.0

	w.wizard.Dx = 0
	w.wizard.Dy = 0

	w.wizard.ActiveAnimation = w.wizard.GetActiveAnimation()
	w.wizard.ActiveAnimation.Update()

	w.wizard.NameTag.X = w.wizard.X + shared.TileSize/2
	w.wizard.NameTag.Y = w.wizard.Y + shared.TileSize + 2

	w.camera.FollowTarget(
		w.wizard.X,
		w.wizard.Y,
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)
	w.camera.Constrain(
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.startNewRound()
	}
}

func (w *WizArena) updateRoundOver() {
	w.NewScoreboard()

	// add a fade effect when players die
	if w.wizard.Combat.Dead {
		if !w.wizard.Combat.Fell {
			if w.wizard.Y > -shared.HalfTile {
				w.wizard.Y -= 3.0
			}
		}
	}
	if w.wizard.Alpha > 0.0 {
		w.wizard.Alpha -= .02 //fade speed
	} else {
		w.wizard.Alpha = 0.0
	}

	w.wizard.ActiveAnimation = w.wizard.GetActiveAnimation()
	w.wizard.ActiveAnimation.Update()

	w.wizard.NameTag.X = w.wizard.X + shared.TileSize/2
	w.wizard.NameTag.Y = w.wizard.Y + shared.TileSize + 2

	for _, enemy := range w.enemies {
		if enemy.Combat.Dead {
			enemy.DieAnim = true
		}
		enemy.ActiveAnimation = enemy.GetActiveAnimation()
		enemy.ActiveAnimation.Update()

		enemy.X += enemy.Dx / 2 // slow motion woahhhh
		enemy.Y += enemy.Dy / 2

		if enemy.Alpha > 0.0 {
			enemy.Alpha -= .01 //fade speed
		}

	}

	for _, projectile := range w.projectiles {
		projectile.ActiveAnimation.Update()

		projectile.X += projectile.Dx / 2 // slow motion woahhhh
		projectile.Y += projectile.Dy / 2

		if projectile.Alpha > 0.0 {
			projectile.Alpha -= .01 //fade speed
		}
	}

	for _, wizard := range w.wizards {
		wizard.Reflect.ActiveAnimation = wizard.Reflect.GetActiveAnimation()

		if wizard.Reflect.ActiveAnimation != nil {
			wizard.Reflect.ActiveAnimation.Update()

			wizard.Reflect.X = wizard.X + shared.HalfTile
			wizard.Reflect.Y = wizard.Y + shared.HalfTile
		}
	}

	w.UpdateTraps()

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		if w.round > Rounds {
			w.phase = GameOver
			w.phaseTimer = time.Now().Add(GameOverTime * time.Second)
			w.gameOver()
			return
		}

		w.phase = RoundStart
		w.phaseTimer = time.Now().Add(RoundStartTime * time.Second)
	}
}

func (w *WizArena) gameStart() {
	w.AssignPlayerSpawns()

	w.camera.FollowTarget(
		w.wizard.X,
		w.wizard.Y,
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)
	w.camera.Constrain(
		float64(w.tilemapJSON.Layers[0].Width)*16.0,
		float64(w.tilemapJSON.Layers[0].Height)*16.0,
	)

	w.backgroundAlpha = 255
	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.startNewRound()
	}
}

func (w *WizArena) gameOver() {
	w.NewScoreboard()
	if time.Until(w.phaseTimer) <= 2*time.Second {
		_ = sound.FadeOut(w.audioPlayer)
	}

	if w.phaseTimer.IsZero() || time.Until(w.phaseTimer) <= 0 {
		w.Conn.WriteMsg(
			protocol.SceneChange,
			&protocol.SceneChangeData{
				SceneType: protocol.LobbyScene,
			},
		)
	}
}
