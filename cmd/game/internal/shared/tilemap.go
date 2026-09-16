package shared

import (
	"encoding/json"
	"fmt"
	"image"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const TileSize = 16
const HalfTile = 8

type Tile struct {
	Img   *ebiten.Image
	Flips struct {
		HorizontalFlip bool
		VerticalFlip   bool
		DiagonalFlip   bool
	}
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

func GetTilemapIndex(id int, tilemapJSON *TilemapJSON) int {
	for i := range tilemapJSON.Tilesets {
		if id < tilemapJSON.Tilesets[i].Firstgid {
			return i - 1
		}
	}

	return len(tilemapJSON.Tilesets) - 1
}

func newTilemapImgList(tilemapJSON *TilemapJSON) ([]*ebiten.Image, error) {
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

func NewTileCache(tilemapJSON *TilemapJSON) (map[int]*Tile, error) {
	imgMap := make(map[int]*Tile)
	tilemapImgList, err := newTilemapImgList(tilemapJSON)
	for _, layer := range tilemapJSON.Layers {
		for _, id := range layer.Data {
			if id == 0 {
				continue
			}

			id &= ^(int(FlagFlippedHorizontally) |
				int(FlagFlippedVertically) |
				int(FlagFlippedDiagonally) |
				int(FlagRotatedHexagonal120))

			if _, ok := imgMap[id]; !ok {
				tilemapImgIndex := GetTilemapIndex(int(id), tilemapJSON)
				if err != nil {
					return nil, err
				}

				tileImg := tilemapImgList[tilemapImgIndex]

				srcX := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) % tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns
				srcY := (id - tilemapJSON.Tilesets[tilemapImgIndex].Firstgid) / tilemapJSON.Tilesets[tilemapImgIndex].Data.Columns

				srcX *= TileSize
				srcY *= TileSize

				imgMap[id] = &Tile{
					Img: tileImg.SubImage(image.Rect(srcX, srcY, srcX+TileSize, srcY+TileSize)).(*ebiten.Image),
				}
			}
		}
	}
	return imgMap, nil
}
