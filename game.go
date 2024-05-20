package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameState int

const (
	Playing GameState = iota
	GameOver
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
}

func (g *Game) Update() error {
	if g.state == Playing {
		g.updateCount++
		if g.updateCount >= g.updateInterval {
			err := g.gridManager.Update(g)
			if err != nil {
				g.state = GameOver
				g.scoreAdded = false
			}
			g.updateCount = 0
		}
	} else if g.state == GameOver {
		if !g.scoreAdded {
			g.AddScore(g.score, g.playerName)
			g.scoreAdded = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.gridManager = NewGrid(cellSize)
			g.score = 0
			g.state = Playing
			g.scoreAdded = false
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}

	g.updateInterval = 5 - g.score/10
	if g.updateInterval < 1 {
		g.updateInterval = 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gridManager.Draw(screen)
	scoreText := "Score: " + strconv.Itoa(g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 10)

	if g.state == GameOver {

		gridX := (screenWidth - gridWidth) / 2
		gridY := (screenHeight - gridHeight) / 2

		// Image de fond
		gameOverBackground := ebiten.NewImage(gridWidth, gridHeight)
		gameOverBackground.Fill(color.RGBA{51, 79, 255, 255}) // Couleur bleue

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(gridX), float64(gridY))
		screen.DrawImage(gameOverBackground, op)

		textX := gridX + 20
		textY := gridY + 20

		// Affichage du texte de score total
		numberOfScoresText := fmt.Sprintf("Nombre de scores total: %d", len(g.scores))
		ebitenutil.DebugPrintAt(screen, numberOfScoresText, textX, textY)

		textY += 20

		// Affichage des meilleurs scores
		highScoreText := "Meilleurs Scores:\n"
		ebitenutil.DebugPrintAt(screen, highScoreText, textX, textY)

		for _, s := range g.scores {
			textY += 20
			scoreText := fmt.Sprintf("%s: %d\n", s.Name, s.Value)
			ebitenutil.DebugPrintAt(screen, scoreText, textX, textY)
		}
		textY += 20

		msg := "Game Over! Score: " + strconv.Itoa(g.score) + "\nR pour Recommencer, Echap pour quitter"
		ebitenutil.DebugPrintAt(screen, msg, textX, textY)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
