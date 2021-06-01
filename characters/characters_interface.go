package characters

import (
	"../animations"
	"github.com/faiface/pixel"
)

/*
 _______________________________________________
< Implementiert von Rayk von Ende               >
< Ergänzung der pixel.Matrix sowie der          >
< Methoden GetMatrix() und Move(): Thomas Grell >
 -----------------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

****************************************************************************************

Konstruktoren:

Vor.: -
Eff.: Ein neuer Bomberman vom Typ t wird zurückgegeben.
      Die Typen sind in constants.go definiert.
      Der Bomberman ist mit Standardwerten initialisiert.
NewPlayer(t uint8) Player

Vor.: -
Eff.: Ein neues Monster vom Typ t wird zurückgegeben.
      Die Typen sind in constants.go definiert.
      Das Monster ist mit Standardwerten initialisiert.
NewEnemy(t uint8) Enemy

****************************************************************************************

*/

// Alle Methoden, die nur Monster besitzen.
type Enemy interface {
	Character

	// Vor.: -
	// Eff.: Je höher der gelieferte Wert, desto unwahrscheinlicher ist ein Richtungswechsel des Enemy.
	GetBehaviour() uint8

	// Vor.: -
	// Eff.: Liefert true, wenn das Monster die Spieler verfolgen kann, sonst false.
	IsFollowing() bool

	// Vor.: val > 0
	// Eff.: Setzt einen Wert für die Wahrscheinlichkeit eines Richtungswechsels.
	//       Je höher, desto unwahrscheinlicher ist ein Wechsel.
	SetBehaviour(val uint8)
}

// Alle Methoden, die nur Spieler besitzen.
type Player interface {
	Character

	// Vor.: -
	// Eff.: Erhöht den Punktestand des Spielers um den Wert points.
	AddPoints(points uint32)

	// Vor.: -
	// Eff.: Reduziert den Zähler für aktuell gelegte Bomben.
	DecBombs()

	// Vor.: -
	// Eff.: Reduziert den Explosionsradius von Bomben.
	DecPower()

	// Vor.: -
	// Eff.: Liefert die Anzahl gelegter Bomben.
	GetBombs() uint8

	// Vor.: -
	// Eff.: Liefert die maximale Anzahl an Bomben, die der Spieler legen kann.
	GetMaxBombs() uint8

	// Vor.: -
	// Eff.: Liefert den Sprengradius der Bomben eines Spielers
	GetPower() uint8

	// Vor.: -
	// Eff.: Liefert die Anzahl an Siegen eines Spielers im Multiplayer-Modus.
	GetWins() uint8

	// Vor.: -
	// Eff.: Gibt Auskunft, ob der Bomberman eine Fernzündung durchführen darf
	HasRemote() bool

	// Vor.: -
	// Eff.: Erhöht den Zähler für aktuell gelegte Bomben.
	IncBombs()

	// Vor.: -
	// Eff.: Erhöht die verbleibenden Leben eines Spielers um 1.
	IncLife()

	// Vor.: -
	// Eff.: Erhöht die maximale Anzahl legbarer Bomben eines Spielers um 1.
	IncMaxBombs()

	// Vor.: -
	// Eff.: Erhöht den Sprengradius der Bomben eines Spielers um 1.
	IncPower()

	// Vor.: -
	// Eff.: Erhöht die Anzahl an Siegen eines Spielers im Multiplayer-Modus um 1.
	IncWins()

	// Vor.: -
	// Eff.: Setzt den Bomberman auf Standardeinstellungen zurück.
	Reset()

	// Vor.: -
	// Eff.: Setzt die Anzahl der Siege auf 0 zurück.
	ResetWins()

	// Vor.: -
	// Eff.: Setzt die Anzahl der verbleibenden Leben eines Spielers auf den Wert lifes.
	SetLife(lifes uint8)

	// Vor.: -
	// Eff.: Setzt die maximale Anzahl legbarer Bomben eines Spielers auf den Wert maxBombs.
	SetMaxBombs(maxBombs uint8)

	// Vor.: -
	// Eff.: Legt fest, ob der Spieler sterblich ist (true) oder unsterblich (false).
	SetMortal(bool)

	// Vor.: -
	// Eff.: Legt den Explosionsradius der Bomben fest.
	SetPower(uint8)

	// Vor.: -
	// Eff.: Legt fest, ob der Bomberman seine Bomben fernzünden kann.
	SetRemote(bool)

	// Vor.: -
	// Eff.: Legt fest, ob der Spieler durch Wände laufen kann (true) oder nicht (false).
	SetWallghost(bool)
}

