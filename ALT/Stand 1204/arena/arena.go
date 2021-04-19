package arena

import (
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

//const TileSize float64 = 16
const WallWidth float64 = 48
const WallHeight float64 = 13

type data struct {
	canvas           *pixelgl.Canvas
	destroyableTiles [2][35]int
	matrix           pixel.Matrix
	permTiles        [2][36]int
	//tiles           [11][13]int
	lowerLeft   pixel.Vec // linke untere Spielfeldecke
	w, h        int
	passability []bool // Slice showing passability for each tile
}

func NewArena(width, height int) *data {
	var a *data = new(data)
	a.w = width
	a.h = height
	a.setPermTiles()
	a.setDestroyableTiles()
	a.passability = make([]bool, width*height)

	//fmt.Println(a.passability)
	for i := range a.passability {
		a.passability[i] = true
	}
	for i := range a.passability {
		for j := 0; j < len(a.permTiles[0]); j++ {
			if i/a.w == a.permTiles[1][j]-2 && i%a.w == a.permTiles[0][j]-2 {
				a.passability[i] = false
			}
		}
		for j := 0; j < len(a.destroyableTiles[0]); j++ {
			if i/a.w == a.destroyableTiles[1][j]-2 && i%a.w == a.destroyableTiles[0][j]-2 {
				a.passability[i] = false
			}
		}
	}
	a.lowerLeft = pixel.V(24, 6)
	a.matrix = pixel.IM
	a.matrix = a.matrix.Moved(pixel.V((float64(width)*TileSize+WallWidth)/2, (float64(height)*TileSize+WallHeight)/2))
	a.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(width)*TileSize+WallWidth, float64(height)*TileSize+WallHeight))
	a.drawWallsAndGround()
	a.drawPermTiles()
	a.drawDestroyableTiles()
	return a
}

func (a *data) GetPassability() []bool {
	return a.passability
}
func (a *data) GetCanvas() *pixelgl.Canvas {
	return a.canvas
}
func (a *data) GetFieldCoord(v pixel.Vec) (x, y int) {
	x = int(math.Trunc((v.X-a.lowerLeft.X)/TileSize))%(a.w+1) + 2
	y = int(math.Trunc((v.Y-a.lowerLeft.Y)/TileSize))%(a.h+1) + 2
	return
}
func (a *data) GetHeight() int {
	return a.h
}
func (a *data) GetLowerLeft() pixel.Vec {
	return a.lowerLeft
}
func (a *data) GetMatrix() *pixel.Matrix {
	return &(a.matrix)
}
func (a *data) GetPermTiles() [2][36]int {
	return a.permTiles
}
func (a *data) GetTileSize() float64 {
	return TileSize
}
func (a *data) GetWidth() int {
	return a.w
}
func (a *data) GrantedDirections(posBox pixel.Rect) [4]bool { // {links,rechts,oben,unten}
	var grDir [4]bool
	var x1, x2, y1, y2 int
	x1 = int(math.Trunc((posBox.Min.X-a.lowerLeft.X)/TileSize)) % (a.w + 1)
	y1 = int(math.Trunc((posBox.Min.Y-a.lowerLeft.Y)/TileSize)) % (a.h + 1)
	x2 = int(math.Trunc((posBox.Max.X-a.lowerLeft.X)/TileSize)) % (a.w + 1)
	y2 = int(math.Trunc((posBox.Max.Y-a.lowerLeft.Y)/TileSize)) % (a.h + 1)
	if posBox.Min.X-1 > a.lowerLeft.X {
		if !a.passability[(y1)*a.w+x2-1] || !a.passability[(y2)*a.w+x2-1] { // if a unpassable field is left of the posBox
			if posBox.Min.X-1 > a.lowerLeft.X+float64(x2)*TileSize {
				grDir[0] = true

			} else {
				grDir[0] = false
			}
		} else {
			grDir[0] = true
		}
	} else {
		grDir[0] = false
	}
	if posBox.Max.X+1 < a.lowerLeft.X+float64(a.w)*TileSize {
		if !a.passability[((y1)*a.w+x1+1)%143] || !a.passability[((y2)*a.w+x1+1)%143] { // if a unpassable field is left of the posBox
			if posBox.Max.X+1 < a.lowerLeft.X+float64(x1+1)*TileSize {
				grDir[1] = true
			} else {
				grDir[1] = false
			}
		} else {
			grDir[1] = true
		}
	} else {
		grDir[1] = false
	}
	if posBox.Max.Y+1 < a.lowerLeft.Y+float64(a.h)*TileSize {
		if !a.passability[((y1+1)*a.w+x1)%143] || !a.passability[((y1+1)*a.w+x2)%143] { // if a unpassable field is left of the posBox
			if posBox.Max.Y+1 < a.lowerLeft.Y+float64(y2+1)*TileSize {
				grDir[2] = true
			} else {
				grDir[2] = false
			}
		} else {
			grDir[2] = true
		}
	} else {
		grDir[2] = false
	}
	if posBox.Min.Y-1 > a.lowerLeft.Y {
		if !a.passability[modulus((y2-1)*a.w+x1)] || !a.passability[modulus((y2-1)*a.w+x2)] { // if a unpassable field is left of the posBox
			if posBox.Min.Y-1 > a.lowerLeft.Y+float64(y2)*TileSize {
				grDir[3] = true
			} else {
				grDir[3] = false
			}
		} else {
			grDir[3] = true
		}
	} else {
		grDir[3] = false
	}
	return grDir
	//return [4]bool{grDir[0],grDir[1],true,true}
}
func (a *data) RemoveTiles(x, y int) {
	k := a.checkCoordsOfDestroyables(x, y)
	if k != -1 { // 42 als Fehlerfall: diese Koordinaten wurden nicht gefunden
		a.destroyableTiles[0][k] = -1 // -1 als "nil-Koordinate"
		a.destroyableTiles[1][k] = -1
		a.passability[(y-2)*a.w+(x-2)] = true
		a.drawWallsAndGround()
		a.drawPermTiles()
		a.drawDestroyableTiles()
	}
}

