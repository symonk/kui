package terminal

import (
	"os"

	"golang.org/x/term"
)

// size returns the height and width of the current
// terminal.
func Size() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	return width, height
}
