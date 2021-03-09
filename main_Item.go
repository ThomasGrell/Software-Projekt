package main

import (
	"./arena"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	//"./animations"
	//"./characters"
	. "./constants"
	"./items"
	//"image/png"
	"os"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func fun() {
	var winSizeX float64 = 816
	var winSizeY float64 = 720
	//var tileSize float64 = 8

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	//charactersPic, err := loadPicture("graphics/animations.png")
	if err != nil {
		panic(err)
	}
	
	whiteBomberman := characters.NewPlayer(WhiteBomberman)
	
	var item = items.NewItem(PowerItem,pixel.V(winSizeX/2, winSizeY/2))
	
	//itemSprite := (item.Ani()).GetSprite()
	
	//whiteBomberman := pixel.NewSprite(charactersPic, pixel.R(0, 361, 16, 385))
	//mat := pixel.IM
	//mat = mat.Moved(pixel.V(winSizeX/2, winSizeY/2))
	//mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(3.3, 3.3))

	win.Clear(colornames.Whitesmoke)
	arena.Draw(win)
	//item.Ani().Update()
	
	//arena.Draw(win)
	//item.Ani().Update()
	//item.Ani().GetSprite().Draw(win, mat)
	//win.Update()
	item.Draw(win)

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		
		win.Update()
	}
	
}

func main() {
	pixelgl.Run(fun)
}
