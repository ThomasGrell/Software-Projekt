package character

import (
	//	"image"
	"github.com/faiface/pixel"
	"image/png"
	"log"
	"os"
	"time"
	//	"github.com/faiface/pixel/pixelgl"
	//	"golang.org/x/image/colornames"
)

/*

Der Ursprung des Koordinatensystems von "pixel" ist unten links.
Monster werden stets animiert.
Bombermen sind nur bei Bewegung animiert.
Animationen sind teilweise von Bewegungsrichtung abhängig.
*/

// Definition der Charaktertypen
const (
	WhiteBomberman = 1 // Spielfiguren im Single Player Mode
	BlackBomberman = 2
	BlueBomberman  = 3
	RedBomberman   = 4
	WhiteBattleman = 5 // Spielfiguren im Multi Player Mode
	BlackBattleman = 6
	BlueBattleman  = 7
	RedBattleman   = 8
	Balloon        = 9
	Teddy          = 10
	Ghost          = 11
	Drop           = 12
	Pinky          = 13
	BluePopEye     = 14
	Jellyfish      = 15
	Snake          = 16
	Spinner        = 17
	YellowPopEye   = 18
	Snowy          = 19
)

// Definition der Bewegungsrichtungen
const (
	Stay  = 0
	Up    = 1
	Down  = 2
	Left  = 3
	Right = 4
	Dead  = 5
)

// Bild, welches die Sprites aller Charaktere enthält
var characterImage *pixel.PictureData
var bm, mo *character

type Character interface {
	DecLife()
	IncLife()
	IncSpeed()
	DecSpeed()
	GetSpeed() float64
	SetMinPos(v pixel.Vec)
	GetMinPos() pixel.Vec
	GetMaxPos() pixel.Vec
	GetCenterPos() pixel.Vec
	IsAlife() bool
	IsMortal() bool
	SetMortal(bool)
	GetSpriteCoords() pixel.Rect
	GetSprite() *pixel.Sprite
	Update()

	// Move() legt die Bewegungsrichtung des Charakters fest.
	// Übergeben wird eine der definierten Konstanten stay, left, right, up oder down
	// stay deaktiviert die Animation. Der Charakter schaut in Richtung seiner letzten Bewegung.
	Direction(direction int)
}

type character struct {

	// Fähigkeiten:

	points   int     // Punkte für den Multi-Player-Modus
	life     int     // verbleibende Anzahl der Leben
	maxBombs int     // maximale Anzahl der legbaren Bomben
	power    int     // Wirkungsradius der Bomben
	bombs    int     // Anzahl der aktuell gelegten Bomben
	speed    float64 // max. Bewegungsgeschwindigkeit in Pixel pro Sekunde

	kick      bool // kann Bomben wegkicken
	mortal    bool // Sterblichkeit
	wallghost bool // kann durch Wände laufen
	bombghost bool // kann durch Bomben laufen
	follow    bool // Folgt einem Spieler
	invisible bool // Unsichtbarer Sprite

	// Position/Bewegung:

	pos   pixel.Vec // pixelgenaue Position der linken unteren Ecke der Kollisionsbox im Spielfeld
	width pixel.Vec // Breite und Höhe der Kollisionsbox
	field int       // Nummer auf dem Spielfeld, wo sich der Charakter gerade aufhält. Hängt von x und y ab.

	direction int // Bewegungsrichtung (siehe oben definierte Konstanten)

	// Eigenschaften der Animation:

	sprite *pixel.Sprite // Zugehöriger Sprite

	lastUpdate int64 // time.Now().UnixNano() des letzten Spritewechsels
	interval   int64 // Dauer in Nanosekunden bis zum nächsten Spritewechsel

	cpos   pixel.Vec // Verschiebungsvektor des Sprites in Bezug zur Kollisionsbox
	cwidth pixel.Vec // Breite und Höhe des Charakter-Sprites
	count  int       // Nummer des zuletzt gezeichneten Sprites der Animationssequenz beginnend bei 1
	delta  int       // entweder -1 oder 1 für Animationsreihenfolge vor und zurück
	seesaw bool      // Bei true geht die Animation hin und her (delta=+1/-1). Bei false werden
	// die Sprites immer in derselben Reihenfolge durchlaufen (delta=1)
	hasIntro bool // True, wenn das Monster eine Animation zum Erscheinen hat.

	spos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für den ruhenden Charakter

	lpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach links bewegenden Charakter
	ln   int       // Anzahl der Sprites für Animationseffekte

	rpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach rechts bewegenden Charakter
	rn   int       // Anzahl der Sprites für Animationseffekte

	// Für rechtsläufige Figuren wird die linksläufige Animation gespiegelt.

	upos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach oben bewegenden Charakter
	un   int       // Anzahl der Sprites für Animationseffekte

	dpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach unten bewegenden Charakter
	dn   int       // Anzahl der Sprites für Animationseffekte

	ipos pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	in   int       // Anzahl der Sprites für die Erscheinungssequenz

	kpos pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	kn   int       // Anzahl der Sprites für die Todessequenz
}

