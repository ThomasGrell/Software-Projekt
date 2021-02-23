package animations

import "github.com/faiface/pixel"

/*
Vor.: Die Nummer der Animation wird als uint8 übergeben.
Eff.: Ein Animationsobjekt wird geliefert, dessen Aussehen nicht bewegungsrichtungsabhängig ist. Gut geeignet für
	  einfache Monster und Items.
	  *basicAnimation erfüllt das Interface Animation
NewBasicAnimation() *basicAnimation

Vor.: Die Nummer der Animation wird als uint8 übergeben.
Eff.: Ein Animationsobjekt wird geliefert, dessen Aussehen bewegungsrichtungsabhängig ist. Gut geeignet für
	  komplexe Monster und Bomberman.
	  *enhancedAnimation erfüllt das Interface Animation
NewEnhancedAnimation() *enhancedAnimation
*/

type Animation interface {
	Die()
	GetCenter() pixel.Vec // Koordinaten der Mitte der Kollisionsbox im Spielfeld
	GetWidth() pixel.Vec  // rechte obere Ecke der Kollisionsbox im Spielfeld
	GetSprite() *pixel.Sprite
	GetSpriteCoords() pixel.Rect // Rechteckkoordinaten des Sprites im Spriteimage
	IntroFinished() bool         //
	IsVisible() bool             // Gestorbene Charaktere werden unsichtbar gesetzt.
	SetDirection(uint8)
	SetIntervall(int64)
	SetVisible(bool)
	Update()
}
