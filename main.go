package main

import (
	"./animations"
	"./arena"
	"./characters"
	. "./constants"
	"./items"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

func sun() {
	var winSizeX float64 = 768
	var winSizeY float64 = 672
	var stepSize float64 = 1
	var bombs []items.Bombe
	var turfNtreesArena arena.Arena
	var tempAniSlice [][]interface{} // [Animation][Matrix]

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	turfNtreesArena = arena.NewArena(13, 11)

	whiteBomberman := characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()
	whiteBomberman.IncPower()
	whiteBomberman.IncPower()

	// Put character at free space with at least two free neighbours in a row
A:
	for i := 2 * turfNtreesArena.GetWidth(); i < len(turfNtreesArena.GetPassability())-2*turfNtreesArena.GetWidth(); i++ { // Einschränkung des Wertebereichs von i um index out of range Probleme zu vermeiden
		if turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-1] && turfNtreesArena.GetPassability()[i-2] || // checke links, rechts, oben, unten
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+1] && turfNtreesArena.GetPassability()[i+2] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i+2*turfNtreesArena.GetWidth()] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i-2*turfNtreesArena.GetWidth()] {
			whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i%turfNtreesArena.GetWidth())*
				TileSize, float64(i/turfNtreesArena.GetWidth())*TileSize)))
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
			bombs = append(bombs, item)

		}

		win.Clear(colornames.Whitesmoke)
		turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
		for index, item := range bombs {
			if ((item).GetTimeStamp()).Before(time.Now()) {
				if len(bombs) == 1 {
					bombs = bombs[:0]
				} else {
					fmt.Println(bombs)
					bombs = append(bombs[0:index], bombs[index+1:]...)
				}
				//destTiles := turfNtreesArena.GetDestroyableTiles()
				x, y := turfNtreesArena.GetFieldCoord(item.GetPos())
				power := int(item.GetPower())
				l, r, u, d := power, power, power, power
				if 2+l > x {
					l = x - 2
				}
				if 14-r < x {
					r = 14 - x
				}
				if 2+d > y {
					d = y - 2
				}
				if 12-u < y {
					u = 12 - y
				}

				if x > 2 {
					for i := 1; i <= int(l); i++ {
						if turfNtreesArena.RemoveTiles(x-i, y) {
							l = i
							break
						}
					}
				}
				if x < 14 {
					for i := 1; i <= int(r); i++ {
						if turfNtreesArena.RemoveTiles(x+i, y) {
							r = i
							break
						}
					}
				}
				if y < 12 {
					for i := 1; i <= int(u); i++ {
						if turfNtreesArena.RemoveTiles(x, y+i) {
							u = i
							break
						}
					}
				}
				if y > 2 {
					for i := 1; i <= int(d); i++ {
						if turfNtreesArena.RemoveTiles(x, y-i) {
							d = i
							break
						}
					}
				}

				ani := animations.NewExplosion(uint8(l), uint8(r), uint8(u), uint8(d))
				ani.Show()
				tempAni := make([]interface{}, 2)
				tempAni[0] = ani
				tempAni[1] = (item.GetMatrix()).Moved(ani.ToCenter())
				tempAniSlice = append(tempAniSlice, tempAni)

				turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
			}
			item.Draw(win)
		}
		for _, a := range tempAniSlice {
			ani := (a[0]).(animations.Animation)
			ani.Update()
			mtx := (a[1]).(pixel.Matrix)
			(ani.GetSprite()).Draw(win, mtx)
		}
		whiteBomberman.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
