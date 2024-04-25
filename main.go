package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridSize     = 20
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Go")

	game := &Game{
		gridManager: NewGrid(screenWidth/gridSize, screenHeight/gridSize),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
