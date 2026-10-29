package wizarena

import (
	"math/rand"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
)

const (
	PlayerHealth      = wizServer.PlayerHealth
	PlayerAttackPower = 1.0
	PlayerMoveSpeed   = 2.0
	ReflectCooldown   = 180
	IFrames           = 90
	FlickerFrames     = 5
	HurtboxRadius     = 8
	ReflectFrames     = wizServer.ReflectFrames
)

// boostable stats
const (
	BaseProjectileSpeed = 3.0
	BaseProjectileScale = 0.5
	BaseKnockback       = 3.0
	BaseAttackCooldown  = 60.0 //determines how many ticks projectile will persist

	PSpeedMult    = .35
	PScaleMult    = .15
	KnockbackMult = .15
	AttackCDMult  = .20
	StandardMult  = 0.0

	PSpeedCap    = 35
	PScaleCap    = 15
	KnockbackCap = 20
	AttackCDCap  = 30
)

const (
	KillEnemyBoost   = 1.0
	KillPlayerBoost  = 3.0
	ChestBoost       = 3.0
	WinRoundBoost    = wizServer.WinRoundBoost
	StatTheftDivisor = 2
)

const (
	EnemyMoveSpeed      = wizServer.EnemyMoveSpeed
	EnemyHealth         = 2.0
	EnemyAttackPower    = 1.0
	EnemyAttackCooldown = 60
	EnemyKnockBack      = 1.2
)

const TrapDamage = 1.0

type Combat interface {
	Health() int
	AttackPower() int
	Damage(amount int)
	Attacking() bool
	Attack() bool
	Update()
}

type BasicCombat struct {
	health          float64
	attackPower     float64
	moveSpeed       float64
	projectileSpeed float64
	projectileScale float64
	knockback       float64
	attackCooldown  float64
	timeSinceAttack float64
	attacking       bool
}

type EnemyCombat struct {
	*BasicCombat
	Dead bool
}

func (b *BasicCombat) Update() {
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) AttackPower() float64 {
	return b.attackPower
}

func (b *BasicCombat) Health() float64 {
	return b.health
}

func (b *BasicCombat) SetHealth(amount float64) {
	b.health = amount
}

func (b *BasicCombat) Damage(amount float64) {
	b.health -= amount
}

func (b *BasicCombat) MoveSpeed() float64 {
	return b.moveSpeed
}

func (b *BasicCombat) SetMoveSpeed(newMoveSpeed float64) {
	b.moveSpeed = newMoveSpeed
}

func (b *BasicCombat) ProjectileSpeed() float64 {
	return b.projectileSpeed
}

func (b *BasicCombat) BoostProjectileSpeed(amount float64, mult float64) {
	if mult == StandardMult {
		mult = PSpeedMult
	}
	b.projectileSpeed += amount * mult
}

func (b *BasicCombat) ProjectileScale() float64 {
	return b.projectileScale
}

func (b *BasicCombat) BoostProjectileScale(amount float64, mult float64) {
	if mult == StandardMult {
		mult = PScaleMult
	}
	b.projectileScale += amount * mult
}

func (b *BasicCombat) Knockback() float64 {
	return b.knockback
}

func (b *BasicCombat) BoostKnockback(amount float64, mult float64) {
	if mult == StandardMult {
		mult = KnockbackMult
	}
	b.knockback += amount * mult
}

func NewBasicCombat(health, attackPower, moveSpeed, projectileSpeed, projectileScale, knockback, attackCooldown, timeSinceAttack float64) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		moveSpeed,
		projectileSpeed,
		projectileScale,
		knockback,
		attackCooldown,
		0,
		false,
	}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.attacking = false
	e.timeSinceAttack += 1
}

func NewEnemyCombat(attackCooldown, health, attackPower, moveSpeed, projectileSpeed, projectileScale, knockback float64) *EnemyCombat {
	return &EnemyCombat{
		NewBasicCombat(
			health,
			attackPower,
			moveSpeed,
			projectileSpeed,
			projectileScale,
			knockback,
			attackCooldown,
			0,
		),
		false,
	}
}

type WizardCombat struct {
	*BasicCombat
	LastHitBy        string
	reflectCooldown  int
	timeSinceAttack  int
	timeSinceReflect int
	iFrames          int
	reflectFrames    int
	Dead             bool
	Fell             bool
	reflecting       bool
}

func (wc *WizardCombat) SetKnockback(amount float64) {
	wc.knockback = amount
}

func (wc *WizardCombat) SetProjectileScale(amount float64) {
	wc.projectileScale = amount
}

func (wc *WizardCombat) SetProjectileSpeed(amount float64) {
	wc.projectileSpeed = amount
}

func (wc *WizardCombat) SetAttackCooldown(amount float64) {
	wc.attackCooldown = amount
}

func (wc *WizardCombat) GetTotalPower() int {
	return int((wc.projectileSpeed*(1/PSpeedMult) - (1 / PSpeedMult * BaseProjectileSpeed)) +
		(wc.projectileScale*(1/PScaleMult) - (1 / PScaleMult * BaseProjectileScale)) +
		(wc.knockback*(1/KnockbackMult) - (1 / KnockbackMult * BaseKnockback)) +
		((60 - wc.attackCooldown) * 1 / AttackCDMult))
}

func (wc *WizardCombat) AttackCooldown() float64 {
	return wc.attackCooldown
}

func (wc *WizardCombat) BoostAttackCooldown(amount float64, mult float64) {
	if mult == StandardMult {
		mult = AttackCDMult
	}
	wc.attackCooldown -= amount * mult
}

