package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type GameState int

const (
	Menu GameState = iota
	NameInput
	DifficultySelection
	Playing
	GameOver
)

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

var (
	currentSelection Difficulty
	lastMenuUpdate   time.Time
	lastEnterPress   time.Time
	backgroundImage  *ebiten.Image
)

type Score struct {
	Value int
	Name  string
}

type Game struct {
	gridManager    GridManager
	scores         []Score
	score          int
	state          GameState
	updateCount    int
	updateInterval int
	scoreAdded     bool
	playerName     string
	difficulty     Difficulty
}

func init() {
	img, _, err := ebitenutil.NewImageFromFile("assets/menu_background.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage = img
}

func (g *Game) Update() error {
	switch g.state {
	case Menu:
		g.state = NameInput
	case NameInput:
		return g.updateNameInput()
	case DifficultySelection:
		return g.updateDifficultySelection()
	case Playing:
		return g.updatePlaying()
	case GameOver:
		return g.updateGameOver()
	}
	return nil
}

func (g *Game) updateNameInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.playerName) > 0 {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.state = DifficultySelection
			lastEnterPress = time.Now()
		}
	}
	g.handleNameInput()
	return nil
}

func (g *Game) updateDifficultySelection() error {
	if time.Since(lastMenuUpdate) > 200*time.Millisecond {
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			if currentSelection > Easy {
				currentSelection--
			}
			lastMenuUpdate = time.Now()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			if currentSelection < Hard {
				currentSelection++
			}
			lastMenuUpdate = time.Now()
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			if time.Since(lastEnterPress) > 500*time.Millisecond {
				g.state = Playing
				g.difficulty = currentSelection
				g.startGame()
				lastEnterPress = time.Now()
				backgroundPlayer.Rewind()
				backgroundPlayer.Play()
			}
		}
	}
	return nil
}

func (g *Game) updatePlaying() error {
	g.updateCount++
	if g.updateCount >= g.updateInterval {
		err := g.gridManager.Update(g)
		if err != nil {
			g.state = GameOver
			g.scoreAdded = false
		}
		g.updateCount = 0
	}
	return nil
}

func (g *Game) updateGameOver() error {
	if !g.scoreAdded {
		g.AddScore(g.score, g.playerName)
		g.scoreAdded = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.gridManager = NewGrid(cellSize)
		g.score = 0
		g.state = Playing
		g.scoreAdded = false
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.state = Menu
			lastEnterPress = time.Now()
		}
	}
	return nil
}

func (g *Game) startGame() {
	switch g.difficulty {
	case Easy:
		g.updateInterval = 15
	case Medium:
		g.updateInterval = 10
	case Hard:
		g.updateInterval = 5
	}
	g.gridManager = NewGrid(cellSize)
	g.score = 0
	g.state = Playing
	g.updateCount = 0
}

func (g *Game) Draw(screen *ebiten.Image) {
	if backgroundImage != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(screenWidth)/float64(backgroundImage.Bounds().Dx()), float64(screenHeight)/float64(backgroundImage.Bounds().Dy()))
		screen.DrawImage(backgroundImage, opts)
	}

	switch g.state {
	case Menu:
	case NameInput:
		g.renderNameInput(screen)
	case DifficultySelection:
		g.renderDifficultySelection(screen)
	case Playing:
		g.gridManager.Draw(screen)
		ebitenutil.DebugPrintAt(screen, "Score: "+strconv.Itoa(g.score), 10, 10)
	case GameOver:
		g.renderGameOver(screen)
	}
}

func (g *Game) renderNameInput(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	msg := "Veuillez entrer votre nom: " + g.playerName
	text.Draw(screen, msg, fontFace, screenWidth/2-100, screenHeight/2, textColor)
}

func (g *Game) renderDifficultySelection(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Choix de la Difficulte", fontFace, screenWidth/2-50, screenHeight/2-100, textColor)
	text.Draw(screen, "Facile", fontFace, screenWidth/2-50, screenHeight/2-50, textColor)
	text.Draw(screen, "Normal", fontFace, screenWidth/2-50, screenHeight/2, textColor)
	text.Draw(screen, "Difficile", fontFace, screenWidth/2-50, screenHeight/2+50, textColor)
	switch currentSelection {
	case Easy:
		text.Draw(screen, ">", fontFace, screenWidth/2-70, screenHeight/2-50, textColor)
	case Medium:
		text.Draw(screen, ">", fontFace, screenWidth/2-70, screenHeight/2, textColor)
	case Hard:
		text.Draw(screen, ">", fontFace, screenWidth/2-70, screenHeight/2+50, textColor)
	}

	text.Draw(screen, "Meilleurs scores", fontFace, screenWidth/2-50, screenHeight/2+100, textColor)
	for i, score := range g.scores {
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, score.Name, score.Value), fontFace, screenWidth/2-50, screenHeight/2+120+(i*20), textColor)
	}
}

func (g *Game) renderGameOver(screen *ebiten.Image) {
	gameArea := ebiten.NewImage(gridWidth, gridHeight)
	gameArea.Fill(color.Black)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64((screenWidth-gridWidth)/2), float64((screenHeight-gridHeight)/2))
	screen.DrawImage(gameArea, opts)

	textColor := color.RGBA{255, 255, 255, 255}
	fontFace := basicfont.Face7x13

	msg := fmt.Sprintf("Game Over! Score: %d\nAppuyez sur R pour recommencer\nAppuyez sur Entree pour acceder au Menu", g.score)
	text.Draw(screen, msg, fontFace, (screenWidth-gridWidth)/2+20, (screenHeight-gridHeight)/2+20, textColor)

	text.Draw(screen, "Meilleurs scores", fontFace, (screenWidth-gridWidth)/2+20, (screenHeight-gridHeight)/2+60, textColor)
	for i, score := range g.scores {
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, score.Name, score.Value), fontFace, (screenWidth-gridWidth)/2+20, (screenHeight-gridHeight)/2+80+(i*20), textColor)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) handleNameInput() {
	for _, r := range ebiten.InputChars() {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			g.playerName += string(r)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.playerName) > 0 {
		g.playerName = g.playerName[:len(g.playerName)-1]
	}
}
