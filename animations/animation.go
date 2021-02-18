package animations

import (
	. "../constants"
	"github.com/faiface/pixel"
	"image/png"
	"log"
	"os"
	"time"
)

var bm, en *animation
var characterImage *pixel.PictureData

type animation struct {
	// Position/Bewegung:

	pos   pixel.Vec // pixelgenaue Position der linken unteren Ecke der Kollisionsbox im Spielfeld
	width pixel.Vec // Breite und Höhe der Kollisionsbox
	field uint16    // Nummer auf dem Spielfeld, wo sich der Charakter gerade aufhält. Hängt von x und y ab.

	// Eigenschaften der Animation:

	sprite *pixel.Sprite // Zugehöriger Sprite

	lastUpdate int64 // time.Now().UnixNano() des letzten Spritewechsels
	intervall  int64 // Dauer in Nanosekunden bis zum nächsten Spritewechsel

	cpos   pixel.Vec // Verschiebungsvektor des Sprites in Bezug zur Kollisionsbox
	cwidth pixel.Vec // Breite und Höhe des Charakter-Sprites
	count  int8      // Nummer des zuletzt gezeichneten Sprites der Animationssequenz beginnend bei 1
	delta  int8      // entweder -1 oder 1 für Animationsreihenfolge vor und zurück
	seesaw bool      // Bei true geht die Animation hin und her (delta=+1/-1). Bei false werden
	// die Sprites immer in derselben Reihenfolge durchlaufen (delta=1)
	hasIntro bool // True, wenn das Monster eine Animation zum Erscheinen hat.
	visible  bool // Sichtbarer Sprite?

	direction uint8 // Bewegungsrichtung (siehe oben definierte Konstanten)
	dn        uint8 // Anzahl der Sprites für Abwärtsbewegung
	in        uint8 // Anzahl der Sprites für Erscheinungssequenz
	kn        uint8 // Anzahl der Sprites für Todessequenz
	ln        uint8 // Anzahl der Sprites für Linksbewegung
	rn        uint8 // Anzahl der Sprites für Rechtsbewegung
	sn        uint8 // Anzahl der Sprites für unbewegte Figur
	un        uint8 // Anzahl der Sprites für Aufwärtsbewegung

	dpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach unten bewegenden Charakter
	ipos pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	kpos pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	lpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach links bewegenden Charakter
	rpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach rechts bewegenden Charakter
	spos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für den ruhenden Charakter
	upos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach oben bewegenden Charakter

	iwidth pixel.Vec // Spritegröße für Introanimation
	kwidth pixel.Vec // Spritegröße für Todesanimation

	child *animation // Für den Endgegner, den langen Drachen, welcher aus mehreren Segmenten besteht
}

