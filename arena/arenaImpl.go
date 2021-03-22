package arena

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const tileSize float64 = 16

type Arena struct {
	canvas           *pixelgl.Canvas
	destroyableTiles [2][35]int
	matrix           pixel.Matrix
	kollisionsKarte  [15][17]pixel.Rect
	permTiles        [2][36]int
	//tiles           [11][13]int
	untenLinks     pixel.Vec    // linke untere Spielfeldecke
	wahrheitsKarte [15][17]bool // [Zeilen][Spalten]
	w, h           float64
}

func NewArena(width, heigth float64) *Arena {
	var a *Arena = new(Arena)
	a.w = width
	a.h = heigth
	a.permTiles = setCabins()
	for i := 2; i < 13; i++ {
		for j := 2; j < 15; j++ {
			a.wahrheitsKarte[i][j] = true
		}
	}
	for i := 0; i < 36; i++ {
		a.wahrheitsKarte[a.permTiles[1][i]][a.permTiles[0][i]] = false
	}
	a.destroyableTiles = setStub()
	for i := 0; i < 35; i++ {
		a.wahrheitsKarte[a.destroyableTiles[1][i]][a.destroyableTiles[0][i]] = false
	}
	//for i := 14; i >= 0; i-- {
	//	fmt.Println(a.wahrheitsKarte[i])
	//}
	a.untenLinks = pixel.V(24, 6)
	a.matrix = pixel.IM
	a.matrix = a.matrix.Moved(pixel.V(width/2, heigth/2))
	a.canvas = pixelgl.NewCanvas(pixel.R(0, 0, width, heigth))
	drawWallsNturf(a.canvas)
	drawCabin(a.canvas, a)
	a.drawStub(a.canvas)
	putWallsOnMap(a) // schreibt Rechtecke in die virtuelle Karte für Kollisionen
	//for i := 14; i >= 0; i-- {
	//fmt.Println(a.kollisionsKarte[i])
	//}
	return a
}

