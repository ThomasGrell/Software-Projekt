package gameStat

import ( //
	"../arena"
	//"../tiles"
	"github.com/faiface/pixel"
)

// Konstruktor für ein zufälliges GemeStat Objekt für ein Spielfeld mit den Maßen width*height und für eine Spielerzahl von "anzPlayer". 
// NewRandomGameStat (width, height int, anzPlayer uint8) *gs

// Konstruktor für ein GameStat Objekt für ein Spielfeld, das durch das dem Kondtruktor übergebene level Objekt definiert ist. 
// Die Anzahl der Spieler wird mit anzPLayer festegelgt.
// NewGameStat(lv level.Level, anzPlayer uint8) *gs

// gs erfüllt das Interface GameStat

type GameStat interface {

	// Vor.: -
	// Eff.: Liefert die zum GameStat bzw. Spielfeld gehörende Arena
	A() arena.Arena

	// Vor.: -
	// Eff.: Wenn sich an der Stelle x,y im Spielfeld ein Item befindet,
	//       wird dessen Typ sowie b=true zurückgegeben und das Item gelöscht.
	//       Befindet sich kein Item an der Stelle, so werden typ=0 und b=false
	//       zurückgegeben.
	CollectItem(x, y int) (typ uint8, b bool)

	// Vor.: -
	// Eff.: Die Zeile y des Spielfeldes mit allen zerstörbaren Teilen und Items ist gezeichnet.
	DrawColumn(y int, win pixel.Target)

	// Vor.: -
	// Eff.: Die Positionskoordinaten (xx,yy) des nächsten Zerstörbaren Teils in Richtug dir
	//       ausgehend von der aktuellen Position (x,y) ist geliefert, falls es ein solches 
	//       gibt. In dem Fall wird true , xx , yy zurück gegeben. Falls es kein zerstörbarens
	//       Teil gibt, wird false, -1 ,-1 zurück gegeben.
	GetPosOfNextTile(x, y int, dir pixel.Vec) (b bool, xx, yy int)

	// Vor.: -
	// Eff.: true wird genau dann zurück gegeben, wenn sich an den übergebenen Koordinaten (x,y) 
	//       ein zerstörbares Teil befindet. Falls nicht wird false zurück gegeben.
	IsDestroyableTile(x, y int) bool

	// Vor.: -
	// Eff.: true wird genau dann zurück gegeben, wenn sich an den übergebenen Koordinaten (x,y) 
	//       ein unzerstörbares Teil befindet. Falls nicht wird false zurück gegeben.
	IsUndestroyableTile(x, y int) bool

	// Vor.: -
	// Eff.: true wird genau dann zurück gegeben, wenn sich an den übergebenen Koordinaten (x,y) 
	//       ein zerstörbares oder unzerstörbares Teil befindet. Falls nicht wird false zurück gegeben.
	IsTile(x, y int) bool

	// Vor.: -
	// Eff.: Falls ein Item an den übergebenen Koordinaten lag, ist dieses nun zerstört und die 
	//       entsprechnde Animation beginnt.
	RemoveItems(x, y int, dir pixel.Vec)

	// Vor.: -
	// Eff.: Falls ein zerstörbares Teil an den übergebenen Koordinaten lag, ist dieses nun zerstört und die 
	//       entsprechnde Animation beginnt.
	RemoveTile(x, y int)

	// Vor.: -
	// Eff.: Die Breite und Höhe des durch GameStat defierten Spielfeldes ist als Paar Breite, Höhe geliefert.
	GetBounds() (int, int)

	// Vor.: -
	// Eff.: Das Gamestat Objekt ist auf seinen Startzustand zurücg gesetzt: Wurden im verlaufe des Spiels 
	//       Items und Teile zerstört, sind diese wieder da. Die Items liegen aber nicht mehr an der selben Stelle.
	Reset ()
}
