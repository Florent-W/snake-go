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
	gridManager GridManager
	score       int
	state       GameState
}

func (g *Game) Update() error {
	if g.state == Playing {
		err := g.gridManager.Update(g)
		if err != nil {
			g.state = GameOver
			return nil
		}
	} else if g.state == GameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.gridManager = NewGrid(screenWidth/gridSize, screenHeight/gridSize)
			g.score = 0
			g.state = Playing
		} else if ebiten.IsKeyPressed(ebiten.KeyQ) {
			os.Exit(0)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gridManager.Draw(screen)
	scoreText := "Score: " + strconv.Itoa(g.score)
	ebitenutil.DebugPrint(screen, scoreText)

	// Affichage du Game Over
	if g.state == GameOver {
		msg := "Game Over! Score: " + strconv.Itoa(g.score) + "\nR pour Recommencer, Q pour quitter"
		x := screenWidth/2 - 200
		y := screenHeight / 2
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
