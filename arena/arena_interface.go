package arena

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// NewArena(width,heigth float64) *data

type Arena interface {

	//Vor.: -
	//Erg.: Zu dem Spielfeld mit den Koordinaten x,y wird der Feldmittelpunkt
	//      als Vektor geliefert.
	CoordToVec(x, y int) pixel.Vec

	//Vor.: /
	//Erg.: Die grafische Darstellung des Spielfelds inklusive Umrandung, Untergrund, permanenten und
	//      zertörbaren Mauern ist geliefert.
	GetCanvas() *pixelgl.Canvas

	//Vor.: /
	//Erg.: Die Koordinaten des Feldes, in dem sich der Punkt v befindet, sind geliefert.
	GetFieldCoord(v pixel.Vec) (x, y int)

	//Vor.: /
	//Erg.: Die Höhe des Spielfelds = Anzahl der Felder in senkrechter Richtung ist geliefert.
	GetHeight() int

	//Vor.: /
	//Erg.: Die Pixelkoordinaten der linken unteren Spielfeldecke sind geliefert.
	GetLowerLeft() pixel.Vec

	//Vor.: /
	//Erg.: Der Zeiger auf die Matrix, die der Arena zugrunde liegt, ist geliefert.
	GetMatrix() *pixel.Matrix

	//Vor.: /
	//Erg.: Ein Slice der Länge (Spielfeldbreite x Spielfeldhöhe) ist geliefert. Jeder Eintrag repräsentiert die
	//		Betretbarkeit einer Spielfeldkachel.
	GetPassability() []bool

	//Vor.: /
	//Erg.: Das Array, das die Koordinaten der unzerstörbaren Kacheln enthält ist geliefert.
	GetPermTiles() [2][]int

	//Vor.: /
	//Erg.: Die Breite des Spielfelds = Anzahl der Felder in waagerechter Richtung ist geliefert.
	GetWidth() int

	//Vor.: /
	//Erg.: Ein Feld mit 4 Wahrheitswerten ist geliefert. Jeder der Werte gibt für eine Richtung an, ob Laufen
	//		in diese Richtung erlaubt ist. Die Reihenfolge der Richtungen ist: links - rechts - oben - unten.
	GrantedDirections(posBox pixel.Rect) [4]bool

	//Vor.: /
	//Erg.: Falls auf der Kachel mit den Koordinaten x,y ein permanentes Hindernis steht, so ist true geliefert, sonst
	//		false.
	IsTile(x, y int) bool

	//Vor.: /
	//Erg.: Falls auf der Kachel mit den Koordinaten x,y ein permanentes Hindernis steht, so ist false geliefert, sonst
	//		true.
	IsFreeTile(x, y int) bool

}
