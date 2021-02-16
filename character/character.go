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

// Bild, welches die Sprites aller Charaktere enthält
var characterImage *pixel.PictureData
var bm *player
var en *enemy

type player struct {
	character
	bombs    int  // Anzahl der aktuell gelegten Bomben
	kick     bool // kann Bomben wegkicken
	maxBombs int  // maximale Anzahl der legbaren Bomben
	power    int  // Wirkungsradius der Bomben
	wins     int  // Siege für Multi-Player-Modus
}

type enemy struct {
	character
	follow bool // Folgt einem Spieler
}

type character struct {
	animation
	bombghost bool // kann durch Bomben laufen
	life      int  // verbleibende Anzahl der Leben
	mortal    bool // Sterblichkeit
	points    int  // Punkte
	wallghost bool // kann durch Wände laufen
}

type animation struct {
	visible bool    // Unsichtbarer Sprite
	speed   float64 // max. Bewegungsgeschwindigkeit in Pixel pro Sekunde

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

	upos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach oben bewegenden Charakter
	un   int       // Anzahl der Sprites für Animationseffekte

	dpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach unten bewegenden Charakter
	dn   int       // Anzahl der Sprites für Animationseffekte

	ipos   pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	in     int       // Anzahl der Sprites für die Erscheinungssequenz
	iwidth pixel.Vec // Spritegröße für Introanimation

	kpos   pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	kn     int       // Anzahl der Sprites für die Todessequenz
	kwidth pixel.Vec // Spritegröße für Todesanimation

	child *animation // Für den Endgegner, den langen Drachen, welcher aus mehreren Segmenten besteht
}

func NewPlayer(t int) *player {
	c := new(player)

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

	c.sprite = pixel.NewSprite(characterImage, pixel.R(c.spos.X, c.spos.Y, c.spos.X+c.cwidth.X, c.spos.Y+c.cwidth.Y))

	return c
}

