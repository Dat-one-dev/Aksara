package desktop

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/Dat-one-dev/Aksara/Desktop/wallpaper"
)

type Wallpaper struct {
	img image.Image
}

func Load(path string) (*Wallpaper, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &Wallpaper{img: img}, nil
}

func (w *Wallpaper) Render(width, height int) string {
	if w.img == nil {
		return ""
	}
	return wallpaper.ToASCII(w.img, width, height)
}
