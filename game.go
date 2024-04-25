package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	gridManager GridManager
}

func (g *Game) Update() error {
	g.gridManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gridManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
