# Ragin' Mages

Ragin' Mages is a multiplayer webgame inspired by the flash games of the early 2000s. Up to 8 players can play together in one room. Ragin' Mages is much more fun with friends, but you can play solo too!

# How to Play

1. Go to [raginmages.xyz](https://raginmages.xyz/) in your browser.
2. Create a room, or join a room created by another player with their room code.
3. When everyone is ready, the host can start the game by pressing Enter. 
4. Have fun! :)

## Controls

| Input | Action |
| --- | --- |
| WASD / arrow keys | Move |
| Left click | Shoot fireball |
| Right click | Reflect fireball (short cooldown on failure, no cooldown on success) |
| F3 | Toggle hitbox debug overlay |
| F4 | Toggle FPS/TPS |

##  Mechanics

- Enemies chase players and use a melee attack at close range, which deals damage and applies a small amount of knockback
- Player projectiles deal damage and apply knockback scaling with their knockback stat
- Getting knocked into a hole results in INSTANT DEATH, regardless of remaining health
- Players can reflect projectiles by right clicking right as the projectile is about to hit them; the projectile is reflected in the direction of the reflector's cursor with increased stats
- Stat increases are gained by killing skeletons, opening chests, killing other players, and winning rounds.
- Killing another player always gives a flat stat point reward. However, if the victim had higher stats than the shooter, the shooter also STEALS a percentage of the victim's stats.
- Win rounds by being the last player standing. If two players are alive at the end of the round, or if the last two players die on the same tick, no stats points are awarded.
- Win the game by having the highest stat total at the end of the third round.

# Roadmap

- add more robust disconnection logic
- move movement and projectile logic to server, which will necessitate client-side prediction and interpolation
- more games! The original idea was to have many mini/micro games à la Mario Party or WarioWare
- optional bot players
- WizArena:
  - investigate potential performance improvements
  - briefly stun enemies hit by projectiles
  - add a minimap
  - add an overhead bar to indicate when the player can attack again
  - give players a visual indication of their stats increasing
  - add walls to the map
  - add new mechanics i.e. invisibility, bouncing projectiles, spell ammo

# How It Works

# Attributions

# Contributing