// Gemeinsame Methoden von Spielern und Monstern
type Character interface {

	// Vor.: -
	// Eff.: Liefert das zugehörige Animationsobjekt des Characters
	Ani() animations.Animation

	// Vor.: -
	// Eff.: Reduziert die verbleibenden Leben eines Characters um 1.
	DecLife()

	// Vor.: -
	// Eff.: Reduziert die Geschwindigkeit eines Characters um einen festen Wert.
	DecSpeed()

	// Vor.: -
	// Eff.: Zeichnet den Sprite in das angegebene Target (Batch, Canvas oder Window).
	//       Dabei sind die Mitten der Grundlinien der Kollisionsbox und des Sprites
	//       aufeinander ausgerichtet.
	Draw(target pixel.Target)

	// Vor.: -
	// Eff.: Die Koordinaten der Mitte der Grundlinie der Kollisionsbox werden geliefert.
	GetBaselineCenter() pixel.Vec

	// Vor.: -
	// Eff.: Liefert die aktuelle Bewegungsrichtung des Characters
	GetDirection() uint8

	// Vor.: -
	// Eff.: Liefert die Nummer des Feldes, auf welchem sich der Character befindet.
	//       Die Nummer berechnet sich aus den Koordianten x,y des Feldes wiefolgt:
	//       x + y * Spielfeldbreite
	GetFieldNo() int

	// Vor.: -
	// Eff.: Liefert die Anzahl der verbleibenden Leben eines Characters.
	GetLife() uint8

	// Vor.: -
	// Eff.: Liefert einen Pointer auf die Lebensvariable, damit dieser an den Titlebar
	//       weitergereicht werden kann.
	GetLifePointer() *uint8

	// Vor.: -
	// Eff.: Liefert die aktuelle Positionsmatrix des Characters im Spielfeld.
	GetMatrix() pixel.Matrix

	// Vor.: -
	// Eff.: Liefert einen Vektor der die Koordinaten für das Zeichnen des Sprites enthält,
	//		 damit die Kollisionsbox und der Sprite bzgl. der Mitte ihrer jeweiligen Grundlinie
	//       ausgerichtet sind.
	GetMovedPos() pixel.Vec

	// Vor.: -
	// Eff.: Bei einem Bomberman werden die gesammelten Punkte geliefert.
	//       Bei einem Enemy werden die Punkte geliefert, die man beim
	//       töten des Enemys verdient.
	GetPoints() uint32

	// Vor.: -
	// Eff.: Bei einem Bomberman wird ein Pointer auf die gesammelten Punkte geliefert.
	//       Bei einem Enemy wird ein Pointer auf die Punkte geliefert, die man beim
	//       töten des Enemys verdient.
	GetPointsPointer() *uint32

	// Vor.: -
	// Eff.: Liefert einen Vektor der die Koordinaten der linken unteren Ecke
	//       der Kollisionsbox enthält
	GetPos() pixel.Vec

	// Vor.: -
	// Eff.: Liefert die Kollisionsbox als pixel.Rect
	GetPosBox() pixel.Rect

	// Vor.: -
	// Eff.: Breite und Höhe der Kollisionsbox werden als pixel.Vec zurückgegeben.
	GetSize() pixel.Vec

	// Vor.: -
	// Eff.: Liefert die aktuelle Geschwindigkeit des Characters
	GetSpeed() float64

	// Vor.: -
	// Eff.: Erhöht die Geschwindigkeit eines Characters um einen festen Wert.
	IncSpeed()

	// Vor.: -
	// Eff.: Liefert true, wenn die Anzahl der verbleibenden Leben größer als 0 ist.
	IsAlife() bool

	// Vor.: -
	// Eff.: Liefert true, wenn der Character durch Bomben gehen kann, sonst false.
	IsBombghost() bool

	// Vor.: -
	// Eff.: Liefert true, wenn der Character aktuell sterblich ist, sonst false.
	IsMortal() bool

	// Vor.: -
	// Eff.: Liefert true, wenn der Character durch Wände laufen kann, sonst false.
	IsWallghost() bool

	// Vor.: -
	// Eff.: Die Kollisionsbox wird um den Vektor delta verschoben.
	Move(delta pixel.Vec)

	// Vor.: -
	// Eff.: Die Kollisionsbox wird an die Position pos verschoben bzgl.
	//       der linken unteren Ecke.
	MoveTo(pos pixel.Vec)

	// Vor.: -
	// Eff.: Legt fest, ob der Character durch Bomben laufen kann (true) oder nicht (false).
	SetBombghost(bool)

	// Vor.: -
	// Eff.: Legt die aktuelle Bewegungsrichtung des Characters fest
	//       0 -> Left, 1 -> Right, 2 -> Up, 3 -> Down
	SetDirection(dir uint8)

	// Vor.: -
	// Eff.: Setzt die Nummer des Feldes, auf welchem sich der Character befindet.
	//       Die Nummer berechnet sich aus den Koordianten x,y des Feldes wiefolgt:
	//       x + y * Spielfeldbreite
	SetFieldNo(no int)
}