func (a *Arena) GetBoolMap() [15][17]bool {
	return a.wahrheitsKarte
}
func (a *Arena) GetCanvas() *pixelgl.Canvas {
	return a.canvas
}
func (a *Arena) GetFieldCoord(v pixel.Vec) (x, y int) {
	var columns int = 13
	var rows int = 11
	x = int(math.Trunc((v.X-a.untenLinks.X)/tileSize))%(columns+1) + 2
	y = int(math.Trunc((v.Y-a.untenLinks.Y)/tileSize))%(rows+1) + 2
	return
}
func (a *Arena) GetMatrix() pixel.Matrix {
	return a.matrix
}
func GetTileSize() float64 {
	return tileSize
}
func (a *Arena) GrantedDirections(posBox pixel.Rect) [4]bool { // {links,rechts,oben,unten}
	var grDir [4]bool
	var x1, x2, y1, y2 int
	var columns int = 13
	var rows int = 11
	//fmt.Println(posVec)
	x1 = int(math.Trunc((posBox.Min.X-a.untenLinks.X)/tileSize))%(columns+1) + 2
	y1 = int(math.Trunc((posBox.Min.Y-a.untenLinks.Y)/tileSize))%(rows+1) + 2 // Eintritt in nächste Kachel oben erst 2 Pixel später (ist schicker)
	x2 = int(math.Trunc((posBox.Max.X-a.untenLinks.X)/tileSize))%(columns+1) + 2
	y2 = int(math.Trunc((posBox.Max.Y-a.untenLinks.Y)/tileSize))%(rows+1) + 2
	//fmt.Println("Collison Box", posBox)
	//fmt.Println("", a.kollisionsKarte[y2][x2-1], "\n", a.kollisionsKarte[y1][x2-1])
	if !pixel.R(posBox.Min.X-1, posBox.Min.Y, posBox.Max.X, posBox.Max.Y).Intersects(a.kollisionsKarte[y1][x2-1]) &&
		!pixel.R(posBox.Min.X-1, posBox.Min.Y, posBox.Max.X, posBox.Max.Y).Intersects(a.kollisionsKarte[y2][x2-1]) { // Left
		grDir[0] = true
	} else {
		grDir[0] = false
	}
	if !pixel.R(posBox.Min.X, posBox.Min.Y, posBox.Max.X+1, posBox.Max.Y+1).Intersects(a.kollisionsKarte[y1][x1+1]) &&
		!pixel.R(posBox.Min.X, posBox.Min.Y, posBox.Max.X+1, posBox.Max.Y+1).Intersects(a.kollisionsKarte[y2][x1+1]) { // Right
		grDir[1] = true
	} else {
		grDir[1] = false
	}
	if !pixel.R(posBox.Min.X, posBox.Min.Y, posBox.Max.X, posBox.Max.Y+1).Intersects(a.kollisionsKarte[y1+1][x1]) &&
		!pixel.R(posBox.Min.X, posBox.Min.Y, posBox.Max.X, posBox.Max.Y+1).Intersects(a.kollisionsKarte[y1+1][x2]) { // Up
		grDir[2] = true
	} else {
		grDir[2] = false
	}
	if !pixel.R(posBox.Min.X, posBox.Min.Y-1, posBox.Max.X, posBox.Max.Y).Intersects(a.kollisionsKarte[y2-1][x1]) &&
		!pixel.R(posBox.Min.X, posBox.Min.Y-1, posBox.Max.X, posBox.Max.Y).Intersects(a.kollisionsKarte[y2-1][x2]) { // Down
		grDir[3] = true
	} else {
		grDir[3] = false
	}
	return grDir
	//return [4]bool{grDir[0],grDir[1],true,true}
}
func (a *Arena) RemoveTile(x, y int) {
	//for i := 14; i >= 0; i-- {
	//	fmt.Println(a.wahrheitsKarte[i])
	//}
	//fmt.Println("")
	k := checkCoordsOfDestroyables(x, y, a.destroyableTiles)
	if k != 42 { // 42 als Fehlerfall: diese Koordinaten wurden nicht gefunden
		a.destroyableTiles[0][k] = 0
		a.destroyableTiles[1][k] = 0
		a.kollisionsKarte[y][x] = pixel.R(0, 0, 0, 0)
		a.wahrheitsKarte[y][x] = true
		//for i := 14; i >= 0; i-- {
		//	fmt.Println(a.wahrheitsKarte[i])
		//}
		drawWallsNturf(a.canvas)
		drawCabin(a.canvas, a)
		a.drawStub(a.canvas)
	}
}

//------------------------- Hilfsfunktionen ---------------------------------

func checkCoordsOfDestroyables(x, y int, locations [2][35]int) int {
	var j int
	for i := 0; i < 35; i++ {
		if locations[0][i] == x && locations[1][i] == y {
			j = i
			return j
			//fmt.Println(j)
		}
	}
	j = 42
	return j
}

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

func checkFreeTiles(x, y int, locations [2][35]int) (u int, v int, w bool) {
	var cabinTiles [2][36]int = setCabins()
	w = true
	u = x
	v = y
	for i := 0; i < 36; i++ {
		if cabinTiles[0][i] == x && cabinTiles[1][i] == y {
			w = false
		}
	}
	for i := 0; i < 35; i++ {
		if locations[0][i] == x && locations[1][i] == y {
			w = false
		}
	}
	return
}

