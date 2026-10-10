package wizarena

import (
	"github.com/ben-rw/ragin-mages/internal/protocol"
)

func (w *WizArena) WriteWizardMovementUpdate() {
	w.Conn.WriteMsg(protocol.WizardMovementUpdate, &protocol.WizardMovementUpdateData{
		X:        w.wizard.X,
		Y:        w.wizard.Y,
		Dx:       w.wizard.Y,
		Dy:       w.wizard.Y,
		Username: w.wizard.Data.Username,
	})
}

// send when w.wizard kills an enemy, hits a player, opens a chest, or wins a round
func (w *WizArena) WriteWizardStatUpdate() {
	w.Conn.WriteMsg(protocol.WizardStatUpdate, &protocol.WizardStatUpdateData{
		Score:           w.wizard.Data.Score,
		Health:          w.wizard.Combat.Health(),
		AttackCooldown:  w.wizard.Combat.AttackCooldown(),
		ProjectileSpeed: w.wizard.Combat.ProjectileSpeed(),
		ProjectileScale: w.wizard.Combat.ProjectileScale(),
		Knockback:       w.wizard.Combat.Knockback(),
		Username:        w.wizard.Data.Username,
	})
}

func (w *WizArena) WriteNewProjectile(cx, cy int) {
	w.Conn.WriteMsg(protocol.NewProjectile, &protocol.NewProjectileData{
		CasterName: w.wizard.Data.Username,
		CursorX:    cx,
		CursorY:    cy,
	})
}

func (w *WizArena) WriteReflecting(cx, cy int) {
	w.Conn.WriteMsg(protocol.Reflecting, &protocol.ReflectingData{
		ReflectorName: w.wizard.Data.Username,
		CursorX:       cx,
		CursorY:       cy,
	})
}

func (w *WizArena) WriteProjectileReflected(reflectorName string, projectileID uint16) {
	w.Conn.WriteMsg(protocol.ProjectileReflected, &protocol.ProjectileReflectedData{
		ReflectorName: reflectorName,
		ProjectileID:  projectileID,
	})
}

func (w *WizArena) WriteProjectileHitWizard(victim *WizardPlayer, victimDx, victimDy float64) {
	w.Conn.WriteMsg(protocol.ProjectileHitWizard, &protocol.ProjectileHitWizardData{
		ShooterStatData: &protocol.WizardStatUpdateData{
			Score:           w.wizard.Data.Score,
			Health:          w.wizard.Combat.Health(),
			AttackCooldown:  w.wizard.Combat.AttackCooldown(),
			ProjectileSpeed: w.wizard.Combat.ProjectileSpeed(),
			ProjectileScale: w.wizard.Combat.ProjectileScale(),
			Knockback:       w.wizard.Combat.Knockback(),
			Username:        w.wizard.Data.Username,
		},
		VictimStatData: &protocol.WizardStatUpdateData{
			Score:           victim.Data.Score,
			Health:          victim.Combat.Health(),
			AttackCooldown:  victim.Combat.AttackCooldown(),
			ProjectileSpeed: victim.Combat.ProjectileSpeed(),
			ProjectileScale: victim.Combat.ProjectileScale(),
			Knockback:       victim.Combat.Knockback(),
			Username:        victim.Data.Username,
		},
		VictimDx:     victimDx,
		VictimDy:     victimDy,
		VictimNoclip: true,
	})
}

func (w *WizArena) WriteProjectileHitEnemy(victim *Enemy, victimDx, victimDy float64) {
	w.Conn.WriteMsg(protocol.ProjectileHitEnemy, &protocol.ProjectileHitEnemyData{
		ShooterStatData: &protocol.WizardStatUpdateData{
			Score:           w.wizard.Data.Score,
			Health:          w.wizard.Combat.Health(),
			AttackCooldown:  w.wizard.Combat.AttackCooldown(),
			ProjectileSpeed: w.wizard.Combat.ProjectileSpeed(),
			ProjectileScale: w.wizard.Combat.ProjectileScale(),
			Knockback:       w.wizard.Combat.Knockback(),
			Username:        w.wizard.Data.Username,
		},
		VictimStatData: &protocol.EnemyStatUpdateData{
			Health: victim.Combat.Health(),
			ID:     victim.ID,
		},
		VictimDx: victimDx,
		VictimDy: victimDy,
	})
}
