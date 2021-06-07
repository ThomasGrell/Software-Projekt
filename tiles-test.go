package main

import (
	"./animations"
	. "./constants"
	"./tiles"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"./characters"
	"golang.org/x/image/colornames"
	"fmt"
	"time"
)


func sun() {

	wincfg := pixelgl.WindowConfig{
		Title:  "GameStat Test",
		Bounds: pixel.R(-100,-100, 400, 400),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	mod := 1
	win.SetMatrix(pixel.IM)//.Scaled(pixel.V(0, 0), 3))
	win.Clear(colornames.Blue)
	win.Update()
	
	wB := characters.NewPlayer(WhiteBomberman)

	itemBatch := pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)
	var ti []tiles.Tile
	var it []tiles.Item
	var bs []tiles.Bombe
	
	i:= tiles.NewItem(SkullItem,pixel.V(0, 0))
	fmt.Println(i.GetType())
	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		itemBatch.Clear()
		if mod%2==0 {
			for _,w := range ti {
				w.SetVisible(true)
			}
			for _,w := range it {
				w.SetVisible(false)
				w.SetDestroyable(false)
				w.SetTimeStamp(time.Now())
			}
		} else {
			for _,w := range ti {
				w.SetVisible(false)
			}
			for _,w := range it {
				w.SetVisible(true)
				w.SetDestroyable(true)
			}
		}
		if win.JustPressed(pixelgl.MouseButton1){				// Move Tile
			t:= tiles.NewTile(House,win.MousePosition())
			ti = append(ti,t)
			fmt.Println("Zeichenmatrix: ",t.GetMatrix())
			fmt.Println("Position: ",t.GetPos())
			fmt.Println("Ist sichtbar: ",t.IsVisible())
			fmt.Println(t.GetType())
			mod++
		} else if win.JustPressed(pixelgl.MouseButton2){				// Move Item
			i:=tiles.NewItem(SkullItem,win.MousePosition())//.Scaled(1/3))
			it = append(it,i)
			mod++
		}
		
		if win.JustPressed(pixelgl.KeySpace) {							// stats
			for _,w := range ti {
				fmt.Println("Tile Animation sichtbar? ",w.IsVisible())
			}
			for _,w := range it {
				fmt.Println("Item Animation sichtbar? ",w.IsVisible())
				fmt.Println("Item zerstörbar? ",w.IsDestroyable())
				fmt.Println("Item TimeStamp: ",w.GetTimeStamp())
			}
			for _,w:=range bs {
				fmt.Println("Bombenpower: ",w.GetPower())
				bb,cc := w.Owner()
				fmt.Println("Bombenpower: ",bb,cc)
			}
		}
		
		if win.JustPressed(pixelgl.KeyB) {	
			b:= tiles.NewBomb(wB,win.MousePosition())
			b.SetPower(1)
			bs = append(bs,b)
		}
		
		for _,w:=range ti {
			w.Draw(itemBatch)
		}
		for _,w:=range it {
			w.Draw(itemBatch)
		}
		for _,w:=range bs {
			w.Draw(itemBatch)
		}
		//win.Clear(colornames.Blue)
		itemBatch.Draw(win)
		win.Update()
	}
 
 }


func main (){
	pixelgl.Run(sun)
}
