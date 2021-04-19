package items

import (
	"../animations"
	"../characters"
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"time"
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
	timeStamp 	time.Time
	itemType uint8
}

type bombe struct {
	data
	owner characters.Player // Bombenbesitzer

	// Weshalb float64 und nicht uint8 ???
	power float64 // Wirkungsradius der Bomben
	
	// Der Bombe fehlt ein Zeitstempel, wann sie gelegt wurde,
	// damit sie nach Ablauf einer Zeitspanne explodieren kann.
	// Auch die Zeitspanne sollte man setzen können, um ggf
	// später eine Fernzündung zu ermöglichen.
	// Fernzündung -> lange Zeitspanne einstellen
	// Normalzündung -> mittlere Zeitspanne einstellen
	// Totenkopf -> kurze Zeitspanne einstellen? (ist nur so ne Idee)
}

func NewItem(t uint8, pos pixel.Vec) *data {
	var item = new(data)
	(*item).itemType = t
	(*item).ani = animations.NewAnimation(t)
	((*item).ani).Show()
	(*item).matrix = pixel.IM.Moved(pos)
	(*item).pos = pos
	d,_:= time.ParseDuration("100m")
	(*item).timeStamp = (time.Now()).Add(d)
	return item
}

func NewBomb(p characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).itemType = Bomb
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	//fmt.Println(p.GetPosBox().Min)
	(*bomb).pos = pixel.Vec{math.Round(p.GetPosBox().Center().X/16) * 16, math.Round(p.GetPosBox().Center().Y/16) * 16}
	(*bomb).matrix = pixel.IM.Moved(bomb.pos)
	d,_:= time.ParseDuration("3s")
	(*bomb).timeStamp = (time.Now()).Add(d)
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

func (item *data) GetTimeStamp () time.Time {
	return (*item).timeStamp
}

func (item *data) GetMatrix() pixel.Matrix {
	return (*item).matrix
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

func (item *bombe) GetPower() float64 {
	return (*item).power
}

func (item *bombe) SetAnimation (newAni animations.Animation) {
	(*item).ani = newAni
	((*item).ani).Update()
}
