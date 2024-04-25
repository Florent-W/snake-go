package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	cells [][]bool
}

func NewGrid(width, height int) *Grid {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &Grid{cells: cells}
}

func (g *Grid) Update() {

}

func (g *Grid) Draw(screen *ebiten.Image) {

}
