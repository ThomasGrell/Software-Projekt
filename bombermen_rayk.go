package main

import (
	"./animations"
	"./characters"
	. "./constants"
	"./sounds"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/png"
	"os"
	"time"
)

func loadPic(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func showIntro(win *pixelgl.Window) {

	pic, err := loadPic("graphics/bomberman.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	win.Clear(colornames.Darkblue)
	win.SetSmooth(true)

	// Startbild: Zoom in
	for i := float64(0); i <= 0.5; i = i + 0.01 {
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
		win.Update()
	}

	// Startbild: Rotate
	for i := float64(0); i <= 6.282; i = i + 0.3141 {
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Rotated(pixel.ZV, i))
		win.Update()
	}
}

func fadeOut(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	imd.SetColorMask(pixel.Alpha(0.05))
	imd.Push(pixel.V(-win.Bounds().W()/2, -win.Bounds().H()/2))
	imd.Push(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2))
	imd.Rectangle(0)
	for i := 0; i < 100; i++ {
		imd.Draw(win)
		win.Update()
		time.Sleep(50 * time.Millisecond)
	}
}

func die(win *pixelgl.Window, wB characters.Enemy) {

	wB.DecLife()
	wB.DecLife()
	wB.DecLife()
	last := time.Now()
	test := animations.NewExplosion(4, 4, 4, 4)
	test.Show()

	for time.Since(last).Seconds() < 5 {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}
		win.Clear(colornames.Black)
		wB.Ani().Update()
		test.Update()
		test.GetSprite().Draw(win, pixel.IM.Moved(wB.GetBaselineCenter()).Moved(test.ToBaseline()))

		wB.Ani().GetSprite().Draw(win, pixel.IM.Moved(wB.GetMovedPos()))
		win.Update()
	}
}

func walkIn(win *pixelgl.Window, wB characters.Enemy) {

	//Figur kommt Betrachter entgegen und wird größer
	for i := 0.1; i <= 2; i += 0.01 {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			win.Destroy()
			break
		}
		win.Clear(colornames.Black)
		wB.Ani().Update()
		wB.Ani().GetSprite().Draw(win, pixel.IM.Moved(wB.GetMovedPos()))
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV, i*i).Moved(win.Bounds().Center()))
		win.Update()
	}
}

func stay(win *pixelgl.Window, wB characters.Enemy) {

	wB.Ani().SetView(Stay)
	wB.Ani().Update()
	win.Clear(colornames.Black)
	wB.SetPos(pixel.V(0, 0))
	wB.Ani().GetSprite().Draw(win, pixel.IM.Moved(wB.GetMovedPos()))
	win.Update()
	for !win.Pressed(pixelgl.KeyEnter) {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			win.Destroy()
			break
		}
		win.Update()
	}

}

func walkAway(win *pixelgl.Window, wB characters.Enemy) {
	var t int

	// Figur läuft nach links
	t = time.Now().Second()
	wB.Ani().SetView(Left)
	var dt float64
	last := time.Now()
	dt = time.Since(last).Seconds()
	for time.Now().Second()-t < 4 {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			win.Destroy()
			break
		}
		wB.Ani().Update()
		win.Clear(colornames.Black)
		v := wB.GetPos()
		dt = time.Since(last).Seconds()
		last = time.Now()
		wB.SetPos(v.Sub(pixel.V(wB.GetSpeed()*dt, 0)))
		wB.Ani().GetSprite().Draw(win, pixel.IM.Moved(wB.GetMovedPos()))
		win.Update()
	}
	/*
		// Figur läuft weg und wird kleiner
		wB.SetDirection(character.Up)
		for i := 3.0; i > 0.1; i -= 0.02 {
			if win.Closed() {
				break
			}
			if win.Pressed(pixelgl.KeyEscape) {
				win.Destroy()
				break
			}
			wB.Update()
			win.Clear(colornames.Black)
			wB.GetSprite().Draw(win, pixel.IM.Moved(wB.GetMinPos()))
			win.SetMatrix(pixel.IM.Scaled(pixel.ZV, i*i))
			win.Update()
		}*/
}

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
	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	s1 := sounds.NewSound(ThroughSpace)
	go s1.PlaySound()
	/*
		win.Update()
		time.Sleep(10 * time.Second)
		showIntro(win)


		time.Sleep(3 * time.Second)

		fadeOut(win)
	*/
	wB := characters.NewEnemy(Fireball)
	wB.Ani().SetView(Intro)
	wB.Ani().Show()

	win.SetSmooth(false)
	walkIn(win, wB)
	stay(win, wB)
	walkAway(win, wB)
	s2 := sounds.NewSound(Deathflash)
	go s2.PlaySound()
	die(win, wB)
	win.Destroy()
}

func main() {
	// Hier darf nichts weiter stehen als die folgende Anweisung:
	pixelgl.Run(run)
}