// pass StandardMult unless stats have already had their corresponding mult applied
func (wc *WizardCombat) RandomBoost(amount float64, mult float64) {
	boosts := make([]int, 0, 4)

	if wc.ProjectileSpeed() < PSpeedCap {
		boosts = append(boosts, 0)
	}
	if wc.ProjectileScale() < PScaleCap {
		boosts = append(boosts, 1)
	}
	if wc.Knockback() < KnockbackCap {
		boosts = append(boosts, 2)
	}
	if wc.AttackCooldown() > AttackCDCap {
		boosts = append(boosts, 3)
	}

	if len(boosts) > 0 {
		i := rand.Intn(len(boosts))
		switch boosts[i] {
		case 0:
			wc.BoostProjectileSpeed(amount, mult)
		case 1:
			wc.BoostProjectileScale(amount, mult)
		case 2:
			wc.BoostKnockback(amount, mult)
		case 3:
			wc.BoostAttackCooldown(amount, mult)
		default:
		}
	}
}

func (wc *WizardCombat) Attack() bool {
	if wc.timeSinceAttack >= int(wc.attackCooldown) {
		wc.attacking = true
		wc.timeSinceAttack = 0
		return true
	}
	return false
}

func (wc *WizardCombat) Reflect() bool {
	if wc.timeSinceReflect >= int(wc.reflectCooldown) {
		wc.reflecting = true
		wc.reflectFrames = ReflectFrames
		wc.timeSinceReflect = 0
		return true
	}
	return false
}

func (wc *WizardCombat) Reflecting() bool {
	return wc.reflecting
}

// reset cooldown on successful reflect
func (wc *WizardCombat) ResetReflectCooldown() {
	wc.timeSinceReflect = ReflectCooldown
}

func (wc *WizardCombat) Update() {
	wc.attacking = false
	wc.timeSinceAttack += 1
	wc.timeSinceReflect += 1
	if wc.reflectFrames > 0 {
		wc.reflectFrames -= 1
		if wc.reflectFrames <= 0 {
			wc.reflecting = false
		}
	}
	if wc.iFrames > 0 {
		wc.iFrames -= 1
	}
}

func (wc *WizardCombat) Damage(amount float64) {
	wc.health -= amount
	wc.iFrames += IFrames
}

func (wc *WizardCombat) IFrames() int {
	return wc.iFrames
}

func (wc *WizardCombat) SetIFrames(amount int) {
	wc.iFrames = amount
}

type WizardPlayer struct {
	*shared.Player
	Combat        *WizardCombat
	Reflect       *Reflect
	FlickerFrames int
	HurtboxRadius float64
	DeathX        float64
	DeathY        float64
}

func (w *WizardPlayer) IFrameFlicker() {
	w.FlickerFrames -= 1
	if w.FlickerFrames <= 0 {
		if w.Combat.IFrames() < FlickerFrames*2 {
			w.Alpha = 1.0
			w.FlickerFrames = FlickerFrames
		} else if w.Alpha == 0.1 {
			w.Alpha = 0.5
			w.FlickerFrames = FlickerFrames
		} else {
			w.Alpha = 0.1
			w.FlickerFrames = FlickerFrames
		}
	}
}

func NewWizard(player *shared.Player) *WizardPlayer {
	return &WizardPlayer{
		player,
		&WizardCombat{
			BasicCombat: NewBasicCombat(
				PlayerHealth,
				PlayerAttackPower,
				PlayerMoveSpeed,
				BaseProjectileSpeed,
				BaseProjectileScale,
				BaseKnockback,
				BaseAttackCooldown,
				0,
			),
			LastHitBy:        "",
			reflectCooldown:  ReflectCooldown,
			timeSinceAttack:  BaseAttackCooldown, // set timeSinceAttack so player can attack immediately
			timeSinceReflect: ReflectCooldown,    // set timeSinceReflect so player can reflect immediately
			iFrames:          0,
			reflectFrames:    0,
			Dead:             false,
			Fell:             false,
			reflecting:       false,
		},
		NewReflect(player.X, player.Y),
		FlickerFrames,
		HurtboxRadius,
		0,
		0,
	}
}

func StealStats(thief *WizardCombat, victim *WizardCombat) (float64, float64, float64, float64) {
	var stolenPScale float64 = 0
	var stolenPSpeed float64 = 0
	var stolenKnockback float64 = 0
	var stolenAttackCooldown float64 = 0
	if thief.ProjectileScale() < victim.ProjectileScale() {
		stolenPScale = (victim.ProjectileScale() - thief.ProjectileScale()) / StatTheftDivisor
		thief.BoostProjectileScale(stolenPScale, 1) // 1 for all mults as existing player stats have had already had mult applied
		victim.BoostProjectileScale(-stolenPScale, 1)
	}
	if thief.ProjectileSpeed() < victim.ProjectileSpeed() {
		stolenPSpeed = (victim.ProjectileSpeed() - thief.ProjectileSpeed()) / StatTheftDivisor
		thief.BoostProjectileSpeed(stolenPSpeed, 1)
		victim.BoostProjectileSpeed(-stolenPSpeed, 1)
	}
	if thief.Knockback() < victim.Knockback() {
		stolenKnockback = (victim.Knockback() - thief.Knockback()) / StatTheftDivisor
		thief.BoostKnockback(stolenKnockback, 1)
		victim.BoostKnockback(-stolenKnockback, 1)
	}
	if thief.AttackCooldown() > victim.AttackCooldown() {
		stolenAttackCooldown = (thief.AttackCooldown() - victim.AttackCooldown()) / StatTheftDivisor
		thief.BoostAttackCooldown(float64(stolenAttackCooldown), 1)
		victim.BoostAttackCooldown(float64(-stolenAttackCooldown), 1)
	}

	return stolenPScale, stolenPSpeed, stolenKnockback, stolenAttackCooldown
}
