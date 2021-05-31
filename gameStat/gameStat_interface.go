package gameStat

import ( //
	"../arena"
	//"../tiles"
	"github.com/faiface/pixel"
)

// Konstruktor für ein leeres Levelobjekt
// NewBlankLevel ()

type GameStat interface {

	// Vor.: -
	// Eff.: Liefert die zum Level gehörende Arena
	A() arena.Arena

	// Vor.: -
	// Eff.: Wenn sich an der Stelle x,y im Spielfeld ein Item befindet,
	//       wird dessen Typ sowie b=true zurückgegeben und das Item gelöscht.
	//       Befindet sich kein Item an der Stelle, so werden typ=0 und b=false
	//       zurückgegeben.
	CollectItem(x, y int) (typ uint8, b bool)

	// Vor.: -
	// Eff.:
	DrawColumn(y int, win pixel.Target)

	// Vor.: -
	// Eff.:
	GetPosOfNextTile(x, y int, dir pixel.Vec) (b bool, xx, yy int)

	// Vor.: -
	// Eff.:
	IsDestroyableTile(x, y int) bool

	// Vor.: -
	// Eff.:
	IsUndestroyableTile(x, y int) bool

	// Vor.: -
	// Eff.:
	IsTile(x, y int) bool

	// Vor.: -
	// Eff.:
	RemoveItems(x, y int, dir pixel.Vec)

	// Vor.: -
	// Eff.:
	RemoveTile(x, y int)

	// Vor.: -
	// Eff.:
	GetBounds() (int, int)

	// Vor.: -
	// Eff.: 
	Reset ()
}
