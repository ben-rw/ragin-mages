package shared

import (
	"math/rand"
)

const (
	DefaultPlayerHealth      = 3.0
	DefaultPlayerAttackPower = 1.0
	DefaultPlayerMoveSpeed   = 2.0
	DefaultProjectileSpeed   = 3.0
	DefaultProjectileScale   = 0.5
	DefaultPlayerKnockback   = 3.0
	DefaultAttackCooldown    = 60 //determines how many ticks projectile will persist
	DefaultIFrames           = 90
	DefaultFlickerFrames     = 5
	DefaultHurtboxRadius     = 8
	EnemyMoveSpeed           = 1.5
	EnemyHealth              = 2.0
	EnemyAttackPower         = 1.0
	EnemyAttackCooldown      = 60
	EnemyKnockBack           = 1.2
	KillEnemyBoost           = 1.0
	KillPlayerBoost          = 1.0
	ChestBoost               = 3.0
	WinRoundBoost            = 10.0
	TrapDamage               = 1.0
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
	health          float64
	attackPower     float64
	moveSpeed       float64
	projectileSpeed float64
	projectileScale float64
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

func (b *BasicCombat) BoostProjectileSpeed(amount float64) {
	b.projectileSpeed += amount
}

func (b *BasicCombat) ProjectileScale() float64 {
	return b.projectileScale
}

func (b *BasicCombat) BoostProjectileScale(amount float64) {
	b.projectileScale += amount
}

func (b *BasicCombat) Knockback() float64 {
	return b.knockback
}

func (b *BasicCombat) BoostKnockback(amount float64) {
	b.knockback += amount
}

func NewBasicCombat(health, attackPower, moveSpeed, projectileSpeed, projectileScale, knockback float64) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		moveSpeed,
		projectileSpeed,
		projectileScale,
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

func NewEnemyCombat(attackCooldown int, health, attackPower, moveSpeed, projectileSpeed, projectileScale, knockback float64) *EnemyCombat {
	return &EnemyCombat{
		NewBasicCombat(
			health,
			attackPower,
			moveSpeed,
			projectileSpeed,
			projectileScale,
			knockback,
		),
		attackCooldown,
		0,
	}
}

type WizardCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
	Dead            bool
	Fell            bool
	iFrames         int
}

func (wc *WizardCombat) AttackCooldown() int {
	return wc.attackCooldown
}

func (wc *WizardCombat) BoostAttackCooldown(amount int) {
	wc.attackCooldown += amount
}

func (wc *WizardCombat) RandomBoost(amount float64) {
	boosts := map[int]func(amount float64){
		0: wc.BoostProjectileSpeed,
		1: wc.BoostProjectileScale,
		2: wc.BoostKnockback,
		3: func(amount float64) {
			wc.BoostAttackCooldown(int(amount))
		},
	}

	boosts[rand.Intn(len(boosts))](amount)
}

func (wc *WizardCombat) Attack() bool {
	if wc.timeSinceAttack >= int(wc.attackCooldown) {
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

func (wc *WizardCombat) Damage(amount float64) {
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
	HurtboxRadius float64
}

func (w *WizardPlayer) IFrameFlicker() {
	w.FlickerFrames -= 1
	if w.FlickerFrames <= 0 {
		if w.Combat.IFrames() < DefaultFlickerFrames*2 {
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
	spawnIndex := rand.Intn(8)
	player.X = PlayerSpawns[spawnIndex].X * TileSize
	player.Y = PlayerSpawns[spawnIndex].Y * TileSize

	return &WizardPlayer{
		player,
		&WizardCombat{
			NewBasicCombat(
				DefaultPlayerHealth,
				DefaultPlayerAttackPower,
				DefaultPlayerMoveSpeed,
				DefaultProjectileSpeed,
				DefaultProjectileScale,
				DefaultPlayerKnockback,
			),
			DefaultAttackCooldown,
			60, // set timeSinceAttack so player can attack immediately
			false,
			false,
			0,
		},
		1.0,
		DefaultFlickerFrames,
		DefaultHurtboxRadius,
	}
}
