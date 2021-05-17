package characters

import (
	//	"image"
	//"github.com/faiface/pixel"
	//"image/png"
	//"log"
	//"os"
	//"time"
	"../animations"
	. "../constants"
	"github.com/faiface/pixel"
	//	"golang.org/x/image/colornames"
)

/*
Der Ursprung des Koordinatensystems von "pixel" ist unten links.
Monster werden stets animiert.
Bombermen sind nur bei Bewegung animiert.
Animationen sind teilweise von Bewegungsrichtung abhängig.
*/

// Bild, welches die Sprites aller Charaktere enthält
var bm *player
var en *enemy

const cBoxSize = 15

type player struct {
	character
	bombs    uint8 // Anzahl der aktuell gelegten Bomben
	kick     bool  // kann Bomben wegkicken
	maxBombs uint8 // maximale Anzahl der legbaren Bomben
	power    uint8 // Wirkungsradius der Bomben
	wins     uint8 // Siege für Multi-Player-Modus
}
type enemy struct {
	character
	follow bool // Folgt einem Spieler
}
type character struct {
	ani          animations.Animation
	bombghost    bool       // kann durch Bomben laufen
	collisionbox pixel.Rect // Kollisionsbox
	direction    uint8
	fieldNo      int
	life         uint8 // verbleibende Anzahl der Leben
	matrix       pixel.Matrix
	mortal       bool   // Sterblichkeit
	points       uint32 // Punkte
	speed        float64
	wallghost    bool // kann durch Wände laufen
}

func NewPlayer(t uint8) *player {
	c := new(player)
	switch t {
	case WhiteBomberman, BlackBomberman, BlueBomberman, RedBomberman:
		*c = *bm
		c.collisionbox = pixel.R(-cBoxSize/2, -cBoxSize*3/4, cBoxSize/2, cBoxSize/4)
	case WhiteBattleman, BlackBattleman, BlueBattleman, RedBattleman:
		*c = *bm
		c.life = 1
	default:
		panic("Unknown Player")
	}
	c.ani = animations.NewAnimation(t)
	c.matrix = pixel.IM.Moved(c.GetMovedPos())
	//	(*c).scale = 1.0
	return c
}
func NewEnemy(t uint8) *enemy {
	c := new(enemy)
	*c = *en
	c.ani = animations.NewAnimation(t)
	c.collisionbox = pixel.R(-cBoxSize/2, -cBoxSize/2, cBoxSize/2, cBoxSize/2)
	c.matrix = pixel.IM.Moved(c.GetMovedPos())
	switch t {
	case Balloon:
		c.speed = 10
		c.points = 100
	case Teddy:
		c.speed = 30
		c.points = 100
	case Ghost:
		c.speed = 30
		c.points = 100
	case Drop:
		c.speed = 20
		c.points = 100
	case Pinky:
		c.points = 100
	case BluePopEye:
		c.points = 100
	case Jellyfish:
		c.points = 100
	case Snake:
		c.points = 100
	case Spinner:
		c.points = 100
	case YellowPopEye:
		c.points = 100
	case YellowBubble:
		c.points = 100
	case PinkPopEye:
		c.points = 100
	case Fire:
		c.points = 100
	case Crocodile:
		c.points = 100
	case Coin:
		c.points = 100
	case Puddle:
		c.points = 100
	case PinkCyclops:
		c.points = 100
	case RedCyclops:
		c.points = 100
	case PinkFlower:
		c.points = 100
	case Fireball:
		c.ani.SetView(Intro)
		c.points = 100
	case Snowy:
		c.points = 100
	case BlueRabbit:
		c.points = 100
	case PinkDevil:
		c.points = 100
	case Penguin:
		c.points = 100
	case BlueCyclops:
		c.points = 100
	default:
		panic("Unknown Enemy")
	}
	return c
}

func (c *enemy) IsFollowing() bool {
	return c.follow
}

