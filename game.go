package main

import (
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

type Game struct {
	gridManager    GridManager
	score          int
	state          GameState
	updateCount    int
	updateInterval int
}

func (g *Game) Update() error {
	if g.state == Playing {
		g.updateCount++
		if g.updateCount >= g.updateInterval {
			err := g.gridManager.Update(g)
			if err != nil {
				g.state = GameOver
			}
			g.updateCount = 0
		}
	} else if g.state == GameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.gridManager = NewGrid(cellSize)
			g.score = 0
			g.state = Playing
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}

	// Mise à jour de l'intervalle en fonction du score
	g.updateInterval = 5 - g.score/10
	if g.updateInterval < 1 {
		g.updateInterval = 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gridManager.Draw(screen)
	scoreText := "Score: " + strconv.Itoa(g.score)
	ebitenutil.DebugPrint(screen, scoreText)

	// Affichage du Game Over
	if g.state == GameOver {
		msg := "Game Over! Score: " + strconv.Itoa(g.score) + "\nR pour Recommencer, Echap pour quitter"
		x := screenWidth/2 - 200
		y := screenHeight / 2
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
