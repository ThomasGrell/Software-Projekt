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

func getRectForChar(c rune) pixel.Rect {
	x := float64(c) - 32
	if x < 0 || x > 96 {
		return pixel.R(0, 0, 16, 32)
	}
	x = x * 16
	return pixel.R(x, 0, x+16, 32)
}

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
	fontImage = pixel.PictureDataFromImage(img)
	sprite = pixel.NewSprite(fontImage, pixel.R(0, 0, 16, 32))
}
