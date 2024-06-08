package game

import (
	"image/color"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"snake-go/src/audio"
	"snake-go/src/constants"
	utils "snake-go/src/input"
	"snake-go/src/ui"
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
	Facile Difficulty = iota
	Normal
	Difficile
)

var (
	currentSelection Difficulty
	lastMenuUpdate   time.Time
	lastEnterPress   time.Time
	backgroundImage  *ebiten.Image
	heartImage       *ebiten.Image
	snakeSprite      *ebiten.Image
	backgroundArea   *ebiten.Image
)

type Score struct {
	Value int
	Name  string
}

type Game struct {
	GridManager    GridManager
	Scores         []Score
	Score          int
	State          GameState
	UpdateCount    int
	UpdateInterval int
	ScoreAdded     bool
	PlayerName     string
	Difficulty     Difficulty
	Mode           string
	Lives          int
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

	snakeImg, _, err := ebitenutil.NewImageFromFile("assets/snake-sprite.png")
	if err != nil {
		log.Fatal(err)
	}
	snakeSprite = snakeImg

	areaImg, _, err := ebitenutil.NewImageFromFile("assets/bg-area.png")
	if err != nil {
		log.Fatal(err)
	}
	backgroundArea = areaImg
}

func (g *Game) Update() error {
	switch g.State {
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
		g.State = NameInput
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.State = Credits
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		os.Exit(0)
	}
	return nil
}

func (g *Game) updateNameInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.PlayerName) > 0 {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.State = ModeSelection
			lastEnterPress = time.Now()
		}
	}
	utils.HandleNameInput(&g.PlayerName)
	return nil
}

func (g *Game) updateModeSelection() error {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.Mode = "Classique"
		g.State = DifficultySelection
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.Mode = "Challenge"
		g.State = DifficultySelection
	}
	return nil
}

func (g *Game) updateDifficultySelection() error {
	if time.Since(lastMenuUpdate) > 200*time.Millisecond {
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			if currentSelection > Facile {
				currentSelection--
			}
			lastMenuUpdate = time.Now()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			if currentSelection < Difficile {
				currentSelection++
			}
			lastMenuUpdate = time.Now()
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			if time.Since(lastEnterPress) > 500*time.Millisecond {
				g.State = Playing
				g.Difficulty = currentSelection
				g.startGame()
				lastEnterPress = time.Now()
				audio.BackgroundPlayer.Rewind()
				audio.BackgroundPlayer.Play()
			}
		}
	}
	return nil
}

func (g *Game) updatePlaying() error {
	g.UpdateCount++
	if g.UpdateCount >= g.UpdateInterval {
		err := g.GridManager.Update(g)
		if err != nil {
			if g.Lives > 1 {
				g.Lives--
				g.GridManager = NewGridWithObstacles(constants.CellSize, g.Difficulty)
			} else {
				g.State = GameOver
				g.ScoreAdded = false
				audio.LoseSoundPlayer.Rewind()
				audio.LoseSoundPlayer.Play()
			}
		}
		g.UpdateCount = 0
	}
	return nil
}

func (g *Game) updateGameOver() error {
	if !g.ScoreAdded {
		g.AddScore(g.Score, g.PlayerName)
		g.ScoreAdded = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.startGame()
		g.ScoreAdded = false
		audio.BackgroundPlayer.Rewind()
		audio.BackgroundPlayer.Play()
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.State = Menu
			lastEnterPress = time.Now()
			audio.BackgroundPlayer.Rewind()
			audio.BackgroundPlayer.Play()
		}
	}
	return nil
}

func (g *Game) updateCredits() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.State = Menu
		audio.BackgroundPlayer.Rewind()
		audio.BackgroundPlayer.Play()
	}
	return nil
}

func (g *Game) startGame() {
	switch g.Difficulty {
	case Facile:
		g.UpdateInterval = 15
		if g.Mode == "Challenge" {
			g.Lives = 3
		}
	case Normal:
		g.UpdateInterval = 10
		if g.Mode == "Challenge" {
			g.Lives = 2
		}
	case Difficile:
		g.UpdateInterval = 5
		if g.Mode == "Challenge" {
			g.Lives = 1
		}
	}

	if g.Mode == "Challenge" {
		g.GridManager = NewGridWithObstacles(constants.CellSize, g.Difficulty)
	} else {
		g.GridManager = NewGrid(constants.CellSize)
		g.Lives = 1
	}

	g.Score = 0
	g.State = Playing
	g.UpdateCount = 0
	audio.BackgroundPlayer.Rewind()
	audio.BackgroundPlayer.Play()
}

func (g *Game) AddScore(newScore int, newName string) {
	g.Scores = append(g.Scores, Score{Value: newScore, Name: newName})

	sort.Slice(g.Scores, func(i, j int) bool {
		return g.Scores[i].Value > g.Scores[j].Value
	})

	if len(g.Scores) > 10 {
		g.Scores = g.Scores[:10]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if backgroundImage != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(constants.ScreenWidth)/float64(backgroundImage.Bounds().Dx()), float64(constants.ScreenHeight)/float64(backgroundImage.Bounds().Dy()))
		screen.DrawImage(backgroundImage, opts)
	}

	switch g.State {
	case Menu:
		ui.RenderMenu(screen)
	case NameInput:
		ui.RenderNameInput(screen, g.PlayerName)
	case ModeSelection:
		ui.RenderModeSelection(screen)
	case DifficultySelection:
		ui.RenderDifficultySelection(screen, int(currentSelection), convertScores(g.Scores))
	case Playing:
		g.GridManager.Draw(screen)
		text.Draw(screen, "Score: "+strconv.Itoa(g.Score), basicfont.Face7x13, 10, 20, color.Black)
		g.drawLives(screen)
	case GameOver:
		ui.RenderGameOver(screen, g.Score, convertScores(g.Scores))
	case Credits:
		ui.RenderCredits(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func (g *Game) drawLives(screen *ebiten.Image) {
	if heartImage == nil {
		return
	}

	textColor := color.RGBA{255, 0, 0, 255}
	fontFace := basicfont.Face7x13
	text.Draw(screen, "Vies :", fontFace, 10, 50, textColor)

	for i := 0; i < g.Lives; i++ {
		opts := &ebiten.DrawImageOptions{}
		scale := 0.02
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(55+i*int(float64(heartImage.Bounds().Dx())*scale)), 40)
		screen.DrawImage(heartImage, opts)
	}
}

func convertScores(scores []Score) []ui.Score {
	converted := make([]ui.Score, len(scores))
	for i, score := range scores {
		converted[i] = ui.Score{
			Value: score.Value,
			Name:  score.Name,
		}
	}
	return converted
}
