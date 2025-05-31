package main

import (
	"encoding/json"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// first tile ID constants for each tileset
const groundFirstTileID = 1
const waterFirstTileID = 257

// Constants for tile flipping
const (
	FlippedHorizontally = 0x80000000
	FlippedVertically   = 0x40000000
	FlippedDiagonally   = 0x20000000
	FlipMask            = FlippedHorizontally | FlippedVertically | FlippedDiagonally
)


// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
	Name   string `json:"name"`
}

type TileImageMap struct {
	Images map[int]*ebiten.Image
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

func (t TilemapJSON) LoadTiles() TileImageMap {
	tileMap := TileImageMap{
		Images: make(map[int]*ebiten.Image),
	}

	// collect all unique tile IDs from all layers
	uniqueTileIDs := make(map[int]bool)
	for _, layer := range t.Layers {
		for _, tileID := range layer.Data {
			if tileID != 0 { // 0 typically represents empty/no tile
				uniqueTileIDs[tileID] = true
			}
		}
	}

	// load image for each unique tile ID
	for tileID := range uniqueTileIDs {
		i,err := getTileImage(tileID)
		if err != nil{
			panic(err)
		}
		tileMap.Images[tileID] = i
	}

	return tileMap
}

// opens the file, parses it, and returns the json object + potential error
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}

// getTileImage returns the ebiten image for a given tile ID
func getTileImage(tileID int) (*ebiten.Image, error) {
    flippedH := (tileID & FlippedHorizontally) != 0
    flippedV := (tileID & FlippedVertically) != 0
    flippedD := (tileID & FlippedDiagonally) != 0

    // Get the actual tile ID without flip flags
    actualTileID := tileID &^ FlipMask // Use defined FlipMask

    var tilesetImage *ebiten.Image
    var err error
    var localTileID int
    var tileWidth, tileHeight int = 64, 64 // Standard tile size
    var tilesPerRow int

    if actualTileID >= groundFirstTileID && actualTileID < waterFirstTileID {
        // Ground tileset
        tilesetImage, _, err = ebitenutil.NewImageFromFile("assets/Grass Tileset.png")
        if err != nil {
            return nil, err
        }
        localTileID = actualTileID - groundFirstTileID // Use actualTileID, not tileID
        tilesPerRow = 16 // Ground tileset has 16 columns
    } else if actualTileID >= waterFirstTileID { // Use actualTileID, not tileID
        // Water tileset
        tilesetImage, _, err = ebitenutil.NewImageFromFile("assets/Animated water tiles.png")
        if err != nil {
            return nil, err
        }
        localTileID = actualTileID - waterFirstTileID // Use actualTileID, not tileID
        tilesPerRow = 70 // Water tileset has 70 columns
    } else {
        return nil, nil // Invalid tile ID
    }

    // Calculate tile position in the tileset
    tileX := (localTileID % tilesPerRow) * tileWidth
    tileY := (localTileID / tilesPerRow) * tileHeight

    // Extract the tile from the tileset
    tileRect := image.Rect(tileX, tileY, tileX+tileWidth, tileY+tileHeight)
    tileImage := tilesetImage.SubImage(tileRect).(*ebiten.Image)

    // Apply flips if needed
    if flippedH || flippedV || flippedD {
        tileImage = applyFlips(tileImage, flippedH, flippedV, flippedD)
    }
    
    return tileImage, nil
}

func applyFlips(img *ebiten.Image, flippedH, flippedV, flippedD bool) *ebiten.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Create a new image for the flipped result
	flippedImg := ebiten.NewImage(w, h)

	opts := &ebiten.DrawImageOptions{}

	// Handle diagonal flip (90° rotation) first
	if flippedD {
		// 90° clockwise rotation + horizontal flip = diagonal flip in Tiled
		opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		opts.GeoM.Rotate(3.14159 / 2) // 90 degrees in radians
		opts.GeoM.Scale(-1, 1)        // Horizontal flip
		opts.GeoM.Translate(float64(h)/2, float64(w)/2)
	} else {
		// Handle horizontal flip
		if flippedH {
			opts.GeoM.Scale(-1, 1)
			opts.GeoM.Translate(float64(w), 0)
		}

		// Handle vertical flip
		if flippedV {
			opts.GeoM.Scale(1, -1)
			opts.GeoM.Translate(0, float64(h))
		}
	}

	flippedImg.DrawImage(img, opts)
	return flippedImg
}


