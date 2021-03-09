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
	owner characters.Player	// Bombenbesitzer
	power float64				// Wirkungsradius der Bomben
}

func NewItem (t uint8) *data {
	var item = new(data)
	(*item).ani = animations.NewAnimation(t)		
	((*item).ani).Show()				
	return item
}

func NewBomb (p characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	return bomb
}

func (item *data) SetDestroyable (b bool) {
	(*item).destroyable = b
}

func (item *data) IsDestroyable () bool {
	return (*item).destroyable
}

func (item *data) Ani() animations.Animation {
	return (*item).ani
}

func (item *data) SetVisible(b bool) {
	((*item).ani).SetVisible(b)
}

func (item *data) IsVisible() bool {
	return ((*item).ani).IsVisible()
}


//------------------ Funktionen für Bomben -----------------------------

func (item *bombe) Owner () (bool,characters.Player) {
	return (*item).owner!=nil,(*item).owner
}

func (item *bombe) SetOwner (player characters.Player) {
	(*item).owner = player
}

