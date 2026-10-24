package wizarena

import (
	"errors"
	"time"

	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizArena) HandleWizardMovementUpdate(message *protocol.Message) error {
	wizardMovementUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := wizardMovementUpdateData.(*protocol.WizardMovementUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if _, ok := w.wizards[data.Username]; ok {
		w.wizards[data.Username].X = data.X
		w.wizards[data.Username].Y = data.Y
		w.wizards[data.Username].Dx = data.Dx
		w.wizards[data.Username].Dy = data.Dy
	}

	return nil
}

func (w *WizArena) HandleEnemyMovementUpdate(message *protocol.Message) error {
	enemyMovementUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := enemyMovementUpdateData.(*protocol.EnemyMovementUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if _, ok := w.enemies[data.ID]; ok {
		w.enemies[data.ID].X = data.X
		w.enemies[data.ID].Y = data.Y
		w.enemies[data.ID].Dx = data.Dx
		w.enemies[data.ID].Dy = data.Dy
	}

	return nil
}

func (w *WizArena) HandleWizArenaMovementUpdate(message *protocol.Message) error {
	wizArenaMovementData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := wizArenaMovementData.(*protocol.WizArenaMovementData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	for _, wizardData := range data.WizardsMovementData {
		if _, ok := w.wizards[wizardData.Username]; ok {
			if w.wizards[wizardData.Username].Data.Username == w.wizard.Data.Username &&
				w.phase == wizServer.Playing {
				continue
			}
			w.wizards[wizardData.Username].X = wizardData.X
			w.wizards[wizardData.Username].Y = wizardData.Y
			w.wizards[wizardData.Username].Dx = wizardData.Dx
			w.wizards[wizardData.Username].Dy = wizardData.Dy
		}
	}

	for _, enemyData := range data.EnemiesMovementData {
		if _, ok := w.enemies[enemyData.ID]; ok {
			w.enemies[enemyData.ID].X = enemyData.X
			w.enemies[enemyData.ID].Y = enemyData.Y
			w.enemies[enemyData.ID].Dx = enemyData.Dx
			w.enemies[enemyData.ID].Dy = enemyData.Dy
		} else if !ok {
			// add newly spawned enemies
			w.enemies[enemyData.ID] = NewEnemy(
				w.enemyImageCache[Skeleton],
				Skeleton,
				true,
				enemyData.X,
				enemyData.Y,
				enemyData.ID,
			)
		}
	}

	return nil
}

func (w *WizArena) HandleWizardStatUpdate(message *protocol.Message) error {
	wizardStatUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := wizardStatUpdateData.(*protocol.WizardStatUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if wizard, ok := w.wizards[data.Username]; ok {
		if wizard.Combat.Health() > data.Health {
			wizard.Combat.Damage(wizard.Combat.Health() - data.Health)
			if wizard.Combat.Health() <= 0 {
				wizard.Combat.Dead = true
				wizard.Alpha = 1.0
				wizard.Combat.SetIFrames(0)
				wizard.DeathX = wizard.X
				wizard.DeathY = wizard.Y
				wizard.DieAnim = true //starts death animation
			}
		}
		wizard.Data.Score = data.Score
		wizard.Combat.SetAttackCooldown(data.AttackCooldown)
		wizard.Combat.SetProjectileSpeed(data.ProjectileSpeed)
		wizard.Combat.SetProjectileScale(data.ProjectileScale)
		wizard.Combat.SetKnockback(data.Knockback)
		wizard.Combat.SetHealth(data.Health)
	}

	w.scoresUpdatedCount.Add(1)

	return nil
}

func (w *WizArena) HandleEnemyStatUpdate(message *protocol.Message) error {
	enemyStatUpdate, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := enemyStatUpdate.(*protocol.EnemyStatUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if enemy, ok := w.enemies[data.ID]; ok {
		enemy.Combat.Damage(enemy.Combat.Health() - data.Health)
	}

	return nil
}

func (w *WizArena) HandleNewProjectile(message *protocol.Message) error {
	newProjectileData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := newProjectileData.(*protocol.NewProjectileData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if wizard, ok := w.wizards[data.CasterName]; ok {
		projectile := wizard.NewProjectile(
			w.projectileImageCache[Fireball],
			float64(data.CursorX),
			float64(data.CursorY),
			Fireball,
		)
		w.projectiles[projectile.ID] = projectile
		wizard.Combat.Attack()
		wizard.AttackAnim = true
	}

	return nil
}

func (w *WizArena) HandleReflecting(message *protocol.Message) error {
	reflectingData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := reflectingData.(*protocol.ReflectingData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if wizard, ok := w.wizards[data.ReflectorName]; ok {
		wizard.Combat.Reflect()
		wizard.ReflectAnim = true
		wizard.Reflect.Active = true
		wizard.Reflect.cursorX = float64(data.CursorX)
		wizard.Reflect.cursorY = float64(data.CursorY)
	}

	return nil
}

func (w *WizArena) HandleProjectileReflected(message *protocol.Message) error {
	projectileReflectedData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileReflectedData.(*protocol.ProjectileReflectedData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if wizard, ok := w.wizards[data.ReflectorName]; ok {
		wizard.ReflectProjectile(
			w.projectiles[data.ProjectileID],
		)
		wizard.Combat.ResetReflectCooldown()
	}

	return nil
}

func (w *WizArena) HandleProjectileDestroyed(message *protocol.Message) error {
	projectileDestroyedData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileDestroyedData.(*protocol.ProjectileDestroyedData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if _, ok := w.projectiles[data.ProjectileID]; ok {
		delete(w.projectiles, w.projectiles[data.ProjectileID].ID)
	}

	return nil
}

func (w *WizArena) HandleProjectileHitWizard(message *protocol.Message) error {
	projectileHitWizardData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileHitWizardData.(*protocol.ProjectileHitWizardData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if shooter, ok := w.wizards[data.ShooterStatData.Username]; ok {
		shooter.Data.Score = data.ShooterStatData.Score
		shooter.Combat.SetAttackCooldown(data.ShooterStatData.AttackCooldown)
		shooter.Combat.SetProjectileSpeed(data.ShooterStatData.ProjectileSpeed)
		shooter.Combat.SetProjectileScale(data.ShooterStatData.ProjectileScale)
		shooter.Combat.SetKnockback(data.ShooterStatData.Knockback)
	}

	if victim, ok := w.wizards[data.VictimStatData.Username]; ok {
		if victim.Combat.Health() > data.VictimStatData.Health {
			victim.Combat.Damage(victim.Combat.Health() - data.VictimStatData.Health)
			if victim.Combat.Health() <= 0 {
				victim.Combat.Dead = true
				victim.Alpha = 1.0
				victim.Combat.SetIFrames(0)
				victim.DeathX = victim.X
				victim.DeathY = victim.Y
				victim.DieAnim = true //starts death animation
			}
		}

		victim.Data.Score = data.VictimStatData.Score
		victim.Combat.SetAttackCooldown(data.VictimStatData.AttackCooldown)
		victim.Combat.SetProjectileSpeed(data.VictimStatData.ProjectileSpeed)
		victim.Combat.SetProjectileScale(data.VictimStatData.ProjectileScale)
		victim.Combat.SetKnockback(data.VictimStatData.Knockback)
		victim.Combat.SetHealth(data.VictimStatData.Health)

		victim.KnockbackDx = data.VictimDx
		victim.KnockbackDy = data.VictimDy

		victim.Combat.LastHitBy = data.ShooterStatData.Username
	}

	return nil
}

func (w *WizArena) HandleProjectileHitEnemy(message *protocol.Message) error {
	projectileHitEnemyData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileHitEnemyData.(*protocol.ProjectileHitEnemyData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if shooter, ok := w.wizards[data.ShooterStatData.Username]; ok {
		shooter.Data.Score = data.ShooterStatData.Score
		shooter.Combat.SetAttackCooldown(data.ShooterStatData.AttackCooldown)
		shooter.Combat.SetProjectileSpeed(data.ShooterStatData.ProjectileSpeed)
		shooter.Combat.SetProjectileScale(data.ShooterStatData.ProjectileScale)
		shooter.Combat.SetKnockback(data.ShooterStatData.Knockback)
	}

	if victim, ok := w.enemies[data.VictimStatData.ID]; ok {
		victim.Combat.SetHealth(data.VictimStatData.Health)
	}

	return nil
}

func (w *WizArena) HandleChestUpdate(message *protocol.Message) error {
	chestUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := chestUpdateData.(*protocol.ChestUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if data.X == 0 && data.Y == 0 {
		clear(w.openedChests)
	} else {
		w.openedChests[OpenedChest{X: data.X, Y: data.Y}] = struct{}{}
	}

	return nil
}

func (w *WizArena) HandleRoundTimerUpdate(message *protocol.Message) error {
	roundTimerUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := roundTimerUpdateData.(*protocol.RoundTimerUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.phaseTimer = time.Now().Add(time.Duration(data.TimeRemaining) * time.Millisecond)

	return nil
}

func (w *WizArena) HandlePhaseUpdate(message *protocol.Message) error {
	phaseUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := phaseUpdateData.(*protocol.PhaseUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.phase = wizServer.RoundPhase(data.Phase)

	if w.phase == wizServer.SyncScores {
		w.WriteWizardStatUpdate()
		w.scoresUpdatedCount.Store(0)
	}

	return nil
}

func (w *WizArena) HandleGivePoints(message *protocol.Message) error {
	givePointsData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := givePointsData.(*protocol.GivePointsData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	if wizard, ok := w.wizards[data.RecipientName]; ok {
		wizard.Combat.RandomBoost(float64(data.Random), StandardMult)
		wizard.Combat.BoostProjectileScale(float64(data.ProjectileScale), 1)
		wizard.Combat.BoostProjectileSpeed(float64(data.ProjectileSpeed), 1)
		wizard.Combat.BoostKnockback(float64(data.Knockback), 1)
		wizard.Combat.BoostAttackCooldown(float64(data.AttackCooldown), 1)
	}

	return nil
}
