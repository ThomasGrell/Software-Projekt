package text

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/png"
	"log"
	"os"
)

var fontImage *pixel.PictureData
var sprite *pixel.Sprite

// Vor.: -
// Eff.: getRectForChar liefert zu einem Zeichen das zugehörige Rechteck innerhalb des font.png
func getRectForChar(c rune) pixel.Rect {
	x := float64(c) - 32
	if x < 0 || x > 96 {
		return pixel.R(0, 0, 16, 32)
	}
	x = x * 16
	return pixel.R(x, 0, x+16, 32)
}

// Vor.: -
// Eff.: Print schreibt den Text auf ein Canvas-Element und gibt dieses zurück.
func Print(text string) *pixelgl.Canvas {
	canvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(16*len(text)), 32))
	for i, val := range text {
		sprite.Set(fontImage, getRectForChar(val))
		sprite.Draw(canvas, pixel.IM.Moved(pixel.V(8, 16)).Moved(pixel.V(float64(i*16), 0)))
	}
	return canvas
}

func init() {
	file, err := os.Open("./text/font.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Quelle: https://github.com/Lallassu/gizmo/blob/master/assets/mixed/font.png
	fontImage = pixel.PictureDataFromImage(img)
	sprite = pixel.NewSprite(fontImage, pixel.R(0, 0, 16, 32))
}
