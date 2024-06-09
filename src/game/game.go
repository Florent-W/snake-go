package game

import (
	"image/color"
	"sort"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"snake-go/src/audio"
	"snake-go/src/constants"
	"snake-go/src/resources"
	"snake-go/src/ui"
)

// Etats dans le jeu, on peut être dans le menu, en train de jouer, en train de choisir le mode de jeu, etc.
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

// Déclaration des niveaux de difficulté
type Difficulty int

const (
	Facile Difficulty = iota
	Normal
	Difficile
)

// Variables globales
var (
	currentSelection Difficulty
	lastEnterPress   time.Time
)

type Score struct {
	Value int
	Name  string
}

// Structure représentant l'état du jeu
type Game struct {
	GridManager       GridManager
	Scores            []Score
	Score             int
	State             GameState
	UpdateCount       int
	UpdateInterval    int
	ScoreAdded        bool
	PlayerName        string
	Difficulty        Difficulty
	Mode              string
	Lives             int
	LastSpeedIncrease int
}

// Fonction principale de mise à jour du jeu, appel les méthodes selon l'état du jeu
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

// Mise à jour de l'état de jeu pendant la partie
func (g *Game) updatePlaying() error {
	g.UpdateCount++
	if g.UpdateCount >= g.UpdateInterval {
		if g.Score > 0 && g.Score%5 == 0 && g.Score != g.LastSpeedIncrease { // Ma vitesse sera augmentée à chaque fois que 5 pommes sont mangées
			g.UpdateInterval = max(3, g.UpdateInterval-1) // Réduire l'intervalle de mise à jour mais pas en dessous de 3
			g.LastSpeedIncrease = g.Score                 // Permet d'enregistrer le score où la vitesse a été augmentée comme ça on ne l'augmente qu'une seule fois par 5 points
		}
		err := g.GridManager.Update(g)
		if err != nil {
			if g.Lives > 1 { // Si on a plus d'une vie (dans le mode challenge), on perd une vie et on recommence tout en gardant le score
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

// Mise à jour de l'état de jeu lors du game over
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
			g.State = Menu // Retourau menu
			lastEnterPress = time.Now()
			audio.BackgroundPlayer.Rewind()
			audio.BackgroundPlayer.Play()
		}
	}
	return nil
}

// Mise à jour de l'état de jeu lors des crédits, ça permet de revenir au menu principal
func (g *Game) updateCredits() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.State = Menu
		audio.BackgroundPlayer.Rewind()
		audio.BackgroundPlayer.Play()
	}
	return nil
}

// Initialisation des paramètres de jeu selon la difficulté choisie
// Dans le mode classique, on a une seule vie qu'importe la difficulté
// Dans le mode challenge, on a 3 vies en facile, 2 en normal et 1 en difficile, la vitesse de départ change et il y a des obstacles en mode challenge qui diffèrennt selon la difficulté ainsi que plus de vies
func (g *Game) startGame() {
	switch g.Difficulty {
	case Facile:
		g.UpdateInterval = 15 // Le serpent bouge lentement au début
		g.Lives = 1
		if g.Mode == "Challenge" {
			g.Lives = 3
		}
	case Normal:
		g.UpdateInterval = 10
		g.Lives = 1
		if g.Mode == "Challenge" {
			g.Lives = 2 // Mode challenge : 2 vies
		}
	case Difficile:
		g.UpdateInterval = 5
		g.Lives = 1
	}

	// Initialisation de la grille en fonction du mode de jeu pour savoir si il y a des obstacles ou non
	if g.Mode == "Challenge" {
		g.GridManager = NewGridWithObstacles(constants.CellSize, g.Difficulty)
	} else {
		g.GridManager = NewGrid(constants.CellSize)
	}

	// Réinitialisation des autres paramètres de jeu
	g.Score = 0
	g.State = Playing
	g.UpdateCount = 0
	g.LastSpeedIncrease = 0
	audio.BackgroundPlayer.Rewind()
	audio.BackgroundPlayer.Play()
}

// Ajout d'un nouveau score à la liste des scores
func (g *Game) AddScore(newScore int, newName string) {
	g.Scores = append(g.Scores, Score{Value: newScore, Name: newName})

	sort.Slice(g.Scores, func(i, j int) bool {
		return g.Scores[i].Value > g.Scores[j].Value
	})

	if len(g.Scores) > 10 {
		g.Scores = g.Scores[:10]
	}
}

// Dessin des éléments à l'écran selon l'état du jeu
func (g *Game) Draw(screen *ebiten.Image) {
	if resources.BackgroundImage != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(constants.ScreenWidth)/float64(resources.BackgroundImage.Bounds().Dx()), float64(constants.ScreenHeight)/float64(resources.BackgroundImage.Bounds().Dy()))
		screen.DrawImage(resources.BackgroundImage, opts)
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

// Dessin des vies restantes à l'écran
func (g *Game) drawLives(screen *ebiten.Image) {
	if resources.HeartImage == nil {
		return
	}

	textColor := color.RGBA{255, 0, 0, 255}
	fontFace := basicfont.Face7x13
	text.Draw(screen, "Vies :", fontFace, 10, 50, textColor)

	for i := 0; i < g.Lives; i++ {
		opts := &ebiten.DrawImageOptions{}
		scale := 0.02
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(55+i*int(float64(resources.HeartImage.Bounds().Dx())*scale)), 40)
		screen.DrawImage(resources.HeartImage, opts)
	}
}

// Conversion des scores pour affichage
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
