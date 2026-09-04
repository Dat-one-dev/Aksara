package desktop

import (
	"strings"

	"github.com/Dat-one-dev/Aksara/Desktop/panel"
)

func Panel(width int) string {
	left := " AKSARA"
	center := panel.Battery()
	right := panel.Time()

	contentWidth := len(left) + len(center) + len(right)

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