func NewAnimation(t uint8) *animation {
	c := new(animation)
	c.lastUpdate = time.Now().UnixNano()
	c.intervall = 2e8

	switch t {
	case Balloon:
		*c = *en
	case Teddy:
		*c = *en
		c.spos.Y = 352
		c.upos.Y = 352
		c.dpos.Y = 352
		c.lpos.Y = 352
		c.rpos.Y = 352
		c.kpos.Y = 352
	case Ghost:
		*c = *en
		c.spos.Y = 336
		c.upos.Y = 336
		c.dpos.Y = 336
		c.lpos.Y = 336
		c.rpos.Y = 336
		en.kpos.Y = 21 * 16
		en.kn = 9
	case Drop:
		*c = *en
		c.spos.Y = 20 * 16
		c.upos.Y = 20 * 16
		c.dpos.Y = 20 * 16
		c.lpos.Y = 20 * 16
		c.rpos.Y = 20 * 16
		c.kpos.Y = 20 * 16
	case Pinky:
		*c = *en
		c.spos.Y = 19 * 16
		c.upos.Y = 19 * 16
		c.dpos.Y = 19 * 16
		c.lpos.Y = 19 * 16
		c.rpos.Y = 19 * 16
		c.kpos.Y = 19 * 16
	case BluePopEye:
		*c = *en
		c.spos.Y = 18 * 16
		c.upos.Y = 18 * 16
		c.dpos.Y = 18 * 16
		c.lpos.Y = 18 * 16
		c.rpos.Y = 18 * 16
		c.kpos.Y = 18 * 16
		c.kn = 9
	case Jellyfish:
		*c = *en
		c.spos.Y = 17 * 16
		c.upos.Y = 17 * 16
		c.dpos.Y = 17 * 16
		c.lpos.Y = 17 * 16
		c.rpos.Y = 17 * 16
		c.kpos.Y = 17 * 16
	case Snake:
		*c = *en
		c.spos.Y = 16 * 16
		c.upos.Y = 16 * 16
		c.dpos.Y = 16 * 16
		c.lpos.Y = 16 * 16
		c.rpos.Y = 16 * 16
	case Spinner:
		*c = *en
		c.seesaw = false
		c.dn = 4
		c.un = 4
		c.ln = 4
		c.rn = 4
		c.spos.X = 20 * 16
		c.spos.Y = 15 * 16
		c.upos.X = 19 * 16
		c.upos.Y = 15 * 16
		c.dpos.X = 19 * 16
		c.dpos.Y = 15 * 16
		c.lpos.X = 19 * 16
		c.lpos.Y = 15 * 16
		c.rpos.X = 19 * 16
		c.rpos.Y = 15 * 16
		c.kpos.X = 304 + 4*16
		c.kpos.Y = 15 * 16
	case YellowPopEye:
		*c = *en
		c.spos.Y = 13 * 16
		c.upos.Y = 13 * 16
		c.dpos.Y = 13 * 16
		c.lpos.Y = 13 * 16
		c.rpos.Y = 13 * 16
		c.kpos.Y = 13 * 16
	case Snowy:
		*c = *en
		c.spos.X = 224
		c.spos.Y = 224
		c.upos.X = 256
		c.upos.Y = 224
		c.dpos.X = 208
		c.dpos.Y = 224
		c.lpos.X = 304
		c.lpos.Y = 224
		c.ln = 2
		c.rpos.X = 336
		c.rpos.Y = 224
		c.rn = 2
		c.kpos.X = 336 + 2*16
		c.kpos.Y = 224
	}
	return c
}

