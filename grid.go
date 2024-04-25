package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	cells [][]bool
	snake []Position
}

type Position struct {
	X, Y int
}

func NewGrid(width, height int) *Grid {
	grid := &Grid{
		cells: make([][]bool, height),
		snake: []Position{{X: width / 2, Y: height / 2}},
	}
	for i := range grid.cells {
		grid.cells[i] = make([]bool, width)
	}
	grid.cells[grid.snake[0].Y][grid.snake[0].X] = true
	return grid
}

func (g *Grid) Update() {

}

func (g *Grid) Draw(screen *ebiten.Image) {
	for _, pos := range g.snake {
		rect := image.Rect(pos.X*10, pos.Y*10, (pos.X+1)*10, (pos.Y+1)*10)
		snakePart := ebiten.NewImage(10, 10)
		snakePart.Fill(color.White)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		screen.DrawImage(snakePart, opts)
	}
}
