package arena

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

/* NewArena(width,heigth float64) *data */

type Arena interface {

	GetBoolMap() [15][17]bool

	GetCanvas() *pixelgl.Canvas

	GetFieldCoord(v pixel.Vec) (x, y int)

	GetMatrix() *pixel.Matrix

	GeTileSize() float64

	GrantedDirections(posBox pixel.Rect) [4]bool

	RemoveTiles(x, y int)

}
