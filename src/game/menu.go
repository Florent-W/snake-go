package game

import (
	"os"
	"time"

	"snake-go/src/audio"
	"snake-go/src/input"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	lastMenuUpdate time.Time
)

// Cette fonction gère les input dans le menu principal
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

// Gère la saisie du nom du joueur
func (g *Game) updateNameInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.PlayerName) > 0 {
		if time.Since(lastEnterPress) > 500*time.Millisecond {
			g.State = ModeSelection
			lastEnterPress = time.Now()
		}
	}
	input.HandleNameInput(&g.PlayerName)
	return nil
}

// Gère la sélection du mode de jeu
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

// Gère la sélection de la difficulté
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
