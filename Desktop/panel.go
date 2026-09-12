package desktop

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Dat-one-dev/Aksara/Desktop/panel"
)

var (
	aksaraStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	systemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	timeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("81"))
)

func Panel(width int) string {
	left := aksaraStyle.Render(" AKSARA")

	center := systemStyle.Render(
		fmt.Sprintf(
			"[Battery: %s || RAM: %s]",
			panel.Battery(),
			panel.RAM(),
		),
	)

	right := timeStyle.Render(panel.Time())

	contentWidth := lipgloss.Width(left) +
		lipgloss.Width(center) +
		lipgloss.Width(right)

	if contentWidth >= width {
		return strings.Repeat(" ", width)
	}

	remaining := width - contentWidth
	leftGap := remaining / 2
	rightGap := remaining - leftGap

	return left +
		strings.Repeat(" ", leftGap) +
		center +
		strings.Repeat(" ", rightGap) +
		right
}
