package tiles

import (
	"../animations"
	"../characters"
	"github.com/faiface/pixel"
	"time"
)

/*	Vor.: "constant" muss ein gültiger Bezeichner für ein Item sein
	Eff.: Ein "zerstörbares Tile"-Objekt mit dem Bezeichner "constant",
	das weder ein Item noch eine Bombe ist, ist an Position pos geliefert.
	*data erfüllt das Interface Item

	NewTile(constant uint8, pos pixel.Vec) *data
*/

/*	Vor.: "constant" muss ein gültiger Bezeichner für ein Item sein
	Eff.: Ein Item-Objekt mit dem Bezeichner "constant",
	das keine Bombe ist, ist an Position pos geliefert.
	*data erfüllt das Interface Item

	NewItem(constant uint8, pos pixel.Vec) *data
*/

/*	Vor.: -
	Eff.: Ein Bomben-Objekt ist geliefert.
	*bombe erfüllt das Interface Item

	NewBomb() *bombe
*/

type Tile interface {
	/*	Vor.: -
	 * 	Eff.: Die zum Objekt gehörenden Animation ist geliefert,
	 * 	in der Animation ist auch die Position des Items
	 * 	als pixel.Vec gespeichert.
	 */
	Ani() animations.Animation

	/*	Vor.: win darf nicht nil sein - das fenser win ist geöffnet
	 * 	Eff.: Das Item ist gezeichnet
	 */
	Draw(win pixel.Target)

	/*	Vor.: -
	 * 	Eff.: Die ZeichenMatrix, des Items ist geliefert
	 */
	GetMatrix() pixel.Matrix

	/*	Vor.: -
	 * 	Eff.: Der Vektor der Position das Items ist geliefert
	 */
	GetPos() pixel.Vec

	/*	Vor.: -
	 * 	Eff.: Der Typ das Items ist geliefert (siehe constants.go)
	 */
	GetType() uint8

	//GetIndexPos() (x, y int)

	/*	Vor.: -
	 * 	Eff.: true ist genu dann geliefert, wenn das Objekt sichtbar ist
	 */
	IsVisible() bool

	/*	Vor.: -
	 * 	Eff.: Das Item wird an die Position nach Vektor pos gesetzt,
	 *  ist aber noch nicht dort gezeichnet: Aufruf von Draw nötig!
	 */
	SetPos(pos pixel.Vec)

	/*	Vor.: -
	 * 	Eff.: Das Item ist nun sichtbar genau dann, wenn b=true ist
	 */
	SetVisible(b bool)
}

type Item interface {
	Tile

	/*	Vor.: -
	 * 	Eff.: Der Zeitpunkt, ab dem das Item nicht mehr existieren sollte ist geliefert.
	 */
	GetTimeStamp() time.Time

	/*	Vor.: -
	 * 	Eff.: true ist genu dann geliefert, wenn das Objekt zerstörbar ist
	 */
	IsDestroyable() bool

	/*	Vor.: -
	 * 	Eff.: Das Objekt ist entweder zerstörbar, falls true übergeben
	 * 	wurde oder unzerstörbar, falls false übergeben wurde.
	 */
	SetDestroyable(b bool)

	/* Vor.: -
	 * Eff.: Der Zeitpunkt, ab dem das Item nicht mehr existiert ist nun geändert
	 */
	SetTimeStamp(time.Time)
}

type Bombe interface {
	Item

	/*	Vor.: -
	 * 	Eff.: Wirkungsradius der Bombe ist geliefert
	 */
	GetPower() float64

	/*	Vor.: -
	 * 	Eff.: Falls das Objekt keinen Besitzer hat ist das Tupel false,nil geliefert.
	 * 	Falls das Objekt einen Besitzer hat, ist true und ein Zeiger auf den Besitzer
	 * 	geliefert.
	 */
	Owner() (bool, characters.Player)

	/*	Vor.: -
	 * 	Eff.: Die Animation der Bombe ist nun verändert
	 */
	SetAnimation(animations.Animation)

	/*	Vor.: Das Objekt hat noch keinen Besitzer.
	 * 	Eff.: Der Besitzer des Objekt ist nun der übergebene player.
	 */
	SetOwner(player characters.Player)
}
