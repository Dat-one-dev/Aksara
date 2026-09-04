package panel

import (
	"os"
	"strings"
)

func Battery() string {
	data, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return "N/A"
	}

	return strings.TrimSpace(string(data)) + "%"
}
