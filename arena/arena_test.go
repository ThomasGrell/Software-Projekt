package arena

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"testing"
	"time"
)

func cun() {
	const winWidth, winHeight, zoomFactor = 1024, 768, 3
	wincfg := pixelgl.WindowConfig{
		Title: "Arena Test",
		Bounds: pixel.R(0,0, winWidth,winHeight),
		VSync: true,
	}

	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {panic(err)}
	defer win.Destroy()

	arena := NewArena(0,13,11)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(0,0),basicAtlas)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, -winHeight+25), 2))
	win.Update()
	fmt.Fprintln(basicTxt, "Test einiger Methoden des Pakets 'arena':")
	basicTxt.Draw(win, pixel.IM)
	win.Update()
	basicTxt = text.New(pixel.V(0,winHeight-150),basicAtlas)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 1))
	win.Update()
	fmt.Fprintln(basicTxt, "arena.CoordToVec(12,10):", arena.CoordToVec(12,10))
	x,y := arena.GetFieldCoord(pixel.V(224,176))
	fmt.Fprintln(basicTxt, "arena.GetFieldCoord( pixel.V(224,176) ):",x,y)
	fmt.Fprintln(basicTxt, "arena.GetWidth(), arena.GetHeight():",arena.GetWidth(),", ",arena.GetHeight())
	fmt.Fprintln(basicTxt, "arena.GetLowerLeft():",arena.GetLowerLeft())
	fmt.Fprintln(basicTxt, "arena.GetPermTiles[0][0]:",arena.GetPermTiles()[0][0],"arena.GetPermTiles[1][0]:",arena.GetPermTiles()[1][0])
	fmt.Fprintln(basicTxt, "arena.IsTile(0,2):", arena.IsTile(0,2))
	fmt.Fprintln(basicTxt, "arena.IsFreeTile(12,10)",arena.IsFreeTile(12,10))
	fmt.Fprintln(basicTxt, "arena.GetPassability()[10*arena.GetWidth()+12]:",arena.GetPassability()[10*arena.GetWidth()+12])
	fmt.Fprintln(basicTxt, "arena.GetPassability()[2*arena.GetWidth()+0]:",arena.GetPassability()[2*arena.GetWidth()+0])
	//win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 1))
	//win.Update()
	basicTxt.Draw(win, pixel.IM)

	//win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 3))
	win.Update()
	basicTxt = text.New(pixel.V(0,winHeight-475),basicAtlas)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 2))
	win.Update()
	fmt.Fprintln(basicTxt, "Press Enter to swap arena!")
	basicTxt.Draw(win, pixel.IM)
	win.Update()

	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 3))

	var i int
	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		if i % 3 == 0 {
			arena = NewArena(0,13,11)
			arena.GetCanvas().Draw(win, *(arena.GetMatrix()))
			win.Update()
		}else if i % 3 == 1{
			arena = NewArena(1,13,11)
			arena.GetCanvas().Draw(win, *(arena.GetMatrix()))
			win.Update()
		}else if i % 3 == 2{
			arena = NewArena(2,13,11)
			arena.GetCanvas().Draw(win, *(arena.GetMatrix()))
			win.Update()
		}
		if win.Pressed(pixelgl.KeyEnter) {
			i++
			time.Sleep(1e8)
		}
		time.Sleep(1e7)
		win.Update()
	}
}

func TestMain(*testing.M) {
	pixelgl.Run(cun)
}
//func main() {
//	pixelgl.Run(cun)
//}