//------------------------- Hilfsfunktionen ---------------------------------

func (a *data) checkCoordsOfDestroyables(x, y int) int {
	var j int
	for i := 0; i < 35; i++ {
		if a.destroyableTiles[0][i] == x && a.destroyableTiles[1][i] == y {
			j = i
			return j
		}
	}
	j = -1 // -1 als nil-Wert
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

func (a *data) checkTileStatus(x, y int) (int, int, bool) { // checkt, ob die Kachel (x,y) belegt ist
	var w bool = true
	for i := 0; i < 36; i++ {
		if a.permTiles[0][i] == x && a.permTiles[1][i] == y {
			w = false
		}
	}
	for i := 0; i < 35; i++ {
		if a.destroyableTiles[0][i] == x && a.destroyableTiles[1][i] == y {
			w = false
		}
	}
	return x, y, w
}

// Berechnet den Betrag eines int-Wertes
func modulus(x int) int {
	if x >= 0 {
		return x
	} else {
		return x * (-1)
	}
}

// Erzeugt x-y-Koordinaten der Häuser für eine Spielfeldmatrix (13x11 Matrix)
func (a *data) setPermTiles() {
	//rand.Seed(time.Now().UnixNano())
	//for i := range locations[0] {
	//	a.permTiles[0][i] = rand.Intn(13)+2
	//	a.permTiles[1][i] = rand.Intn(11)+2
	//}
	a.permTiles = [2][36]int{
		{3, 5, 7, 9, 11, 13, 10, 11, 12, 3, 5, 7, 9, 11, 12, 13, 3, 5, 7, 9, 11, 13, 3, 5, 7, 9, 11, 13, 4, 3, 5, 7, 9, 11, 13, 9},
		{3, 3, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 9, 9, 9, 9, 9, 9, 10, 11, 11, 11, 11, 11, 11, 12}}
}

// Erzeugt x-y-Koordinaten der Baumstümpfe für eine Spielfeldmatrix (13x11 Matrix)
func (a *data) setDestroyableTiles() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 35; i++ {
		for {
			x, y, w := a.checkTileStatus(rand.Intn(13)+2, rand.Intn(11)+2)
			if w {
				a.destroyableTiles[0][i] = x
				a.destroyableTiles[1][i] = y
				break
			}
		}
	}
}

