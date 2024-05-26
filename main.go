package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const sampleRate = 44100

var (
	moveSoundPlayer   *audio.Player
	eatSoundPlayer    *audio.Player
	loseSoundPlayer   *audio.Player
	backgroundContext *audio.Context
	backgroundPlayer  *audio.Player
)

const (
	screenWidth     = 1280
	screenHeight    = 720
	gridWidth       = 500
	gridHeight      = 500
	cellSize        = 15
	borderThickness = 5
)

func initAudio() {
	backgroundContext = audio.NewContext(sampleRate)

	moveSoundPlayer = loadAudioPlayer(backgroundContext, "./assets/move.mp3")
	eatSoundPlayer = loadAudioPlayer(backgroundContext, "./assets/eating.mp3")
	loseSoundPlayer = loadAudioPlayer(backgroundContext, "./assets/lose.mp3")
	backgroundPlayer = loadAudioPlayer(backgroundContext, "./assets/HeatleyBros - HeatleyBros II - 06 8 Bit Adventure.mp3")
}

func loadAudioPlayer(ctx *audio.Context, filename string) *audio.Player {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Impossible de lire le fichier audio: %v", err)
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

func loadIcon() image.Image {
	icon, _, err := ebitenutil.NewImageFromFile("assets/icone.png")
	if err != nil {
		log.Fatal(err)
	}

	return icon
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Go")

	icon := loadIcon()
	ebiten.SetWindowIcon([]image.Image{icon})

	initAudio()

	game := &Game{
		gridManager:    NewGrid(cellSize),
		score:          0,
		updateInterval: 3,
		scores: []Score{
			{Value: 30, Name: "Joueur1"},
			{Value: 23, Name: "Joueur2"},
			{Value: 12, Name: "Joueur3"},
		},
		state:      Menu,
		playerName: "",
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
