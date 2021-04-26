package arena

import (
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
	"math"
	"os"
)

type data struct {
	canvas      *pixelgl.Canvas
	lowerLeft pixel.Vec // linke untere Spielfeldecke für korrekte Positionsbestimmung
	matrix    pixel.Matrix
	typ       int
	passably  []bool // Slice showing passably for each tile
	permTiles [2][]int
	w, h      int
}

func NewArena(typ, width, height int) *data {
	var a = new(data)
	a.w = width
	a.h = height
	a.typ = typ
	a.setPermTiles()
	a.passably = make([]bool, width*height)
	for i := range a.passably {
		a.passably[i] = true
	}
	for i := 0; i < len(a.permTiles[0]); i++ {
		a.passably[a.w * a.permTiles[1][i] + a.permTiles[0][i]] = false
	}
	switch a.typ {
	case MfS:
		a.lowerLeft = pixel.V(24,8)
	case TurfNtrees:
		a.lowerLeft = pixel.V(24, 8)
	case Castle:
		a.lowerLeft = pixel.V(24, 8)
	}
	a.matrix = pixel.IM
	a.matrix = a.matrix.Moved(pixel.V((float64(width)*TileSize+WallWidth)/2-TileSize/4, (float64(height)*TileSize+WallHeight)/2-TileSize/2))
	a.canvas = pixelgl.NewCanvas(pixel.R(-2*TileSize, -2*TileSize, float64(width)*TileSize+WallWidth + TileSize/2, float64(height)*TileSize+WallHeight))
	a.drawWallsAndGround()
	a.drawPermTiles()
	return a
}

