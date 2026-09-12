package wallpaper

import (
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/nfnt/resize"
)

var asciiChars = []rune(" .':-=+x*#%@")

func ToASCII(img image.Image, w, h int) string {
	resizedImg := resize.Resize(uint(w), uint(h), img, resize.Bilinear)
	bounds := resizedImg.Bounds()

	var builder strings.Builder
	builder.Grow(w * h * 20)

	rampLen := float64(len(asciiChars) - 1)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := resizedImg.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			gray := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)

			if gray < 15.0 {
				builder.WriteRune(' ')
				continue
			}

			norm := math.Pow(gray/255.0, 0.9)
			charIdx := int(norm * rampLen)

			if charIdx < 0 {
				charIdx = 0
			} else if charIdx > int(rampLen) {
				charIdx = int(rampLen)
			}

			fmt.Fprintf(&builder, "\x1b[38;2;%d;%d;%dm%c\x1b[0m", r8, g8, b8, asciiChars[charIdx])
		}
		if y < bounds.Max.Y-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}
