package tiles

import (
	"../animations"
	"../characters"
	. "../constants"
	"github.com/faiface/pixel"
	"math"
	"time"
)

/*
 * Es wird zwischen Bomben und anderen Items unterschieden
 * Hauptunterschied: Normale Items benötigen keinen Besitzer,
 * denn noemale Items verbessern die Werte eines Players beim Einsammeln
 */

type tile struct {
	tileType uint8
	ani      animations.Animation
	matrix   pixel.Matrix
	pos      pixel.Vec
	x, y     int
}

type item struct {
	tile
	destroyable bool
	timeStamp   time.Time
}

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

func NewItem(t uint8, loleft pixel.Vec, x, y int) *item {
	var it = new(item)
	(*it).tileType = t
	(*it).ani = animations.NewAnimation(t)
	((*it).ani).Show()
	(*it).x = x
	(*it).y = y
	(*it).pos = loleft.Add(pixel.V(float64(x)*TileSize+TileSize/2, float64(y)*TileSize+TileSize/2))
	(*it).matrix = pixel.IM.Moved(it.pos)
	d, _ := time.ParseDuration("100m")
	(*it).timeStamp = (time.Now()).Add(d)
	return it
}

func NewBomb(p characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).tileType = Bomb
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	//fmt.Println(p.GetPosBox().Min)
	(*bomb).pos = pixel.Vec{math.Round(p.GetPosBox().Center().X/TileSize) * TileSize, math.Round(p.GetPosBox().Center().Y/TileSize) * TileSize}
	(*bomb).matrix = pixel.IM.Moved(bomb.pos)
	d, _ := time.ParseDuration("3s")
	(*bomb).timeStamp = (time.Now()).Add(d)
	return bomb
}

func NewTile(t uint8, loleft pixel.Vec, x, y int) *tile {
	var nt = new(tile)
	(*nt).tileType = t
	(*nt).ani = animations.NewAnimation(t)
	((*nt).ani).Show()
	(*nt).x = x
	(*nt).y = y
	(*nt).pos = loleft.Add(pixel.V(float64(x)*TileSize+TileSize/2, float64(y)*TileSize+TileSize/2))
	if nt.ani.GetSize().Y > TileSize {
		(*nt).matrix = pixel.IM.Moved(nt.pos.Add(pixel.V(0, ((*nt).ani.GetSize().Y-TileSize)/2)))
	} else {
		(*nt).matrix = pixel.IM.Moved(nt.pos)
	}

	return nt
}

func (it *tile) GetIndexPos() (x, y int) {
	return it.x, it.y
}

func (it *item) SetDestroyable(b bool) {
	(*it).destroyable = b
}

func (it *item) IsDestroyable() bool {
	return (*it).destroyable
}

func (it *tile) Ani() animations.Animation {
	return (*it).ani
}

func (it *tile) SetVisible(b bool) {
	((*it).ani).SetVisible(b)
}

func (it *tile) IsVisible() bool {
	return ((*it).ani).IsVisible()
}

func (it *tile) Draw(win pixel.Target) {
	((*it).ani).Update()
	if it.ani.IsVisible() {
		(((*it).ani).GetSprite()).Draw(win, (*it).matrix)
	}
}

func (it *tile) SetPos(pos pixel.Vec) {
	(*it).pos = pos
	(*it).matrix = pixel.IM
	(*it).matrix = ((*it).matrix).Moved(pos)
}

func (it *tile) GetPos() pixel.Vec {
	return (*it).pos
}

func (it *item) GetTimeStamp() time.Time {
	return (*it).timeStamp
}

func (it *item) SetTimeStamp(t time.Time) {
	(*it).timeStamp = t
}

func (it *tile) GetMatrix() pixel.Matrix {
	return (*it).matrix
}

//------------------ Funktionen für Bomben -----------------------------

func (it *bombe) Owner() (bool, characters.Player) {
	return (*it).owner != nil, (*it).owner
}

func (it *bombe) SetOwner(player characters.Player) {
	(*it).owner = player
}

func (it *bombe) SetPower(newPower float64) {
	(*it).power = newPower
	(*it).matrix = ((*it).matrix).ScaledXY((*it).pos, pixel.V(newPower*3.3, newPower*3.3))
}

func (it *bombe) GetPower() float64 {
	return (*it).power
}

func (it *bombe) SetAnimation(newAni animations.Animation) {
	(*it).ani = newAni
	((*it).ani).Update()
}