func (a *data) GetCanvas() *pixelgl.Canvas {
	return a.canvas
}
func (a *data) GetFieldCoord(v pixel.Vec) (x, y int) {
	x = int(math.Trunc((v.X - a.lowerLeft.X)/TileSize))%(a.w+1)
	y = int(math.Trunc((v.Y - a.lowerLeft.Y)/TileSize))%(a.h+1)
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
func (a *data) GetPassability() []bool {
	return a.passably
}
func (a *data) GetPermTiles() [2][]int {
	return a.permTiles
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
		if !a.passably[(y1)*a.w+x2-1] || !a.passably[(y2)*a.w+x2-1] { // if a unpassable field is left of the posBox
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
		if !a.passably[((y1)*a.w+x1+1)% (a.w * a.h)] || !a.passably[((y2)*a.w+x1+1)% (a.w * a.h)] { // if a unpassable field is left of the posBox
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
		if !a.passably[((y1+1)*a.w+x1)% (a.w * a.h)] || !a.passably[((y1+1)*a.w+x2)% (a.w * a.h)] { // if a unpassable field is left of the posBox
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
		if !a.passably[modulus((y2-1)*a.w+x1)] || !a.passably[modulus((y2-1)*a.w+x2)] { // if a unpassable field is left of the posBox
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
func (a *data) IsFreeTile(x,y int) bool {
	return a.passably[a.w*y+x]
}
func (a *data) IsTile(x, y int) bool {
	return !a.IsFreeTile(x,y)
}

//------------------------- Hilfsfunktionen ---------------------------------

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

// Berechnet den Betrag eines int-Wertes
func modulus(x int) int {
	if x >= 0 {
		return x
	} else {
		return x * (-1)
	}
}

func (a *data) setPermTiles() {
	var permTilesDefault [2][]int
	switch a.typ {
	case MfS:
		permTilesDefault = [2][]int{
			{0,1,1,1,1,1,1,2,3,3,3,3,3,3,3,3,4,5,5,5,5,5,5,5,7,7,7,7,7,7,7,9,9,9,9,9,9,10},
			{2,1,3,5,7,9,11,11,0,1,3,5,7,9,10,11,9,1,2,3,5,7,8,9,11,1,3,5,7,8,9,11,1,3,5,7,9,11,10}}
	case TurfNtrees:
		permTilesDefault = [2][]int{
			{1, 3, 5, 7, 9, 11, 8, 9, 10, 1, 3, 5, 7, 9, 10, 11, 1, 3, 5, 7, 9, 11, 1, 3, 5, 7, 9, 11, 2, 1, 3, 5, 7, 9, 11, 7},
			{1, 1, 1, 1, 1, 1, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 5, 5, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 8, 9, 9, 9, 9, 9, 9, 10}}
	case Castle:
		permTilesDefault = [2][]int{
		{1,3,5,7,9,11,2,1,3,5,7,9,11,12,1,3,5,6,7,8,9,11,2,1,3,5,7,9,11,5,11,1,3,5,7,9,11},
		{1,1,1,1,1,1,2,3,3,3,3,3,3,4,5,5,5,5,5,5,5,5,6,7,7,7,7,7,7,8,8,9,9,9,9,9,9}}
	}
	// Feld kleiner oder gleich groß wie das Standardfeld
	for j := range permTilesDefault[0] {
		if permTilesDefault[0][j] < a.w && permTilesDefault[1][j] < a.h {
			a.permTiles[0] = append(a.permTiles[0], permTilesDefault[0][j])
			a.permTiles[1] = append(a.permTiles[1], permTilesDefault[1][j])
		}
	}
	if a.w > 13 || a.h > 11 { // Feld größer als das Standardfeld
		j := 36
		for {
			if a.permTiles[0][j-36] + 12 < a.w - 1 {
				a.permTiles[0] = append(a.permTiles[0], a.permTiles[0][j-36]+12) // permTiles rechts vom Standardfeld
				a.permTiles[1] = append(a.permTiles[1], a.permTiles[1][j-36])
			}
			if a.permTiles[1][j-36] + 10 < a.h - 1 {
				a.permTiles[0] = append(a.permTiles[0], a.permTiles[0][j-36]) // permTiles oberhalb vom Standardfeld
				a.permTiles[1] = append(a.permTiles[1], a.permTiles[1][j-36]+10)
			}
			if a.permTiles[0][j-36] + 12 < a.w - 1 && a.permTiles[1][j-36] + 10 < a.h - 1 {
				a.permTiles[0] = append(a.permTiles[0], a.permTiles[0][j-36]+12) // permTiles rechts oberhalb vom Standardfeld
				a.permTiles[1] = append(a.permTiles[1], a.permTiles[1][j-36]+10)
				j++
			}else{
				j++
			}
			if j >= len(a.permTiles[0])+36{break}
		}
	}
}

func (a *data) drawPermTiles() {
	var permSprite *pixel.Sprite
	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	switch a.typ {
	case MfS:
		permSprite = pixel.NewSprite(tilesPic, pixel.R(8*TileSize, 18*TileSize, 9*TileSize, 19*TileSize))
	case TurfNtrees:
		permSprite = pixel.NewSprite(tilesPic, pixel.R(4*TileSize, 18*TileSize, 5*TileSize, 19*TileSize))
	case Castle:
		permSprite = pixel.NewSprite(tilesPic, pixel.R(28*TileSize, 9*TileSize, 29*TileSize, 10*TileSize))
	}
	permMat := pixel.IM
	permMat = permMat.Moved(pixel.V(TileSize + TileSize/2, TileSize/2))
	for i := range a.permTiles[0] {
		permMat = permMat.Moved(pixel.V(float64(a.permTiles[0][i])*TileSize, float64(a.permTiles[1][i])*TileSize))
		permSprite.Draw(a.canvas, permMat)
		permMat = permMat.Moved(pixel.V(-float64(a.permTiles[0][i])*TileSize, -float64(a.permTiles[1][i])*TileSize))
	}
}

func (a *data) drawWallsAndGround() { // baut Arena spaltenweise auf, beginnt unten links
	var edgeLowLeft, wallLeft, edgeHiLeft, hiWall, edgeHiRight, wallRight, edgeLowRight, loWall, ground *pixel.Sprite
	tilesPic, err := loadPicture("graphics/tiles.png")
	if err != nil {
		panic(err)
	}
	switch a.typ {
	case MfS:
		edgeLowLeft = pixel.NewSprite(tilesPic, pixel.R(18*TileSize, 1*TileSize, 20*TileSize, 3*TileSize)) // Default-Sprites
		wallLeft = pixel.NewSprite(tilesPic, pixel.R(18*TileSize, 3*TileSize, 20*TileSize, 4*TileSize))
		edgeHiLeft = pixel.NewSprite(tilesPic, pixel.R(18*TileSize, 3*TileSize, 20*TileSize, 5*TileSize))
		hiWall = pixel.NewSprite(tilesPic, pixel.R(20*TileSize, 4*TileSize, 21*TileSize, 5*TileSize))
		edgeHiRight = pixel.NewSprite(tilesPic, pixel.R(21*TileSize, 3*TileSize, 23*TileSize, 5*TileSize))
		wallRight = pixel.NewSprite(tilesPic, pixel.R(21*TileSize, 3*TileSize, 23*TileSize, 4*TileSize))
		edgeLowRight = pixel.NewSprite(tilesPic, pixel.R(21*TileSize, 1*TileSize, 23*TileSize, 3*TileSize))
		loWall = pixel.NewSprite(tilesPic, pixel.R(20*TileSize, 1*TileSize, 21*TileSize, 2*TileSize))
		ground = pixel.NewSprite(tilesPic, pixel.R(10*TileSize, 18*TileSize, 11*TileSize, 19*TileSize))
	case TurfNtrees:
		edgeLowLeft = pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 3*TileSize, 26*TileSize, 5*TileSize)) // Default-Sprites
		wallLeft = pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 5*TileSize, 26*TileSize, 6*TileSize))
		edgeHiLeft = pixel.NewSprite(tilesPic, pixel.R(24*TileSize, 6*TileSize, 26*TileSize, 8*TileSize))
		hiWall = pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 7*TileSize, 27*TileSize, 8*TileSize))
		edgeHiRight = pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 6*TileSize, 29*TileSize, 8*TileSize))
		wallRight = pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 5*TileSize, 29*TileSize, 6*TileSize))
		edgeLowRight = pixel.NewSprite(tilesPic, pixel.R(27*TileSize, 3*TileSize, 29*TileSize, 5*TileSize))
		loWall = pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 3*TileSize, 27*TileSize, 4*TileSize))
		ground = pixel.NewSprite(tilesPic, pixel.R(7*TileSize, 18*TileSize, 8*TileSize, 19*TileSize))
	case Castle:
		edgeLowLeft = pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 8*TileSize, 28*TileSize, 10*TileSize)) // Default-Sprites
		wallLeft = pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 10*TileSize, 28*TileSize, 11*TileSize))
		edgeHiLeft = pixel.NewSprite(tilesPic, pixel.R(26*TileSize, 11*TileSize, 28*TileSize, 13*TileSize))
		hiWall = pixel.NewSprite(tilesPic, pixel.R(28*TileSize, 12*TileSize, 29*TileSize, 13*TileSize))
		edgeHiRight = pixel.NewSprite(tilesPic, pixel.R(29*TileSize, 11*TileSize, 31*TileSize, 13*TileSize))
		wallRight = pixel.NewSprite(tilesPic, pixel.R(29*TileSize, 10*TileSize, 31*TileSize, 11*TileSize))
		edgeLowRight = pixel.NewSprite(tilesPic, pixel.R(29*TileSize, 8*TileSize, 31*TileSize, 10*TileSize))
		loWall = pixel.NewSprite(tilesPic, pixel.R(28*TileSize, 8*TileSize, 29*TileSize, 9*TileSize))
		ground = pixel.NewSprite(tilesPic, pixel.R(28*TileSize, 11*TileSize, 29*TileSize, 12*TileSize))
	}
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
				ground.Draw(a.canvas, drawMat)
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
