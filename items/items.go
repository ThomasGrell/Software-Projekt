package items

import (
	"../animations"
	"../characters"
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

/*
 * Es wird zwischen Bomben und anderen Items unterschieden
 * Hauptunterschied: Normale Items benötigen keinen Besitzer,
 * denn noemale Items verbessern die Werte eines Players beim Einsammeln
 */

type data struct {
	destroyable bool
	ani         animations.Animation
	matrix      pixel.Matrix
	pos         pixel.Vec
}

type bombe struct {
	data
	owner characters.Player // Bombenbesitzer
	power float64           // Wirkungsradius der Bomben
}

func NewItem(t uint8, pos pixel.Vec) *data {
	var item = new(data)
	(*item).ani = animations.NewAnimation(t)
	((*item).ani).Show()
	mat := pixel.IM
	mat = mat.Moved(pos)
	(*item).matrix = mat
	(*item).pos = pos
	return item
}

func NewBomb(p characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	mat := pixel.IM
	mat = mat.Moved(pixel.V(p.GetPos().X, p.GetPos().Y))
	//fmt.Println(p.GetPos())
	(*bomb).matrix = mat
	(*bomb).pos = p.GetPos()
	return bomb
}

func (item *data) SetDestroyable(b bool) {
	(*item).destroyable = b
}

func (item *data) IsDestroyable() bool {
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

func (item *data) Draw(win *pixelgl.Window) {
	((*item).ani).Update()
	(((*item).ani).GetSprite()).Draw(win, (*item).matrix)
}

func (item *data) SetPos(pos pixel.Vec) {
	(*item).pos = pos
	(*item).matrix = pixel.IM
	(*item).matrix = ((*item).matrix).Moved(pos)
}

func (item *data) GetPos() pixel.Vec {
	return (*item).pos
}

//------------------ Funktionen für Bomben -----------------------------

func (item *bombe) Owner() (bool, characters.Player) {
	return (*item).owner != nil, (*item).owner
}

func (item *bombe) SetOwner(player characters.Player) {
	(*item).owner = player
}

func (item *bombe) SetPower(newPower float64) {
	(*item).power = newPower
	(*item).matrix = ((*item).matrix).ScaledXY((*item).pos, pixel.V(newPower*3.3, newPower*3.3))
}
