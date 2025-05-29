package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(outerScreenWidth, outerScreenHeight int) (ScreenWidth, ScreenHeight int) {
	return 640, 480
}

func (g *Game) Update() error {
	return nil
}

func main() {
	g := &Game{}

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
