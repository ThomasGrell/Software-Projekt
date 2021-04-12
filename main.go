package main

import (
	"./arena"
	"./animations"
	"./characters"
	. "./constants"
	"./items"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
	"fmt"
)

func sun() {
	var winSizeX float64 = 768
	var winSizeY float64 = 672
	var stepSize float64 = 1
	var slice []items.Bombe
	var tempslice []animations.Animation
	var tempMatrix []pixel.Matrix
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
				turfNtreesArena.GetTileSize(), float64(i / turfNtreesArena.GetWidth())*turfNtreesArena.GetTileSize())))
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
		if win.JustPressed(pixelgl.KeyB) {
			var item items.Bombe
			item = items.NewBomb(characters.Player(whiteBomberman))
			slice = append(slice, item)
			
		}

		win.Clear(colornames.Whitesmoke)
		turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
		for index,item := range slice {
			if ((item).GetTimeStamp()).Before(time.Now()) {
				if len(slice)==1 {
					slice = slice[:0]
				} else {
					fmt.Println(slice)
					slice = append(slice[0:index],slice[index+1:]...)
				}
				var l,r,u,d uint8
				l = 3
				r = 3
				u = 3
				d = 3
				ani := animations.NewExplosion(l,r,u,d)
				ani.Show()
				tempslice = append(tempslice,ani)
				tempMatrix = append(tempMatrix,item.GetMatrix())
				
				x, y := turfNtreesArena.GetFieldCoord(item.GetPos())
				if x > 2 {
					for i:=0; i<int(l); i++ {
						turfNtreesArena.RemoveTiles(x-i, y)
					}
				}
				if x < 14 {
					for i:=0; i<int(r); i++ {
						turfNtreesArena.RemoveTiles(x+i, y)
					}
				}
				if y < 12 {
					for i:=0; i<int(u); i++ {
						turfNtreesArena.RemoveTiles(x, y+i)
					}
				}
				if y > 2 {
					for i:=0; i<int(d); i++ {
						turfNtreesArena.RemoveTiles(x,y-i)
					}
				}
				turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
			}
			item.Draw(win)
		}
		for index,a := range(tempslice) {
			a.Update()
			(a.GetSprite()).Draw(win,tempMatrix[index])
		}
		whiteBomberman.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
