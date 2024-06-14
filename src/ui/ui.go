package ui

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"snake-go/src/constants"
	"snake-go/src/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

func loadFont(size float64) font.Face {
    fontBytes, err := ioutil.ReadFile("assets/upheavtt.ttf")
    if err != nil {
        log.Fatalf("failed to read font file: %v", err)
    }

    tt, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatalf("failed to parse font: %v", err)
    }

    const dpi = 72
    face, err := opentype.NewFace(tt, &opentype.FaceOptions{
        Size:    size,
        DPI:     dpi,
        Hinting: font.HintingFull,
    })
    if err != nil {
        log.Fatalf("failed to create font face: %v", err)
    }

    return face
}


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

	borderColor := color.RGBA{223, 173, 59, 255}
	borderImage := ebiten.NewImage(constants.GridWidth+2*constants.BorderThickness, constants.GridHeight+2*constants.BorderThickness)
	borderImage.Fill(borderColor)
	borderOpts := &ebiten.DrawImageOptions{}
	borderOpts.GeoM.Translate(float64(gridX-constants.BorderThickness), float64(gridY-constants.BorderThickness))
	screen.DrawImage(borderImage, borderOpts)

	backgroundColor := color.RGBA{R: 26, G: 26, B: 26, A: 255}
	gameArea := ebiten.NewImage(constants.GridWidth, constants.GridHeight)
	gameArea.Fill(backgroundColor)
	gameAreaOpts := &ebiten.DrawImageOptions{}
	gameAreaOpts.GeoM.Translate(float64(gridX), float64(gridY))
	screen.DrawImage(gameArea, gameAreaOpts)

	// TITLE TEXT 

    titleText := "Game Over !"
    // Obtenez les dimensions du texte 
	titleColor := color.RGBA{223, 173, 59, 255}
    fontTitle := loadFont(70) // Charger la police avec une taille spécifique
    bounds := text.BoundString(fontTitle, titleText)
    textWidth := bounds.Dx()
    textHeight := bounds.Dy()

    // Calculer les positions pour centrer le texte
    x := gridX + (constants.GridWidth-textWidth)/2
    y := gridY + 20

    // Dessiner le texte centré
    text.Draw(screen, titleText, fontTitle, x, y+textHeight, titleColor)


	// SCORE TEXT

	scoreText := fmt.Sprintf("Score: %d", score)
	
	// Obtenez les dimensions du texte
	scoreColor := color.RGBA{255, 255, 255, 255}
	fontScore := loadFont(30) // Charger la police avec une taille spécifique
	bounds = text.BoundString(fontScore, scoreText)
	textWidth = bounds.Dx()
	textHeight = bounds.Dy()

	// Calculer les positions pour centrer le texte
	x = gridX + (constants.GridWidth-textWidth)/2
	y = gridY + 100

	text.Draw(screen, scoreText, fontScore, x, y+textHeight, scoreColor)

    // COLUMN TITLES
    columnTitleFont := loadFont(25)
    playerTitle := "JOUEUR"
    scoreTitle := "SCORE"
    playerColumnWidth := constants.GridWidth * 0.5
    scoreColumnWidth := constants.GridWidth * 0.2
    columnGap := (constants.GridWidth - int(playerColumnWidth) - int(scoreColumnWidth)) / 2
    playerTitleX := gridX + columnGap
    scoreTitleX := playerTitleX + int(playerColumnWidth) 

    text.Draw(screen, playerTitle, columnTitleFont, playerTitleX, gridY+160, color.RGBA{173, 216, 230, 255})
    text.Draw(screen, scoreTitle, columnTitleFont, scoreTitleX, gridY+160, color.RGBA{173, 216, 230, 255})

    // SCORE LIST TEXT
    playerFont := loadFont(20)
    scoreFont := loadFont(20)
    column1X := playerTitleX
    column2X := scoreTitleX
    columnY := gridY + 200

    for i, score := range scores {
        text.Draw(screen, fmt.Sprintf("%s", score.Name), playerFont, column1X, columnY+(i*30), color.RGBA{255, 255, 204, 255})
        text.Draw(screen, fmt.Sprintf("%d", score.Value), scoreFont, column2X, columnY+(i*30), color.RGBA{255, 255, 204, 255})
    }

    // RELAUNCH TEXT
    relaunchText1 := "RECOMMENCER"
    relaunchText2 := "MENU"
    relaunchFont := loadFont(20) // Charger la police avec une taille spécifique

    // Obtenez les dimensions du texte
    bounds1 := text.BoundString(relaunchFont, relaunchText1)

    textHeight1 := bounds1.Dy()
    bounds2 := text.BoundString(relaunchFont, relaunchText2)

    textHeight2 := bounds2.Dy()

    // Obtenez les dimensions des images
    rKeyImageWidth, rKeyImageHeight := resources.RKeyImage.Size()
	enterKeyImageWidth, _ := resources.EnterKeyImage.Size()

    // Positionner les images et les textes en bas à gauche de la grille
    x1 := gridX + 20
    y1 := gridY + constants.GridHeight - rKeyImageHeight - textHeight1 - 50
    x2 := gridX + 20
    y2 := y1 + rKeyImageHeight + textHeight1 + 10

    // Dessiner les images et les textes
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(x1), float64(y1))
    screen.DrawImage(resources.RKeyImage, op)
    text.Draw(screen, relaunchText1, relaunchFont, x1+rKeyImageWidth+10, y1+textHeight1+5, color.RGBA{255, 255, 255, 255})

    op = &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(x2), float64(y2))
    screen.DrawImage(resources.EnterKeyImage, op)
    text.Draw(screen, relaunchText2, relaunchFont, x2+enterKeyImageWidth+10, y2+textHeight2+5, color.RGBA{255, 255, 255, 255})
}

// Dessine l'écran des crédits
func RenderCredits(screen *ebiten.Image) {
	textColor := color.RGBA{0, 0, 0, 255}
	fontFace := basicfont.Face7x13

	text.Draw(screen, "Credits", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-100, textColor)
	text.Draw(screen, "Developpe par Florent Weltmann, Dantin Durand", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2-50, textColor)
	text.Draw(screen, "Appuyez sur Echap pour revenir au menu", fontFace, constants.ScreenWidth/2-50, constants.ScreenHeight/2, textColor)
}

type Score struct {
	Value int
	Name  string
}
