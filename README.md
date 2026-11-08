# Ragin' Mages

Ragin' Mages is a multiplayer webgame inspired by the flash games of the early 2000s. Up to 8 players can play together in one room. Ragin' Mages is much more fun with friends, but you can play solo too!

## How to Play

1. Go to [raginmages.xyz](https://raginmages.xyz/) in your browser.
2. Create a room, or join a room created by another player with their room code.
3. When everyone is ready, the host can start the game by pressing Enter.
4. Have fun! :)

### Controls

| Input | Action |
| --- | --- |
| WASD / arrow keys | Move |
| Left click | Shoot fireball |
| Right click | Reflect fireball (short cooldown on failure, no cooldown on success) |
| F3 | Toggle hitbox debug overlay |
| F4 | Toggle FPS/TPS |

### Mechanics

- Enemies chase players and use a melee attack at close range, which deals damage and applies a small amount of knockback
- Player projectiles deal damage and apply knockback scaling with their knockback stat
- Getting knocked into a hole results in INSTANT DEATH, regardless of remaining health
- Players can reflect projectiles by right clicking right as the projectile is about to hit them; the projectile is reflected in the direction of the reflector's cursor with increased stats
- Stat increases are gained by killing skeletons, opening chests, killing other players, and winning rounds
- Killing another player always gives a flat stat point reward. However, if the victim had higher stats than the shooter, the shooter also STEALS a percentage of the victim's stats
- Win rounds by being the last player standing. If two players are alive at the end of the round, or if the last two players die on the same tick, no stats points are awarded
- Win the game by having the highest stat total at the end of the third round

## How It Works

Made with Ebitengine, uses coder/websocket

## Dependencies

- [Ebitengine](https://github.com/hajimehoshi/ebiten/v2) for the game client
- [coder/websocket](https://github.com/coder/websocket) for multiplayer networking
- [google/uuid](https://github.com/google/uuid) for generating unique IDs
- [godotenv](https://github.com/joho/godotenv) for loading `.env`

## Roadmap

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

## Attributions

- Lobby Music - Amalfi Coast by [Albert Behar](https://albertbehar.com/)
- WizArena Music - Void Construct by [Alex Grohl](https://pixabay.com/users/alexgrohl-25289918/)
- Character, Map, and UI Assets - [Ninja Adventure Asset Pack](https://pixel-boy.itch.io/ninja-adventure-asset-pack) by 
[pixel-boy](https://itch.io/profile/pixel-boy)
- Projectile & Reflection Assets - [Warped Shooting Fx Asset Pack](https://ansimuz.itch.io/warped-shooting-fx) by [ansimuz](https://itch.io/profile/ansimuz)
- Game engine - [Ebitengine](https://ebitengine.org/)
- Ebitengine tutorials - [RPG in Go](https://www.youtube.com/playlist?list=PLvN4CrYN-8i7xnODFyCMty6ossz4eW0Cn) by [Coding with Sphere](https://www.youtube.com/@codingwithsphere)
- Ebitengine tutorials - [Trevor's Tutorials](https://trevors-tutorials.com/0000-introduction/)
- Landing page font - [PressStart2p](https://github.com/codeman38/PressStart2P/tree/master) by [codeman38](https://github.com/codeman38)
- In-game font - PICO-8 mono by [lexaloffle and RhythmLynx](https://www.lexaloffle.com/bbs/?tid=3760)
- Go performance tips(most of which have not yet been applied!) - [goperf.dev](goperf.dev)
- Programming knowledge and experience - [boot.dev](boot.dev)

## Contributing

Contributions are welcome! If you'd like to contribute:

- Open an issue to discuss your idea or change.
- Fork the repository and open a pull request to the `main` branch.
- Please make sure the game works before sending the PR. Thanks!

Create a `.env` file:

```text
PORT="8080"
FILEPATH_ROOT="."
URL_ROOT="http://localhost"
```

To compile both binaries and start the server:
> $ GOOS=js GOARCH=wasm go build -o frontend/static/ragin-mages.wasm ./cmd/game && go build -o ragin-mages-server ./cmd/server/ && ./ragin-mages-server

Alternatively, you can use the `watch` command I have in `.watchexec` to have the binaries recompiled every time you save a go, html, css, or javascript file. Requires [watchexec](https://github.com/watchexec/watchexec).
> $ source .watchexec
> $ $watch
