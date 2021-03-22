package items

import (
	"../animations"
	"../characters"
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"fmt"
)

/* 
 * Es wird zwischen Bomben und anderen Items unterschieden
 * Hauptunterschied: Normale Items benötigen keinen Besitzer, 
 * denn noemale Items verbessern die Werte eines Players beim Einsammeln
 */

type data struct {
	destroyable bool
	ani animations.Animation
	matrix pixel.Matrix
	pos pixel.Vec
	t uint8
}

type bombe struct {
	data
	owner characters.Player		// Bombenbesitzer
	power float64				// Wirkungsradius der Bomben
}

func NewItem (t uint8, pos pixel.Vec) *data {
	var item = new(data)
	(*item).ani = animations.NewAnimation(t)		
	((*item).ani).Show()	
	mat := pixel.IM
	mat = mat.Moved(pos)
	mat = mat.ScaledXY(pos, pixel.V(3.3, 3.3))
	(*item).matrix = mat
	(*item).pos = pos	
	(*item).t = t		
	return item
}

func NewBomb (p characters.Player) *bombe {
	var bomb = new(bombe)
	(*bomb).owner = p
	(*bomb).power = float64(p.GetPower())
	(*bomb).ani = animations.NewAnimation(Bomb)
	((*bomb).ani).Show()
	mat := pixel.IM
	mat = mat.Moved(p.GetPos())
	mat = mat.ScaledXY(p.GetPos(), pixel.V(3.3, 3.3))
	(*bomb).matrix = mat	
	(*bomb).pos = p.GetPos()
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

func (item *data) Draw (win *pixelgl.Window) {
	((*item).ani).Update()
	(((*item).ani).GetSprite()).Draw(win,(*item).matrix)
}

func (item *data) SetPos (pos pixel.Vec) {
	(*item).pos = pos
	(*item).matrix = pixel.IM
	(*item).matrix = ((*item).matrix).Moved(pos)
}

func (item *data) GetPos () pixel.Vec {
	return (*item).pos
}

func (item *data) Use (p characters.Player) {
	switch (*item).t {
		case Bomb:
			fmt.Println("Bomben sind nicht zum verzehren da.")
		case PowerItem:
			// p.IncBombPower wäre gut
		case BombItem:
			p.IncMaxBombs()
		case PunchItem:
			// keine Ahnung -.-
		case HeartItem:
			// p.IncLife() ?
		case RollerbladeItem:
			// p.IncSpeed() ?
		case SkullItem:
			// p.DecLife() ?
		case WallghostItem:
			p.SetWallghost(true)			// timer nötig?
		case BombghostItem:
			p.SetBombghost(true)			// timer nötig?
		case LifeItem:
			p.IncLife()
		case KickItem:
			// keine Ahnung -.-
		default:
			fmt.Println("Hier kst was schiefgelaufen: ",(*item).t," ist kein Gültiger Item-Wert!")
	}
}


//------------------ Funktionen für Bomben -----------------------------

func (item *bombe) Owner () (bool,characters.Player) {
	return (*item).owner!=nil,(*item).owner
}

func (item *bombe) SetOwner (player characters.Player) {
	(*item).owner = player
}

func (item *bombe) SetPower (newPower float64) {
	(*item).power = newPower
	(*item).matrix = ((*item).matrix).ScaledXY((*item).pos, pixel.V(newPower*3.3, newPower*3.3))
}

