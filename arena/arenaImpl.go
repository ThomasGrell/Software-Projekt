package arena

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
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

func Draw(win *pixelgl.Window) {
	//var winSizeX float64 = 816
	//var winSizeY float64 = 720
	var linksUnten pixel.Vec = pixel.V(0, 0)
	var zoomFactor float64 = 3
	var kurzeSeiteTeile = 9
	var langeSeiteTeile = 13
	var edgeLowLeftCenterX, edgeLowLeftCenterY,
		wallLeftCenterX, wallLeftCenterY,
		edgeHiLeftCenterX, edgeHiLeftCenterY,
		hiWallCenterX, hiWallCenterY,
		edgeHiRightCenterX, edgeHiRightCenterY,
		wallRightCenterX, wallRightCenterY,
		edgeLowRightCenterX, edgeLowRightCenterY,
		loWallCenterX, loWallCenterY,
		turfCenterX, turfCenterY float64

	//wincfg := pixelgl.WindowConfig{
	//	Title: "Test Map",
	//	Bounds: pixel.R(0,0,winSizeX, winSizeY),
	//	VSync: true,
	//}
	//win, err := pixelgl.NewWindow(wincfg)
	//if err != nil {
	//	panic(err)
	//}
	//win.Clear(colornames.Whitesmoke)

	TilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	edgeLowLeft := pixel.NewSprite(TilesPic, pixel.R(288, 81, 312, 113))
	wallLeft := pixel.NewSprite(TilesPic, pixel.R(288, 114, 312, 130))
	edgeHiLeft := pixel.NewSprite(TilesPic, pixel.R(288, 114, 312, 144))
	hiWall := pixel.NewSprite(TilesPic, pixel.R(312, 136, 328, 144))
	edgeHiRight := pixel.NewSprite(TilesPic, pixel.R(344, 114, 368, 144))
	wallRight := pixel.NewSprite(TilesPic, pixel.R(344, 112, 368, 128))
	edgeLowRight := pixel.NewSprite(TilesPic, pixel.R(344, 81, 368, 115))
	loWall := pixel.NewSprite(TilesPic, pixel.R(312, 81, 328, 87))
	turf := pixel.NewSprite(TilesPic, pixel.R(112, 288, 128, 304))
	edgeLowLeftCenterX = (312 - 288) / 2
	edgeLowLeftCenterY = (112 - 81) / 2
	wallLeftCenterX = edgeLowLeftCenterX
	wallLeftCenterY = (130 - 114) / 2
	edgeHiLeftCenterX = edgeLowLeftCenterX
	edgeHiLeftCenterY = (144 - 114) / 2
	hiWallCenterX = (328 - 312) / 2
	hiWallCenterY = (144 - 136) / 2
	edgeHiRightCenterX = (368 - 344) / 2
	edgeHiRightCenterY = edgeHiLeftCenterY
	wallRightCenterX = edgeHiRightCenterX
	wallRightCenterY = (128 - 112) / 2
	edgeLowRightCenterX = edgeHiRightCenterX
	edgeLowRightCenterY = edgeLowLeftCenterY
	loWallCenterX = hiWallCenterX
	loWallCenterY = (87 - 81) / 2
	turfCenterX = (128 - 112) / 2
	turfCenterY = (304 - 288) / 2

	edgeLowLeftMat := pixel.IM
	edgeLowLeftMat = edgeLowLeftMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	edgeLowLeftMat = edgeLowLeftMat.Moved(pixel.V(edgeLowLeftCenterX*zoomFactor, (edgeLowLeftCenterY+1)*zoomFactor))
	// Moved verschiebt den MatrixMITTELPUNKT, +1 in der y-Komponente, weil in tiles.png etwas mehr als 3 tiles in die Mitte passen
	wallLeftMat := pixel.IM
	wallLeftMat = wallLeftMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	wallLeftMat = wallLeftMat.Moved(pixel.V(wallLeftCenterX*zoomFactor, (2*edgeLowLeftCenterY+wallLeftCenterY+1)*zoomFactor)) // +1 in der y-Komponente, weil in tiles.png etwas mehr als 3 tiles in die Mitte passen
	edgeHiLeftMat := pixel.IM
	edgeHiLeftMat = edgeHiLeftMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	edgeHiLeftMat = edgeHiLeftMat.Moved(pixel.V(edgeHiLeftCenterX*zoomFactor, (2*edgeLowLeftCenterY+
		2*float64(kurzeSeiteTeile)*wallLeftCenterY+edgeHiLeftCenterY+1)*zoomFactor))
	hiWallMat := pixel.IM
	hiWallMat = hiWallMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	hiWallMat = hiWallMat.Moved(pixel.V((2*edgeHiLeftCenterX+hiWallCenterX)*zoomFactor, (2*edgeLowLeftCenterY+
		2*wallRightCenterY*float64(kurzeSeiteTeile)+2*edgeHiLeftCenterY-hiWallCenterY+1)*zoomFactor))
	edgeHiRightMat := pixel.IM
	edgeHiRightMat = edgeHiRightMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	edgeHiRightMat = edgeHiRightMat.Moved(pixel.V((2*edgeHiLeftCenterX+2*hiWallCenterX*float64(langeSeiteTeile+1)+
		edgeHiRightCenterX)*zoomFactor, (2*edgeLowRightCenterY+2*wallRightCenterY*float64(kurzeSeiteTeile)+
		edgeHiRightCenterY+1)*zoomFactor))
	wallRightMat := pixel.IM
	wallRightMat = wallRightMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	wallRightMat = wallRightMat.Moved(pixel.V((2*edgeLowLeftCenterX+2*loWallCenterX*float64(langeSeiteTeile+1)+
		wallRightCenterX)*zoomFactor, (2*edgeLowRightCenterY+wallRightCenterY)*zoomFactor))
	edgeLowRightMat := pixel.IM
	edgeLowRightMat = edgeLowRightMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	edgeLowRightMat = edgeLowRightMat.Moved(pixel.V((2*edgeLowLeftCenterX+2*loWallCenterX*float64(langeSeiteTeile+1)+
		edgeLowRightCenterX)*zoomFactor, (edgeLowRightCenterY+2)*zoomFactor))
	loWallMat := pixel.IM
	loWallMat = loWallMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	loWallMat = loWallMat.Moved(pixel.V((2*edgeLowLeftCenterX+loWallCenterX)*zoomFactor, loWallCenterY*zoomFactor))
	turfMat := pixel.IM
	turfMat = turfMat.ScaledXY(linksUnten, pixel.V(zoomFactor, zoomFactor))
	turfMat = turfMat.Moved(pixel.V((2*wallLeftCenterX+turfCenterX)*zoomFactor, (2*loWallCenterY+turfCenterY)*zoomFactor))

	edgeLowLeft.Draw(win, edgeLowLeftMat)
	wallLeft.Draw(win, wallLeftMat)
	edgeHiLeft.Draw(win, edgeHiLeftMat)
	hiWall.Draw(win, hiWallMat)
	edgeHiRight.Draw(win, edgeHiRightMat)
	wallRight.Draw(win, wallRightMat)
	edgeLowRight.Draw(win, edgeLowRightMat)
	loWall.Draw(win, loWallMat)

	wallLeftShift := 2 * wallLeftCenterY * zoomFactor
	for i := 0; i < kurzeSeiteTeile; i++ { // draws left wall
		wallLeftMat = wallLeftMat.Moved(pixel.V(0, wallLeftShift))
		wallLeft.Draw(win, wallLeftMat)
		wallRightMat = wallRightMat.Moved(pixel.V(0, wallLeftShift))
		wallRight.Draw(win, wallRightMat)
	}
	hiWallShift := 2 * hiWallCenterX * zoomFactor
	for i := 0; i < langeSeiteTeile; i++ {
		hiWallMat = hiWallMat.Moved(pixel.V(hiWallShift, 0))
		hiWall.Draw(win, hiWallMat)
		loWallMat = loWallMat.Moved(pixel.V(hiWallShift, 0))
		loWall.Draw(win, loWallMat)
	}
	turfRightShift := 2 * turfCenterX * zoomFactor
	turfUpShift := 2 * turfCenterY * zoomFactor
	for i := 0; i <= kurzeSeiteTeile+2; i++ {
		turf.Draw(win, turfMat)
		for j := 0; j < langeSeiteTeile; j++ {
			turfMat = turfMat.Moved(pixel.V(turfRightShift, 0))
			turf.Draw(win, turfMat)
		}
		turfMat = turfMat.Moved(pixel.V(float64(-langeSeiteTeile)*turfRightShift, turfUpShift))
	}

	//for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
	//	win.Update()
	//}

}

//func main () {
//	pixelgl.Run(run)
//}
