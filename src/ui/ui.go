package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"snake-go/src/constants"
)

// Dessine le menu principal
func RenderMenu(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Menu Principal", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-100, textColor)
	text.Draw(screen, "1. Commencer le jeu", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-50, textColor)
	text.Draw(screen, "2. Credits", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2, textColor)
	text.Draw(screen, "3. Quitter", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2+50, textColor)
}

// Dessine l'écran de saisie du nom du joueur
//
// playerName: le nom actuellement saisi par le joueur
func RenderNameInput(screen *ebiten.Image, playerName string) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	msg := "Veuillez entrer votre nom: " + playerName
	text.Draw(screen, msg, fontFace, constants.ScreenWidth/2-100, constants.ScreenHeight/2, textColor)
}

// Dessine l'écran de sélection du mode de jeu
func RenderModeSelection(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Choisissez le mode de jeu", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-50, textColor)
	text.Draw(screen, "1. Mode Classique", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2, textColor)
	text.Draw(screen, "2. Mode Challenge", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2+50, textColor)
}

// Dessine l'écran de sélection de la difficulté
//
// currentSelection: la difficulté actuellement sélectionnée
// scores: la liste des meilleurs scores à afficher
func RenderDifficultySelection(screen *ebiten.Image, currentSelection int, scores []Score) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Choix de la Difficulte", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-100, textColor)
	text.Draw(screen, "Facile", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-50, textColor)
	text.Draw(screen, "Normal", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2, textColor)
	text.Draw(screen, "Difficile", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2+50, textColor)
	switch currentSelection {
	case 0:
		text.Draw(screen, ">", fontFace, constants.ScreenWidth/2-70, constants.ScreenHeight/2-50, textColor)
	case 1:
		text.Draw(screen, ">", fontFace, constants.ScreenWidth/2-70, constants.ScreenHeight/2, textColor)
	case 2:
		text.Draw(screen, ">", fontFace, constants.ScreenWidth/2-70, constants.ScreenHeight/2+50, textColor)
	}

	text.Draw(screen, "Meilleurs scores", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2+100, textColor)
	for i, score := range scores {
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, score.Name, score.Value), fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2+120+(i*20), textColor)
	}
}

// Dessine l'écran de fin de partie avec le score final et les meilleurs scores
//
// score: le score final
// scores: la liste des meilleurs scores
func RenderGameOver(screen *ebiten.Image, score int, scores []Score) {
	gridX := (constants.ScreenWidth - constants.GridWidth) / 2
	gridY := (constants.ScreenHeight - constants.GridHeight) / 2

	borderColor := color.RGBA{R: 113, G: 105, B: 66, A: 255}
	borderImage := ebiten.NewImage(constants.GridWidth+2*constants.BorderThickness, constants.GridHeight+2*constants.BorderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-constants.BorderThickness), float64(gridY-constants.BorderThickness))
	screen.DrawImage(borderImage, borderOpts)

	backgroundColor := color.RGBA{R: 140, G: 130, B: 81, A: 255}
	gameArea := ebiten.NewImage(constants.GridWidth, constants.GridHeight)
	gameArea.Fill(backgroundColor)
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	textColor := color.RGBA{255, 255, 255, 255}
	fontFace := basicfont.Face7x13

	msg := fmt.Sprintf("Game Over! Score: %d\nAppuyez sur R pour recommencer\nAppuyez sur Entree pour acceder au Menu", score)
	text.Draw(screen, msg, fontFace, gridX+20, gridY+20, textColor)

	text.Draw(screen, "Meilleurs scores", fontFace, gridX+20, gridY+60, textColor)
	for i, score := range scores {
		text.Draw(screen, fmt.Sprintf("%d. %s: %d", i+1, score.Name, score.Value), fontFace, gridX+20, gridY+80+(i*20), textColor)
	}
}

// Dessine l'écran des crédits
func RenderCredits(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Credits", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-100, textColor)
	text.Draw(screen, "Developpe par Florent Weltmann, Dantin Durand, William Girard-Reydet", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-50, textColor)
	text.Draw(screen, "Appuyez sur Echap pour revenir au menu", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2, textColor)
}

type Score struct {
	Value int
	Name  string
}
