package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	level *TilemapJSON
	images TileImageMap
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.level == nil {
		return
	}
	
	const tileSize = 64 // Standard tile size
	
	// Draw layers in reverse order (last layer first)
	//for i := len(g.level.Layers) - 1; i >= 0; i-- {
	for i := 0; i < len(g.level.Layers); i++ {

		layer := g.level.Layers[i]
		
		// Draw each tile in the layer
		for y := 0; y < layer.Height; y++ {
			for x := 0; x < layer.Width; x++ {
				index := y*layer.Width + x
				if index < len(layer.Data) {
					tileID := layer.Data[index]
					
					// Skip empty tiles (ID 0)
					if tileID == 0 {
						continue
					}
					
					// Get the tile image from the map
					if tileImage, exists := g.images.Images[tileID]; exists && tileImage != nil {
						// Calculate screen position
						screenX := float64(x * tileSize)
						screenY := float64(y * tileSize)
						
						// Draw the tile
						opts := &ebiten.DrawImageOptions{}
						opts.GeoM.Translate(screenX, screenY)
						screen.DrawImage(tileImage, opts)
					}
				}
			}
		}
	}
}

func (g *Game) Layout(outerScreenWidth, outerScreenHeight int) (ScreenWidth, ScreenHeight int) {
	return 1920, 1280
}

func (g *Game) Update() error {
	return nil
}

func main() {
	g := &Game{}
	t, err := NewTilemapJSON("assets/level.tmj")
	if err != nil {
		panic(err)
	}
	g.level = t
	g.images = t.LoadTiles()

	err = ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
