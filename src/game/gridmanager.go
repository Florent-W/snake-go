package game

import "github.com/hajimehoshi/ebiten/v2"

// GridManager définit les méthodes nécessaires pour gérer et dessiner une grille.
type GridManager interface {
	Update(game *Game) error
	Draw(screen *ebiten.Image)
}