func (a *data) drawDestroyableTiles() {
	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	destrSprite := pixel.NewSprite(tilesPic, pixel.R(80, 304, 96, 288))
	destrMat := pixel.IM
	destrMat = destrMat.Moved(pixel.V(TileSize/2, TileSize/2).Add(a.lowerLeft))
	for i := range a.destroyableTiles[0] {
		if a.destroyableTiles[0][i] != -1 {
			destrMat = destrMat.Moved(pixel.V(float64(a.destroyableTiles[0][i]-2)*TileSize, float64(a.destroyableTiles[1][i]-2)*TileSize))
			destrSprite.Draw(a.canvas, destrMat)
			destrMat = destrMat.Moved(pixel.V(-float64(a.destroyableTiles[0][i]-2)*TileSize, -float64(a.destroyableTiles[1][i]-2)*TileSize))
		}
	}
}

func (a *data) drawPermTiles() {
	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	permSprite := pixel.NewSprite(tilesPic, pixel.R(64, 304, 80, 288))
	permMat := pixel.IM
	permMat = permMat.Moved(pixel.V(TileSize/2, TileSize/2).Add(a.lowerLeft))
	for i := range a.permTiles[0] {
		permMat = permMat.Moved(pixel.V(float64(a.permTiles[0][i]-2)*TileSize, float64(a.permTiles[1][i]-2)*TileSize))
		permSprite.Draw(a.canvas, permMat)
		permMat = permMat.Moved(pixel.V(-float64(a.permTiles[0][i]-2)*TileSize, -float64(a.permTiles[1][i]-2)*TileSize))
	}
}

func (a *data) drawWallsAndGround() { // zeichnet die Umrandung und die Wiese
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

	edgeLowLeft.Draw(a.canvas, edgeLowLeftMat)
	wallLeft.Draw(a.canvas, wallLeftMat)
	edgeHiLeft.Draw(a.canvas, edgeHiLeftMat)
	hiWall.Draw(a.canvas, hiWallMat)
	edgeHiRight.Draw(a.canvas, edgeHiRightMat)
	wallRight.Draw(a.canvas, wallRightMat)
	edgeLowRight.Draw(a.canvas, edgeLowRightMat)
	loWall.Draw(a.canvas, loWallMat)

	wallLeftShift := 2 * wallLeftCenterY
	for i := 0; i < shortSideWallParts; i++ { // draws left wall
		wallLeftMat = wallLeftMat.Moved(pixel.V(0, wallLeftShift))
		wallLeft.Draw(a.canvas, wallLeftMat)
		wallRightMat = wallRightMat.Moved(pixel.V(0, wallLeftShift))
		wallRight.Draw(a.canvas, wallRightMat)
	}
	hiWallShift := 2 * hiWallCenterX
	for i := 0; i < longSideWallParts; i++ {
		hiWallMat = hiWallMat.Moved(pixel.V(hiWallShift, 0))
		hiWall.Draw(a.canvas, hiWallMat)
		loWallMat = loWallMat.Moved(pixel.V(hiWallShift, 0))
		loWall.Draw(a.canvas, loWallMat)
	}
	turfRightShift := 2 * turfCenterX
	turfUpShift := 2 * turfCenterY
	for i := 0; i <= shortSideWallParts+2; i++ { // es sind 2 Wandteile weniger als Kacheln
		turf.Draw(a.canvas, turfMat)
		for j := 0; j < longSideWallParts; j++ { // one is already drawn in the line before
			turfMat = turfMat.Moved(pixel.V(turfRightShift, 0))
			turf.Draw(a.canvas, turfMat)
		}
		turfMat = turfMat.Moved(pixel.V(float64(-(longSideWallParts))*turfRightShift, turfUpShift))
	}
}
