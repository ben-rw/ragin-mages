package wizarena

//
// import (
// 	"errors"
//
// 	"github.com/ben-rw/ragin-mages/internal/protocol"
// )
//
// func (w *WizArena) HandleWizardMovementUpdate(message *protocol.Message) error {
// 	wizardMovementUpdateData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := wizardMovementUpdateData.(*protocol.WizardMovementUpdateData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleEnemyMovementUpdate(message *protocol.Message) error {
// 	enemyMovementUpdateData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := enemyMovementUpdateData.(*protocol.EnemyMovementUpdateData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleWizArenaMovementUpdate(message *protocol.Message) error {
// 	wizArenaMovementData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := wizArenaMovementData.(*protocol.WizArenaMovementData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleWizardStatUpdate(message *protocol.Message) error {
// 	wizardStatUpdateData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := wizardStatUpdateData.(*protocol.WizardStatUpdateData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleNewProjectile(message *protocol.Message) error {
// 	newProjectileData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := newProjectileData.(*protocol.NewProjectileData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleProjectileReflected(message *protocol.Message) error {
// 	projectileReflectedData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := projectileReflectedData.(*protocol.ProjectileReflectedData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleProjectileHitWizard(message *protocol.Message) error {
// 	projectileHitWizardData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := projectileHitWizardData.(*protocol.ProjectileHitWizardData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
//
// func (w *WizArena) HandleProjectileHitEnemy(message *protocol.Message) error {
// 	projectileHitEnemyData, err := message.UnmarshalMessageData()
// 	if err != nil {
// 		return err
// 	}
// 	data, ok := projectileHitEnemyData.(*protocol.ProjectileHitEnemyData)
// 	if !ok {
// 		err = errors.New("failed type assertion")
// 		return err
// 	}
// 	return nil
// }
