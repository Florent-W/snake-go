package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	cells [][]bool
	snake []Position
	food  Position
}

type Position struct {
	X, Y int
}

func NewGrid(width, height int) *Grid {
	rand.Seed(time.Now().UnixNano())
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
	// segments du serpent
	for _, pos := range g.snake {
		snakePart := ebiten.NewImage(10, 10)
		snakePart.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(pos.X*10), float64(pos.Y*10))
		screen.DrawImage(snakePart, opts)
	}

	foodPart := ebiten.NewImage(10, 10)
	foodPart.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	foodOpts := &ebiten.DrawImageOptions{}
	foodOpts.GeoM.Translate(float64(g.food.X*10), float64(g.food.Y*10))
	screen.DrawImage(foodPart, foodOpts)
}
