package main

import (
	"fmt"
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
	screenWidth     = 1280
	screenHeight    = 720
	gridWidth       = 500
	gridHeight      = 500
	cellSize        = 15
	borderThickness = 5
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

	var playerName string
	for {
		fmt.Println("Veuillez entrer votre nom:")
		fmt.Scanln(&playerName)
		if len(playerName) > 0 {
			break
		}
		fmt.Println("Le nom doit contenir au moins un caractère. Réessayez.")
	}
	game := &Game{
		gridManager:    NewGrid(cellSize),
		score:          0,
		updateInterval: 3,
		scores: []Score{
			{Value: 30, Name: "Joueur1"},
			{Value: 23, Name: "Joueur2"},
			{Value: 12, Name: "Joueur3"},
		},
		playerName: playerName,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
