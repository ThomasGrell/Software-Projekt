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
//const WallWidth float64 = 48
//const WallHeight float64 = 13

type data struct {
	canvas           *pixelgl.Canvas
	destroyableTiles [][]int
	matrix           pixel.Matrix
	permTiles        [2][36]int
	//tiles           [11][13]int
	lowerLeft   pixel.Vec // linke untere Spielfeldecke für korrekte Positionsbestimmung
	w, h        int
	passability []bool // Slice showing passability for each tile
}

func NewArena(width, height int) *data {
	var a *data = new(data)
	a.destroyableTiles = make([][]int, 2)
	for i := 0; i < 2; i++ {
		a.destroyableTiles[i] = make([]int, 35)
	}

	a.w = width
	a.h = height
	a.setPermTiles()
	a.setDestroyableTiles() // Reihenfolge!!!
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
	a.matrix = a.matrix.Moved(pixel.V((float64(width)*TileSize+WallWidth)/2-TileSize/4, (float64(height)*TileSize+WallHeight)/2-TileSize/2))
	a.canvas = pixelgl.NewCanvas(pixel.R(-2*TileSize, -2*TileSize, float64(width)*TileSize+WallWidth + TileSize/2, float64(height)*TileSize+WallHeight))
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
func (a *data) GetDestroyableTiles() [][]int {
	return a.destroyableTiles[:][:]
}
func (a *data) GetFieldCoord(v pixel.Vec) (x, y int) {
	x = int(math.Trunc((v.X - a.lowerLeft.X)/TileSize))%(a.w+1) + 2
	y = int(math.Trunc((v.Y - a.lowerLeft.Y)/TileSize))%(a.h+1) + 2
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
	//return [4]bool{true,true,true,true}
}
func (a *data) IsTile(x, y int) bool {
	var w bool
	for i := 0; i < 36; i++ {
		if a.permTiles[0][i] == x && a.permTiles[1][i] == y {
			return true
		}
	}
	for i := 0; i < 35; i++ {
		if a.destroyableTiles[0][i] == x && a.destroyableTiles[1][i] == y {
			return true
		}
	}
	return w
}
func (a *data) RemoveTiles(x, y int) bool {
	var b bool
	k := a.checkCoordsOfDestroyables(x, y)
	if k != -1 { // -1 als Fehlerfall: diese Koordinaten wurden nicht gefunden
		a.destroyableTiles[0][k] = -1 // -1 als "nil-Koordinate"
		a.destroyableTiles[1][k] = -1
		a.passability[(y-2)*a.w+(x-2)] = true
		a.drawWallsAndGround()
		a.drawPermTiles()
		a.drawDestroyableTiles()
		b = true
	}
	return b
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
	for i := 0; i < 36; i++ {
		if a.permTiles[0][i] == x && a.permTiles[1][i] == y {
			return x, y, false
		}
	}
	for i := 0; i < 35; i++ {
		if a.destroyableTiles[0][i] == x && a.destroyableTiles[1][i] == y {
			return x, y, false
		}
	}
	return x, y, true
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
				//v := make([]int,2)
				//v[0] = x
				//v[1] = y
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
	destrMat = destrMat.Moved(pixel.V(TileSize + TileSize/2, TileSize/2))
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
	permMat = permMat.Moved(pixel.V(TileSize + TileSize/2, TileSize/2))
	for i := range a.permTiles[0] {
		permMat = permMat.Moved(pixel.V(float64(a.permTiles[0][i]-2)*TileSize, float64(a.permTiles[1][i]-2)*TileSize))
		permSprite.Draw(a.canvas, permMat)
		permMat = permMat.Moved(pixel.V(-float64(a.permTiles[0][i]-2)*TileSize, -float64(a.permTiles[1][i]-2)*TileSize))
	}
}

func (a *data) drawWallsAndGround() { // baut Arena spaltenweise auf, beginnt unten links
	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	edgeLowLeft := pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 3*TileSize, 26*TileSize, 5*TileSize))
	wallLeft := pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 5*TileSize, 26*TileSize, 6*TileSize))
	edgeHiLeft := pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 6*TileSize, 26*TileSize, 8*TileSize))
	hiWall := pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 7*TileSize, 27*TileSize, 8*TileSize))
	edgeHiRight := pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 6*TileSize, 29*TileSize, 8*TileSize))
	wallRight := pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 5*TileSize, 29*TileSize, 6*TileSize))
	edgeLowRight := pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 3*TileSize, 29*TileSize, 5*TileSize))
	loWall := pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 3*TileSize, 27*TileSize, 4*TileSize))
	turf := pixel.NewSprite(tilesPic, pixel.R(7*TileSize, 18*TileSize, 8*TileSize, 19*TileSize))
	drawMat := pixel.IM
	edgeLowLeft.Draw(a.canvas, /*edgeLowLeftMat*/ drawMat)
	drawMat = drawMat.Moved(pixel.V(0,TileSize + TileSize/2))
	wallLeft.Draw(a.canvas, drawMat)
	for i:=0; i<a.h-3; i++ { // -3 weil die beiden Ecken schon Felder sind und ein Wandstück bereits gezeichnet wurde
		drawMat = drawMat.Moved(pixel.V(0,TileSize))
		wallLeft.Draw(a.canvas, drawMat)
	}
	drawMat = drawMat.Moved(pixel.V(0,TileSize+TileSize/2))
	edgeHiLeft.Draw(a.canvas, drawMat)
	for j:=0; j<a.w; j++ {
		if j==0 {
			drawMat = drawMat.Moved(pixel.V(TileSize+TileSize/2, -(TileSize*float64(a.h) + TileSize/2)))
		}else{
			drawMat = drawMat.Moved(pixel.V(TileSize, -(TileSize*float64(a.h) + TileSize)))
		}
		for i := 0; i < a.h+2; i++ { // +2 weil oben und unten Wände sind
			if i == 0 {
				loWall.Draw(a.canvas, drawMat)
			} else if i < a.h+1 {
				drawMat = drawMat.Moved(pixel.V(0, TileSize))
				turf.Draw(a.canvas, drawMat)
			} else {
				drawMat = drawMat.Moved(pixel.V(0, TileSize))
				hiWall.Draw(a.canvas, drawMat)
			}

		}
	}
	drawMat = drawMat.Moved(pixel.V(TileSize+TileSize/2, -(TileSize*float64(a.h) + TileSize/2)))
	edgeLowRight.Draw(a.canvas,drawMat)
	drawMat = drawMat.Moved(pixel.V(0,TileSize + TileSize/2))
	wallRight.Draw(a.canvas, drawMat)
	for i:=0; i<a.h-3; i++ { // -3 weil die beiden Ecken schon Felder sind und ein Wandstück bereits gezeichnet wurde
		drawMat = drawMat.Moved(pixel.V(0,TileSize))
		wallRight.Draw(a.canvas, drawMat)
	}
	drawMat = drawMat.Moved(pixel.V(0,TileSize+TileSize/2))
	edgeHiRight.Draw(a.canvas, drawMat)

}
