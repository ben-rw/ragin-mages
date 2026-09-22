package wizarena

import (
	"image"
	"log"
	"time"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// wrap ws.Connection to allow for testing by avoiding importing syscall/js
type MessageConn interface {
	Check() []protocol.Message
	WriteMsg(mt protocol.MessageType, data any)
}

const (
	tilemapPath = "assets/maps/ninja_dungeon.json"
	heartPath   = "assets/images/ninja_adventure/Ui/Receptacle/IconHeart.png"
	songPath    = "assets/audio/music/void-construct-loop-0.6.ogg"
	introLen    = 1 // 1 second song intro
)

type WizArena struct {
	shared.Roster
	Conn                 *MessageConn
	Sprites              []*shared.Sprite
	wizard               *shared.WizardPlayer
	wizards              map[string]*shared.WizardPlayer
	wizardImgCache       map[int]*ebiten.Image
	wizardFrameCache     map[int][]*ebiten.Image
	enemyFrameCache      map[shared.EnemyType][]*ebiten.Image
	projectileFrameCache map[shared.ProjectileType][]*ebiten.Image
	enemies              []*shared.Enemy
	enemyRespawnTimer    time.Time
	roundTimer           time.Time
	statText             map[Stat]string
	tilemapJSON          *shared.TilemapJSON
	tiles                map[string][]*shared.Tile
	camera               *shared.Camera
	colliders            []image.Rectangle
	chests               []image.Rectangle
	chestRespawnTimer    time.Time
	openedChests         map[OpenedChest]struct{}
	holes                []image.Rectangle
	traps                []image.Rectangle
	trapsUp              bool
	audioPlayer          *audio.Player
	projectiles          []*shared.Projectile
	deadProjectiles      []*shared.Projectile
	projectileImageCache map[shared.ProjectileType]*ebiten.Image
	enemyImageCache      map[shared.EnemyType]*ebiten.Image
	heartImage           *ebiten.Image
	debug                bool
}

func NewWizArena(c *MessageConn) *WizArena {
	log.Println("scene changed to Wizards")

	tilemap, err := shared.NewTilemapJSON(tilemapPath)
	if err != nil {
		log.Printf("couldn't load tilemap: %v", err)
	}

	tiles, err := shared.NewOrderedTileList(tilemap)
	if err != nil {
		log.Printf("couldn't build tile cache: %v", err)
	}

	audioPlayer, err := sound.NewAudioPlayer(songPath, true, introLen)
	if err != nil {
		log.Printf("couldn't create audio player: %v", err)
	}

	projectileImgCache, err := shared.NewProjectileImageCache([]shared.ProjectileType{shared.Fireball})
	if err != nil {
		log.Printf("couldn't load projectile image: %v", err)
	}

	heartImg, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, heartPath)
	if err != nil {
		log.Printf("couldn't load image: %v", err)
	}

	enemyImgCache, err := shared.NewEnemyImageCache([]shared.EnemyType{shared.Skeleton})
	if err != nil {
		log.Fatal(err)
	}

	w := &WizArena{
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8), // max 8 players
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:                 c,
		Sprites:              []*shared.Sprite{},
		wizards:              make(map[string]*shared.WizardPlayer, 8),                       // max 8 wizards
		enemies:              make([]*shared.Enemy, 0, 32),                                   //16 enemies spawn at at time, doubled for headroom
		wizardImgCache:       make(map[int]*ebiten.Image, 8),                                 // 8 players
		wizardFrameCache:     make(map[int][]*ebiten.Image, 8),                               // 8 player sprites
		enemyFrameCache:      make(map[shared.EnemyType][]*ebiten.Image, 1),                  // 1 enemy type
		projectileFrameCache: make(map[shared.ProjectileType][]*ebiten.Image, 1),             //1 projectile type
		enemyRespawnTimer:    time.Now().Add(shared.RoundStartEnemySpawnTimer * time.Second), // wait 5 seconds before spawning first group of enemies
		projectiles:          make([]*shared.Projectile, 0, 16),                              // there shouldn't ever be more than 16 projectiles
		deadProjectiles:      make([]*shared.Projectile, 0, 16),                              // alive or dead at one time
		projectileImageCache: projectileImgCache,
		enemyImageCache:      enemyImgCache,
		tilemapJSON:          tilemap,
		tiles:                tiles,
		camera:               nil,
		colliders:            []image.Rectangle{},
		chests:               make([]image.Rectangle, 8), // 8 chests on the map
		openedChests:         make(map[OpenedChest]struct{}, 8),
		holes:                make([]image.Rectangle, 2092), //2091 holes on map
		traps:                make([]image.Rectangle, 84),   //traps on map
		trapsUp:              false,
		audioPlayer:          audioPlayer,
		heartImage:           heartImg,
		debug:                false,
	}

	w.enemyFrameCache[shared.Skeleton] = spritesheet.LoadFrames(
		shared.EnemySpriteSheet,
		enemyImgCache[shared.Skeleton],
	)
	w.projectileFrameCache[shared.Fireball] = spritesheet.LoadFrames(
		shared.ProjectileSpriteSheet,
		projectileImgCache[shared.Fireball],
	)
	for i, path := range shared.PlayerSpriteIndex {
		img, _, err := ebitenutil.NewImageFromFileSystem(shared.AssetsFS, path)
		if err != nil {
			log.Printf("wizarena init: couldnt' create wizard subimage: %v", err)
		}
		w.wizardImgCache[i] = img
		w.wizardFrameCache[i] = spritesheet.LoadFrames(shared.PlayerSpriteSheet, w.wizardImgCache[i])
	}

	w.TileBounds()
	return w
}

func (w *WizArena) TileBounds() {
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

			tilemapIndex := shared.GetTilesetIndex(id, w.tilemapJSON)
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
}
