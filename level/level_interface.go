package level

import(	//"../arena" 
		"../tiles"
		"github.com/faiface/pixel/pixelgl")

// Konstruktor für ein leeres Levelobjekt
// NewBlankLevel ()

type Level interface {

GetTiles () []tiles.Tile

SetRandomTiles (number int)//, ar arena.Arena)

SetRandomItems (number int)

DrawTiles (win *pixelgl.Window)

DrawItems (win *pixelgl.Window)

IsTile (x,y int) bool

RemoveTile(x,y int) bool

}
