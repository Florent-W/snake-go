package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Grid struct {
	cells         [][]bool
	snake         []Position
	food          Position
	direction     Direction
	width, height int
}

type Position struct {
	X, Y int
}

func NewGrid(width, height int) *Grid {
	rand.Seed(time.Now().UnixNano())
	grid := &Grid{
		cells:     make([][]bool, height),
		snake:     []Position{{X: width / 2, Y: height / 2}},
		direction: Right,
		width:     width,
		height:    height,
	}
	for i := range grid.cells {
		grid.cells[i] = make([]bool, width)
	}
	grid.cells[grid.snake[0].Y][grid.snake[0].X] = true
	grid.placeFood()
	return grid
}

func (g *Grid) placeFood() {
	if g.width <= 0 || g.height <= 0 {
		fmt.Println("Erreur de taille de la grille lors du placement de la nourriture")
		return
	}
	// Pour éviter que la nourriture soit placé trop près du bord
	margin := 1

	foodX := rand.Intn(g.width-2*margin) + margin
	foodY := rand.Intn(g.height-2*margin) + margin
	g.food = Position{X: foodX, Y: foodY}

	// Pour éviter que la nourriture soit placé sur le snake
	for _, pos := range g.snake {
		if pos == g.food {
			g.placeFood()
			return
		}
	}
}

func (g *Grid) Update(game *Game) error {
	// Changement de direction
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction != Down {
		g.direction = Up
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction != Up {
		g.direction = Down
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction != Right {
		g.direction = Left
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction != Left {
		g.direction = Right
	}

	head := g.snake[0]
	newHead := head
	switch g.direction {
	case Up:
		newHead.Y--
	case Down:
		newHead.Y++
	case Left:
		newHead.X--
	case Right:
		newHead.X++
	}

	// Vérifie les collisions avec les murs
	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		return fmt.Errorf("game over: collision avec un mur")
	}

	// Vérifie les collisions avec lui-même
	for i, segment := range g.snake[1:] {
		if newHead == segment {
			return fmt.Errorf("game over: collision avec soi-même à l'index %d", i+1)
		}
	}

	// Vérifie si la nourriture est mangée
	if newHead == g.food {
		game.score++
		g.snake = append([]Position{newHead}, g.snake...)
		g.placeFood()
	} else {
		g.snake = append([]Position{newHead}, g.snake[:len(g.snake)-1]...)
	}

	return nil
}

func (g *Grid) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()
	gridWidth := int(float64(screenWidth) * 1)
	gridHeight := int(float64(screenHeight) * 1)
	borderThickness := 5

	gridX := (screenWidth - gridWidth) / 2
	gridY := (screenHeight - gridHeight) / 2

	// Bordure
	borderColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	borderImage := ebiten.NewImage(gridWidth+2*borderThickness, gridHeight+2*borderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-borderThickness), float64(gridY-borderThickness))
	screen.DrawImage(borderImage, borderOpts)

	gameArea := ebiten.NewImage(gridWidth, gridHeight)
	gameArea.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	cellSize := 10
	for _, pos := range g.snake {
		snakePart := ebiten.NewImage(cellSize, cellSize)
		snakePart.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(gridX+pos.X*cellSize), float64(gridY+pos.Y*cellSize))
		screen.DrawImage(snakePart, opts)
	}

	foodPart := ebiten.NewImage(cellSize, cellSize)
	foodPart.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	foodOpts := &ebiten.DrawImageOptions{}
	foodOpts.GeoM.Translate(float64(gridX+g.food.X*cellSize), float64(gridY+g.food.Y*cellSize))
	screen.DrawImage(foodPart, foodOpts)
}
