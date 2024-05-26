package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"snake-go/src/audio"
	"snake-go/src/constants"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Position struct {
	X, Y int
}

type Grid struct {
	cells         [][]bool
	snake         []Position
	food          Position
	obstacles     []Position
	direction     Direction
	nextDirection Direction
	width, height int
}

func NewGrid(cellSize int) *Grid {
	width := constants.GridWidth / cellSize
	height := constants.GridHeight / cellSize

	rand.Seed(time.Now().UnixNano())
	initialDirection := Right

	grid := &Grid{
		cells:         make([][]bool, height),
		snake:         []Position{{X: width / 2, Y: height / 2}},
		direction:     initialDirection,
		nextDirection: initialDirection,
		width:         width,
		height:        height,
	}
	for i := range grid.cells {
		grid.cells[i] = make([]bool, width)
	}
	grid.cells[grid.snake[0].Y][grid.snake[0].X] = true
	grid.placeFood()
	return grid
}

func NewGridWithObstacles(cellSize int, difficulty Difficulty) *Grid {
	width := constants.GridWidth / cellSize
	height := constants.GridHeight / cellSize

	rand.Seed(time.Now().UnixNano())
	initialDirection := Right

	grid := &Grid{
		cells:         make([][]bool, height),
		snake:         []Position{{X: width / 2, Y: height / 2}},
		direction:     initialDirection,
		nextDirection: initialDirection,
		width:         width,
		height:        height,
	}
	for i := range grid.cells {
		grid.cells[i] = make([]bool, width)
	}
	grid.cells[grid.snake[0].Y][grid.snake[0].X] = true
	grid.placeFood()
	grid.placeObstacles(difficulty)
	return grid
}

func (g *Grid) placeFood() {
	margin := 1

	foodX := rand.Intn(g.width-2*margin) + margin
	foodY := rand.Intn(g.height-2*margin) + margin
	g.food = Position{X: foodX, Y: foodY}

	for _, pos := range g.snake {
		if pos == g.food {
			g.placeFood()
			return
		}
	}

	for _, pos := range g.obstacles {
		if pos == g.food {
			g.placeFood()
			return
		}
	}
}

func (g *Grid) placeObstacles(difficulty Difficulty) {
	var obstacleCount int
	switch difficulty {
	case Facile:
		obstacleCount = 2
	case Normal:
		obstacleCount = 3
	case Difficile:
		obstacleCount = 5
	}

	margin := 1

	for i := 0; i < obstacleCount; i++ {
		obstacleX := rand.Intn(g.width-2*margin) + margin
		obstacleY := rand.Intn(g.height-2*margin) + margin
		obstacle := Position{X: obstacleX, Y: obstacleY}

		for g.cells[obstacle.Y][obstacle.X] {
			obstacleX = rand.Intn(g.width-2*margin) + margin
			obstacleY = rand.Intn(g.height-2*margin) + margin
			obstacle = Position{X: obstacleX, Y: obstacleY}
		}

		g.obstacles = append(g.obstacles, obstacle)
		g.cells[obstacle.Y][obstacle.X] = true
	}
}

func (g *Grid) Update(game *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction != Down {
		g.nextDirection = Up
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction != Up {
		g.nextDirection = Down
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction != Right {
		g.nextDirection = Left
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction != Left {
		g.nextDirection = Right
	}

	if g.nextDirection != g.direction {
		g.direction = g.nextDirection
		audio.MoveSoundPlayer.Rewind()
		audio.MoveSoundPlayer.Play()
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

	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		audio.LoseSoundPlayer.Rewind()
		audio.LoseSoundPlayer.Play()
		return fmt.Errorf("game over: collision avec un mur")
	}

	for _, segment := range g.snake[1:] {
		if newHead == segment {
			audio.LoseSoundPlayer.Rewind()
			audio.LoseSoundPlayer.Play()
			return fmt.Errorf("game over: collision avec soi-même")
		}
	}

	for _, obstacle := range g.obstacles {
		if newHead == obstacle {
			audio.LoseSoundPlayer.Rewind()
			audio.LoseSoundPlayer.Play()
			return fmt.Errorf("game over: collision avec un obstacle")
		}
	}

	if newHead == g.food {
		game.Score++
		audio.EatSoundPlayer.Rewind()
		audio.EatSoundPlayer.Play()
		g.snake = append([]Position{newHead}, g.snake...)
		g.placeFood()
	} else {
		g.snake = append([]Position{newHead}, g.snake[:len(g.snake)-1]...)
	}

	return nil
}

func (g *Grid) Draw(screen *ebiten.Image) {
	gridX := (constants.ScreenWidth - constants.GridWidth) / 2
	gridY := (constants.ScreenHeight - constants.GridHeight) / 2

	borderColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	borderImage := ebiten.NewImage(constants.GridWidth+2*constants.BorderThickness, constants.GridHeight+2*constants.BorderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-constants.BorderThickness), float64(gridY-constants.BorderThickness))
	screen.DrawImage(borderImage, borderOpts)

	gameArea := ebiten.NewImage(constants.GridWidth, constants.GridHeight)
	gameArea.Fill(color.Black)
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	for _, pos := range g.snake {
		snakePart := ebiten.NewImage(constants.CellSize, constants.CellSize)
		snakePart.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(gridX+pos.X*constants.CellSize), float64(gridY+pos.Y*constants.CellSize))
		screen.DrawImage(snakePart, opts)
	}

	foodPart := ebiten.NewImage(constants.CellSize, constants.CellSize)
	foodPart.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	foodOpts := &ebiten.DrawImageOptions{}
	foodOpts.GeoM.Translate(float64(gridX+g.food.X*constants.CellSize), float64(gridY+g.food.Y*constants.CellSize))
	screen.DrawImage(foodPart, foodOpts)

	for _, pos := range g.obstacles {
		obstaclePart := ebiten.NewImage(constants.CellSize, constants.CellSize)
		obstaclePart.Fill(color.RGBA{R: 128, G: 128, B: 128, A: 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(gridX+pos.X*constants.CellSize), float64(gridY+pos.Y*constants.CellSize))
		screen.DrawImage(obstaclePart, opts)
	}
}
