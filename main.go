package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 44100

var (
	moveSoundPlayer *audio.Player
	eatSoundPlayer  *audio.Player
	loseSoundPlayer *audio.Player
)

const (
    screenWidth  = 1280
    screenHeight = 720
    gridWidth    = 500  // Taille en pixels de la largeur de la grille
    gridHeight   = 500  // Taille en pixels de la hauteur de la grille
    cellSize     = 15   // Taille en pixels de chaque cellule de la grille
    borderThickness = 5 // Épaisseur de la bordure autour de la grille
)

func initAudio() {
	var audioContext = audio.NewContext(sampleRate)

	moveSoundPlayer = loadAudioPlayer(audioContext, "./assets/move.mp3")
	eatSoundPlayer = loadAudioPlayer(audioContext, "./assets/eating.mp3")
	loseSoundPlayer = loadAudioPlayer(audioContext, "./assets/lose.mp3")
}

func loadAudioPlayer(ctx *audio.Context, filename string) *audio.Player {
	file, err := os.Open(filename)
    if err != nil {
        log.Printf("failed to open audio file: %v", err)
        return nil
    }

	d, err := mp3.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		log.Fatal(err)
	}

	p, err := ctx.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func main() {
	initAudio()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Go")

	game := &Game{
		gridManager:    NewGrid(cellSize), // Utilisez cellSize ici pour la création de la grille
		score:          0,
		updateInterval: 5,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
