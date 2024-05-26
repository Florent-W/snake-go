package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"snake-go/src/audio"
	"snake-go/src/constants"
	"snake-go/src/game"
)

func loadIcon() image.Image {
	icon, _, err := ebitenutil.NewImageFromFile("assets/icone.png")
	if err != nil {
		log.Fatal(err)
	}

	return icon
}

func main() {
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Snake Go")

	icon := loadIcon()
	ebiten.SetWindowIcon([]image.Image{icon})

	audio.InitAudio()

	g := &game.Game{
		GridManager:    game.NewGrid(constants.CellSize),
		Score:          0,
		UpdateInterval: 3,
		Scores: []game.Score{
			{Value: 30, Name: "Joueur1"},
			{Value: 23, Name: "Joueur2"},
			{Value: 12, Name: "Joueur3"},
		},
		State:      game.Menu,
		PlayerName: "",
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
