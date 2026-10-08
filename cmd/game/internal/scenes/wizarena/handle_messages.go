package wizarena

import (
	"errors"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
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

	if _, ok := w.wizards[data.Name]; ok {
		w.wizards[data.Name].X = data.X
		w.wizards[data.Name].Y = data.Y
		w.wizards[data.Name].Dx = data.Dx
		w.wizards[data.Name].Dy = data.Dy
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
		if _, ok := w.wizards[wizardData.Name]; ok {
			w.wizards[wizardData.Name].X = wizardData.X
			w.wizards[wizardData.Name].Y = wizardData.Y
			w.wizards[wizardData.Name].Dx = wizardData.Dx
			w.wizards[wizardData.Name].Dy = wizardData.Dy
		}
	}

	for _, enemyData := range data.EnemiesMovementData {
		if _, ok := w.enemies[enemyData.ID]; ok {
			w.enemies[enemyData.ID].X = enemyData.X
			w.enemies[enemyData.ID].Y = enemyData.Y
			w.enemies[enemyData.ID].Dx = enemyData.Dx
			w.enemies[enemyData.ID].Dy = enemyData.Dy
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

	if _, ok := w.wizards[data.Name]; ok {
		w.wizards[data.Name].Data.Score = data.Score
		w.wizards[data.Name].Combat.SetAttackCooldown(data.AttackCooldown)
		w.wizards[data.Name].Combat.SetProjectileSpeed(data.ProjectileSpeed)
		w.wizards[data.Name].Combat.SetProjectileScale(data.ProjectileScale)
		w.wizards[data.Name].Combat.SetKnockback(data.Knockback)
		w.wizards[data.Name].Combat.SetHealth(data.Health)
	}

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

	if _, ok := w.enemies[data.ID]; ok {
		w.enemies[data.ID].Combat.SetHealth(data.Health)
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
		wizard.NewProjectile(
			w.projectileImageCache[shared.Fireball],
			float64(data.CursorX),
			float64(data.CursorY),
			shared.Fireball,
		)
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
			float64(data.CursorX),
			float64(data.CursorY),
		)
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

	if shooter, ok := w.wizards[data.ShooterStatData.Name]; ok {
		shooter.Data.Score = data.ShooterStatData.Score
		shooter.Combat.SetAttackCooldown(data.ShooterStatData.AttackCooldown)
		shooter.Combat.SetProjectileSpeed(data.ShooterStatData.ProjectileSpeed)
		shooter.Combat.SetProjectileScale(data.ShooterStatData.ProjectileScale)
		shooter.Combat.SetKnockback(data.ShooterStatData.Knockback)
	}

	if victim, ok := w.wizards[data.VictimStatData.Name]; ok {
		victim.Data.Score = data.VictimStatData.Score
		victim.Combat.SetAttackCooldown(data.VictimStatData.AttackCooldown)
		victim.Combat.SetProjectileSpeed(data.VictimStatData.ProjectileSpeed)
		victim.Combat.SetProjectileScale(data.VictimStatData.ProjectileScale)
		victim.Combat.SetKnockback(data.VictimStatData.Knockback)
		victim.Combat.SetHealth(data.VictimStatData.Health)

		victim.KnockbackDx = data.VictimDx
		victim.KnockbackDy = data.VictimDy
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

	if shooter, ok := w.wizards[data.ShooterStatData.Name]; ok {
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