// Erzeugt x-y-Koordinaten der Häuser für eine Spielfeldmatrix (13x11 Matrix)
func setCabins() [2][36]int {
	var locations [2][36]int
	//rand.Seed(time.Now().UnixNano())
	//for i := range locations[0] {
	//	locations[0][i] = rand.Intn(13)+2
	//	locations[1][i] = rand.Intn(11)+2
	//}
	locations[0][0] = 3
	locations[0][1] = 5
	locations[0][2] = 7
	locations[0][3] = 9
	locations[0][4] = 11
	locations[0][5] = 13
	locations[1][0] = 3
	locations[1][1] = 3
	locations[1][2] = 3
	locations[1][3] = 3
	locations[1][4] = 3
	locations[1][5] = 3
	locations[0][6] = 10
	locations[0][7] = 11
	locations[0][8] = 12
	locations[0][9] = 3
	locations[0][10] = 5
	locations[0][11] = 7
	locations[1][6] = 4
	locations[1][7] = 4
	locations[1][8] = 4
	locations[1][9] = 5
	locations[1][10] = 5
	locations[1][11] = 5
	locations[0][12] = 9
	locations[0][13] = 11
	locations[0][14] = 12
	locations[0][15] = 13
	locations[0][16] = 3
	locations[0][17] = 5
	locations[1][12] = 5
	locations[1][13] = 5
	locations[1][14] = 5
	locations[1][15] = 5
	locations[1][16] = 7
	locations[1][17] = 7
	locations[0][18] = 7
	locations[0][19] = 9
	locations[0][20] = 11
	locations[0][21] = 13
	locations[0][22] = 3
	locations[0][23] = 5
	locations[1][18] = 7
	locations[1][19] = 7
	locations[1][20] = 7
	locations[1][21] = 7
	locations[1][22] = 9
	locations[1][23] = 9
	locations[0][24] = 7
	locations[0][25] = 9
	locations[0][26] = 11
	locations[0][27] = 13
	locations[0][28] = 4
	locations[0][29] = 3
	locations[1][24] = 9
	locations[1][25] = 9
	locations[1][26] = 9
	locations[1][27] = 9
	locations[1][28] = 10
	locations[1][29] = 11
	locations[0][30] = 5
	locations[0][31] = 7
	locations[0][32] = 9
	locations[0][33] = 11
	locations[0][34] = 13
	locations[0][35] = 9
	locations[1][30] = 11
	locations[1][31] = 11
	locations[1][32] = 11
	locations[1][33] = 11
	locations[1][34] = 11
	locations[1][35] = 12
	return locations
}

// Erzeugt x-y-Koordinaten der Baumstümpfe für eine Spielfeldmatrix (13x11 Matrix)
func setStub() [2][35]int {
	var locations [2][35]int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 35; i++ {
		for {
			x, y, w := checkFreeTiles(rand.Intn(13)+2, rand.Intn(11)+2, locations)
			if w {
				locations[0][i] = x
				locations[1][i] = y
				break
			}
		}
	}
	return locations
}

func putWallsOnMap(a *Arena) { // setzt die Kollisionsboxen in die kollisionKarte
	for i := 1; i < 14; i++ {
		for j := 1; j < 16; j++ {
			if i < 2 || i > 12 { // first and last row
				a.kollisionsKarte[i][j] = pixel.R(
					tileSize*float64(j-2)+24,
					tileSize*float64(i-2)+6,
					tileSize*float64(j-1)+24,
					tileSize*float64(i-1)+6)
			}
			if j < 2 || j > 14 {
				a.kollisionsKarte[i][j] = pixel.R(
					tileSize*float64(j-2)+24,
					tileSize*float64(i-2)+6,
					tileSize*float64(j-1)+24,
					tileSize*float64(i-1)+6)
			}
		}
	}
}

func (a *Arena) drawStub(can *pixelgl.Canvas) { // and SetStubs on kollisionsKarte !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	const numbOfCabs = 35
	var stubRow [numbOfCabs]int
	var stubColumn [numbOfCabs]int
	stubCoord := a.destroyableTiles
	stubRow = stubCoord[1]
	stubColumn = stubCoord[0]

	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	stub := pixel.NewSprite(tilesPic, pixel.R(80, 304, 96, 288))
	stubMat := pixel.IM
	stubMat = stubMat.Moved(pixel.V(tileSize/2, tileSize/2).Add(a.untenLinks))

	for i := 0; i < 35; i++ {
		if stubColumn[i] != 0 {
			stubMat = stubMat.Moved(pixel.V(float64(stubColumn[i]-2)*tileSize, float64(stubRow[i]-2)*tileSize))
			stub.Draw(can, stubMat)
			a.kollisionsKarte[stubRow[i]][stubColumn[i]] = pixel.R(
				tileSize*float64(stubColumn[i]-2)+a.untenLinks.X,
				tileSize*float64(stubRow[i]-2)+a.untenLinks.Y,
				tileSize*float64(stubColumn[i]-1)+a.untenLinks.X-1,
				tileSize*float64(stubRow[i]-1)+a.untenLinks.Y-1)

			stubMat = stubMat.Moved(pixel.V(-float64(stubColumn[i]-2)*tileSize, -float64(stubRow[i]-2)*tileSize))
		}
	}
}

