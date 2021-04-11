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
const wallWidth float64 = 48
const wallHeight float64 = 13

type data struct {
	canvas           *pixelgl.Canvas
	destroyableTiles [2][35]int
	matrix           pixel.Matrix
	permTiles        [2][36]int
	//tiles           [11][13]int
	lowerLeft     pixel.Vec    // linke untere Spielfeldecke
	w, h           int
	passability []bool // Slice showing passability for each tile
}

func NewArena(width, height int) *data {
	var a *data = new(data)
	var s []bool = make([]bool,width*height)
	a.w = width
	a.h = height
	a.permTiles = setCabins()
	a.destroyableTiles = setStub()
	a.passability = s
	//fmt.Println(a.passability)
	for i := range a.passability { a.passability[i] = true }
	for i := range a.passability {
		for j := 0; j < len(a.permTiles[0]); j++ {
			if i / a.w == a.permTiles[1][j] -2 && i % a.w == a.permTiles[0][j] -2 {
				a.passability[i] = false
			}
		}
		for j := 0; j < len(a.destroyableTiles[0]); j++ {
			if i / a.w == a.destroyableTiles[1][j] -2 && i % a.w == a.destroyableTiles[0][j] -2 {
				a.passability[i] = false
			}
		}
	}
	a.lowerLeft = pixel.V(24, 6)
	a.matrix = pixel.IM
	a.matrix = a.matrix.Moved(pixel.V((float64(width)*tileSize+wallWidth)/2, (float64(height)*tileSize+wallHeight)/2))
	a.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(width)*tileSize+wallWidth, float64(height)*tileSize+wallHeight))
	drawWallsNturf(a.canvas)
	drawCabin(a.canvas, a)
	a.drawStub(a.canvas)
	return a
}

func (a *data) GetBoolMap() []bool {
	return a.passability
}
func (a *data) GetCanvas() *pixelgl.Canvas {
	return a.canvas
}
func (a *data) GetFieldCoord(v pixel.Vec) (x, y int) {
	x = int(math.Trunc((v.X-a.lowerLeft.X)/tileSize))%(a.w+1) + 2
	y = int(math.Trunc((v.Y-a.lowerLeft.Y)/tileSize))%(a.h+1) + 2
	return
}
func (a *data) GetHeight() int {
	return a.h
}
func (a *data) GetMatrix() *pixel.Matrix {
	return &(a.matrix)
}
func (a *data) GetTileSize() float64 {
	return tileSize
}
func (a *data) GetWidth() int {
	return a.w
}
func (a *data) GrantedDirections(posBox pixel.Rect) [4]bool { // {links,rechts,oben,unten}
	var grDir [4]bool
	var x1, x2, y1, y2 int
	x1 = int(math.Trunc((posBox.Min.X-a.lowerLeft.X)/tileSize))%(a.w+1)
	y1 = int(math.Trunc((posBox.Min.Y-a.lowerLeft.Y)/tileSize))%(a.h+1)
	x2 = int(math.Trunc((posBox.Max.X-a.lowerLeft.X)/tileSize))%(a.w+1)
	y2 = int(math.Trunc((posBox.Max.Y-a.lowerLeft.Y)/tileSize))%(a.h+1)
	if posBox.Min.X-1 > a.lowerLeft.X {
		if !a.passability[(y1)*a.w+x2-1] || !a.passability[(y2)*a.w+x2-1] { // if a unpassable field is left of the posBox
			if posBox.Min.X-1 > a.lowerLeft.X+float64(x2)*tileSize {
				grDir[0] = true

			} else {
				grDir[0] = false
			}
		}else{
			grDir[0] = true
		}
	} else {
		grDir[0] = false
	}
	if posBox.Max.X+1 < a.lowerLeft.X + float64(a.w)*tileSize {
		if !a.passability[((y1)*a.w+x1+1) % 143] || !a.passability[((y2)*a.w+x1+1) % 143] { // if a unpassable field is left of the posBox
			if posBox.Max.X+1 < a.lowerLeft.X+float64(x1+1)*tileSize {
				grDir[1] = true
			} else {
				grDir[1] = false
			}
		}else{
			grDir[1] = true
		}
	} else {
		grDir[1] = false
	}
	if posBox.Max.Y+1 < a.lowerLeft.Y + float64(a.h)*tileSize {
		if !a.passability[((y1+1)*a.w+x1) % 143] || !a.passability[((y1+1)*a.w+x2) % 143] { // if a unpassable field is left of the posBox
			if posBox.Max.Y+1 < a.lowerLeft.Y+float64(y2+1)*tileSize {
				grDir[2] = true
			} else {
				grDir[2] = false
			}
		}else{
			grDir[2] = true
		}
	} else {
		grDir[2] = false
	}
	if posBox.Min.Y-1 > a.lowerLeft.Y {
		if !a.passability[modulus((y2-1)*a.w+x1)] || !a.passability[modulus((y2-1)*a.w+x2)] { // if a unpassable field is left of the posBox
			if posBox.Min.Y-1 > a.lowerLeft.Y+float64(y2)*tileSize {
				grDir[3] = true
			} else {
				grDir[3] = false
			}
		}else{
			grDir[3] = true
		}
	} else {
		grDir[3] = false
	}
	return grDir
	//return [4]bool{grDir[0],grDir[1],true,true}
}
func (a *data) RemoveTiles(x, y int) {
	k := checkCoordsOfDestroyables(x, y, a.destroyableTiles)
	if k != -1 { // 42 als Fehlerfall: diese Koordinaten wurden nicht gefunden
		a.destroyableTiles[0][k] = -1	// -1 als "nil-Koordinate"
		a.destroyableTiles[1][k] = -1
		a.passability[(y-2)*a.w + (x-2)] = true
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
		}
	}
	j = -1	// -1 als nil-Wert
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
// Berechnet den Betrag eines int-Wertes
func modulus (x int) int {
	if x >= 0 {
		return x
	}else{
		return x * (-1)
	}
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

func (a *data) drawStub(can *pixelgl.Canvas) { // and SetStubs on kollisionsKarte !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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
	stubMat = stubMat.Moved(pixel.V(tileSize/2, tileSize/2).Add(a.lowerLeft))
	for i := 0; i < 35; i++ {
		if stubColumn[i] != -1 {
			stubMat = stubMat.Moved(pixel.V(float64(stubColumn[i]-2)*tileSize, float64(stubRow[i]-2)*tileSize))
			stub.Draw(can, stubMat)
			stubMat = stubMat.Moved(pixel.V(-float64(stubColumn[i]-2)*tileSize, -float64(stubRow[i]-2)*tileSize))
		}
	}
}

func drawCabin(can *pixelgl.Canvas, a *data) {
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
	cabinMat = cabinMat.Moved(pixel.V(tileSize/2, tileSize/2).Add(a.lowerLeft))
	for i := range cabinRow {
		cabinMat = cabinMat.Moved(pixel.V(float64(cabinColumn[i]-2)*tileSize, float64(cabinRow[i]-2)*tileSize))
		cabin.Draw(can, cabinMat)
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
