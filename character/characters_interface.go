package character

import "github.com/faiface/pixel"

// Definition der Charaktertypen
const (
	WhiteBomberman = 1 // Spielfiguren im Single Player Mode
	BlackBomberman = 2
	BlueBomberman  = 3
	RedBomberman   = 4
	WhiteBattleman = 5 // Spielfiguren im Multi Player Mode
	BlackBattleman = 6
	BlueBattleman  = 7
	RedBattleman   = 8
	Balloon        = 9
	Teddy          = 10
	Ghost          = 11
	Drop           = 12
	Pinky          = 13
	BluePopEye     = 14
	Jellyfish      = 15
	Snake          = 16
	Spinner        = 17
	YellowPopEye   = 18
	Snowy          = 19
)

// Definition der Bewegungsrichtungen
const (
	Stay  = 0
	Up    = 1
	Down  = 2
	Left  = 3
	Right = 4
	Dead  = 5
	Intro = 6
)

type Enemy interface {
	// Alle Methoden eines Monsters
	Character
	IsFollowing() bool
}

type Player interface {
	// Alle Methoden eines Spielers
	Character
	AddPoints(int)
	GetMaxBombs() int
	GetWins() int
	IncLife()
	IncMaxBombs()
	IncWins()
	ResetWins()
	SetLife(int)
	SetMaxBombs(int)
	SetMortal(bool)
	SetWallghost(bool)
}

type Character interface {
	// Gemeinsame Methoden von Spielern und Monstern
	Animation
	DecLife()
	GetPoints() int
	IsAlife() bool
	IsBombghost() bool
	IsMortal() bool
	IsWallghost() bool
	SetBombghost(bool)
}

type Animation interface {
	DecSpeed()
	GetCenterPos() pixel.Vec // Koordinaten der Mitte der Kollisionsbox im Spielfeld
	GetMaxPos() pixel.Vec    // rechte obere Ecke der Kollisionsbox im Spielfeld
	GetMinPos() pixel.Vec    // linke untere Ecke der Kollisionsbox im Spielfeld
	GetOffset() pixel.Vec    // liefert Verschiebevektor zwischen Sprite und Kollisionsbox
	GetSpeed() float64
	GetSprite() *pixel.Sprite
	GetSpriteCoords() pixel.Rect // Rechteckkoordinaten des Sprites im Spriteimage
	IncSpeed()
	IsVisible() bool
	SetDirection(direction int)
	SetMinPos(v pixel.Vec) // Platziert den Sprite auf dem Spielfeld
	SetVisible(bool)
	Update()
}