func drawCabin(can *pixelgl.Canvas, a *Arena) { // and SetCabins on kollisionsKarte !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	const numbOfCabs = 36
	var cabinRow [numbOfCabs]int
	var cabinColumn [numbOfCabs]int
	cabinCoord := a.permTiles
	cabinRow = cabinCoord[1]
	cabinColumn = cabinCoord[0]

	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	cabin := pixel.NewSprite(tilesPic, pixel.R(64, 304, 80, 288))
	cabinMat := pixel.IM
	cabinMat = cabinMat.Moved(pixel.V(tileSize/2, tileSize/2).Add(a.untenLinks))

	for i := range cabinRow {
		cabinMat = cabinMat.Moved(pixel.V(float64(cabinColumn[i]-2)*tileSize, float64(cabinRow[i]-2)*tileSize))
		cabin.Draw(can, cabinMat)
		a.kollisionsKarte[cabinRow[i]][cabinColumn[i]] = pixel.R(
			tileSize*float64(cabinColumn[i]-2)+24,
			tileSize*float64(cabinRow[i]-2)+6,
			tileSize*float64(cabinColumn[i]-1)+24-1,
			tileSize*float64(cabinRow[i]-1)+6-1)

		cabinMat = cabinMat.Moved(pixel.V(-float64(cabinColumn[i]-2)*tileSize, -float64(cabinRow[i]-2)*tileSize))

	}
}

