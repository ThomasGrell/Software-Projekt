package animations

import (
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"testing"
	"time"
)

func run() {
	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(-640, -480, 640, 480),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	win.SetMatrix(pixel.IM.Moved(pixel.V(-5*32, -32)).Scaled(pixel.V(0, 0), 3))

	var ani [33]Animation
	for i := uint8(1); i <= 10; i++ {
		ani[3*i-1] = NewAnimation(3 * i)
		ani[3*i-1].SetView(Down)
		ani[3*i-1].Show()
		ani[3*i] = NewAnimation(3*i + 1)
		ani[3*i].SetView(Down)
		ani[3*i].Show()
		ani[3*i+1] = NewAnimation(3*i + 2)
		ani[3*i+1].SetView(Down)
		ani[3*i+1].Show()
	}

	last := time.Now()
	for time.Since(last) < time.Second*10 {
		win.Clear(color.Black)
		for i := uint8(1); i <= 10; i++ {
			ani[3*i-1].Update()
			ani[3*i-1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 32)))
			ani[3*i].Update()
			ani[3*i].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 64)))
			ani[3*i+1].Update()
			ani[3*i+1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 98)))
		}
		win.Update()
	}
	for i := uint8(1); i <= 10; i++ {
		ani[3*i-1].SetView(Up)
		ani[3*i].SetView(Up)
		ani[3*i+1].SetView(Up)
	}

	last = time.Now()
	for time.Since(last) < time.Second*10 {
		win.Clear(color.Black)
		for i := uint8(1); i <= 10; i++ {
			ani[3*i-1].Update()
			ani[3*i-1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 32)))
			ani[3*i].Update()
			ani[3*i].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 64)))
			ani[3*i+1].Update()
			ani[3*i+1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 98)))
		}
		win.Update()
	}
	bang := NewExplosion(5, 5, 5, 5)
	bang.Show()

	for i := uint8(1); i <= 10; i++ {
		ani[3*i-1].Die()
		ani[3*i].Die()
		ani[3*i+1].Die()
	}

	last = time.Now()
	for time.Since(last) < time.Second*10 {
		win.Clear(color.Black)
		for i := uint8(1); i <= 10; i++ {
			ani[3*i-1].Update()
			ani[3*i-1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 32)))
			ani[3*i].Update()
			ani[3*i].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 64)))
			ani[3*i+1].Update()
			ani[3*i+1].GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(i)*32, 98)))
		}
		bang.Update()
		bang.GetSprite().Draw(win, pixel.IM.Moved(pixel.V(float64(5)*32, 64)))
		win.Update()
	}

}

func TestMain(*testing.M) {
	pixelgl.Run(run)
}
