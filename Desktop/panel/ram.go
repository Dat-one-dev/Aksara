package panel

import (
	"os"
	"strconv"
	"strings"
)

func RAM() string {
	data, err := os.ReadFile("/proc/meminfo")
	var total, avail uint64
	if err != nil {
		return "N/A"
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			total = value
		case "MemAvailable:":
			avail = value
		}
	}

	if total == 0 {
		return "N/A"
	}

	used := total - avail
	percentage := (used * 100) / total

	return strconv.FormatUint(percentage, 10) + "%"
}
