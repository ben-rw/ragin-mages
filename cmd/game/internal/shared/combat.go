package shared

import (
	"math/rand"
)

const (
	DefaultPlayerHealth      = 3.0
	DefaultPlayerAttackPower = 1.0
	DefaultPlayerMoveSpeed   = 2.0
	DefaultProjectileSpeed   = 3.0
	DefaultProjectileSize    = 1.0
	DefaultPlayerKnockback   = 3.0
	DefaultAttackRange       = 60 //determines how many ticks projectile will persist
	DefaultAttackCooldown    = 60
	DefaultIFrames           = 90
	DefaultFlickerFrames     = 5
	EnemyMoveSpeed           = 1.5
	EnemyHealth              = 1.0
	EnemyAttackPower         = 1.0
	EnemyAttackCooldown      = 60
	EnemyKnockBack           = 1.2
	KillEnemyBoost           = 1.0
	KillPlayerBoost          = 1.0
	ChestBoost               = 3.0
	WinRoundBoost            = 10.0
)

type Combat interface {
	Health() int
	AttackPower() int
	Damage(amount int)
	Attacking() bool
	Attack() bool
	Update()
}

type BasicCombat struct {
	health          int
	attackPower     int
	moveSpeed       float64
	projectileSpeed float64
	projectileSize  float64
	knockback       float64
	attacking       bool
}

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
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

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Damage(amount int) {
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

func (b *BasicCombat) BoostProjectileSpeed(amount float64) {
	b.projectileSpeed += amount
}

func (b *BasicCombat) ProjectileSize() float64 {
	return b.projectileSize
}

func (b *BasicCombat) BoostProjectileSize(amount float64) {
	b.projectileSize += amount
}

func (b *BasicCombat) Knockback() float64 {
	return b.knockback
}

func (b *BasicCombat) BoostKnockback(amount float64) {
	b.knockback += amount
}

func NewBasicCombat(health, attackPower int, moveSpeed, projectileSpeed, projectileSize, knockback float64) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		moveSpeed,
		projectileSpeed,
		projectileSize,
		knockback,
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
	e.timeSinceAttack += 1
}

func NewEnemyCombat(health, attackPower, attackCooldown int, moveSpeed, projectileSpeed, projectileSize, knockback float64) *EnemyCombat {
	return &EnemyCombat{
		NewBasicCombat(
			health,
			attackPower,
			moveSpeed,
			projectileSpeed,
			projectileSize,
			knockback,
		),
		attackCooldown,
		0,
	}
}

type WizardCombat struct {
	*BasicCombat
	attackRange     float64
	attackCooldown  int
	timeSinceAttack int
	Dead            bool
	iFrames         int
}

func (wc *WizardCombat) AttackRange() float64 {
	return wc.attackRange
}

func (wc *WizardCombat) BoostAttackRange(amount float64) {
	wc.attackRange += amount
}

func (wc *WizardCombat) RandomBoost(amount float64) {
	boosts := map[int]func(amount float64){
		0: wc.BoostProjectileSpeed,
		1: wc.BoostProjectileSize,
		2: wc.BoostKnockback,
		3: wc.BoostAttackRange,
	}

	boosts[rand.Intn(len(boosts))](amount)
}

func (wc *WizardCombat) Attack() bool {
	if wc.timeSinceAttack >= wc.attackCooldown {
		wc.attacking = true
		wc.timeSinceAttack = 0
		return true
	}
	return false
}

func (wc *WizardCombat) Update() {
	wc.timeSinceAttack += 1
	if wc.iFrames > 0 {
		wc.iFrames -= 1
	}
}

func (wc *WizardCombat) Damage(amount int) {
	wc.health -= amount
	wc.iFrames += DefaultIFrames
}

func (wc *WizardCombat) IFrames() int {
	return wc.iFrames
}

type WizardPlayer struct {
	*Player
	Combat        *WizardCombat
	Alpha         float32
	FlickerFrames int
}

func (w *WizardPlayer) IFrameFlicker() {
	w.FlickerFrames -= 1
	if w.FlickerFrames <= 0 {
		if w.Combat.IFrames() < DefaultFlickerFrames {
			w.Alpha = 1.0
			w.FlickerFrames = DefaultFlickerFrames
		} else if w.Alpha == 0.1 {
			w.Alpha = 0.5
			w.FlickerFrames = DefaultFlickerFrames
		} else {
			w.Alpha = 0.1
			w.FlickerFrames = DefaultFlickerFrames
		}
	}
}

func NewWizard(player *Player) *WizardPlayer {
	return &WizardPlayer{
		player,
		&WizardCombat{
			NewBasicCombat(
				DefaultPlayerHealth,
				DefaultPlayerAttackPower,
				DefaultPlayerMoveSpeed,
				DefaultProjectileSpeed,
				DefaultProjectileSize,
				DefaultPlayerKnockback,
			),
			DefaultAttackRange,
			DefaultAttackCooldown,
			0,
			false,
			0,
		},
		1.0,
		DefaultFlickerFrames,
	}
}
