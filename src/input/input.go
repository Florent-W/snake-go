package input

import (
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
)

// gère la saisie du nom du joueur
func HandleNameInput(playerName *string) {
	for _, r := range ebiten.InputChars() {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			*playerName += string(r)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(*playerName) > 0 {
		*playerName = (*playerName)[:len(*playerName)-1]
	}
}
