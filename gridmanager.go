package main

import "github.com/hajimehoshi/ebiten/v2"

type GridManager interface {
	Update()
	Draw(screen *ebiten.Image)
}
