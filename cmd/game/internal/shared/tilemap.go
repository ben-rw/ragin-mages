package shared

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"path"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const TileSize = 16
const HalfTile = 8

const NumberOfTileImgs = 46

type Tile struct {
	Img   *ebiten.Image
	Flips struct {
		HorizontalFlip bool
		VerticalFlip   bool
		DiagonalFlip   bool
	}
	X, Y, ImgId int
}

const (
	FlagFlippedHorizontally uint32 = 1 << 31
	FlagFlippedVertically   uint32 = 1 << 30
	FlagFlippedDiagonally   uint32 = 1 << 29
	FlagRotatedHexagonal120 uint32 = 1 << 28 //not used. tiled docs say to unset it anyway
)

type TilemapJSON struct {
	Layers   []*TilemapLayerJSON `json:"layers"`
	Tilesets []*Tileset          `json:"tilesets"`
}

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type Tileset struct {
	Firstgid int    `json:"firstgid"`
	Source   string `json:"source"`
	Data     struct {
		Columns   int    `json:"columns"`
		ImagePath string `json:"image"`
	}
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	data, err := AssetsFS.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(data, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	for _, tileset := range tilemapJSON.Tilesets {
		data, err := AssetsFS.ReadFile(fmt.Sprintf("assets/maps/%v", tileset.Source))
		var tilesetData Tileset
		err = json.Unmarshal(data, &tilesetData.Data)
		if err != nil {
			return nil, err
		}
		tileset.Data = tilesetData.Data
	}

	return &tilemapJSON, nil
}

func GetTilesetIndex(id int, tilemapJSON *TilemapJSON) int {
	for i := range tilemapJSON.Tilesets {
		if id < tilemapJSON.Tilesets[i].Firstgid {
			return i - 1
		}
	}

	return len(tilemapJSON.Tilesets) - 1
}

func newTilesetImageList(tilemapJSON *TilemapJSON) ([]*ebiten.Image, error) {
	imgList := make([]*ebiten.Image, len(tilemapJSON.Tilesets))
	for i := range tilemapJSON.Tilesets {
		img, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, (fmt.Sprintf("assets%v", path.Clean("/"+tilemapJSON.Tilesets[i].Data.ImagePath))))
		if err != nil {
			return nil, err
		}
		imgList[i] = img
	}

	return imgList, nil
}

func NewTileImageCache(tilemapJSON *TilemapJSON) (map[int]*ebiten.Image, error) {
	imgMap := make(map[int]*ebiten.Image, NumberOfTileImgs)
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
				id &= ^(int(FlagFlippedHorizontally) |
					int(FlagFlippedVertically) |
					int(FlagFlippedDiagonally) |
					int(FlagRotatedHexagonal120))

				tilemapImgIndex := GetTilesetIndex(int(id), tilemapJSON)

				tileImg := tilemapImgList[tilemapImgIndex]

				srcX := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) % tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns
				srcY := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) / tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns

				srcX *= TileSize
				srcY *= TileSize

				imgMap[id] = tileImg.SubImage(image.Rect(srcX, srcY, srcX+TileSize, srcY+TileSize)).(*ebiten.Image)
			}
		}
	}
	return imgMap, nil
}

func FixRotatedTile(tile *Tile, opts *ebiten.DrawImageOptions) {
	switch {
	case tile.Flips.HorizontalFlip && tile.Flips.VerticalFlip:
		opts.GeoM.Scale(-1, -1)
		opts.GeoM.Translate(TileSize, 0)
		opts.GeoM.Translate(0, TileSize)
	case tile.Flips.DiagonalFlip && tile.Flips.HorizontalFlip:
		opts.GeoM.Translate(-TileSize/2, -TileSize/2)
		opts.GeoM.Rotate(math.Pi / 2)
		opts.GeoM.Translate(TileSize/2, TileSize/2)
	case tile.Flips.HorizontalFlip:
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(TileSize, 0)
	case tile.Flips.DiagonalFlip && tile.Flips.VerticalFlip:
		opts.GeoM.Translate(-TileSize/2, -TileSize/2)
		opts.GeoM.Rotate(3 * math.Pi / 2)
		opts.GeoM.Translate(TileSize/2, TileSize/2)
	case tile.Flips.VerticalFlip:
		opts.GeoM.Scale(1, -1)
		opts.GeoM.Translate(0, TileSize)
	default:
	}
}

func NewOrderedTileList(tilemapJSON *TilemapJSON) (map[string][]*Tile, error) {
	tileImgCache, err := NewTileImageCache(tilemapJSON)
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
			if id&int(FlagFlippedHorizontally) != 0 {
				tile.Flips.HorizontalFlip = true
			}
			if id&int(FlagFlippedVertically) != 0 {
				tile.Flips.VerticalFlip = true
			}
			if id&int(FlagFlippedDiagonally) != 0 {
				tile.Flips.DiagonalFlip = true
			}

			id &= ^(int(FlagFlippedHorizontally) |
				int(FlagFlippedVertically) |
				int(FlagFlippedDiagonally) |
				int(FlagRotatedHexagonal120))

			x := i % layer.Width
			y := i / layer.Width

			x *= TileSize
			y *= TileSize

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
