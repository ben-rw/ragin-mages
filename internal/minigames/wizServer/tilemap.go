package wizServer

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/ben-rw/ragin-mages/internal/gameAssets"
)

const TilemapPath = "assets/maps/ninja_dungeon.json"

const TileSize = 16
const HalfTile = 8
const NumberOfTileImgs = 46

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
	Name   string `json:"name"`
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Tileset struct {
	Data struct {
		Columns   int    `json:"columns"`
		ImagePath string `json:"image"`
	}
	Source   string `json:"source"`
	Firstgid int    `json:"firstgid"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	data, err := gameAssets.AssetsFS.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(data, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	for i := range tilemapJSON.Tilesets {
		data, err := gameAssets.AssetsFS.ReadFile(fmt.Sprintf("assets/maps/%v", tilemapJSON.Tilesets[i].Source))
		var tilesetData Tileset
		err = json.Unmarshal(data, &tilesetData.Data)
		if err != nil {
			return nil, err
		}
		tilemapJSON.Tilesets[i].Data = tilesetData.Data
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

func (w *WizardServer) TileBounds() {
	for _, layer := range w.tilemapJSON.Layers {
		for i, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := i % layer.Width
			y := i / layer.Width

			x *= TileSize
			y *= TileSize

			id &= ^(int(FlagFlippedHorizontally) |
				int(FlagFlippedVertically) |
				int(FlagFlippedDiagonally) |
				int(FlagRotatedHexagonal120))

			tilemapIndex := GetTilesetIndex(id, w.tilemapJSON)
			src := w.tilemapJSON.Tilesets[tilemapIndex].Source

			if src == "TilesetHole.json" {
				w.holes = append(w.holes, image.Rect(
					x,
					y,
					x+TileSize,
					y+TileSize,
				))
				// if id != 387 {
				// 	w.colliders = append(w.colliders, image.Rect(
				// 		x,
				// 		y,
				// 		x+TileSize,
				// 		y+TileSize,
				// 	))
				// }
			}
			// else if layer.Name == "traps_up" {
			// 	w.traps = append(w.traps, image.Rect(
			// 		x,
			// 		y,
			// 		x+TileSize,
			// 		y+TileSize,
			// 	))
			// } else if layer.Name == "chests" {
			// 	w.chests = append(w.chests, image.Rect(
			// 		x,
			// 		y,
			// 		x+TileSize,
			// 		y+TileSize,
			// 	))
			// }
		}
	}
}
