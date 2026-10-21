package wizarena

import (
	"image"
	"log"
	"sync/atomic"
	"time"

	"github.com/ben-rw/ragin-mages/internal/gameAssets"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/sound"
	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared/spritesheet"
	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// wrap ws.Connection to avoid syscall/js import
// this allows for benchmarking with go test -bench
type MessageConn interface {
	Check() []*protocol.Message
	WriteMsg(mt protocol.MessageType, data any)
}

const (
	heartPath = "assets/images/ninja_adventure/Ui/Receptacle/IconHeart.png"
	songPath  = "assets/audio/music/void-construct-loop-0.6.ogg"
	introLen  = 1 // 1 second song intro
)

type WizArena struct {
	phaseTimer time.Time
	shared.Roster
	Conn                   MessageConn
	wizard                 *WizardPlayer
	wizards                map[string]*WizardPlayer
	wizardImgCache         map[int]*ebiten.Image
	wizardFrameCache       map[int][]*ebiten.Image
	enemyFrameCache        map[EnemyType][]*ebiten.Image
	projectileFrameCache   map[ProjectileType][]*ebiten.Image
	reflectFrameCache      map[ReflectType][]*ebiten.Image
	dynamicText            map[DynamicText]string
	dynamicTextImageCache  map[DynamicText]*ebiten.Image
	tilemapJSON            *wizServer.TilemapJSON
	tiles                  map[string][]*shared.Tile
	staticTilemapTrapsUp   *ebiten.Image
	staticTilemapTrapsDown *ebiten.Image
	camera                 *shared.Camera
	openedChests           map[OpenedChest]struct{}
	audioPlayer            *audio.Player
	projectileImageCache   map[ProjectileType]*ebiten.Image
	enemyImageCache        map[EnemyType]*ebiten.Image
	heartImage             *ebiten.Image
	reflectImage           *ebiten.Image
	enemies                map[uint16]*Enemy
	sortedEnemies          []*Enemy
	sortedWizards          []*WizardPlayer
	colliders              []image.Rectangle
	chests                 []image.Rectangle
	holes                  []image.Rectangle
	traps                  []image.Rectangle
	trapTimer              int
	projectiles            map[uint16]*Projectile
	sortedProjectiles      []*Projectile
	scoreboard             string
	roundTimerString       string
	roomID                 string
	phase                  wizServer.RoundPhase
	screenAlpha            float32
	backgroundAlpha        uint8
	lastFormattedSecond    int
	round                  uint8
	trapsUp                bool
	debug                  bool
	displayFPS             bool
	scoreboardUpdated      bool
	scoresUpdated          bool
	scoresUpdatedCount     atomic.Int32
}

func NewWizArena(c MessageConn, roomID string) *WizArena {
	tilemap, err := wizServer.NewTilemapJSON(wizServer.TilemapPath)
	if err != nil {
		log.Printf("couldn't load tilemap: %v\n", err)
	}

	tiles, err := shared.NewOrderedTileList(tilemap)
	if err != nil {
		log.Printf("couldn't build tile cache: %v\n", err)
	}

	audioPlayer, err := sound.NewAudioPlayer(songPath, true, introLen)
	if err != nil {
		log.Printf("couldn't create audio player: %v\n", err)
	}

	projectileImgCache, err := NewProjectileImageCache([]ProjectileType{Fireball})
	if err != nil {
		log.Printf("couldn't load projectile image: %v\n", err)
	}

	heartImg, _, err := ebitenutil.NewImageFromFileSystem(gameAssets.AssetsFS, heartPath)
	if err != nil {
		log.Printf("couldn't load image: %v\n", err)
	}

	reflectImg, _, err := ebitenutil.NewImageFromFileSystem(gameAssets.AssetsFS, ReflectPath)
	if err != nil {
		log.Printf("couldn't load image: %v\n", err)
	}

	enemyImgCache, err := NewEnemyImageCache([]EnemyType{Skeleton})
	if err != nil {
		log.Printf("couldn't load enemy image cache: %v\n", err)
	}

	w := &WizArena{
		phase: wizServer.WaitingForClients,
		Roster: shared.Roster{
			Players: make(map[string]*shared.Player, 8), // max 8 players
			Player:  shared.NewPlayer(&protocol.PlayerData{}, 0),
		},
		Conn:                 c,
		wizards:              make(map[string]*WizardPlayer, 8),           // max 8 wizards
		enemies:              make(map[uint16]*Enemy, 32),                 // 16 enemies spawn at at time, doubled for headroom
		wizardImgCache:       make(map[int]*ebiten.Image, 8),              // 8 players
		wizardFrameCache:     make(map[int][]*ebiten.Image, 8),            // 8 player sprites
		enemyFrameCache:      make(map[EnemyType][]*ebiten.Image, 1),      // 1 enemy type
		projectileFrameCache: make(map[ProjectileType][]*ebiten.Image, 1), // 1 projectile type
		reflectFrameCache:    make(map[ReflectType][]*ebiten.Image, 1),    // 1 reflect img
		projectiles:          make(map[uint16]*Projectile, 16),            // there shouldn't ever be more than 16 projectiles at one time
		projectileImageCache: projectileImgCache,
		enemyImageCache:      enemyImgCache,
		tilemapJSON:          tilemap,
		tiles:                tiles,
		camera:               nil,
		colliders:            make([]image.Rectangle, 0, 972), // 972 colliders on map
		chests:               make([]image.Rectangle, 0, 8),   // 8 chests on the map
		openedChests:         make(map[OpenedChest]struct{}, 8),
		holes:                make([]image.Rectangle, 0, 2092), // 2092 holes on map
		traps:                make([]image.Rectangle, 84),      // traps on map
		trapsUp:              false,
		audioPlayer:          audioPlayer,
		heartImage:           heartImg,
		reflectImage:         reflectImg,
		roomID:               roomID,
		backgroundAlpha:      255,
		round:                1,
		trapTimer:            100, // traps start up for 100 ticks
		debug:                false,
		displayFPS:           false,
		scoreboardUpdated:    false,
		scoresUpdated:        false,
	}

	// preload entity images
	w.enemyFrameCache[Skeleton] = spritesheet.LoadFrames(
		EnemySpriteSheet,
		enemyImgCache[Skeleton],
	)
	w.projectileFrameCache[Fireball] = spritesheet.LoadFrames(
		ProjectileSpriteSheet,
		projectileImgCache[Fireball],
	)
	w.reflectFrameCache[Circle] = spritesheet.LoadFrames(
		ReflectSpriteSheet,
		reflectImg,
	)
	for i, path := range shared.PlayerSpriteIndex {
		img, _, err := ebitenutil.NewImageFromFileSystem(gameAssets.AssetsFS, path)
		if err != nil {
			log.Printf("wizarena init: couldnt' create wizard subimage: %v\n", err)
		}
		w.wizardImgCache[i] = img
		w.wizardFrameCache[i] = spritesheet.LoadFrames(shared.PlayerSpriteSheet, w.wizardImgCache[i])
	}

	// set up tile collisions
	w.TileBounds()

	// preload map images
	w.staticTilemapTrapsUp = w.NewStaticTilemapTrapsUp()
	w.staticTilemapTrapsDown = w.NewStaticTilemapTrapsDown()

	// tell server this client is ready
	w.Conn.WriteMsg(protocol.ClientLoaded, &protocol.ClientLoadedData{
		Loaded: true,
	})

	return w
}

func (w *WizArena) NewStaticTilemapTrapsUp() *ebiten.Image {
	var opts ebiten.DrawImageOptions
	img := ebiten.NewImage(
		w.tilemapJSON.Layers[0].Width*wizServer.TileSize,
		w.tilemapJSON.Layers[0].Height*wizServer.TileSize,
	)

	for _, layer := range w.tilemapJSON.Layers {
		if layer.Name == "traps_down" {
			continue
		}
		if layer.Name == "chests" {
			continue
		}

		for _, tile := range w.tiles[layer.Name] {
			if tile.Flips.HorizontalFlip ||
				tile.Flips.VerticalFlip ||
				tile.Flips.DiagonalFlip {
				shared.FixRotatedTile(tile, &opts)
			}

			opts.GeoM.Translate(float64(tile.X), float64(tile.Y))

			img.DrawImage(
				tile.Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}
	return img
}

func (w *WizArena) NewStaticTilemapTrapsDown() *ebiten.Image {
	var opts ebiten.DrawImageOptions
	img := ebiten.NewImage(
		w.tilemapJSON.Layers[0].Width*wizServer.TileSize,
		w.tilemapJSON.Layers[0].Height*wizServer.TileSize,
	)

	for _, layer := range w.tilemapJSON.Layers {
		if layer.Name == "traps_up" {
			continue
		}
		if layer.Name == "chests" {
			continue
		}

		for _, tile := range w.tiles[layer.Name] {
			if tile.Flips.HorizontalFlip ||
				tile.Flips.VerticalFlip ||
				tile.Flips.DiagonalFlip {
				shared.FixRotatedTile(tile, &opts)
			}

			opts.GeoM.Translate(float64(tile.X), float64(tile.Y))

			img.DrawImage(
				tile.Img,
				&opts,
			)

			opts.GeoM.Reset()
		}
	}
	return img
}

func (w *WizArena) TileBounds() {
	for _, layer := range w.tilemapJSON.Layers {
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := i % layer.Width
			y := i / layer.Width

			x *= wizServer.TileSize
			y *= wizServer.TileSize

			id &= ^(int(wizServer.FlagFlippedHorizontally) |
				int(wizServer.FlagFlippedVertically) |
				int(wizServer.FlagFlippedDiagonally) |
				int(wizServer.FlagRotatedHexagonal120))

			tilemapIndex := wizServer.GetTilesetIndex(id, w.tilemapJSON)
			src := w.tilemapJSON.Tilesets[tilemapIndex].Source

			if src == "TilesetHole.json" {
				w.holes = append(w.holes, image.Rect(
					x,
					y,
					x+wizServer.TileSize,
					y+wizServer.TileSize,
				))
				if id != 387 {
					w.colliders = append(w.colliders, image.Rect(
						x,
						y,
						x+wizServer.TileSize,
						y+wizServer.TileSize,
					))
				}
			} else if layer.Name == "traps_up" {
				w.traps = append(w.traps, image.Rect(
					x,
					y,
					x+wizServer.TileSize,
					y+wizServer.TileSize,
				))
			} else if layer.Name == "chests" {
				w.chests = append(w.chests, image.Rect(
					x,
					y,
					x+wizServer.TileSize,
					y+wizServer.TileSize,
				))
			}
		}
	}
}
