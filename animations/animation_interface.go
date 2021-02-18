package animations

import "github.com/faiface/pixel"

// Definition der Bewegungsrichtungen

type Animation interface {
	Die()
	GetCenterPos() pixel.Vec // Koordinaten der Mitte der Kollisionsbox im Spielfeld
	GetMaxPos() pixel.Vec    // rechte obere Ecke der Kollisionsbox im Spielfeld
	GetMinPos() pixel.Vec    // linke untere Ecke der Kollisionsbox im Spielfeld
	GetOffset() pixel.Vec    // liefert Verschiebevektor zwischen Sprite und Kollisionsbox
	GetSprite() *pixel.Sprite
	GetSpriteCoords() pixel.Rect // Rechteckkoordinaten des Sprites im Spriteimage
	IsVisible() bool
	SetDirection(uint8)
	SetIntervall(int64)
	SetMinPos(pixel.Vec) // Platziert den Sprite auf dem Spielfeld
	SetVisible(bool)
	Update()
}
