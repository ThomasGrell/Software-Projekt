package main

import (
	"./arena"
	"./characters"
	. "./constants"
	"./items"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func sun() {
	var winSizeX float64 = 768
	var winSizeY float64 = 672
	var stepSize float64 = 1
	var slice []items.Bombe
	var turfNtreesArena arena.Arena

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	turfNtreesArena = arena.NewArena(13,11 )

	whiteBomberman := characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()

// Put character at free space with at least two free neighbours in a row
A:	for i := 2 * turfNtreesArena.GetWidth(); i < len(turfNtreesArena.GetPassability()) - 2 * turfNtreesArena.GetWidth(); i++ {	// Einschränkung des Wertebereichs von i um index out of range Probleme zu vermeiden
		if turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-1] && turfNtreesArena.GetPassability()[i-2] ||	// checke links, rechts, oben, unten
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+1] && turfNtreesArena.GetPassability()[i+2] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i + turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i+ 2 * turfNtreesArena.GetWidth()] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i - turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i - 2 * turfNtreesArena.GetWidth()] {
			whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i % turfNtreesArena.GetWidth())*
				TileSize, float64(i / turfNtreesArena.GetWidth())*TileSize)))
			break A
		}
	}

	win.Clear(colornames.Whitesmoke)
	turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
	whiteBomberman.Ani().Update()
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 3))
	win.Update()

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		grDir := turfNtreesArena.GrantedDirections(whiteBomberman.GetPosBox()) // [4]bool left-right-up-down granted?
		if win.Pressed(pixelgl.KeyLeft) && grDir[0] {
			whiteBomberman.Move(pixel.V(-stepSize, 0))
			whiteBomberman.Ani().SetView(Left)
		}
		if win.Pressed(pixelgl.KeyRight) && grDir[1] {
			whiteBomberman.Move(pixel.V(stepSize, 0))
			whiteBomberman.Ani().SetView(Right)
		}
		if win.Pressed(pixelgl.KeyUp) && grDir[2] {
			whiteBomberman.Move(pixel.V(0, stepSize))
			whiteBomberman.Ani().SetView(Up)
		}
		if win.Pressed(pixelgl.KeyDown) && grDir[3] {
			whiteBomberman.Move(pixel.V(0, -stepSize))
			whiteBomberman.Ani().SetView(Down)
		}
		if win.Pressed(pixelgl.KeyB) {
			var item items.Bombe
			item = items.NewBomb(characters.Player(whiteBomberman))
			slice = append(slice, item)
			x, y := turfNtreesArena.GetFieldCoord(item.GetPos())
			if x > 2 {
				turfNtreesArena.RemoveTiles(x-1, y)
			}
			if x < 14 {
				turfNtreesArena.RemoveTiles(x+1, y)
			}
			if y < 12 {
				turfNtreesArena.RemoveTiles(x, y+1)
			}
			if y > 2 {
				turfNtreesArena.RemoveTiles(x, y-1)
			}
			turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
		}

		win.Clear(colornames.Whitesmoke)
		turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
		for _, item := range slice {
			item.(items.Bombe).Draw(win)
		}
		whiteBomberman.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
