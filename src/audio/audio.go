package audio

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"snake-go/src/constants"
)

var (
	MoveSoundPlayer   *audio.Player
	EatSoundPlayer    *audio.Player
	LoseSoundPlayer   *audio.Player
	BackgroundContext *audio.Context
	BackgroundPlayer  *audio.Player
)

func InitAudio() {
	BackgroundContext = audio.NewContext(constants.SampleRate)

	MoveSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/move.mp3")
	EatSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/eating.mp3")
	LoseSoundPlayer = loadAudioPlayer(BackgroundContext, "assets/lose.mp3")
	BackgroundPlayer = loadAudioPlayer(BackgroundContext, "assets/HeatleyBros - HeatleyBros II - 06 8 Bit Adventure.mp3")

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
	}
}

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
