package characters

import (
	"../animations"
	"github.com/faiface/pixel"
)

type Enemy interface {
	// Alle Methoden eines Monsters
	Character
	IsFollowing() bool
}

type Player interface {
	Character // Alle Methoden eines Spielers
	AddPoints(uint32)
	GetMaxBombs() uint8
	GetPower() uint8 //NEU NEU NEU
	GetWins() uint8
	IncLife()
	IncMaxBombs()
	IncPower()
	IncWins()
	ResetWins()
	SetLife(uint8)
	SetMaxBombs(uint8)
	SetMortal(bool)
	SetWallghost(bool)
}

type Character interface {
	// Gemeinsame Methoden von Spielern und Monstern

	// Vor.: keine
	// Eff.: Liefert das zugehörige Animationsobjekt des Characters
	Ani() animations.Animation

	DecLife()
	DecSpeed()

	// Vor.: keine
	// Eff.: Zeichnet den Sprite in das Fenster. Dabei sind die Mitten der Grundlinien
	//       der Kollisionsbox und des Sprites aufeinander ausgerichtet.
	Draw(target pixel.Target)

	// Vor.: keine
	// Eff.: Die Koordinaten der Mitte der Grundlinie der Kollisionsbox werden geliefert.
	GetBaselineCenter() pixel.Vec

	GetDirection() uint8

	GetLife() uint8

	GetLifePointer() *uint8

	GetMatrix() pixel.Matrix

	// Vor.: keine
	// Eff.: Liefert einen Vektor der die Koordinaten für das Zeichnen des Sprites enthält,
	//		 damit die Kollisionsbox und der Sprite bzgl. der Mitte ihrer jeweilien Grundlinie
	//       ausgerichtet sind.
	GetMovedPos() pixel.Vec

	// Vor.: keine
	// Eff.: Bei einem Bomberman werden die gesammelten Punkte geliefert.
	//       Bei einem Enemy werden die Punkte geliefert, die man beim
	//       töten des Enemys verdient.
	GetPoints() uint32

	// Vor.: keine
	// Eff.: Bei einem Bomberman wird ein Pointer auf die gesammelten Punkte geliefert.
	//       Bei einem Enemy wird ein Pointer auf die Punkte geliefert, die man beim
	//       töten des Enemys verdient.
	GetPointsPointer() *uint32

	// Vor.: keine
	// Eff.: Liefert einen Vektor der die Koordinaten der linken unteren Ecke
	//       der Kollisionsbox enthält
	GetPos() pixel.Vec

	// Vor.: keine
	// Eff.: Liefert die Kollisionsbox als pixel.Rect
	GetPosBox() pixel.Rect

	// Vor.: keine
	// Eff.: Breite und Höhe der Kollisionsbox werden als pixel.Vec zurückgegeben.
	GetSize() pixel.Vec

	GetSpeed() float64
	IncSpeed()
	IsAlife() bool
	IsBombghost() bool
	IsMortal() bool
	IsWallghost() bool

	// Vor.: keine
	// Eff.: Die Kollisionsbox ist um den Vektor delta verschoben worden.
	Move(delta pixel.Vec)

	// Vor.: keine
	// Eff.: Die Kollisionsbox wurde an die Position pos verschoben bzgl.
	//       der linken unteren Ecke
	MoveTo(pos pixel.Vec)

	SetBombghost(bool)

	SetDirection(dir uint8)
}
