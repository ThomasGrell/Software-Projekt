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

func fun() {
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

	turfNtreesArena = arena.NewArena(0, winSizeX, winSizeY)

	whiteBomberman := characters.NewPlayer(WhiteBomberman)
	whiteBomberman.SetScale(1)
	whiteBomberman.Ani().Show()

A:
	for i := 0; i < 15; i++ {
		for j := 0; j < 17; j++ {
			if turfNtreesArena.GetPassableTiles()[i][j] {
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
		grDir := turfNtreesArena.GrantedDirection(whiteBomberman.GetPosBox(), whiteBomberman.GetPos()) // [4]bool left-right-up-down granted?

		if win.Pressed(pixelgl.KeyLeft) && grDir[0] {
			whiteBomberman.MoveTo(pixel.V(-stepSize, 0))
			whiteBomberman.Ani().SetView(Left)
		} else if win.Pressed(pixelgl.KeyRight) && grDir[1] {
			whiteBomberman.MoveTo(pixel.V(stepSize, 0))
			whiteBomberman.Ani().SetView(Right)
		} else if win.Pressed(pixelgl.KeyUp) && grDir[2] {
			whiteBomberman.MoveTo(pixel.V(0, stepSize))
			whiteBomberman.Ani().SetView(Up)
		} else if win.Pressed(pixelgl.KeyDown) && grDir[3] {
			whiteBomberman.MoveTo(pixel.V(0, -stepSize))
			whiteBomberman.Ani().SetView(Down)
		}
		if win.Pressed(pixelgl.KeyB) {
			var item items.Bombe
			item = items.NewBomb(characters.Player(whiteBomberman))
			slice = append(slice, item)
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
	pixelgl.Run(fun)
}
