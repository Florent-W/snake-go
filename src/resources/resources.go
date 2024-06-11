package resources

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Variables globales pour les ressources
var (
	BackgroundImage *ebiten.Image
	HeartImage      *ebiten.Image
	SnakeSprite     *ebiten.Image
	RKeyImage	   *ebiten.Image
	EnterKeyImage  *ebiten.Image
)

func init() {
	var err error

	BackgroundImage, _, err = ebitenutil.NewImageFromFile("assets/menu_background.png")
	if err != nil {
		log.Fatal(err)
	}

	HeartImage, _, err = ebitenutil.NewImageFromFile("assets/coeur.png")
	if err != nil {
		log.Fatal(err)
	}

	SnakeSprite, _, err = ebitenutil.NewImageFromFile("assets/snake-sprite.png")
	if err != nil {
		log.Fatal(err)
	}

	RKeyImage, _, err = ebitenutil.NewImageFromFile("assets/press-r.png")
	if err != nil {
		log.Fatal(err)
	}

	EnterKeyImage, _, err = ebitenutil.NewImageFromFile("assets/press-enter.png")
	if err != nil {
		log.Fatal(err)
	}

}
