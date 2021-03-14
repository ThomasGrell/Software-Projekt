package characters

import (
	"../animations"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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
	IncWins()
	ResetWins()
	SetLife(uint8)
	SetMaxBombs(uint8)
	SetMortal(bool)
	SetWallghost(bool)
}

type Character interface {
	// Gemeinsame Methoden von Spielern und Monstern
	Ani() animations.Animation
	DecLife()
	DecSpeed()
	Draw(win *pixelgl.Window) //NEU NEU NEU
	GetBaselineCenter() pixel.Vec
	GetPoints() uint32
	GetPosBox() pixel.Rect
	GetPos() pixel.Vec
	GetMovedPos() pixel.Vec
	GetSize() pixel.Vec
	GetSpeed() float64
	IncSpeed()
	IsAlife() bool
	IsBombghost() bool
	IsMortal() bool
	IsWallghost() bool
	SetBombghost(bool)
	MoveTo(pixel.Vec)
	SetScale(s float64)
}
