package main

import (
	"./titlebar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	wincfg := pixelgl.WindowConfig{
		Title:  "Bombermen 2021",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))

	bar := titlebar.New(16 * 16)
	bar.SetPlayers(4)
	bar.SetLife(1, 2, 3, 4)
	bar.SetPoints(12345)
	go bar.Manager()
	bar.Draw(win)
	bar.SetSeconds(3 * 60)
	bar.StartCountdown()
	win.Update()
	for {
		bar.Draw(win)
		win.Update()
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}
	}
}

func main() {
	// Hier darf nichts weiter stehen als die folgende Anweisung:
	pixelgl.Run(run)
}
