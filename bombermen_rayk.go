package main

import (
	"./character"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/png"
	"log"
	"os"
	"time"
)

func playSound(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

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

	var mat pixel.Matrix

	pic, err := loadPic("graphics/bomberman.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	win.Clear(colornames.Darkblue)
	win.SetSmooth(true)

	// Startbild: Zoom in
	for i := float64(0); i <= 0.5; i = i + 0.01 {
		mat = pixel.IM
		mat = mat.Scaled(pixel.ZV, i)
		mat = mat.Moved(win.Bounds().Center())
		sprite.Draw(win, mat)
		win.Update()
	}

	// Startbild: Rotate
	for i := float64(0); i <= 6.282; i = i + 0.3141 {
		mat = pixel.IM
		mat = mat.Scaled(pixel.ZV, 0.5)
		mat = mat.Rotated(pixel.ZV, i)
		mat = mat.Moved(win.Bounds().Center())
		sprite.Draw(win, mat)
		win.Update()
	}
}

func fadeOut(win *pixelgl.Window) {
	for i := 0; i < 100; i++ {
		imd := imdraw.New(nil)
		imd.Color = colornames.Black
		imd.SetColorMask(pixel.Alpha(0.05))
		imd.Push(pixel.V(0, 0))
		imd.Push(pixel.V(win.Bounds().W(), win.Bounds().H()))
		imd.Rectangle(0)
		imd.Draw(win)
		win.Update()
		time.Sleep(20 * time.Millisecond)
	}
}

func walkIn(win *pixelgl.Window, wB character.Character) {
	var t int

	//Figur kommt Betrachter entgegen und wird größer
	for i := 0.1; i <= 3; i += 0.02 {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			win.Destroy()
			break
		}
		win.Clear(colornames.Black)
		wB.Update()
		wB.GetSprite().Draw(win, pixel.IM.Scaled(pixel.ZV, i*i).Moved(win.Bounds().Center()))
		win.Update()
	}

	// Figur steht still da
	wB.Direction(character.Stay)
	wB.Update()
	win.Clear(colornames.Black)
	wB.SetMinPos(win.Bounds().Center())
	wB.GetSprite().Draw(win, pixel.IM.Scaled(pixel.ZV, 9).Moved(wB.GetMinPos()))
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

	// Figur läuft nach links
	t = time.Now().Second()
	wB.Direction(character.Left)
	var dt float64
	last := time.Now()
	dt = time.Since(last).Seconds()
	for time.Now().Second()-t < 3 {
		if win.Closed() {
			break
		}
		if win.Pressed(pixelgl.KeyEscape) {
			win.Destroy()
			break
		}
		wB.Update()
		win.Clear(colornames.Black)
		v := wB.GetMinPos()
		dt = time.Since(last).Seconds()
		last = time.Now()
		wB.SetMinPos(v.Sub(pixel.V(wB.GetSpeed()*dt, 0)))
		wB.GetSprite().Draw(win, pixel.IM.Scaled(pixel.ZV, 9).Moved(wB.GetMinPos()))
		win.Update()
	}

	// Figur läuft weg und wird kleiner
	wB.Direction(character.Up)
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
		wB.GetSprite().Draw(win, pixel.IM.Scaled(pixel.ZV, i*i).Moved(wB.GetMinPos()))
		win.Update()
	}
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

	go playSound("sound/through space.ogg")

	showIntro(win)

	var wB character.Character
	wB = character.NewCharacter(character.RedBomberman)

	time.Sleep(3 * time.Second)

	fadeOut(win)

	win.SetSmooth(false)

	walkIn(win, wB)

}

func main() {
	// Hier darf nichts weiter stehen als die folgende Anweisung:
	pixelgl.Run(run)
}
