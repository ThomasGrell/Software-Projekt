package animations

import (
	. "../constants"
	"github.com/faiface/pixel"
	"image/png"
	"log"
	"os"
	"time"
)

var bm, en *enhancedAnimation
var be *basicAnimation

var characterImage *pixel.PictureData
var itemImage *pixel.PictureData

type basicAnimation struct {
	// Position/Bewegung:
	whatAmI uint8

	// Eigenschaften der Animation:

	sprite *pixel.Sprite // Zugehöriger Sprite

	lastUpdate int64 // time.Now().UnixNano() des letzten Spritewechsels
	intervall  int64 // Dauer in Nanosekunden bis zum nächsten Spritewechsel

	count  int8 // Nummer des zuletzt gezeichneten Sprites der Animationssequenz beginnend bei 1
	delta  int8 // entweder -1 oder 1 für Animationsreihenfolge vor und zurück
	seesaw bool // Bei true geht die Animation hin und her (delta=+1/-1). Bei false werden
	// die Sprites immer in derselben Reihenfolge durchlaufen (delta=1)
	hasIntro      bool // True, wenn das Monster eine Animation zum Erscheinen hat.
	introFinished bool
	visible       bool // Sichtbarer Sprite?

	direction uint8 // Bewegungsrichtung (siehe oben definierte Konstanten)
	in        uint8 // Anzahl der Sprites für Erscheinungssequenz
	kn        uint8 // Anzahl der Sprites für Todessequenz
	n         uint8 // Anzahl der Sprites für unbewegte Figur

	ipos pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	kpos pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	pos  pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für den ruhenden Charakter

	width  pixel.Vec // Breite und Höhe des Charakter-Sprites
	iwidth pixel.Vec // Spritegröße für Introanimation
	kwidth pixel.Vec // Spritegröße für Todesanimation
}
type enhancedAnimation struct {
	// Position/Bewegung:
	whatAmI uint8

	// Eigenschaften der Animation:

	sprite *pixel.Sprite // Zugehöriger Sprite

	lastUpdate int64 // time.Now().UnixNano() des letzten Spritewechsels
	intervall  int64 // Dauer in Nanosekunden bis zum nächsten Spritewechsel

	count  int8 // Nummer des zuletzt gezeichneten Sprites der Animationssequenz beginnend bei 1
	delta  int8 // entweder -1 oder 1 für Animationsreihenfolge vor und zurück
	seesaw bool // Bei true geht die Animation hin und her (delta=+1/-1). Bei false werden
	// die Sprites immer in derselben Reihenfolge durchlaufen (delta=1)
	hasIntro      bool // True, wenn das Monster eine Animation zum Erscheinen hat.
	introFinished bool
	visible       bool // Sichtbarer Sprite?

	direction uint8 // Bewegungsrichtung (siehe oben definierte Konstanten)
	in        uint8 // Anzahl der Sprites für Erscheinungssequenz
	kn        uint8 // Anzahl der Sprites für Todessequenz
	n         uint8 // Anzahl der Sprites für unbewegte Figur
	dn        uint8 // Anzahl der Sprites für Abwärtsbewegung
	ln        uint8 // Anzahl der Sprites für Linksbewegung
	rn        uint8 // Anzahl der Sprites für Rechtsbewegung
	un        uint8 // Anzahl der Sprites für Aufwärtsbewegung

	dpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach unten bewegenden Charakter
	lpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach links bewegenden Charakter
	rpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach rechts bewegenden Charakter
	upos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach oben bewegenden Charakter
	ipos pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	kpos pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	pos  pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für den ruhenden Charakter

	width  pixel.Vec // Breite und Höhe des Charakter-Sprites
	iwidth pixel.Vec // Spritegröße für Introanimation
	kwidth pixel.Vec // Spritegröße für Todesanimation
}

