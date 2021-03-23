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
	var turfNtreesArena *arena.Arena

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	turfNtreesArena = arena.NewArena(winSizeX, winSizeY)

	whiteBomberman := characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()

A:
	for i := 0; i < 15; i++ {
		for j := 0; j < 17; j++ {
			if turfNtreesArena.GetBoolMap()[i][j] {
				whiteBomberman.MoveTo(pixel.V(float64(j)*arena.GetTileSize(), float64(i)*arena.GetTileSize()-14)) // Hier bekommt die Animation ihren Ort.
				//fmt.Println(i,j)
				break A
			}
		}

	}

	win.Clear(colornames.Whitesmoke)
	turfNtreesArena.GetCanvas().Draw(win, turfNtreesArena.GetMatrix())
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
			x, y := turfNtreesArena.GetFieldCoord(whiteBomberman.GetPosBox().Min)
			if x > 2 {
				turfNtreesArena.RemoveTile(x-1, y)
			}
			if x < 14 {
				turfNtreesArena.RemoveTile(x+1, y)
			}
			if y < 12 {
				turfNtreesArena.RemoveTile(x, y+1)
			}
			if y > 2 {
				turfNtreesArena.RemoveTile(x, y-1)
			}
			turfNtreesArena.GetCanvas().Draw(win, turfNtreesArena.GetMatrix())
		}

		//fmt.Println(whiteBomberman.GetPosBox(),"Size",whiteBomberman.GetSize())

		win.Clear(colornames.Whitesmoke)
		turfNtreesArena.GetCanvas().Draw(win, turfNtreesArena.GetMatrix())
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
