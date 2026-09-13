package desktop

import (
	"strings"

	"charm.land/lipgloss/v2"
	prana "github.com/Dat-one-dev/Prana"
)

var (
	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
	searchFocused = lipgloss.Color("205")
	searchBlurred = lipgloss.Color("245")
	continueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))
)

func SearchBar(w int, content string, focused bool) string {
	c := searchBlurred
	if focused {
		c = searchFocused
	}
	return lipgloss.NewStyle().
		Width(w - 2).
		Height(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c).
		Render(content)
}

func Render(w, h int, wp *Wallpaper) string {
	quitTxt := mutedStyle.Render("[q to quit]")
	panel := prana.Box(prana.BoxConfig{
		Content: []string{
			Panel(w - 4),
		},
		Width:         w,
		Height:        3,
		Border:        lipgloss.RoundedBorder(),
		Align:         lipgloss.Center,
		AlignVertical: lipgloss.Center,
	})

	var wallpaperContent string
	if wp != nil {
		wallpaperContent = wp.Render(w-4, h-5)
	} else {
		wallpaperContent = "NO WALLPAPER LOADED"
	}

	desktopBox := prana.Box(prana.BoxConfig{
		Content: []string{
			wallpaperContent,
		},
		Footer:        quitTxt,
		Width:         w,
		Height:        h - 3,
		Border:        lipgloss.RoundedBorder(),
		Align:         lipgloss.Center,
		AlignVertical: lipgloss.Center,
	})

	return desktopBox + "\n" + panel
}

func Output(w, h int, lines []string, off int) string {
	page := h - 5
	if page < 1 {
		page = 1
	}
	maxOff := len(lines) - page
	if maxOff < 0 {
		maxOff = 0
	}
	if off < 0 {
		off = 0
	}
	if off > maxOff {
		off = maxOff
	}
	end := off + page
	if end > len(lines) {
		end = len(lines)
	}
	body := strings.Join(lines[off:end], "\n")

	quitTxt := continueStyle.Render("Press ESC to continue - Q to quit")

	panel := prana.Box(prana.BoxConfig{
		Content: []string{
			Panel(w - 4),
		},
		Width:         w,
		Height:        3,
		Border:        lipgloss.RoundedBorder(),
		Align:         lipgloss.Center,
		AlignVertical: lipgloss.Center,
	})

	desktopBox := prana.Box(prana.BoxConfig{
		Content: []string{
			body,
		},
		Footer:        quitTxt,
		Width:         w,
		Height:        h - 3,
		Border:        lipgloss.RoundedBorder(),
		Align:         lipgloss.Left,
		AlignVertical: lipgloss.Top,
	})

	return desktopBox + "\n" + panel
}