func NewEnemy(t int) *enemy {
	c := new(enemy)

	c.lastUpdate = time.Now().UnixNano()
	c.interval = 2e8

	if t == Balloon {
		*c = *en
	}
	if t == Teddy {
		*c = *en
		c.spos.Y = 352
		c.upos.Y = 352
		c.dpos.Y = 352
		c.lpos.Y = 352
		c.rpos.Y = 352
		c.kpos.Y = 352
		c.follow = true
	}
	if t == Ghost {
		*c = *en
		c.spos.Y = 336
		c.upos.Y = 336
		c.dpos.Y = 336
		c.lpos.Y = 336
		c.rpos.Y = 336
		en.kpos.Y = 21 * 16
		en.kn = 9
		c.wallghost = true
	}
	if t == Drop {
		*c = *en
		c.spos.Y = 20 * 16
		c.upos.Y = 20 * 16
		c.dpos.Y = 20 * 16
		c.lpos.Y = 20 * 16
		c.rpos.Y = 20 * 16
		c.kpos.Y = 20 * 16
	}
	if t == Pinky {
		*c = *en
		c.spos.Y = 19 * 16
		c.upos.Y = 19 * 16
		c.dpos.Y = 19 * 16
		c.lpos.Y = 19 * 16
		c.rpos.Y = 19 * 16
		c.kpos.Y = 19 * 16
	}
	if t == BluePopEye {
		*c = *en
		c.spos.Y = 18 * 16
		c.upos.Y = 18 * 16
		c.dpos.Y = 18 * 16
		c.lpos.Y = 18 * 16
		c.rpos.Y = 18 * 16
		c.kpos.Y = 18 * 16
		c.kn = 9
	}
	if t == Jellyfish {
		*c = *en
		c.spos.Y = 17 * 16
		c.upos.Y = 17 * 16
		c.dpos.Y = 17 * 16
		c.lpos.Y = 17 * 16
		c.rpos.Y = 17 * 16
		c.kpos.Y = 17 * 16
	}
	if t == Snake {
		*c = *en
		c.spos.Y = 16 * 16
		c.upos.Y = 16 * 16
		c.dpos.Y = 16 * 16
		c.lpos.Y = 16 * 16
		c.rpos.Y = 16 * 16
	}
	if t == Spinner {
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
	}
	if t == YellowPopEye {
		*c = *en
		c.spos.Y = 13 * 16
		c.upos.Y = 13 * 16
		c.dpos.Y = 13 * 16
		c.lpos.Y = 13 * 16
		c.rpos.Y = 13 * 16
		c.kpos.Y = 13 * 16
	}
	if t == Snowy {
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

	c.sprite = pixel.NewSprite(characterImage, pixel.R(c.spos.X, c.spos.Y, c.spos.X+c.cwidth.X, c.spos.Y+c.cwidth.Y))

	return c
}

func (c *enemy) IsFollowing() bool {
	return c.follow
}

func (c *player) AddPoints(p int) {
	c.points += p
}
func (c *player) GetMaxBombs() int { return c.maxBombs }
func (c *player) GetWins() int     { return c.wins }
func (c *player) IncLife() {
	c.life++
}
func (c *player) IncMaxBombs() { c.maxBombs++ }
func (c *player) IncWins()     { c.wins++ }
func (c *player) ResetWins()   { c.wins = 0 }
func (c *player) SetLife(l int) {
	if l >= 0 {
		c.life = l
	}
}
func (c *player) SetMaxBombs(b int) {
	if b >= 0 {
		c.maxBombs = b
	}
}
func (c *player) SetMortal(b bool) {
	c.mortal = b
}
func (c *player) SetWallghost(w bool) { c.wallghost = w }

func (c *character) DecLife() {
	if c.life == 0 {
		return
	}
	if c.mortal {
		c.life--
		if c.life == 0 {
			c.die()
		}
	}
}
func (c *character) GetPoints() int { return c.points }
func (c *character) IsAlife() bool {
	return c.life > 0
}
func (c *character) IsBombghost() bool { return c.bombghost }
func (c *character) IsMortal() bool {
	return c.mortal
}
func (c *character) IsWallghost() bool   { return c.wallghost }
func (c *character) SetBombghost(b bool) { c.bombghost = b }

func (c *animation) DecSpeed() {
	if c.speed > 10 {
		c.speed -= 10
	}
}
func (c *animation) die() {
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
func (c *animation) GetSpeed() float64 {
	// GetSpeed() gibt die Geschwindigkeit in Pixel/Sekunde zurück.
	return c.speed
}
func (c *animation) GetSprite() *pixel.Sprite {
	// GetSprite() liefert den aktuell zu zeichnenden Sprite.
	return c.sprite
}
func (c *animation) GetSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n int

	// Wenn die Figur ruht, wird stets derselbe Sprite in Blickrichtung der Figur ausgegeben.
	// Bewegt sie sich, so wird die Animation durchlaufen.

	if !c.visible {
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
		if c.direction == Intro {
			v = c.ipos
			n = c.in
		}

		// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
		if n > 1 {
			timenow := time.Now().UnixNano()
			if timenow-c.lastUpdate > c.interval {
				c.lastUpdate = timenow
				if c.count == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
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
func (c *animation) IncSpeed() {
	c.speed += 10
}
func (c *animation) IsVisible() bool {
	return c.visible
}
func (c *animation) SetDirection(direction int) {
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
func (c *animation) SetVisible(b bool) { c.visible = b }
func (c *animation) Update() {
	c.sprite.Set(characterImage, c.GetSpriteCoords())
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
	bm = new(player)
	bm.lastUpdate = time.Now().UnixNano()
	bm.interval = 3e8
	bm.life = 3
	bm.maxBombs = 1
	bm.power = 1
	bm.speed = 10
	bm.kick = false
	bm.mortal = true
	bm.wallghost = false
	bm.bombghost = false
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
	en = new(enemy)
	en.lastUpdate = time.Now().UnixNano()
	en.interval = 2e8
	en.life = 1
	en.speed = 10
	en.mortal = true
	en.wallghost = false
	en.bombghost = false
	en.follow = false
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
	en.spos.X = 304 + 16
	en.spos.Y = 368
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
