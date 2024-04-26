package main

import (
	"fmt"
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
			g.AddScore(g.score, "Test")
			g.scoreAdded = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.gridManager = NewGrid(screenWidth/gridSize, screenHeight/gridSize)
			g.score = 0
			g.state = Playing
			g.scoreAdded = false
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}

	// Mise à jour de l'intervalle en fonction du score
	g.updateInterval = 3 - g.score/10
	if g.updateInterval < 1 {
		g.updateInterval = 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gridManager.Draw(screen)
	scoreText := "Score: " + strconv.Itoa(g.score)
	ebitenutil.DebugPrint(screen, scoreText)

	if g.state == GameOver {
		numberOfScoresText := fmt.Sprintf("Nombre de scores total: %d", len(g.scores))
		ebitenutil.DebugPrintAt(screen, numberOfScoresText, screenWidth/2-100, screenHeight/2-150)

		highScoreText := "Meilleurs Scores:\n"
		for _, s := range g.scores {
			highScoreText += fmt.Sprintf("%s: %d\n", s.Name, s.Value)
		}
		ebitenutil.DebugPrintAt(screen, highScoreText, screenWidth/2-100, screenHeight/2-100)

		msg := "Game Over! Score: " + strconv.Itoa(g.score) + "\nR pour Recommencer, Echap pour quitter"
		x := screenWidth/2 - 200
		y := screenHeight/2 + 100
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
