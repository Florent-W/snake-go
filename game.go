package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"sort"
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
	ModeSelection
	DifficultySelection
	Playing
	GameOver
	Credits
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
	heartImage       *ebiten.Image
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
	mode           string
	lives          int
}

func init() {
	img, _, err := ebitenutil.NewImageFromFile("assets/menu_background.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage = img

	heartImg, _, err := ebitenutil.NewImageFromFile("assets/coeur.png")
	if err != nil {
		log.Fatal(err)
	}
	heartImage = heartImg
}

func (g *Game) Update() error {
	switch g.state {
	case Menu:
		return g.updateMenu()
	case NameInput:
		return g.updateNameInput()
	case ModeSelection:
		return g.updateModeSelection()
	case DifficultySelection:
		return g.updateDifficultySelection()
	case Playing:
		return g.updatePlaying()
	case GameOver:
		return g.updateGameOver()
	case Credits:
		return g.updateCredits()
	}
	return nil
}

func (g *Game) updateMenu() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.state = NameInput
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.state = Credits
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		os.Exit(0)
	}
	return nil
}

func (g *Game) updateNameInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.playerName) > 0 {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.state = ModeSelection
			lastEnterPress = time.Now()
		}
	}
	g.handleNameInput()
	return nil
}

func (g *Game) updateModeSelection() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.mode = "Classique"
		g.state = DifficultySelection
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.mode = "Challenge"
		g.state = DifficultySelection
	}
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
			if g.lives > 1 {
				g.lives--
				g.gridManager = NewGridWithObstacles(cellSize, g.difficulty)
			} else {
				g.state = GameOver
				g.scoreAdded = false
			}
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
		g.startGame()
		g.scoreAdded = false
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.state = Menu
			lastEnterPress = time.Now()
		}
	}
	return nil
}

func (g *Game) updateCredits() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.state = Menu
	}
	return nil
}

func (g *Game) startGame() {
	switch g.difficulty {
	case Easy:
		g.updateInterval = 15
		if g.mode == "Challenge" {
			g.lives = 3
		}
	case Medium:
		g.updateInterval = 10
		if g.mode == "Challenge" {
			g.lives = 2
		}
	case Hard:
		g.updateInterval = 5
		if g.mode == "Challenge" {
			g.lives = 1
		}
	}

	if g.mode == "Challenge" {
		g.gridManager = NewGridWithObstacles(cellSize, g.difficulty)
	} else {
		g.gridManager = NewGrid(cellSize)
		g.lives = 1
	}

	g.score = 0
	g.state = Playing
	g.updateCount = 0
}

func (g *Game) AddScore(newScore int, newName string) {
	g.scores = append(g.scores, Score{Value: newScore, Name: newName})

	sort.Slice(g.scores, func(i, j int) bool {
		return g.scores[i].Value > g.scores[j].Value
	})

	if len(g.scores) > 10 {
		g.scores = g.scores[:10]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if backgroundImage != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(screenWidth)/float64(backgroundImage.Bounds().Dx()), float64(screenHeight)/float64(backgroundImage.Bounds().Dy()))
		screen.DrawImage(backgroundImage, opts)
	}

	switch g.state {
	case Menu:
		g.renderMenu(screen)
	case NameInput:
		g.renderNameInput(screen)
	case ModeSelection:
		g.renderModeSelection(screen)
	case DifficultySelection:
		g.renderDifficultySelection(screen)
	case Playing:
		g.gridManager.Draw(screen)
		text.Draw(screen, "Score: "+strconv.Itoa(g.score), basicfont.Face7x13, 10, 20, color.Black)
		g.drawLives(screen)
	case GameOver:
		g.renderGameOver(screen)
	case Credits:
		g.renderCredits(screen)
	}
}

func (g *Game) renderMenu(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Menu Principal", fontFace, screenWidth/2-50, screenHeight/2-100, textColor)
	text.Draw(screen, "1. Commencer le jeu", fontFace, screenWidth/2-50, screenHeight/2-50, textColor)
	text.Draw(screen, "2. Credits", fontFace, screenWidth/2-50, screenHeight/2, textColor)
	text.Draw(screen, "3. Quitter", fontFace, screenWidth/2-50, screenHeight/2+50, textColor)
}

func (g *Game) renderNameInput(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	msg := "Veuillez entrer votre nom: " + g.playerName
	text.Draw(screen, msg, fontFace, screenWidth/2-100, screenHeight/2, textColor)
}

func (g *Game) renderModeSelection(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Choisissez le mode de jeu", fontFace, screenWidth/2-50, screenHeight/2-50, textColor)
	text.Draw(screen, "1. Mode Classique", fontFace, screenWidth/2-50, screenHeight/2, textColor)
	text.Draw(screen, "2. Mode Challenge", fontFace, screenWidth/2-50, screenHeight/2+50, textColor)
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
	gridX := (screenWidth - gridWidth) / 2
	gridY := (screenHeight - gridHeight) / 2

	borderColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	borderImage := ebiten.NewImage(gridWidth+2*borderThickness, gridHeight+2*borderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-borderThickness), float64(gridY-borderThickness))
	screen.DrawImage(borderImage, borderOpts)

	gameArea := ebiten.NewImage(gridWidth, gridHeight)
	gameArea.Fill(color.Black)
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	textColor := color.RGBA{255, 255, 255, 255}
	fontFace := basicfont.Face7x13

	msg := fmt.Sprintf("Game Over! Score: %d\nAppuyez sur R pour recommencer\nAppuyez sur Entree pour acceder au Menu", g.score)
	text.Draw(screen, msg, fontFace, gridX+20, gridY+20, textColor)

	text.Draw(screen, "Meilleurs scores", fontFace, gridX+20, gridY+60, textColor)
	for i, score := range g.scores {
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, score.Name, score.Value), fontFace, gridX+20, gridY+80+(i*20), textColor)
	}
}

func (g *Game) renderCredits(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Credits", fontFace, screenWidth/2-50, screenHeight/2-100, textColor)
	text.Draw(screen, "Developpe par Florent Weltmann, Dantin Durand, William Girard-Reydet", fontFace, screenWidth/2-50, screenHeight/2-50, textColor)
	text.Draw(screen, "Appuyez sur Echap pour revenir au menu", fontFace, screenWidth/2-50, screenHeight/2, textColor)
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

func (g *Game) drawLives(screen *ebiten.Image) {
	if heartImage == nil {
		return
	}

	textColor := color.RGBA{255, 0, 0, 255}
	fontFace := basicfont.Face7x13
	text.Draw(screen, "Vies :", fontFace, 10, 50, textColor)

	for i := 0; i < g.lives; i++ {
		opts := &ebiten.DrawImageOptions{}
		scale := 0.02
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(55+i*int(float64(heartImage.Bounds().Dx())*scale)), 40)
		screen.DrawImage(heartImage, opts)
	}
}
