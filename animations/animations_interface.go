package animations

import "github.com/faiface/pixel"

/*

 _________________________________
< Implementiert von Rayk von Ende >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

****************************************************************************************

Konstruktoren:

Vor.: Die Nummer der Animation wird als uint8 übergeben.
Eff.: Ein Animationsobjekt wird geliefert.
NewAnimation() Animation

Vor.: Die Ausdehnung der Explosion in 4 Richtungen wird übergeben.
Eff.: Ein Explosionsobjekt wird geliefert. Die Views Up, Down etc. sind irrelevant. Ebenso die Methoden Die() und SetView().
	  Eine Explosion wird mit Show() gestartet.
NewExplosion(l, r, u, d uint8) Animation

****************************************************************************************
*/

type Animation interface {
	//	Vor.: keine
	//	Eff.: Die Todessequenz wird eingeleited. Ist diese beendet, wird die Animation unsichtbar. ( IsVisible()==false )
	Die()

	//	Vor.: keine
	//	Eff.: Koordinaten der Mitte der Animation werden geliefert. Bei Explosionen ist dies das Zentrum der Explosion.
	ToCenter() pixel.Vec

	//	Vor.: keine
	//	Eff.: Vektor zur Verschiebung des Sprites auf die Mitte der Grundlinie der Animation wird geliefert.
	//		  Bei Explosionen ist dies die Mitte der Grundlinie des Explosionszentrums.
	ToBaseline() pixel.Vec

	//	Vor.: keine
	//	Eff.: Vektor mit der Breite und Höhe der aktuellen Animation wird geliefert.
	//        Nach Aufruf der Methode Update() kann sich dieser aber ändern.
	GetSize() pixel.Vec

	//  Vor.: keine
	//	Eff.: Ein Pointer auf das Sprite der Animation wird geliefert.
	GetSprite() *pixel.Sprite

	//  Vor.: keine
	//	Eff.: Liefert die aktuelle Sicht. (Stay, Left, Right, Up, Down, Dead, Intro)
	GetView() uint8

	//	Vor.: keine
	//	Eff.: Liefert true wenn die Animation sichtbar ist. Am Ende einer Todessequenz wird die Sichtbarkeit stets auf
	//		  false gesetzt. Eine neu erstellte Animation ist ebenfalls nicht sichtbar und muss mit Show() gestartet werden.
	IsVisible() bool

	//  Vor.: keine
	//	Eff.: Durchläuft die Animation gerade eine Intro- oder Todessequenz,
	//		  so liefert die Funktion false, solange die Sequenz nicht beendet ist.
	//		  In allen anderen Fällen wird true geliefert.
	SequenceFinished() bool

	//	Vor.: keine
	//	Eff.: Mittels der Konstanten Intro, Dead, Left, Right, Up, Down und Stay kann das Aussehen der Animation
	//		  festgelegt werden. Viele der Animationen haben kein Intro, in diesem Fall wird der View auf Stay gesetzt.
	//		  Nur für die Bombermen und einige Enemies gibt es unterschiedliche Sprites für die Bewegungsrichtungen
	//		  Up, Down, Left, Right. Bomben haben keine Todessequenz, sondern einen eigenen Konstruktor für die
	//		  Explosion NewExplosion().
	SetView(uint8)

	//	Vor.: keine
	//	Eff.: Legt die Zeit in Nanosekunden fest, die beim Wechsel zwischen zwei Sprites mindestens vergehen muss.
	//  	  Je kürzer die Zeit, desto schneller die Animation.
	SetIntervall(int64)

	//	Vor.: keine
	//	Eff.: Legt fest, ob die Animation sichtbar ist.
	SetVisible(bool)

	//	Vor.: keine
	//	Eff.: Macht eine Animation sichtbar und startet diese. Muss nach dem Anlegen der Animation mittels NewAnimation()
	//		  oder NewExplosion() aufgerufen werden.
	Show()

	//	Vor.: Die Animation wurde mittels Show() gestartet und ist sichtbar.
	//	Eff.: Die Systemzeit wird abgefragt. Ist das Zeitintervall, welches mit SetIntervall() verändert werden kann,
	//		  abgelaufen, wird der Sprite aktualisiert. Kommt eine Animation an das Ende der Introsequenz, wird der View
	//		  auf Stay gesetzt. Am Ende einer Todessequenz wird die Animation unsichtbar gesetzt IsVisible()==false.
	Update()
}
