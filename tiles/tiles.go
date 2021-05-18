package tiles

import (
	"../animations"
	"../characters"
	. "../constants"
	"github.com/faiface/pixel"
	"time"
)

/*
 * Es wird zwischen Bomben und anderen Items unterschieden
 * Hauptunterschied: Normale Items benötigen keinen Besitzer,
 * denn noemale Items verbessern die Werte eines Players beim Einsammeln
 */

type bombe struct {
	item
	owner characters.Player // Bombenbesitzer

	// Weshalb float64 und nicht uint8 ???
	power float64 // Wirkungsradius der Bomben

	// Auch die Zeitspanne sollte man setzen können, um ggf
	// später eine Fernzündung zu ermöglichen.
	// Fernzündung -> lange Zeitspanne einstellen
	// Normalzündung -> mittlere Zeitspanne einstellen
	// Totenkopf -> kurze Zeitspanne einstellen? (ist nur so ne Idee)
}

type item struct {
	tile
	destroyable bool
	timeStamp   time.Time
}

type tile struct {
	tileType uint8
	ani      animations.Animation
	matrix   pixel.Matrix
	pos      pixel.Vec
	//	x, y     int
}

func NewBomb(p characters.Player, pos pixel.Vec) *bombe {
	var bomb = new(bombe)
	(*bomb).tileType = Bomb
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	bomb.pos = pos
	(*bomb).matrix = pixel.IM.Moved(bomb.pos)
	d, _ := time.ParseDuration("3s")
	(*bomb).timeStamp = (time.Now()).Add(d)
	return bomb
}

func NewItem(t uint8, pos pixel.Vec) *item {
	var it = new(item)
	(*it).tileType = t
	(*it).ani = animations.NewAnimation(t)
	((*it).ani).Show()
	//	(*it).x = x
	//	(*it).y = y
	(*it).pos = pos
	(*it).matrix = pixel.IM.Moved(it.pos)
	d, _ := time.ParseDuration("100m")
	(*it).timeStamp = (time.Now()).Add(d)
	return it
}

func NewTile(t uint8, pos pixel.Vec) *tile {
	var nt = new(tile)
	(*nt).tileType = t
	(*nt).ani = animations.NewAnimation(t)
	((*nt).ani).Show()
	//	(*nt).x = x
	//	(*nt).y = y
	(*nt).pos = pos
	if nt.ani.GetSize().Y > TileSize {
		(*nt).matrix = pixel.IM.Moved(nt.pos.Add(pixel.V(0, ((*nt).ani.GetSize().Y-TileSize)/2)))
	} else {
		(*nt).matrix = pixel.IM.Moved(nt.pos)
	}

	return nt
}

/*
Bisher völlig unbenutzt:

func (it *tile) GetIndexPos() (x, y int) {
	return it.x, it.y
}
*/

// ******************* Items ***********************

func (it *item) GetTimeStamp() time.Time {
	return (*it).timeStamp
}

func (it *item) IsDestroyable() bool {
	return (*it).destroyable
}

func (it *item) SetDestroyable(b bool) {
	(*it).destroyable = b
}

func (it *item) SetTimeStamp(t time.Time) {
	(*it).timeStamp = t
}

// ******************* Tiles ***********************

func (it *tile) Ani() animations.Animation {
	return (*it).ani
}

func (it *tile) Draw(win pixel.Target) {
	((*it).ani).Update()
	if !it.ani.IsVisible() {
		return
	}
	if it.ani.GetView() == Dead {
		it.matrix = pixel.IM.Moved(it.pos.Add(pixel.V(0, it.ani.GetSize().Y/4)))
	}
	if it.ani.IsVisible() {
		(((*it).ani).GetSprite()).Draw(win, (*it).matrix)
	}
}

func (it *tile) GetMatrix() pixel.Matrix {
	return (*it).matrix
}

func (it *tile) GetPos() pixel.Vec {
	return (*it).pos
}

func (it *tile) GetType() uint8 {
	return it.tileType
}

func (it *tile) IsVisible() bool {
	return ((*it).ani).IsVisible()
}

func (it *tile) SetVisible(b bool) {
	((*it).ani).SetVisible(b)
}

func (it *tile) SetPos(pos pixel.Vec) {
	(*it).pos = pos
	(*it).matrix = pixel.IM
	(*it).matrix = ((*it).matrix).Moved(pos)
}

//------------------ Funktionen für Bomben -----------------------------

func (it *bombe) GetPower() float64 {
	return (*it).power
}

func (it *bombe) Owner() (bool, characters.Player) {
	return (*it).owner != nil, (*it).owner
}

func (it *bombe) SetAnimation(newAni animations.Animation) {
	(*it).ani = newAni
	((*it).ani).Update()
}

func (it *bombe) SetOwner(player characters.Player) {
	(*it).owner = player
}

func (it *bombe) SetPower(newPower float64) {
	(*it).power = newPower
	(*it).matrix = ((*it).matrix).ScaledXY((*it).pos, pixel.V(newPower*3.3, newPower*3.3))
}
