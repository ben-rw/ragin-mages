package wizServer

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizardServer) handleChestUpdate(message *protocol.Message) error {
	// chestUpdateData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := chestUpdateData.(*protocol.ChestUpdateData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }

	// w.BroadcastMsg(protocol.ChestUpdate, protocol.ChestUpdateData{})
	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleWizardMovementUpdate(message *protocol.Message) error {
	// wizardMovementUpdateData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := wizardMovementUpdateData.(*protocol.WizardMovementUpdateData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if _, ok := w.wizards[data.Username]; ok {
	// 	w.wizards[data.Username].X = data.X
	// 	w.wizards[data.Username].Y = data.Y
	// 	w.wizards[data.Username].Dx = data.Dx
	// 	w.wizards[data.Username].Dy = data.Dy
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleEnemyMovementUpdate(message *protocol.Message) error {
	// enemyMovementUpdateData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := enemyMovementUpdateData.(*protocol.EnemyMovementUpdateData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if _, ok := w.enemies[data.ID]; ok {
	// 	w.enemies[data.ID].X = data.X
	// 	w.enemies[data.ID].Y = data.Y
	// 	w.enemies[data.ID].Dx = data.Dx
	// 	w.enemies[data.ID].Dy = data.Dy
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleWizArenaMovementUpdate(message *protocol.Message) error {
	// wizArenaMovementData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := wizArenaMovementData.(*protocol.WizArenaMovementData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// for _, wizardData := range data.WizardsMovementData {
	// 	if _, ok := w.wizards[wizardData.Username]; ok {
	// 		w.wizards[wizardData.Username].X = wizardData.X
	// 		w.wizards[wizardData.Username].Y = wizardData.Y
	// 		w.wizards[wizardData.Username].Dx = wizardData.Dx
	// 		w.wizards[wizardData.Username].Dy = wizardData.Dy
	// 	}
	// }
	//
	// for _, enemyData := range data.EnemiesMovementData {
	// 	if _, ok := w.enemies[enemyData.ID]; ok {
	// 		w.enemies[enemyData.ID].X = enemyData.X
	// 		w.enemies[enemyData.ID].Y = enemyData.Y
	// 		w.enemies[enemyData.ID].Dx = enemyData.Dx
	// 		w.enemies[enemyData.ID].Dy = enemyData.Dy
	// 	}
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleWizardStatUpdate(message *protocol.Message) error {
	// wizardStatUpdateData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := wizardStatUpdateData.(*protocol.WizardStatUpdateData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if _, ok := w.wizards[data.Username]; ok {
	// 	w.wizards[data.Username].Data.Score = data.Score
	// 	w.wizards[data.Username].Combat.SetAttackCooldown(data.AttackCooldown)
	// 	w.wizards[data.Username].Combat.SetProjectileSpeed(data.ProjectileSpeed)
	// 	w.wizards[data.Username].Combat.SetProjectileScale(data.ProjectileScale)
	// 	w.wizards[data.Username].Combat.SetKnockback(data.Knockback)
	// 	w.wizards[data.Username].Combat.SetHealth(data.Health)
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleEnemyStatUpdate(message *protocol.Message) error {
	// enemyStatUpdate, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := enemyStatUpdate.(*protocol.EnemyStatUpdateData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if _, ok := w.enemies[data.ID]; ok {
	// 	w.enemies[data.ID].Combat.SetHealth(data.Health)
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleNewProjectile(message *protocol.Message) error {
	// newProjectileData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := newProjectileData.(*protocol.NewProjectileData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if wizard, ok := w.wizards[data.CasterName]; ok {
	// 	wizard.NewProjectile(
	// 		w.projectileImageCache[Fireball],
	// 		float64(data.CursorX),
	// 		float64(data.CursorY),
	// 		Fireball,
	// 	)
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleProjectileReflected(message *protocol.Message) error {
	// projectileReflectedData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := projectileReflectedData.(*protocol.ProjectileReflectedData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if wizard, ok := w.wizards[data.ReflectorName]; ok {
	// 	wizard.ReflectProjectile(
	// 		w.projectiles[data.ProjectileID],
	// 		float64(data.CursorX),
	// 		float64(data.CursorY),
	// 	)
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleProjectileHitWizard(message *protocol.Message) error {
	// projectileHitWizardData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := projectileHitWizardData.(*protocol.ProjectileHitWizardData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if shooter, ok := w.wizards[data.ShooterStatData.Username]; ok {
	// 	shooter.Data.Score = data.ShooterStatData.Score
	// 	shooter.Combat.SetAttackCooldown(data.ShooterStatData.AttackCooldown)
	// 	shooter.Combat.SetProjectileSpeed(data.ShooterStatData.ProjectileSpeed)
	// 	shooter.Combat.SetProjectileScale(data.ShooterStatData.ProjectileScale)
	// 	shooter.Combat.SetKnockback(data.ShooterStatData.Knockback)
	// }
	//
	// if victim, ok := w.wizards[data.VictimStatData.Username]; ok {
	// 	victim.Data.Score = data.VictimStatData.Score
	// 	victim.Combat.SetAttackCooldown(data.VictimStatData.AttackCooldown)
	// 	victim.Combat.SetProjectileSpeed(data.VictimStatData.ProjectileSpeed)
	// 	victim.Combat.SetProjectileScale(data.VictimStatData.ProjectileScale)
	// 	victim.Combat.SetKnockback(data.VictimStatData.Knockback)
	// 	victim.Combat.SetHealth(data.VictimStatData.Health)
	//
	// 	victim.KnockbackDx = data.VictimDx
	// 	victim.KnockbackDy = data.VictimDy
	// }

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleProjectileHitEnemy(message *protocol.Message) error {
	// projectileHitEnemyData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := projectileHitEnemyData.(*protocol.ProjectileHitEnemyData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }
	//
	// if shooter, ok := w.wizards[data.ShooterStatData.Username]; ok {
	// 	shooter.Data.Score = data.ShooterStatData.Score
	// 	shooter.Combat.SetAttackCooldown(data.ShooterStatData.AttackCooldown)
	// 	shooter.Combat.SetProjectileSpeed(data.ShooterStatData.ProjectileSpeed)
	// 	shooter.Combat.SetProjectileScale(data.ShooterStatData.ProjectileScale)
	// 	shooter.Combat.SetKnockback(data.ShooterStatData.Knockback)
	// }
	//
	// if victim, ok := w.enemies[data.VictimStatData.ID]; ok {
	// 	victim.Combat.SetHealth(data.VictimStatData.Health)
	// }

	w.Broadcast(message)

	return nil
}
