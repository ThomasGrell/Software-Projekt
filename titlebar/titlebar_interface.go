package titlebar

import (
	"github.com/faiface/pixel"
)

type Titlebar interface {
	DecLife(uint8)

	Draw(target pixel.Target)

	GetSeconds() uint16

	IncLevel()

	IncLife(uint8)

	Manager()

	SetLevel(level uint8)

	SetLife(life ...uint8)

	SetPlayers(uint8)

	SetPoints(points uint32)

	SetSeconds(seconds uint16)

	StartCountdown()

	StopCountdown()

	Update()
}
