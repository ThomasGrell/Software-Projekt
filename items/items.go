package items

import (
	"../animations"
	"../characters"
	. "../constants"
	//"github.com/faiface/pixel"
	//"github.com/faiface/pixel/pixelgl"
)

/* 
 * Es wird zwischen Bomben und anderen Items unterschieden
 * Hauptunterschied: Normale Items benötigen keinen Besitzer, 
 * denn noemale Items verbessern die Werte eines Players beim Einsammeln
 */

type data struct {
	destroyable bool
	ani animations.Animation
}

type bombe struct {
	data
	owner *characters.Player	// Bombenbesitzer
	power float64				// Wirkungsradius der Bomben
}

func NewItem (t uint8) *data {
	var item = new(data)
	(*item).ani = NewAnimation(t)		
	((*item).ani).Show()				
	((*item).ani).SetVisible(false)			// Zu Anfang unsichtbar
	return item
}

func NewBomb (p *characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).owner = p
	(*bomb).power = p.GetPower()
	(*bomb).ani = NewAnimation(Bomb)
	((*bomb).ani).Show()					// Sofort sichtbar
	return bomb
}

func (item *data) Sichtbar 

func (item *data) SetDestroyable (b bool) {
	(*item).destroyable = b
}

func (item *data) IsDestroyable () bool {
	return (*item).destroyable
}

func (item *data) Ani() animations.Animation {
	return (*item).ani
}


//------------------ Funktionen für Bomben -----------------------------

func (item *bombe) Owner () (bool,*characters.Player) {
	return (*item).owner!=nil,(*item).owner
}

func (item *bombe) SetOwner (player *characters.Player) {
	(*item).owner = player
}

