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
	// Alle Methoden eines Spielers
	Character
	AddPoints(uint32)
	GetMaxBombs() uint8
	GetPower () uint8
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
	GetBaselineCenter() pixel.Vec
	GetPoints() uint32
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
	SetPos(pixel.Vec)
}
