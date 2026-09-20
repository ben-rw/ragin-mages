package wizarena

import (
	"image"
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/ws"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tilemapPath = "assets/maps/ninja_dungeon.json"
	heartPath   = "assets/images/ninja_adventure/Ui/Receptacle/IconHeart.png"
	songPath    = "assets/audio/music/void-construct-loop-0.6.ogg"
	introLen    = 1 // 1 second song intro
)

type WizArena struct {
	shared.Roster
	Conn              *ws.Connection
	Sprites           []*shared.Sprite
	wizard            *shared.WizardPlayer
	wizards           map[string]*shared.WizardPlayer
	enemies           []*shared.Enemy
	enemyRespawnTimer time.Time
	roundTimer        time.Time
	tilemapJSON       *shared.TilemapJSON
	tileCache         map[int]*shared.Tile
	camera            *shared.Camera
	colliders         []image.Rectangle
	chests            []image.Rectangle
	chestRespawnTimer time.Time
	openedChests      map[OpenedChest]struct{}
	holes             []image.Rectangle
	traps             []image.Rectangle
	trapsUp           bool
	audioPlayer       *audio.Player
	projectiles       []*shared.Projectile
	deadProjectiles   []*shared.Projectile
	projectileCache   map[shared.ProjectileType]*ebiten.Image
	heartImage        *ebiten.Image
	debug             bool
}

func NewWizArena(c *ws.Connection) *WizArena {
	log.Println("scene changed to Wizards")
	// shared.ScreenHeight = shared.ScreenHeight * 2
	// shared.ScreenWidth = shared.ScreenWidth * 2

	tilemap, err := shared.NewTilemapJSON(tilemapPath)
	if err != nil {
		log.Printf("couldn't load tilemap: %v", err)
	}

	tileCache, err := shared.NewTileCache(tilemap)
	if err != nil {
		log.Printf("couldn't build tile cache: %v", err)
	}

	audioPlayer, err := sound.NewAudioPlayer(songPath, true, introLen)
	if err != nil {
		log.Printf("couldn't create audio player: %v", err)
	}

	fireballImg, err := shared.LoadProjectile(shared.Fireball)
	if err != nil {
		log.Printf("couldn't load projectile image: %v", err)
	}

	heartImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, heartPath)
	if err != nil {
		log.Printf("couldn't load image: %v", err)
	}

	w := &WizArena{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8), // max 8 players
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:              c,
		Sprites:           []*shared.Sprite{},
		wizards:           make(map[string]*shared.WizardPlayer, 8), // max 8 wizards
		enemies:           make([]*shared.Enemy, 32),                //16 enemies spawn at at time, doubled for headroom
		enemyRespawnTimer: time.Time{},
		projectiles:       make([]*shared.Projectile, 16),                   // there shouldn't ever be more than 16 projectiles
		deadProjectiles:   make([]*shared.Projectile, 16),                   // alive or dead at one time
		projectileCache:   make(map[shared.ProjectileType]*ebiten.Image, 1), // number of projectile types
		tilemapJSON:       tilemap,
		tileCache:         tileCache,
		camera:            nil,
		colliders:         []image.Rectangle{},
		chests:            make([]image.Rectangle, 8), // 8 chests on the map
		openedChests:      make(map[OpenedChest]struct{}, 8),
		holes:             []image.Rectangle{}, //TODO: count holes and traps to preallocate slices
		traps:             []image.Rectangle{},
		trapsUp:           false,
		audioPlayer:       audioPlayer,
		heartImage:        heartImg,
		debug:             false,
	}

	w.projectileCache[shared.Fireball] = fireballImg

	for _, layer := range w.tilemapJSON.Layers {
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := i % layer.Width
			y := i / layer.Width

			x *= shared.TileSize
			y *= shared.TileSize

			id &= ^(int(shared.FlagFlippedHorizontally) |
				int(shared.FlagFlippedVertically) |
				int(shared.FlagFlippedDiagonally) |
				int(shared.FlagRotatedHexagonal120))

			tilemapIndex := shared.GetTilemapIndex(id, w.tilemapJSON)
			src := w.tilemapJSON.Tilesets[tilemapIndex].Source

			if src == "TilesetHole.json" {
				w.holes = append(w.holes, image.Rect(
					x,
					y,
					x+shared.TileSize,
					y+shared.TileSize,
				))
			} else if layer.Name == "traps_up" {
				w.traps = append(w.traps, image.Rect(
					x,
					y,
					x+shared.TileSize,
					y+shared.TileSize,
				))
			} else if layer.Name == "chests" {
				w.chests = append(w.chests, image.Rect(
					x,
					y,
					x+shared.TileSize,
					y+shared.TileSize,
				))
			}
		}
	}

	return w
}
