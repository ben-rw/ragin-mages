package wizServer

import (
	"errors"
	"log"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/ben-rw/ragin-mages/internal/room"
)

func (w *WizardServer) handleJoinRequest(message *protocol.Message) error {

	return nil
}

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
	} else if !ok {
		log.Printf("not found in wizards list: %v", data.Username)
	}

	return nil
}

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
	} else if !ok {
		log.Printf("not found in wizards list: %v", data.Username)
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
	} else if !ok {
		log.Printf("not found in enemies list: %v", data.ID)
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

	w.BroadcastExcludeSender(message, data.ShooterName)

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

	// w.wizards[data.ReflectorName].startReflecting()

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
	} else if !ok {
		log.Printf("not found in wizards list: %v", data.VictimStatData.Username)
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

	found := false
	w.Mu.Lock()
	if _, ok := w.enemies[data.VictimStatData.ID]; ok {
		w.enemies[data.VictimStatData.ID].KnockbackDx = data.VictimDx
		w.enemies[data.VictimStatData.ID].KnockbackDy = data.VictimDy
		w.enemies[data.VictimStatData.ID].Health = data.VictimStatData.Health
		if w.enemies[data.VictimStatData.ID].Health <= 0 {
			delete(w.enemies, data.VictimStatData.ID)
		}
		found = true
	} else if !ok {
		log.Printf("not found in enemies list: %v", data.VictimStatData.ID)
	}
	w.Mu.Unlock()

	if found {
		w.BroadcastExcludeSender(message, data.ShooterStatData.Username)
	}

	return nil
}

func (w *WizardServer) handleClientLoaded(message *protocol.Message, c *room.Client) error {
	clientLoaded, err := message.UnmarshalMessageData()
	if err != nil {
		return err
	}
	data, ok := clientLoaded.(*protocol.ClientLoadedData)
	if !ok {
		err = errors.New("failed type assertion")
		return err
	}

	switch data.LoadStage {
	case protocol.InitDone:
		c.Player.ClientInitDone = true
	case protocol.StateUpdated:
		c.Player.ClientStateUpdated = true
	default:
		log.Printf("unrecognized LoadState: %v\n", data.LoadStage)
	}

	return nil
}

func (w *WizardServer) handleGivePoints(message *protocol.Message) error {
	// givePointsData, err := message.UnmarshalMessageData()
	// if err != nil {
	// 	return err
	// }
	// data, ok := givePointsData.(*protocol.GivePointsData)
	// if !ok {
	// 	err = errors.New("failed type assertion")
	// 	return err
	// }

	w.Broadcast(message)

	return nil
}
