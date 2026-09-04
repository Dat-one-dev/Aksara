package desktop

import (
	prana "github.com/Dat-one-dev/Prana"
)

func Render() string {
	w, _ := prana.TerminalSize()

	return Panel(w)
}
