package shared

import (
	"fmt"
	"image"
	"math"
	"path"
	"sort"

	"github.com/ben-rw/ragin-mages/internal/gameAssets"
	"github.com/ben-rw/ragin-mages/internal/minigames/wizServer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tile struct {
	Img   *ebiten.Image
	Flips struct {
		HorizontalFlip bool
		VerticalFlip   bool
		DiagonalFlip   bool
	}
	X, Y, ImgId int
}

func newTilesetImageList(tilemapJSON *wizServer.TilemapJSON) ([]*ebiten.Image, error) {
	imgList := make([]*ebiten.Image, len(tilemapJSON.Tilesets))
	for i := range tilemapJSON.Tilesets {
		img, _, err := ebitenutil.NewImageFromFileSystem(gameAssets.AssetsFS, (fmt.Sprintf("assets%v", path.Clean("/"+tilemapJSON.Tilesets[i].Data.ImagePath))))
		if err != nil {
			return nil, err
		}
		imgList[i] = img
	}

	return imgList, nil
}

func NewTileImageCache(tilemapJSON *wizServer.TilemapJSON, numberOfTileImgs int) (map[int]*ebiten.Image, error) {
	imgMap := make(map[int]*ebiten.Image, numberOfTileImgs)
	tilemapImgList, err := newTilesetImageList(tilemapJSON)
	if err != nil {
		return nil, err
	}
	for _, layer := range tilemapJSON.Layers {
		for _, id := range layer.Data {
			if id == 0 {
				continue
			}

			if _, ok := imgMap[id]; !ok {
				id &= ^(int(wizServer.FlagFlippedHorizontally) |
					int(wizServer.FlagFlippedVertically) |
					int(wizServer.FlagFlippedDiagonally) |
					int(wizServer.FlagRotatedHexagonal120))

				tilemapImgIndex := wizServer.GetTilesetIndex(int(id), tilemapJSON)

				tileImg := tilemapImgList[tilemapImgIndex]

				srcX := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) % tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns
				srcY := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) / tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns

				srcX *= wizServer.TileSize
				srcY *= wizServer.TileSize

				imgMap[id] = tileImg.SubImage(image.Rect(srcX, srcY, srcX+wizServer.TileSize, srcY+wizServer.TileSize)).(*ebiten.Image)
			}
		}
	}
	return imgMap, nil
}

func FixRotatedTile(tile *Tile, opts *ebiten.DrawImageOptions) {
	switch {
	case tile.Flips.HorizontalFlip && tile.Flips.VerticalFlip:
		opts.GeoM.Scale(-1, -1)
		opts.GeoM.Translate(wizServer.TileSize, 0)
		opts.GeoM.Translate(0, wizServer.TileSize)
	case tile.Flips.DiagonalFlip && tile.Flips.HorizontalFlip:
		opts.GeoM.Translate(-wizServer.TileSize/2, -wizServer.TileSize/2)
		opts.GeoM.Rotate(math.Pi / 2)
		opts.GeoM.Translate(wizServer.TileSize/2, wizServer.TileSize/2)
	case tile.Flips.HorizontalFlip:
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(wizServer.TileSize, 0)
	case tile.Flips.DiagonalFlip && tile.Flips.VerticalFlip:
		opts.GeoM.Translate(-wizServer.TileSize/2, -wizServer.TileSize/2)
		opts.GeoM.Rotate(3 * math.Pi / 2)
		opts.GeoM.Translate(wizServer.TileSize/2, wizServer.TileSize/2)
	case tile.Flips.VerticalFlip:
		opts.GeoM.Scale(1, -1)
		opts.GeoM.Translate(0, wizServer.TileSize)
	default:
	}
}

func NewOrderedTileList(tilemapJSON *wizServer.TilemapJSON) (map[string][]*Tile, error) {
	tileImgCache, err := NewTileImageCache(tilemapJSON, wizServer.NumberOfTileImgs)
	if err != nil {
		return nil, err
	}
	tiles := make(map[string][]*Tile, 7000)
	for _, layer := range tilemapJSON.Layers {
		layerTiles := []*Tile{}
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			tile := Tile{}
			if id&int(wizServer.FlagFlippedHorizontally) != 0 {
				tile.Flips.HorizontalFlip = true
			}
			if id&int(wizServer.FlagFlippedVertically) != 0 {
				tile.Flips.VerticalFlip = true
			}
			if id&int(wizServer.FlagFlippedDiagonally) != 0 {
				tile.Flips.DiagonalFlip = true
			}

			id &= ^(int(wizServer.FlagFlippedHorizontally) |
				int(wizServer.FlagFlippedVertically) |
				int(wizServer.FlagFlippedDiagonally) |
				int(wizServer.FlagRotatedHexagonal120))

			x := i % layer.Width
			y := i / layer.Width

			x *= wizServer.TileSize
			y *= wizServer.TileSize

			tile.Img = tileImgCache[id]
			tile.X = x
			tile.Y = y
			tile.ImgId = id
			layerTiles = append(layerTiles, &tile)
		}

		sort.Slice(layerTiles, func(i, j int) bool {
			return layerTiles[i].ImgId < layerTiles[j].ImgId
		})
		tiles[layer.Name] = layerTiles
	}

	return tiles, nil
}