func (c *animation) Die() {
	c.count = 1
	c.delta = 1
	c.direction = Dead
}
func (c *animation) GetCenterPos() (v pixel.Vec) {
	return c.pos.Add(c.width.Scaled(0.5))
}
func (c *animation) GetMaxPos() pixel.Vec {
	return c.pos.Add(c.width)
}
func (c *animation) GetMinPos() pixel.Vec {
	return c.pos
}
func (c *animation) GetOffset() pixel.Vec { return c.cpos }
func (c *animation) GetSprite() *pixel.Sprite {
	// GetSprite() liefert den aktuell zu zeichnenden Sprite.
	return c.sprite
}
func (c *animation) GetSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n uint8

	// Wenn die Figur ruht, wird stets derselbe Sprite in Blickrichtung der Figur ausgegeben.
	// Bewegt sie sich, so wird die Animation durchlaufen.

	if !c.visible {
		return pixel.R(16*16, 22*16, 17*16, 23*16)
	}

	switch c.direction {
	case Stay:
		v = c.spos
		n = c.sn
	case Up:
		v = c.upos
		n = c.un
	case Down:
		v = c.dpos
		n = c.dn
	case Left:
		v = c.lpos
		n = c.ln
	case Right:
		v = c.lpos
		n = c.ln
	case Dead:
		v = c.kpos
		n = c.kn
	case Intro:
		v = c.ipos
		n = c.in
	}
	// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
	if n > 1 {
		timenow := time.Now().UnixNano()
		if timenow-c.lastUpdate > c.intervall {
			c.lastUpdate = timenow
			if uint8(c.count) == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
				if c.direction == Dead {
					c.visible = false
				} else if c.seesaw {
					c.count--
					c.delta = -1
				} else {
					c.count = 1
				}
			} else if c.count == 1 {
				c.count++ // links angekommen in der Bildfolge --> Vorwärtsgang
				c.delta = 1
			} else {
				c.count += c.delta
			}
		}
		switch c.direction {
		case Dead:
			v.X += c.kwidth.X * float64(c.count-1)
		case Intro:
			v.X += c.iwidth.X * float64(c.count-1)
		default:
			v.X += c.cwidth.X * float64(c.count-1)
		}
	}

	return pixel.R(v.X, v.Y, v.X+c.cwidth.X, v.Y+c.cwidth.Y)
}
func (c *animation) IsVisible() bool {
	return c.visible
}
func (c *animation) SetDirection(direction uint8) {
	// SetDirection() setzt die Bewegungsrichtung neu.
	// Mögliche Eingabewerte sind Stay, Left, Right, Up, Down, Dead.
	// Im character.png ist bei animierten Charakteren
	// der zweite Sprite stets für die ruhende Figur.
	// Es muss dann die Charakterbreite addiert werden.
	if direction == Stay {
		switch c.direction {
		case Left, Right:
			c.spos = c.lpos
			if c.ln > 1 {
				c.spos.X += c.cwidth.X
			}
		case Up:
			c.spos = c.upos
			if c.un > 1 {
				c.spos.X += c.cwidth.X
			}
		case Down:
			c.spos = c.dpos
			if c.dn > 1 {
				c.spos.X += c.cwidth.X
			}
		}
	}
	c.direction = direction
}
func (c *animation) SetMinPos(v pixel.Vec) {
	c.pos = v
}
func (c *animation) SetIntervall(i int64) { c.intervall = i }
func (c *animation) SetVisible(b bool)    { c.visible = b }
func (c *animation) Update() {
	c.sprite.Set(characterImage, c.GetSpriteCoords())
}
func init() {
	file, err := os.Open("graphics/characters.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	characterImage = pixel.PictureDataFromImage(img)

	bm = new(animation)
	bm.visible = true
	bm.pos.X = 19
	bm.pos.Y = 19
	bm.width.X = 10
	bm.width.Y = 10
	bm.direction = Down
	bm.seesaw = true
	bm.cwidth.X = 16
	bm.cwidth.Y = 24
	bm.count = 2
	bm.delta = 1
	bm.spos.X = 16
	bm.spos.Y = 360
	bm.sn = 1
	bm.lpos.X = 3 * 16
	bm.lpos.Y = 360
	bm.ln = 3
	bm.rpos.X = 9 * 16
	bm.rpos.Y = 360
	bm.rn = 3
	bm.upos.X = 6 * 16
	bm.upos.Y = 360
	bm.un = 3
	bm.dpos.X = 0
	bm.dpos.Y = 360
	bm.dn = 3
	bm.kpos.X = 12 * 16
	bm.kpos.Y = 360
	bm.kn = 4
	bm.kwidth = bm.cwidth

	// Monster Prototyp
	en = new(animation)
	en.lastUpdate = time.Now().UnixNano()
	en.intervall = 2e8
	en.visible = true
	en.width.X = 10
	en.width.Y = 10
	en.direction = Down
	en.cpos = pixel.V(5, 8)
	en.cwidth.X = 16
	en.cwidth.Y = 16
	en.count = 2
	en.delta = 1
	en.seesaw = true
	en.hasIntro = false
	en.spos.X = 304
	en.spos.Y = 368
	en.sn = 3
	en.lpos.X = 304
	en.lpos.Y = 368
	en.ln = 3
	en.rpos.X = 304
	en.rpos.Y = 368
	en.rn = 3
	en.upos.X = 304
	en.upos.Y = 368
	en.un = 3
	en.dpos.X = 304
	en.dpos.Y = 368
	en.dn = 3
	en.kpos.X = 304 + 3*16
	en.kpos.Y = 23 * 16
	en.kn = 7
	en.kwidth = en.cwidth
}
