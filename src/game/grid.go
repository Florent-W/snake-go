package game

import (
	"fmt"
	"image"
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

	borderColor := color.RGBA{R: 140, G: 130, B: 81, A: 255}
	borderImage := ebiten.NewImage(constants.GridWidth+2*constants.BorderThickness, constants.GridHeight+2*constants.BorderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-constants.BorderThickness), float64(gridY-constants.BorderThickness))
	screen.DrawImage(borderImage, borderOpts)

	backgroundColor := color.RGBA{R: 238, G: 242, B: 121, A: 255}
	gameArea := ebiten.NewImage(constants.GridWidth, constants.GridHeight)
	gameArea.Fill(backgroundColor)
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	for i, pos := range g.snake {
		var segmentType string
		var direction Direction
		var nextDirection Direction

		if i == 0 {
			segmentType = "head"
			direction = g.direction
			if len(g.snake) > 1 {
				nextPos := g.snake[i+1]
				if pos.X < nextPos.X {
					nextDirection = Left
				} else if pos.X > nextPos.X {
					nextDirection = Right
				} else if pos.Y < nextPos.Y {
					nextDirection = Up
				} else {
					nextDirection = Down
				}
			}
		} else if i == len(g.snake)-1 {
			segmentType = "tail"
			prevPos := g.snake[i-1]
			if pos.X < prevPos.X {
				direction = Left
			} else if pos.X > prevPos.X {
				direction = Right
			} else if pos.Y < prevPos.Y {
				direction = Up
			} else {
				direction = Down
			}
		} else {
			segmentType = "body"
			prevPos := g.snake[i-1]
			if pos.X < prevPos.X {
				direction = Left
			} else if pos.X > prevPos.X {
				direction = Right
			} else if pos.Y < prevPos.Y {
				direction = Up
			} else {
				direction = Down
			}

			nextPos := g.snake[i+1]
			if pos.X < nextPos.X {
				nextDirection = Left
			} else if pos.X > nextPos.X {
				nextDirection = Right
			} else if pos.Y < nextPos.Y {
				nextDirection = Up
			} else {
				nextDirection = Down
			}
		}

		snakePart := getSpriteSegment(segmentType, direction, nextDirection)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(constants.CellSize)/64, float64(constants.CellSize)/64)
		opts.GeoM.Translate(float64(gridX+pos.X*constants.CellSize), float64(gridY+pos.Y*constants.CellSize))
		screen.DrawImage(snakePart, opts)
	}

	// Dessiner la pomme
	appleSprite := getAppleSprite()
	appleOpts := &ebiten.DrawImageOptions{}
	appleOpts.GeoM.Scale(float64(constants.CellSize)/64, float64(constants.CellSize)/64)
	appleOpts.GeoM.Translate(float64(gridX+g.food.X*constants.CellSize), float64(gridY+g.food.Y*constants.CellSize))
	screen.DrawImage(appleSprite, appleOpts)

	// Dessiner les obstacles
	obstacleSprite := getObstacleSprite()
	for _, pos := range g.obstacles {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(constants.CellSize)/64, float64(constants.CellSize)/64)
		opts.GeoM.Translate(float64(gridX+pos.X*constants.CellSize), float64(gridY+pos.Y*constants.CellSize))
		screen.DrawImage(obstacleSprite, opts)
	}
}

func getSpriteSegment(segmentType string, direction Direction, nextDirection Direction) *ebiten.Image {
	// Les coordonnées des segments de l'image sprite (colonne, ligne)
	segments := map[string]image.Point{
		"head_up":    {3, 0},
		"head_down":  {4, 1},
		"head_left":  {3, 1},
		"head_right": {4, 0},
		"tail_up":    {4, 3},
		"tail_down":  {3, 2},
		"tail_left":  {4, 2},
		"tail_right": {3, 3},
		"body_v":     {2, 1},
		"body_h":     {1, 0},
		"turn_ur":    {0, 1},
		"turn_ul":    {2, 2},
		"turn_dr":    {0, 0},
		"turn_dl":    {2, 0},
	}

	var segmentKey string

	switch segmentType {
	case "head":
		switch direction {
		case Up:
			segmentKey = "head_up"
		case Down:
			segmentKey = "head_down"
		case Left:
			segmentKey = "head_left"
		case Right:
			segmentKey = "head_right"
		}
	case "tail":
		switch direction {
		case Up:
			segmentKey = "tail_up"
		case Down:
			segmentKey = "tail_down"
		case Left:
			segmentKey = "tail_left"
		case Right:
			segmentKey = "tail_right"
		}
	case "body":
		if direction == Up || direction == Down {
			segmentKey = "body_v"
		} else {
			segmentKey = "body_h"
		}

		switch {
		case direction == Up && nextDirection == Right:
			segmentKey = "turn_dl"
		case direction == Up && nextDirection == Left:
			segmentKey = "turn_dr"
		case direction == Down && nextDirection == Right:
			segmentKey = "turn_ul"
		case direction == Down && nextDirection == Left:
			segmentKey = "turn_ur"
		case direction == Left && nextDirection == Up:
			segmentKey = "turn_dr"
		case direction == Left && nextDirection == Down:
			segmentKey = "turn_ur"
		case direction == Right && nextDirection == Up:
			segmentKey = "turn_dl"
		case direction == Right && nextDirection == Down:
			segmentKey = "turn_ul"
		}
	}

	segmentCoords := segments[segmentKey]
	return snakeSprite.SubImage(image.Rect(segmentCoords.X*64, segmentCoords.Y*64, (segmentCoords.X+1)*64, (segmentCoords.Y+1)*64)).(*ebiten.Image)
}

func getAppleSprite() *ebiten.Image {
	appleCoords := image.Point{X: 0, Y: 3}
	appleSprite := snakeSprite.SubImage(image.Rect(appleCoords.X*64, appleCoords.Y*64, (appleCoords.X+1)*64, (appleCoords.Y+1)*64)).(*ebiten.Image)
	return appleSprite
}

func getObstacleSprite() *ebiten.Image {
	obstacleCoords := image.Point{X: 1, Y: 3}
	obstacleSprite := snakeSprite.SubImage(image.Rect(obstacleCoords.X*64, obstacleCoords.Y*64, (obstacleCoords.X+1)*64, (obstacleCoords.Y+1)*64)).(*ebiten.Image)
	return obstacleSprite
}
