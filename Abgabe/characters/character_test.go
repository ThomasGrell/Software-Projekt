package characters

import (
	. "../constants"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"testing"
	"time"
)

func run() {
	wincfg := pixelgl.WindowConfig{
		Title:  "Character Test",
		Bounds: pixel.R(-640, -480, 640, 480),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 2))

	p := NewPlayer(WhiteBomberman)
	e := NewEnemy(Balloon)
	p.Ani().Show()
	e.Ani().Show()
	p.MoveTo(pixel.V(200, 0))
	p.Ani().SetView(Left)
	e.MoveTo(pixel.V(-200, 0))
	e.Ani().SetView(Right)

	for i := 0; i < 180; i++ {
		win.Clear(color.Black)
		p.Draw(win)
		e.Draw(win)
		win.Update()
		p.Move(pixel.V(-1, 0))
		e.Move(pixel.V(1, 0))
		time.Sleep(time.Millisecond * 2)
	}

	p.Ani().SetView(Down)
	p.Ani().SetView(Stay)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(-300, 0), basicAtlas)
	fmt.Fprintln(basicTxt, "Test einiger Methoden des Pakets 'characters'")
	fmt.Fprintln(basicTxt, "GetLife:", p.GetLife())
	p.IncLife()
	fmt.Fprintln(basicTxt, "IncLife: ", p.GetLife())
	p.DecLife()
	fmt.Fprintln(basicTxt, "DecLife: ", p.GetLife())
	p.SetLife(5)
	fmt.Fprintln(basicTxt, "SetLife(5): ", p.GetLife())
	fmt.Fprintln(basicTxt, "GetSpeed: ", p.GetSpeed())
	p.IncSpeed()
	fmt.Fprintln(basicTxt, "IncSpeed: ", p.GetSpeed())
	p.DecSpeed()
	fmt.Fprintln(basicTxt, "DecSpeed: ", p.GetSpeed())
	fmt.Fprintln(basicTxt, "GetPoints: ", p.GetPoints())
	p.AddPoints(100)
	fmt.Fprintln(basicTxt, "AddPoints(100): ", p.GetPoints())
	fmt.Fprintln(basicTxt, "GetBombs: ", p.GetBombs())
	p.IncBombs()
	fmt.Fprintln(basicTxt, "IncBombs: ", p.GetBombs())
	p.DecBombs()
	fmt.Fprintln(basicTxt, "DecBombs: ", p.GetBombs())
	fmt.Fprintln(basicTxt, "GetPower: ", p.GetPower())
	p.IncPower()
	fmt.Fprintln(basicTxt, "IncPower: ", p.GetPower())
	p.SetPower(3)
	fmt.Fprintln(basicTxt, "SetPower(3): ", p.GetPower())

	win.Clear(color.Black)
	basicTxt.Draw(win, pixel.IM)
	p.Draw(win)
	e.Draw(win)
	win.Update()

	last := time.Now()
	for time.Since(last) < time.Second*10 {
	}

}

func TestMain(*testing.M) {
	pixelgl.Run(run)
}
