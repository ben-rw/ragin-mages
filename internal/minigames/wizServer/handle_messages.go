package wizServer

import (
	"errors"

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
		w.Mu.Lock()
		w.wizards[data.Username].X = data.X
		w.wizards[data.Username].Y = data.Y
		w.wizards[data.Username].Dx = data.Dx
		w.wizards[data.Username].Dy = data.Dy
		w.Mu.Unlock()
	}

	return nil
}

// func (w *WizardServer) handleEnemyMovementUpdate(message *protocol.Message) error {
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

// 	w.Broadcast(message)
//
// 	return nil
// }

// func (w *WizardServer) handleWizArenaMovementUpdate(message *protocol.Message) error {
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

// 	w.Broadcast(message)
//
// 	return nil
// }

func (w *WizardServer) handleWizardStatUpdate(message *protocol.Message) error {
	wizardStatUpdateData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := wizardStatUpdateData.(*protocol.WizardStatUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.Mu.Lock()
	if _, ok := w.wizards[data.Username]; ok {
		w.wizards[data.Username].Health = data.Health
		if w.wizards[data.Username].Health <= 0 {
			w.deadWizards[data.Username] = struct{}{}
		}
	}
	w.Mu.Unlock()

	w.BroadcastExcludeSender(message, data.Username)

	return nil
}

func (w *WizardServer) handleEnemyStatUpdate(message *protocol.Message) error {
	enemyStatUpdate, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := enemyStatUpdate.(*protocol.EnemyStatUpdateData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.Mu.Lock()
	if _, ok := w.enemies[data.ID]; ok {
		w.enemies[data.ID].Health = data.Health
		if w.enemies[data.ID].Health <= 0 {
			delete(w.enemies, data.ID)
		}
	}
	w.Mu.Unlock()

	w.Broadcast(message)

	return nil
}

func (w *WizardServer) handleNewProjectile(message *protocol.Message) error {
	newProjectileData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := newProjectileData.(*protocol.NewProjectileData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.BroadcastExcludeSender(message, data.CasterName)

	return nil
}

func (w *WizardServer) handleProjectileReflected(message *protocol.Message) error {
	projectileReflectedData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileReflectedData.(*protocol.ProjectileReflectedData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.BroadcastExcludeSender(message, data.ReflectorName)

	return nil
}

func (w *WizardServer) handleReflecting(message *protocol.Message) error {
	reflectingData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := reflectingData.(*protocol.ReflectingData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.wizards[data.ReflectorName].startReflecting()

	w.BroadcastExcludeSender(message, data.ReflectorName)

	return nil
}

func (w *WizardServer) handleProjectileDestroyed(message *protocol.Message) error {
	projectileDestroyedData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileDestroyedData.(*protocol.ProjectileDestroyedData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.BroadcastExcludeSender(message, data.CasterName)

	return nil
}

func (w *WizardServer) handleProjectileHitWizard(message *protocol.Message) error {
	projectileHitWizardData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileHitWizardData.(*protocol.ProjectileHitWizardData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.Mu.Lock()
	if _, ok := w.wizards[data.VictimStatData.Username]; ok {
		w.wizards[data.VictimStatData.Username].Health = data.VictimStatData.Health
		if w.wizards[data.VictimStatData.Username].Health <= 0 {
			w.deadWizards[data.VictimStatData.Username] = struct{}{}
		}
	}
	w.Mu.Unlock()

	w.BroadcastExcludeSender(message, data.ShooterStatData.Username)

	return nil
}

func (w *WizardServer) handleProjectileHitEnemy(message *protocol.Message) error {
	projectileHitEnemyData, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := projectileHitEnemyData.(*protocol.ProjectileHitEnemyData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	w.Mu.Lock()
	if _, ok := w.enemies[data.VictimStatData.ID]; ok {
		w.enemies[data.VictimStatData.ID].KnockbackDx = data.VictimDx
		w.enemies[data.VictimStatData.ID].KnockbackDy = data.VictimDy
		w.enemies[data.VictimStatData.ID].Health = data.VictimStatData.Health
		if w.enemies[data.VictimStatData.ID].Health <= 0 {
			delete(w.enemies, data.VictimStatData.ID)
		}
	}
	w.Mu.Unlock()

	if ok {
		w.BroadcastExcludeSender(message, data.ShooterStatData.Username)
	}

	return nil
}

func (w *WizardServer) handleClientLoaded(message *protocol.Message) error {
	w.clientLoadedCount.Add(1)
	return nil
}
