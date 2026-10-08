package shared

import (
	"math/rand"
)

const (
	PlayerHealth      = 3.0
	PlayerAttackPower = 1.0
	PlayerMoveSpeed   = 2.0
	ReflectCooldown   = 180
	IFrames           = 90
	FlickerFrames     = 5
	HurtboxRadius     = 8
	ReflectFrames     = 15
)

// boostable stats
const (
	BaseProjectileSpeed = 3.0
	BaseProjectileScale = 0.5
	BaseKnockback       = 3.0
	BaseAttackCooldown  = 60.0 //determines how many ticks projectile will persist
	PSpeedMult          = .5
	PScaleMult          = .25
	KnockbackMult       = .5
	AttackCDMult        = .5
	StandardMult        = 0.0
)

const (
	KillEnemyBoost   = 1.0
	KillPlayerBoost  = 3.0
	ChestBoost       = 3.0
	WinRoundBoost    = 10.0
	StatTheftDivisor = 4
)

const (
	EnemyMoveSpeed      = 1.5
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
	reflectCooldown  int
	timeSinceAttack  int
	timeSinceReflect int
	iFrames          int
	reflectFrames    int
	totalPower       int
	Dead             bool
	Fell             bool
	reflecting       bool
}

func (wc *WizardCombat) GetTotalPower() int {
	return wc.totalPower
}

func (wc *WizardCombat) AttackCooldown() float64 {
	return wc.attackCooldown
}

func (wc *WizardCombat) BoostAttackCooldown(amount float64, mult float64) {
	if mult == StandardMult {
		mult = AttackCDMult
	}
	wc.attackCooldown -= amount * mult
	if wc.attackCooldown < 12 { // limit projectile spawning to ~5/sec
		wc.attackCooldown = 12 // 12/60 gives a nice even .2 as the attackCD cap
	}
}

// pass StandardMult unless stats have already had their corresponding mult applied
func (wc *WizardCombat) RandomBoost(amount float64, mult float64) {
	wc.totalPower += int(amount)
	boosts := map[int]func(amount float64, mult float64){
		0: wc.BoostProjectileSpeed,
		1: wc.BoostProjectileScale,
		2: wc.BoostKnockback,
		3: func(amount float64, mult float64) {
			wc.BoostAttackCooldown(amount, mult)
		},
	}

	boosts[rand.Intn(len(boosts))](amount, mult)
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
	*Player
	Combat        *WizardCombat
	Reflect       *Reflect
	FlickerFrames int
	HurtboxRadius float64
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

func NewWizard(player *Player) *WizardPlayer {

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
	}
}

func StealStats(thief *WizardCombat, victim *WizardCombat) {
	if thief.ProjectileScale() < victim.ProjectileScale() {
		stolen := (victim.ProjectileScale() - thief.ProjectileScale()) / StatTheftDivisor
		thief.BoostProjectileScale(stolen, 1) // 1 for all mults as existing player stats have had already had mult applied
		victim.BoostProjectileScale(-stolen, 1)
	}
	if thief.ProjectileSpeed() < victim.ProjectileSpeed() {
		stolen := (victim.ProjectileSpeed() - thief.ProjectileSpeed()) / StatTheftDivisor
		thief.BoostProjectileSpeed(stolen, 1)
		victim.BoostProjectileSpeed(-stolen, 1)
	}
	if thief.Knockback() < victim.Knockback() {
		stolen := (victim.Knockback() - thief.Knockback()) / StatTheftDivisor
		thief.BoostKnockback(stolen, 1)
		victim.BoostKnockback(-stolen, 1)
	}
	if thief.AttackCooldown() > victim.AttackCooldown() {
		stolen := (thief.AttackCooldown() - victim.AttackCooldown()) / StatTheftDivisor
		thief.BoostAttackCooldown(float64(stolen), 1)
		victim.BoostAttackCooldown(float64(-stolen), 1)
	}
}
