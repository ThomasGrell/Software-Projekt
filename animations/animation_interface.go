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
	ToCenter() pixel.Vec   // Koordinaten der Mitte der Animation im Spielfeld. Bei Explosionen ist dies das Zentrum der Explosion.
	ToBaseline() pixel.Vec // Vektor zur Verschiebung des Sprites auf die Mitte der Grundlinie der Animation. Bei Explosionen ist dies die Mitte der Grundlinie des Explosionszentrums.
	GetSize() pixel.Vec    // rechte obere Ecke der Kollisionsbox im Spielfeld
	GetSprite() *pixel.Sprite
	IntroFinished() bool //
	IsVisible() bool     // Gestorbene Charaktere werden unsichtbar gesetzt.
	SetDirection(uint8)
	SetIntervall(int64)
	SetVisible(bool)
	Update()
}