func NewBasicAnimation(t uint8) *basicAnimation {
	c := new(basicAnimation)
	c.whatAmI = t
	switch t {
	case Balloon:
		*c = *be
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Teddy:
		*c = *be
		c.pos.Y = 352
		c.kpos.Y = 352
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Ghost:
		*c = *be
		c.pos.Y = 336
		en.kpos.Y = 21 * 16
		en.kn = 9
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Drop:
		*c = *be
		c.pos.Y = 20 * 16
		c.kpos.Y = 20 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Pinky:
		*c = *be
		c.pos.Y = 19 * 16
		c.kpos.Y = 19 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case BluePopEye:
		*c = *be
		c.pos.Y = 18 * 16
		c.kpos.Y = 18 * 16
		c.kn = 9
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Jellyfish:
		*c = *be
		c.pos.Y = 17 * 16
		c.kpos.Y = 17 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Snake:
		*c = *be
		c.pos.Y = 16 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Spinner:
		*c = *be
		c.intervall = 2e8
		c.seesaw = false
		c.pos.X = 20 * 16
		c.pos.Y = 15 * 16
		c.kpos.X = 304 + 4*16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case YellowPopEye:
		*c = *be
		c.pos.Y = 13 * 16
		c.kpos.Y = 13 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case YellowBubble:
		*c = *be
		c.pos.Y = 7 * 16
		c.kpos.Y = 7 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case PinkPopEye:
		*c = *be
		c.pos.Y = 6 * 16
		c.kpos.Y = 6 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Fire:
		*c = *be
		c.pos.Y = 5 * 16
		c.kpos.Y = 5 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Crocodile:
		*c = *be
		c.pos.Y = 4 * 16
		c.kpos.Y = 4 * 16
		c.kn = 9
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Coin:
		*c = *be
		c.seesaw = false
		c.pos.X -= 16
		c.pos.Y = 3 * 16
		c.kpos.Y = 3 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Puddle:
		*c = *be
		c.pos.Y = 2 * 16
		c.kpos.Y = 2 * 16
		c.kn = 6
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case PinkCyclops:
		*c = *be
		c.pos.X = 0
		c.kpos.X = 6 * 16
		c.pos.Y = 15 * 16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case RedCyclops:
		*c = *be
		c.pos.X = 3 * 16
		c.kpos.X = 6 * 16
		c.pos.Y = 15 * 16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case BlueRabbit:
		*c = *be
		c.pos.X = 0
		c.kpos.X = 4 * 16
		c.n = 4
		c.kn = 8
		c.seesaw = false
		c.pos.Y = 13 * 16
		c.kpos.Y = 13 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case PinkFlower:
		*c = *be
		c.direction = Intro
		c.introFinished = false
		c.hasIntro = true
		c.ipos.X = 0
		c.in = 6
		c.pos.X = 6 * 16
		c.kpos.X = 9 * 16
		c.kn = 12
		c.seesaw = false
		c.pos.Y = 12 * 16
		c.kpos.Y = 12 * 16
		c.ipos.Y = 12 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case Fireball:
		*c = *be
		c.direction = Intro
		c.introFinished = false
		c.hasIntro = true
		c.ipos = pixel.V(0, 0)
		c.pos = pixel.V(9*16, 0)
		c.kpos = pixel.V(12*16, 0)
		c.in = 9
		c.kn = 8
		c.width = pixel.V(16, 16)
		c.iwidth = pixel.V(16, 24)
		c.kwidth = pixel.V(16, 16)
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
		c.intervall = 2e8
	case PowerItem:
		c.direction = Stay
		c.intervall = 1e8
		c.hasIntro = false
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(15*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(itemImage, characterImage.Bounds())
		c.seesaw = true
	case BombItem:
		c.direction = Stay
		c.intervall = 1e8
		c.hasIntro = false
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(15*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(itemImage, characterImage.Bounds())
		c.seesaw = true
	default:
		panic("Unknown BasicAnimation")
	}
	c.lastUpdate = time.Now().UnixNano()
	return c
}
func NewEnhancedAnimation(t uint8) *enhancedAnimation {
	c := new(enhancedAnimation)
	c.whatAmI = t

	switch t {
	case WhiteBomberman, WhiteBattleman:
		*c = *bm
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case BlackBomberman, BlackBattleman:
		*c = *bm
		c.pos.Y = bm.pos.Y - 24
		c.upos.Y = bm.upos.Y - 24
		c.dpos.Y = bm.dpos.Y - 24
		c.lpos.Y = bm.lpos.Y - 24
		c.rpos.Y = bm.rpos.Y - 24
		c.kpos.Y = bm.kpos.Y - 24
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case BlueBomberman, BlueBattleman:
		*c = *bm
		c.pos.Y = bm.pos.Y - 24*2
		c.upos.Y = bm.upos.Y - 24*2
		c.dpos.Y = bm.dpos.Y - 24*2
		c.lpos.Y = bm.lpos.Y - 24*2
		c.rpos.Y = bm.rpos.Y - 24*2
		c.kpos.Y = bm.kpos.Y - 24*2
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case RedBomberman, RedBattleman:
		*c = *bm
		c.pos.Y = bm.pos.Y - 24*3
		c.upos.Y = bm.upos.Y - 24*3
		c.dpos.Y = bm.dpos.Y - 24*3
		c.lpos.Y = bm.lpos.Y - 24*3
		c.rpos.Y = bm.rpos.Y - 24*3
		c.kpos.Y = bm.kpos.Y - 24*3
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case Snowy:
		*c = *en
		c.pos.X = 224
		c.pos.Y = 224
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
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case PinkDevil:
		*c = *en
		c.pos.X = 0
		c.dpos.X = 0
		c.lpos.X = 3 * 16
		c.upos.X = 6 * 16
		c.rpos.X = 9 * 16
		c.kpos.X = 12 * 16
		c.pos.Y = 17 * 16
		c.upos.Y = 17 * 16
		c.dpos.Y = 17 * 16
		c.lpos.Y = 17 * 16
		c.rpos.Y = 17 * 16
		c.kpos.Y = 17 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case Penguin:
		*c = *en
		c.pos.X = 16
		c.dpos.X = 0
		c.lpos.X = 6 * 16
		c.upos.X = 3 * 16
		c.rpos.X = 9 * 16
		c.kpos.X = 12 * 16
		c.pos.Y = 16 * 16
		c.upos.Y = 16 * 16
		c.dpos.Y = 16 * 16
		c.lpos.Y = 16 * 16
		c.rpos.Y = 16 * 16
		c.kpos.Y = 16 * 16
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	case BlueCyclops:
		*c = *en
		c.pos.X = 0
		c.dpos.X = 0
		c.lpos.X = 3 * 16
		c.upos.X = 6 * 16
		c.rpos.X = 9 * 16
		c.kpos.X = 0
		c.kn = 14
		c.pos.Y = 10 * 16
		c.upos.Y = 10 * 16
		c.dpos.Y = 10 * 16
		c.lpos.Y = 10 * 16
		c.rpos.Y = 10 * 16
		c.kpos.Y = 8 * 16
		c.width = pixel.V(16, 32)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(characterImage, characterImage.Bounds())
	default:
		panic("Unknown EnhancedAnimation")
	}
	c.lastUpdate = time.Now().UnixNano()
	c.intervall = 2e8
	return c
}

func (c *basicAnimation) Die() {
	c.count = 1
	c.delta = 1
	c.direction = Dead
}
func (c *basicAnimation) GetCenter() (v pixel.Vec) {
	return c.GetWidth().Scaled(0.5)
}
func (c *basicAnimation) GetWidth() (v pixel.Vec) {
	switch c.direction {
	case Dead:
		return c.kwidth
	case Intro:
		return c.iwidth
	default:
		return c.width
	}
}
func (c *basicAnimation) GetSprite() *pixel.Sprite {
	// GetSprite() liefert den aktuell zu zeichnenden Sprite.
	return c.sprite
}
func (c *basicAnimation) GetSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n uint8
	var width pixel.Vec

	// Wenn die Figur ruht, wird stets derselbe Sprite in Blickrichtung der Figur ausgegeben.
	// Bewegt sie sich, so wird die Animation durchlaufen.

	if !c.visible {
		return pixel.R(16*16, 22*16, 17*16, 23*16)
	}

	switch c.direction {
	case Dead:
		v = c.kpos
		n = c.kn
	case Intro:
		v = c.ipos
		n = c.in
	default:
		v = c.pos
		n = c.n
	}
	// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
	if n > 1 {
		timenow := time.Now().UnixNano()
		if timenow-c.lastUpdate > c.intervall {
			c.lastUpdate = timenow
			if uint8(c.count) == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
				if c.direction == Dead {
					c.visible = false
				} else if c.direction == Intro {
					c.introFinished = true
					c.direction = Stay
					c.count = 0
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
			width = c.kwidth
		case Intro:
			v.X += c.iwidth.X * float64(c.count-1)
			width = c.iwidth
		default:
			v.X += c.width.X * float64(c.count-1)
			width = c.width
		}
	}

	return pixel.R(v.X, v.Y, v.X+width.X, v.Y+width.Y)
}
func (c *basicAnimation) IntroFinished() bool { return c.introFinished }
func (c *basicAnimation) IsVisible() bool {
	return c.visible
}
func (c *basicAnimation) SetDirection(direction uint8) {
	// SetDirection() setzt die Bewegungsrichtung neu.
	// Mögliche Eingabewerte sind Stay, Left, Right, Up, Down, Dead.
	c.direction = direction
}
func (c *basicAnimation) SetIntervall(i int64) { c.intervall = i }
func (c *basicAnimation) SetVisible(b bool)    { c.visible = b }
func (c *basicAnimation) Update() {
	if c.whatAmI >= Bomb {
		c.sprite.Set(itemImage, c.GetSpriteCoords())
	} else {
		c.sprite.Set(characterImage, c.GetSpriteCoords())
	}
}

func (c *enhancedAnimation) Die() {
	c.count = 1
	c.delta = 1
	c.direction = Dead
}
func (c *enhancedAnimation) GetCenter() (v pixel.Vec) {
	return c.GetWidth().Scaled(0.5)
}
func (c *enhancedAnimation) GetWidth() (v pixel.Vec) {
	switch c.direction {
	case Dead:
		return c.kwidth
	case Intro:
		return c.iwidth
	default:
		return c.width
	}
}
func (c *enhancedAnimation) GetSprite() *pixel.Sprite {
	// GetSprite() liefert den aktuell zu zeichnenden Sprite.
	return c.sprite
}
func (c *enhancedAnimation) GetSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n uint8
	var width pixel.Vec

	// Wenn die Figur ruht, wird stets derselbe Sprite in Blickrichtung der Figur ausgegeben.
	// Bewegt sie sich, so wird die Animation durchlaufen.

	if !c.visible {
		return pixel.R(16*16, 22*16, 17*16, 23*16)
	}

	switch c.direction {
	case Stay:
		v = c.pos
		n = c.n
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
				} else if c.direction == Intro {
					c.introFinished = true
					c.direction = Stay
					c.count = 0
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
			width = c.kwidth
		case Intro:
			v.X += c.iwidth.X * float64(c.count-1)
			width = c.iwidth
		default:
			v.X += c.width.X * float64(c.count-1)
			width = c.width
		}
	}

	return pixel.R(v.X, v.Y, v.X+width.X, v.Y+width.Y)
}
func (c *enhancedAnimation) IntroFinished() bool { return c.introFinished }
func (c *enhancedAnimation) IsVisible() bool {
	return c.visible
}
func (c *enhancedAnimation) SetDirection(direction uint8) {
	// SetDirection() setzt die Bewegungsrichtung neu.
	// Mögliche Eingabewerte sind Stay, Left, Right, Up, Down, Dead.
	// Im character.png ist bei animierten Charakteren
	// der zweite Sprite stets für die ruhende Figur.
	// Es muss dann die Charakterbreite addiert werden.
	if direction == Stay {
		switch c.direction {
		case Left, Right:
			c.pos = c.lpos
			if c.ln > 1 {
				c.pos.X += c.width.X
			}
		case Up:
			c.pos = c.upos
			if c.un > 1 {
				c.pos.X += c.width.X
			}
		case Down:
			c.pos = c.dpos
			if c.dn > 1 {
				c.pos.X += c.width.X
			}
		}
	}
	c.direction = direction
}
func (c *enhancedAnimation) SetMinPos(v pixel.Vec) {
	c.pos = v
}
func (c *enhancedAnimation) SetIntervall(i int64) { c.intervall = i }
func (c *enhancedAnimation) SetVisible(b bool)    { c.visible = b }
func (c *enhancedAnimation) Update() {
	c.sprite.Set(characterImage, c.GetSpriteCoords())
}
func init() {
	file, err := os.Open("graphics/characters.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	characterImage = pixel.PictureDataFromImage(img)

	file, err = os.Open("graphics/animations.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err = png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	itemImage = pixel.PictureDataFromImage(img)

	bm = new(enhancedAnimation)
	bm.visible = true
	bm.pos.X = 19
	bm.pos.Y = 19
	bm.width.X = 10
	bm.width.Y = 10
	bm.direction = Down
	bm.seesaw = true
	bm.width.X = 16
	bm.width.Y = 24
	bm.count = 2
	bm.delta = 1
	bm.pos.X = 16
	bm.pos.Y = 360
	bm.n = 1
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
	bm.kwidth = bm.width

	// Monster Prototyp
	en = new(enhancedAnimation)
	en.visible = true
	en.width.X = 10
	en.width.Y = 10
	en.direction = Stay
	en.width.X = 16
	en.width.Y = 16
	en.count = 2
	en.delta = 1
	en.seesaw = true
	en.hasIntro = false
	en.pos.X = 304
	en.pos.Y = 368
	en.n = 3
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
	en.kwidth = en.width
	// Monster Prototyp

	be = new(basicAnimation)
	be.visible = true
	be.direction = Stay
	be.width = pixel.V(16, 16)
	be.count = 2
	be.delta = 1
	be.seesaw = true
	be.hasIntro = false
	be.pos = pixel.V(304, 23*16)
	be.n = 3
	be.kpos = pixel.V(304+3*16, 23*16)
	be.kn = 7
	be.kwidth = be.width
}
