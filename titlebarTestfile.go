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
	var playerOneLifes uint8 = 3
	var playerTwoLifes uint8 = 1
	var playerThreeLifes uint8 = 0
	var playerFourLifes uint8 = 5
	bar.SetLifePointers(&playerOneLifes, &playerTwoLifes, &playerThreeLifes, &playerFourLifes)
	var points uint32 = 1243
	bar.SetPointsPointer(&points)
	go bar.Manager()
	bar.SetSeconds(5 * 60)
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