func NewCharacter(t int) *character {

	c := new(character)

	c.lastUpdate = time.Now().UnixNano()
	c.interval = 2e8

	if t == WhiteBomberman {
		*c = *bm
	}
	if t == BlackBomberman {
		*c = *bm
		c.spos.Y = 336
		c.lpos.Y = 336
		c.upos.Y = 336
		c.dpos.Y = 336
		c.kpos.Y = 336
	}
	if t == BlueBomberman {
		*c = *bm
		c.spos.Y = 312
		c.lpos.Y = 312
		c.upos.Y = 312
		c.dpos.Y = 312
		c.kpos.Y = 312
	}
	if t == RedBomberman {
		*c = *bm
		c.spos.Y = 288
		c.lpos.Y = 288
		c.upos.Y = 288
		c.dpos.Y = 288
		c.kpos.Y = 288
	}
	if t == WhiteBattleman {
		c.life = 1
	}
	if t == BlackBattleman {
		*c = *bm
		c.life = 1
		c.spos.Y = 336
		c.lpos.Y = 336
		c.upos.Y = 336
		c.dpos.Y = 336
		c.kpos.Y = 336
	}
	if t == BlueBattleman {
		*c = *bm
		c.life = 1
		c.spos.Y = 312
		c.lpos.Y = 312
		c.upos.Y = 312
		c.dpos.Y = 312
		c.dpos.Y = 312
	}
	if t == RedBattleman {
		*c = *bm
		c.life = 1
		c.spos.Y = 288
		c.lpos.Y = 288
		c.upos.Y = 288
		c.dpos.Y = 288
		c.kpos.Y = 288
	}
	if t == Balloon {
		*c = *mo
	}
	if t == Teddy {
		*c = *mo
		c.spos.Y = 352
		c.upos.Y = 352
		c.dpos.Y = 352
		c.lpos.Y = 352
		c.rpos.Y = 352
		c.kpos.Y = 352
		c.follow = true
	}
	if t == Ghost {
		*c = *mo
		c.spos.Y = 336
		c.upos.Y = 336
		c.dpos.Y = 336
		c.lpos.Y = 336
		c.rpos.Y = 336
		mo.kpos.Y = 21 * 16
		mo.kn = 9
		c.wallghost = true
	}
	if t == Drop {
		*c = *mo
		c.spos.Y = 20 * 16
		c.upos.Y = 20 * 16
		c.dpos.Y = 20 * 16
		c.lpos.Y = 20 * 16
		c.rpos.Y = 20 * 16
		c.kpos.Y = 20 * 16
	}
	if t == Pinky {
		*c = *mo
		c.spos.Y = 19 * 16
		c.upos.Y = 19 * 16
		c.dpos.Y = 19 * 16
		c.lpos.Y = 19 * 16
		c.rpos.Y = 19 * 16
		c.kpos.Y = 19 * 16
	}
	if t == BluePopEye {
		*c = *mo
		c.spos.Y = 18 * 16
		c.upos.Y = 18 * 16
		c.dpos.Y = 18 * 16
		c.lpos.Y = 18 * 16
		c.rpos.Y = 18 * 16
		c.kpos.Y = 18 * 16
		c.kn = 9
	}
	if t == Jellyfish {
		*c = *mo
		c.spos.Y = 17 * 16
		c.upos.Y = 17 * 16
		c.dpos.Y = 17 * 16
		c.lpos.Y = 17 * 16
		c.rpos.Y = 17 * 16
		c.kpos.Y = 17 * 16
	}
	if t == Snake {
		*c = *mo
		c.spos.Y = 16 * 16
		c.upos.Y = 16 * 16
		c.dpos.Y = 16 * 16
		c.lpos.Y = 16 * 16
		c.rpos.Y = 16 * 16
	}
	if t == Spinner {
		*c = *mo
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
	}
	if t == YellowPopEye {
		*c = *mo
		c.spos.Y = 13 * 16
		c.upos.Y = 13 * 16
		c.dpos.Y = 13 * 16
		c.lpos.Y = 13 * 16
		c.rpos.Y = 13 * 16
		c.kpos.Y = 13 * 16
	}
	if t == Snowy {
		*c = *mo
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

	c.sprite = pixel.NewSprite(characterImage, pixel.R(c.spos.X, c.spos.Y, c.spos.X+c.cwidth.X, c.spos.Y+c.cwidth.Y))

	return c
}

func (c *character) DecLife() {
	if c.life == 0 {
		return
	}
	if c.mortal {
		c.life--
		if c.life == 0 {
			c.count = 1
			c.delta = 1
			c.direction = Dead
		}
	}
}

func (c *character) IncLife() {
	c.life++
}

func (c *character) DecSpeed() {
	if c.speed > 10 {
		c.speed -= 10
	}
}

func (c *character) IncSpeed() {
	c.speed += 10
}

// GetSpeed() gibt die Geschwindigkeit in Pixel/Sekunde zurück.
func (c *character) GetSpeed() float64 {
	return c.speed
}

func (c *character) SetMinPos(v pixel.Vec) {
	c.pos = v
}

func (c *character) GetMinPos() pixel.Vec {
	return c.pos
}

func (c *character) GetMaxPos() pixel.Vec {
	return c.pos.Add(c.width)
}

func (c *character) GetCenterPos() (v pixel.Vec) {
	return c.pos.Add(c.width.Scaled(0.5))
}

func (c *character) IsAlife() bool {
	return c.life > 0
}

func (c *character) IsMortal() bool {
	return c.mortal
}

func (c *character) SetMortal(b bool) {
	c.mortal = b
}

func (c *character) GetSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n int

	// Wenn die Figur ruht, wird stets derselbe Sprite in Blickrichtung der Figur ausgegeben.
	// Bewegt sie sich, so wird die Animation durchlaufen.

	if c.invisible {
		return pixel.R(16*16, 22*16, 17*16, 23*16)
	}

	if c.direction == Stay {
		v = c.spos
	} else {

		if c.direction == Up {
			v = c.upos
			n = c.un
		}
		if c.direction == Down {
			v = c.dpos
			n = c.dn
		}
		if c.direction == Left {
			v = c.lpos
			n = c.ln
		}
		if c.direction == Right {
			v = c.lpos
			n = c.ln
		}
		if c.direction == Dead {
			v = c.kpos
			n = c.kn
		}

		// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
		if n > 1 {
			timenow := time.Now().UnixNano()
			if timenow-c.lastUpdate > c.interval {
				c.lastUpdate = timenow
				if c.count == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
					if c.direction == Dead {
						c.invisible = true
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
		}

		v.X += c.cwidth.X * float64(c.count-1)
	}

	return pixel.R(v.X, v.Y, v.X+c.cwidth.X, v.Y+c.cwidth.Y)
}

func (c *character) Update() {
	c.sprite.Set(characterImage, c.GetSpriteCoords())
}

// GetSprite() liefert den aktuell zu zeichnenden Sprite.
func (c *character) GetSprite() *pixel.Sprite {
	return c.sprite
}

// Direction() setzt die Bewegungsrichtung neu.
// Mögliche Eingabewerte sind Stay, Left, Right, Up, Down.
func (c *character) Direction(direction int) {
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

// init() wird beim Import dieses Packets automatisch ausgeführt.
// Es lädt die Bilddatei mit den Charakteren in den Speicher.
// Zugriff über characterImage.
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

	// Bomberman Prototyp
	bm = new(character)
	bm.lastUpdate = time.Now().UnixNano()
	bm.interval = 2e8
	bm.life = 3
	bm.maxBombs = 1
	bm.power = 1
	bm.speed = 100
	bm.kick = false
	bm.mortal = true
	bm.wallghost = false
	bm.bombghost = false
	bm.invisible = false
	bm.pos.X = 19
	bm.pos.Y = 19
	bm.width.X = 10
	bm.width.Y = 10
	bm.direction = Down
	bm.cwidth.X = 16
	bm.cwidth.Y = 24
	bm.count = 2
	bm.delta = 1
	bm.spos.X = 16
	bm.spos.Y = 360
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

	// Monster Prototyp
	mo = new(character)
	mo.lastUpdate = time.Now().UnixNano()
	mo.interval = 2e8
	mo.life = 1
	mo.maxBombs = 0
	mo.power = 0
	mo.speed = 100
	mo.kick = false
	mo.mortal = true
	mo.wallghost = false
	mo.bombghost = false
	mo.follow = false
	mo.invisible = false
	mo.width.X = 10
	mo.width.Y = 10
	mo.direction = Down
	mo.cwidth.X = 16
	mo.cwidth.Y = 16
	mo.count = 2
	mo.delta = 1
	mo.seesaw = true
	mo.hasIntro = false
	mo.spos.X = 304 + 16
	mo.spos.Y = 368
	mo.lpos.X = 304
	mo.lpos.Y = 368
	mo.ln = 3
	mo.rpos.X = 304
	mo.rpos.Y = 368
	mo.rn = 3
	mo.upos.X = 304
	mo.upos.Y = 368
	mo.un = 3
	mo.dpos.X = 304
	mo.dpos.Y = 368
	mo.dn = 3
	mo.kpos.X = 304 + 3*16
	mo.kpos.Y = 23 * 16
	mo.kn = 7
}