func (c *player) AddPoints(p uint32) {
	c.points += p
}
func (c *player) GetMaxBombs() uint8 { return c.maxBombs }
func (c *player) GetWins() uint8     { return c.wins }
func (c *player) GetPower() uint8    { return c.power }
func (c *player) IncLife() {
	c.life++
}
func (c *player) IncMaxBombs()        { c.maxBombs++ }
func (c *player) IncPower()           { c.power++ }
func (c *player) IncWins()            { c.wins++ }
func (c *player) ResetWins()          { c.wins = 0 }
func (c *player) SetLife(l uint8)     { c.life = l }
func (c *player) SetMaxBombs(b uint8) { c.maxBombs = b }
func (c *player) SetMortal(b bool)    { c.mortal = b }
func (c *player) SetWallghost(w bool) { c.wallghost = w }

func (c *character) Ani() animations.Animation { return c.ani }
func (c *character) DecLife() {
	if c.life == 0 {
		return
	}
	if c.mortal {
		c.life--
		if c.life == 0 {
			c.ani.Die()
		}
	}
}
func (c *character) DecSpeed() {
	if c.speed > 10 {
		c.speed -= 10
	}
}
func (c *character) Draw(target pixel.Target) {
	c.ani.Update()
	if c.ani.IsVisible() {
		c.ani.GetSprite().Draw(target, c.matrix)
	}
}
func (c *character) GetBaselineCenter() pixel.Vec {
	return c.collisionbox.Min.Add(pixel.V(c.collisionbox.Size().X/2, 0))
}
func (c *character) GetDirection() uint8    { return c.direction }
func (c *character) GetFieldNo() int        { return c.fieldNo }
func (c *character) GetLife() uint8         { return c.life }
func (c *character) GetLifePointer() *uint8 { return &c.life }
func (c *character) GetMatrix() pixel.Matrix {
	return (*c).matrix
}
func (c *character) GetMovedPos() pixel.Vec {
	return c.GetBaselineCenter().Add(c.ani.ToBaseline())
}
func (c *character) GetPoints() uint32         { return c.points }
func (c *character) GetPointsPointer() *uint32 { return &c.points }
func (c *character) GetPos() pixel.Vec         { return c.collisionbox.Min }
func (c *character) GetPosBox() pixel.Rect {
	return c.collisionbox
}
func (c *character) GetSize() pixel.Vec {
	return c.collisionbox.Size()
}
func (c *character) GetSpeed() float64 { return c.speed }
func (c *character) IncSpeed()         { c.speed += 10 }
func (c *character) IsAlife() bool {
	return c.life > 0
}
func (c *character) IsBombghost() bool { return c.bombghost }
func (c *character) IsMortal() bool {
	return c.mortal
}
func (c *character) IsWallghost() bool { return c.wallghost }
func (c *character) Move(delta pixel.Vec) {
	c.collisionbox = c.collisionbox.Moved(delta)
	c.matrix = c.matrix.Moved(delta)
}
func (c *character) MoveTo(pos pixel.Vec) {
	c.collisionbox = pixel.Rect{pos, pos.Add(c.collisionbox.Size())}
	c.matrix = pixel.IM.Moved(c.GetMovedPos())
}
func (c *character) SetBombghost(b bool) { c.bombghost = b }
func (c *character) SetDirection(dir uint8) {
	if dir >= 0 && dir <= 4 {
		c.direction = dir
	}
}
func (c *character) SetFieldNo(no int) {
	c.fieldNo = no
}

// init() wird beim Import dieses Packets automatisch ausgeführt.
func init() {

	// Bomberman Prototyp
	bm = new(player)
	bm.life = 3
	bm.maxBombs = 1
	bm.power = 1
	bm.speed = 100
	bm.kick = false
	bm.mortal = true
	bm.wallghost = false
	bm.bombghost = false
	bm.collisionbox = pixel.Rect{pixel.Vec{}, pixel.Vec{cBoxSize, cBoxSize}}

	// Monster Prototyp
	en = new(enemy)
	en.life = 1
	en.speed = 30
	en.direction = 1
	en.mortal = true
	en.wallghost = false
	en.bombghost = false
	en.follow = false
	en.collisionbox = pixel.Rect{pixel.Vec{}, pixel.Vec{cBoxSize, cBoxSize}}
}
