package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type MessageType string

const (
	Unset                 MessageType = "Unset"
	JoinRequest           MessageType = "JoinRequest"
	JoinResponse          MessageType = "JoinResponse"
	SceneChange           MessageType = "SceneChange"
	PlayerUpdate          MessageType = "PlayerUpdate"
	WizardMovementUpdate  MessageType = "WizardMovementUpdate"
	WizardStatUpdate      MessageType = "WizardStatUpdate"
	EnemyMovementUpdate   MessageType = "EnemyMovementUpdate"
	EnemyStatUpdate       MessageType = "EnemyStatUpdate"
	WizArenaMovementState MessageType = "WizArenaState"
	NewProjectile         MessageType = "NewProjectile"
	ProjectileHitWizard   MessageType = "ProjectileHitWizard"
	ProjectileHitEnemy    MessageType = "ProjectileHitEnemy"
	ProjectileReflected   MessageType = "ProjectileReflected"
)

type Message struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinRequestData struct {
	RoomID string `json:"room_id"`
}

type PlayerData struct {
	Score       int     `json:"score"`
	SpriteIndex int     `json:"sprite_index"`
	Host        bool    `json:"host"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Name        string  `json:"username"`
}

type JoinResponseData struct {
	PlayerData *PlayerData   `json:"player_data"`
	PlayerList []*PlayerData `json:"player_list"`
}

type PlayerUpdateData struct {
	PlayerData *PlayerData `json:"player_data"`
}

type WizardMovementUpdateData struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Dx   float64 `json:"dx"`
	Dy   float64 `json:"dy"`
	Name string  `json:"username"`
}

type WizardStatUpdateData struct {
	Score           int     `json:"score"`
	AttackCooldown  int     `json:"attack_cooldown"`
	ProjectileSpeed float64 `json:"projectile_speed"`
	ProjectileScale float64 `json:"projectile_scale"`
	Knockback       float64 `json:"knockback"`
	Name            string  `json:"username"`
	Health          float64 `json:"health"`
}

type EnemyMovementUpdateData struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Dx float64 `json:"dx"`
	Dy float64 `json:"dy"`
	ID int     `json:"id"`
}

type EnemyStatUpdateData struct {
	Health float64 `json:"health"`
	ID     int     `json:"id"`
}

// type WizArenaMovementStateData struct {
// 	WizardMovementData  *WizardMovementUpdateData  `json:"wizard_movement_data"`
// 	EnemiesMovementData []*EnemyMovementUpdateData `json:"enemies_movement_data"`
// }

type NewProjectileData struct {
	CasterName string `json:"caster_name"`
	CursorX    int    `json:"cursor_x"`
	CursorY    int    `json:"cursor_y"`
}

type ProjectileReflectedData struct {
	ReflectorName string `json:"reflector_name"`
	CursorX       int    `json:"cursor_x"`
	CursorY       int    `json:"cursor_y"`
	ProjectileID  uint16 `json:"projectile_id"`
}

type ProjectileHitWizardData struct {
	ShooterStatData *WizardStatUpdateData `json:"shooter_stat_data"`
	VictimStatData  *WizardStatUpdateData `json:"victim_stat_data"`
	VictimDx        float64               `json:"victim_dx"`
	VictimDy        float64               `json:"victim_dy"`
	VictimNoclip    bool                  `json:"victim_noclip"`
}

type ProjectileHitEnemyData struct {
	ShooterStatData *WizardStatUpdateData `json:"shooter_stat_data"`
	VictimStatData  *EnemyStatUpdateData  `json:"victim_stat_data"`
	VictimDx        float64               `json:"victim_dx"`
	VictimDy        float64               `json:"victim_dy"`
}

type SceneType int

const (
	LobbyScene SceneType = iota
	RandomScene
	MemoryScene
	WizardsScene
)

type SceneChangeData struct {
	SceneType SceneType `json:"scene_type"`
}

func (m *Message) UnmarshalMessageData() (any, error) {
	var v any

	switch m.Type {
	case JoinRequest:
		v = &JoinRequestData{}
	case JoinResponse:
		v = &JoinResponseData{}
	case SceneChange:
		v = &SceneChangeData{}
	case PlayerUpdate:
		v = &PlayerUpdateData{}
	default:
		return nil, errors.New("unrecognized message format")
	}

	log.Printf("json: %s type: %v\n", m.Data, m.Type)

	err := json.Unmarshal(m.Data, v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return v, nil
}

func MarshalToMessage(mt MessageType, msgData any) (*Message, error) {
	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Type: mt,
		Data: jsonData,
	}

	return msg, nil
}
