package main

import "github.com/hajimehoshi/ebiten/v2"

type GridManager interface {
	Update(game *Game) error
	Draw(screen *ebiten.Image)
}
