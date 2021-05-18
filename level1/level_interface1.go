package level1

import ( //
	"../arena"
	//"../tiles"
	"github.com/faiface/pixel"
)

// Konstruktor für ein leeres Levelobjekt
// NewBlankLevel ()

type Level interface {
	A() arena.Arena

	SetRandomTilesAndItems(numberTiles, numberItems int)

	DrawColumn(y int, win pixel.Target)

	IsTile(x, y int) bool

	IsDestroyableTile(x, y int) bool

	GetPosOfNextTile(x, y int, dir pixel.Vec) (b bool, xx, yy int)

	RemoveTile(x, y int)

	RemoveItems(x, y int, dir pixel.Vec)

	/*
	   GetTiles () []tiles.Tile



	   SetRandomItems (number int)

	   DrawTiles (win *pixelgl.Window)

	   DrawItems (win *pixelgl.Window)

	   IsTile (x,y int) bool

	   RemoveTile(x,y int) bool

	   RemoveItems (x,y,dir,len int)
	*/
}
