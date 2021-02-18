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
	//	"github.com/faiface/pixel/pixelgl"
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
	bombghost bool   // kann durch Bomben laufen
	mortal    bool   // Sterblichkeit
	wallghost bool   // kann durch Wände laufen
	life      uint8  // verbleibende Anzahl der Leben
	points    uint32 // Punkte
	speed     float64
	ani       animations.Animation
}

func NewPlayer(t uint8) *player {
	c := new(player)
	c.ani = animations.NewAnimation(t)
	switch t {
	case WhiteBomberman, BlackBomberman, BlueBomberman, RedBomberman:
		*c = *bm
	case WhiteBattleman, BlackBattleman, BlueBattleman, RedBattleman:
		*c = *bm
		c.life = 1
	}

	return c
}
func NewEnemy(t uint8) *enemy {
	c := new(enemy)
	c.ani = animations.NewAnimation(t)
	switch t {
	case Balloon:
		*c = *en
	case Teddy:
		*c = *en
		c.follow = true
	case Ghost:
		*c = *en
		c.wallghost = true
	case Drop:
		*c = *en
	case Pinky:
		*c = *en
	case BluePopEye:
		*c = *en
	case Jellyfish:
		*c = *en
	case Snake:
		*c = *en
	case Spinner:
		*c = *en
	case YellowPopEye:
		*c = *en
	case Snowy:
		*c = *en
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
func (c *player) IncLife() {
	c.life++
}
func (c *player) IncMaxBombs()        { c.maxBombs++ }
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
func (c *character) GetPoints() uint32 { return c.points }
func (c *character) GetSpeed() float64 { return c.speed }
func (c *character) IncSpeed()         { c.speed += 10 }
func (c *character) IsAlife() bool {
	return c.life > 0
}
func (c *character) IsBombghost() bool { return c.bombghost }
func (c *character) IsMortal() bool {
	return c.mortal
}
func (c *character) IsWallghost() bool   { return c.wallghost }
func (c *character) SetBombghost(b bool) { c.bombghost = b }

// init() wird beim Import dieses Packets automatisch ausgeführt.
func init() {

	// Bomberman Prototyp
	bm = new(player)
	bm.life = 3
	bm.maxBombs = 1
	bm.power = 1
	bm.speed = 10
	bm.kick = false
	bm.mortal = true
	bm.wallghost = false
	bm.bombghost = false

	// Monster Prototyp
	en = new(enemy)
	en.life = 1
	en.speed = 10
	en.mortal = true
	en.wallghost = false
	en.bombghost = false
	en.follow = false
}