func drawWallsNturf(can *pixelgl.Canvas) { // zeichnet die Umrandung und die Wiese
	//var winSizeX float64 = 768 // DIESE FENSTERGRÖẞE WIRD OPTIMAL AUSGEFÜLLT (bei zoomFactor 3)
	//var winSizeY float64 = 672

	var shortSideWallParts = 8
	var longSideWallParts = 12 // Warum???
	var edgeLowLeftCenterX, edgeLowLeftCenterY,
		wallLeftCenterX, wallLeftCenterY,
		edgeHiLeftCenterX, edgeHiLeftCenterY,
		hiWallCenterX, hiWallCenterY,
		edgeHiRightCenterX, edgeHiRightCenterY,
		wallRightCenterX, wallRightCenterY,
		edgeLowRightCenterX, edgeLowRightCenterY,
		loWallCenterX, loWallCenterY,
		turfCenterX, turfCenterY float64

	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	edgeLowLeft := pixel.NewSprite(tilesPic, pixel.R(288, 81, 312, 113))
	wallLeft := pixel.NewSprite(tilesPic, pixel.R(288, 114, 312, 130))
	edgeHiLeft := pixel.NewSprite(tilesPic, pixel.R(288, 114, 312, 144))
	hiWall := pixel.NewSprite(tilesPic, pixel.R(312, 136, 328, 144))
	edgeHiRight := pixel.NewSprite(tilesPic, pixel.R(344, 114, 368, 144))
	wallRight := pixel.NewSprite(tilesPic, pixel.R(344, 112, 368, 128))
	edgeLowRight := pixel.NewSprite(tilesPic, pixel.R(344, 81, 368, 115))
	loWall := pixel.NewSprite(tilesPic, pixel.R(312, 81, 328, 87))
	turf := pixel.NewSprite(tilesPic, pixel.R(112, 288, 128, 304))
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
	edgeLowLeftMat = edgeLowLeftMat.Moved(pixel.V(edgeLowLeftCenterX, edgeLowLeftCenterY+1))
	// Moved verschiebt den MatrixMITTELPUNKT, +1 in der y-Komponente, weil in tiles.png etwas mehr als 3 tiles in die Mitte passen
	wallLeftMat := pixel.IM
	wallLeftMat = wallLeftMat.Moved(pixel.V(wallLeftCenterX, 2*edgeLowLeftCenterY+wallLeftCenterY+1)) // +1 in der y-Komponente, weil in tiles.png etwas mehr als 3 tiles in die Mitte passen
	edgeHiLeftMat := pixel.IM
	edgeHiLeftMat = edgeHiLeftMat.Moved(pixel.V(edgeHiLeftCenterX, 2*edgeLowLeftCenterY+
		2*float64(shortSideWallParts)*wallLeftCenterY+edgeHiLeftCenterY+1))
	hiWallMat := pixel.IM
	hiWallMat = hiWallMat.Moved(pixel.V(2*edgeHiLeftCenterX+hiWallCenterX, 2*edgeLowLeftCenterY+
		2*wallRightCenterY*float64(shortSideWallParts)+2*edgeHiLeftCenterY-hiWallCenterY+1))
	edgeHiRightMat := pixel.IM
	edgeHiRightMat = edgeHiRightMat.Moved(pixel.V(2*edgeHiLeftCenterX+2*hiWallCenterX*float64(longSideWallParts+1)+
		edgeHiRightCenterX, 2*edgeLowRightCenterY+2*wallRightCenterY*float64(shortSideWallParts)+
		edgeHiRightCenterY+1))
	wallRightMat := pixel.IM
	wallRightMat = wallRightMat.Moved(pixel.V(2*edgeLowLeftCenterX+2*loWallCenterX*float64(longSideWallParts+1)+
		wallRightCenterX, 2*edgeLowRightCenterY+wallRightCenterY))
	edgeLowRightMat := pixel.IM
	edgeLowRightMat = edgeLowRightMat.Moved(pixel.V(2*edgeLowLeftCenterX+2*loWallCenterX*float64(longSideWallParts+1)+
		edgeLowRightCenterX, edgeLowRightCenterY+2))
	loWallMat := pixel.IM
	loWallMat = loWallMat.Moved(pixel.V(2*edgeLowLeftCenterX+loWallCenterX, loWallCenterY))
	turfMat := pixel.IM
	turfMat = turfMat.Moved(pixel.V(2*wallLeftCenterX+turfCenterX, 2*loWallCenterY+turfCenterY))

	edgeLowLeft.Draw(can, edgeLowLeftMat)
	wallLeft.Draw(can, wallLeftMat)
	edgeHiLeft.Draw(can, edgeHiLeftMat)
	hiWall.Draw(can, hiWallMat)
	edgeHiRight.Draw(can, edgeHiRightMat)
	wallRight.Draw(can, wallRightMat)
	edgeLowRight.Draw(can, edgeLowRightMat)
	loWall.Draw(can, loWallMat)

	wallLeftShift := 2 * wallLeftCenterY
	for i := 0; i < shortSideWallParts; i++ { // draws left wall
		wallLeftMat = wallLeftMat.Moved(pixel.V(0, wallLeftShift))
		wallLeft.Draw(can, wallLeftMat)
		wallRightMat = wallRightMat.Moved(pixel.V(0, wallLeftShift))
		wallRight.Draw(can, wallRightMat)
	}
	hiWallShift := 2 * hiWallCenterX
	for i := 0; i < longSideWallParts; i++ {
		hiWallMat = hiWallMat.Moved(pixel.V(hiWallShift, 0))
		hiWall.Draw(can, hiWallMat)
		loWallMat = loWallMat.Moved(pixel.V(hiWallShift, 0))
		loWall.Draw(can, loWallMat)
	}
	turfRightShift := 2 * turfCenterX
	turfUpShift := 2 * turfCenterY
	for i := 0; i <= shortSideWallParts+2; i++ { // es sind 2 Wandteile weniger als Kacheln
		turf.Draw(can, turfMat)
		for j := 0; j < longSideWallParts; j++ { // one is already drawn in the line before
			turfMat = turfMat.Moved(pixel.V(turfRightShift, 0))
			turf.Draw(can, turfMat)
		}
		turfMat = turfMat.Moved(pixel.V(float64(-(longSideWallParts))*turfRightShift, turfUpShift))
	}
}
