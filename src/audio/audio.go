package audio

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"snake-go/src/constants"
)

// Variables globales
var (
	MoveSoundPlayer   *audio.Player
	EatSoundPlayer    *audio.Player
	LoseSoundPlayer   *audio.Player
	BackgroundContext *audio.Context
	BackgroundPlayer  *audio.Player
)

// fonction pour initialiser les différents fichiers audio
func InitAudio() {
	BackgroundContext = audio.NewContext(constants.SampleRate)

	// Chargement des bruitages
	MoveSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/move.mp3")
	EatSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/eating.mp3")
	LoseSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/lose.mp3")

	// Chargement de la musique de fond en boucle
	BackgroundPlayer = loadLoopedAudioPlayer(BackgroundContext, "assets/HeatleyBros - HeatleyBros II - 06 8 Bit Adventure.mp3")

	// Réglage du volume pour chaque son
	if MoveSoundPlayer != nil {
		MoveSoundPlayer.SetVolume(constants.MoveVolume)
	}
	if EatSoundPlayer != nil {
		EatSoundPlayer.SetVolume(constants.EatVolume)
	}
	if LoseSoundPlayer != nil {
		LoseSoundPlayer.SetVolume(constants.LoseVolume)
	}
	if BackgroundPlayer != nil {
		BackgroundPlayer.SetVolume(constants.BackgroundVolume)
		BackgroundPlayer.Play()
	}
}

// charge un fichier audio
//
// ctx: le contexte audio utilisé pour lire les sons
// filename: le chemin du fichier audio à charger
// Retourne un pointeur vers un audio.Player, ou nil en cas d'erreur
func loadAudioPlayer(ctx *audio.Context, filename string) *audio.Player {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Impossible de lire le fichier audio: %v", err)
		return nil
	}

	d, err := mp3.DecodeWithSampleRate(constants.SampleRate, file)
	if err != nil {
		log.Fatal(err)
	}

	p, err := ctx.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}

	return p
}

// charge un fichier audio et le pase en boucle
//
// ctx: le contexte audio utilisé pour lire les sons
// filename: le chemin du fichier audio à charger
// Retourne un pointeur vers un audio.Player en boucle, ou nil en cas d'erreur
func loadLoopedAudioPlayer(ctx *audio.Context, filename string) *audio.Player {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Impossible de lire le fichier audio: %v", err)
		return nil
	}

	d, err := mp3.DecodeWithSampleRate(constants.SampleRate, file)
	if err != nil {
		log.Fatal(err)
	}

	loop := audio.NewInfiniteLoop(d, d.Length())

	p, err := ctx.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
	}

	return p
}
