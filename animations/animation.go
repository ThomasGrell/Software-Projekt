package animations

import (
	. "../constants"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

var bm, en *enhancedAnimation
var be *basicAnimation

var CharacterImage *pixel.PictureData
var ItemImage *pixel.PictureData

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
	hasIntro      bool
	introFinished bool
	visible       bool // Sichtbarer Sprite?

	view uint8 // Bewegungsrichtung (siehe oben definierte Konstanten)
	in   uint8 // Anzahl der Sprites für Erscheinungssequenz
	kn   uint8 // Anzahl der Sprites für Todessequenz
	n    uint8 // Anzahl der Sprites für unbewegte Figur

	ipos pixel.Vec // intro position - Pixelgenaue Position des Sprites für Erscheinungsanimation
	kpos pixel.Vec // kill position - Pixelgenaue Position des Sprites für Todessequenz
	pos  pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für den ruhenden Charakter

	width  pixel.Vec // Breite und Höhe des Charakter-Sprites
	iwidth pixel.Vec // Spritegröße für Introanimation
	kwidth pixel.Vec // Spritegröße für Todesanimation
}
type enhancedAnimation struct {
	basicAnimation

	dn uint8 // Anzahl der Sprites für Abwärtsbewegung
	ln uint8 // Anzahl der Sprites für Linksbewegung
	rn uint8 // Anzahl der Sprites für Rechtsbewegung
	un uint8 // Anzahl der Sprites für Aufwärtsbewegung

	dpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach unten bewegenden Charakter
	lpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach links bewegenden Charakter
	rpos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach rechts bewegenden Charakter
	upos pixel.Vec // Pixelgenaue Position des Sprites innerhalb des png für nach oben bewegenden Charakter
}
type explosionAnimation struct {

	// Beim Zünden der Bombe hat sie in die 4 Richtungen ggf. unterschiedliche Ausdehnungen
	lPower uint8
	rPower uint8
	uPower uint8
	dPower uint8

	canvas *pixelgl.Canvas
	batch  *pixel.Batch
	sprite *pixel.Sprite // Zugehöriger Sprite

	lastUpdate int64 // time.Now().UnixNano() des letzten Spritewechsels
	intervall  int64 // Dauer in Nanosekunden bis zum nächsten Spritewechsel

	count int8 // Nummer des zuletzt gezeichneten Sprites der Animationssequenz beginnend bei 0
	delta int8
	// die Sprites immer in derselben Reihenfolge durchlaufen (delta=1)
	visible bool // Sichtbarer Sprite?

	n uint8 // Anzahl der Sprites
}

func NewAnimation(t uint8) Animation {
	switch t {
	case WhiteBomberman, WhiteBattleman:
		c := new(enhancedAnimation)
		c.whatAmI = t
		*c = *bm
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case BlackBomberman, BlackBattleman:
		c := new(enhancedAnimation)
		c.whatAmI = t
		*c = *bm
		c.pos.Y = bm.pos.Y - 24
		c.upos.Y = bm.upos.Y - 24
		c.dpos.Y = bm.dpos.Y - 24
		c.lpos.Y = bm.lpos.Y - 24
		c.rpos.Y = bm.rpos.Y - 24
		c.kpos.Y = bm.kpos.Y - 24
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case BlueBomberman, BlueBattleman:
		c := new(enhancedAnimation)
		c.whatAmI = t
		*c = *bm
		c.pos.Y = bm.pos.Y - 24*2
		c.upos.Y = bm.upos.Y - 24*2
		c.dpos.Y = bm.dpos.Y - 24*2
		c.lpos.Y = bm.lpos.Y - 24*2
		c.rpos.Y = bm.rpos.Y - 24*2
		c.kpos.Y = bm.kpos.Y - 24*2
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case RedBomberman, RedBattleman:
		c := new(enhancedAnimation)
		c.whatAmI = t
		*c = *bm
		c.pos.Y = bm.pos.Y - 24*3
		c.upos.Y = bm.upos.Y - 24*3
		c.dpos.Y = bm.dpos.Y - 24*3
		c.lpos.Y = bm.lpos.Y - 24*3
		c.rpos.Y = bm.rpos.Y - 24*3
		c.kpos.Y = bm.kpos.Y - 24*3
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case Snowy:
		c := new(enhancedAnimation)
		c.whatAmI = t
		*c = *en
		c.pos.X = 208 + 16
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
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case PinkDevil:
		c := new(enhancedAnimation)
		c.whatAmI = t
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
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case Penguin:
		c := new(enhancedAnimation)
		c.whatAmI = t
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
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case BlueCyclops:
		c := new(enhancedAnimation)
		c.whatAmI = t
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
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		c.intervall = 2e8
		return c
	case Balloon:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Teddy:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 352
		c.kpos.Y = 352
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Ghost:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 336
		en.kpos.Y = 21 * 16
		en.kn = 9
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Drop:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 20 * 16
		c.kpos.Y = 20 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Pinky:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 19 * 16
		c.kpos.Y = 19 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case BluePopEye:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 18 * 16
		c.kpos.Y = 18 * 16
		c.kn = 9
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Jellyfish:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 17 * 16
		c.kpos.Y = 17 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Snake:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 16 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Spinner:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.intervall = 2e8
		c.seesaw = false
		c.pos.Y = 15 * 16
		c.n = 4
		c.kpos.X = 304 + 4*16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.lastUpdate = time.Now().UnixNano()
		return c
	case YellowPopEye:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 13 * 16
		c.kpos.Y = 13 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case YellowBubble:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 7 * 16
		c.kpos.Y = 7 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case PinkPopEye:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 6 * 16
		c.kpos.Y = 6 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Fire:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 5 * 16
		c.kpos.Y = 5 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Crocodile:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 4 * 16
		c.kpos.Y = 4 * 16
		c.kn = 9
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Coin:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.seesaw = false
		c.pos.X -= 16
		c.pos.Y = 3 * 16
		c.kpos.Y = 3 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Puddle:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.Y = 2 * 16
		c.kpos.Y = 2 * 16
		c.kn = 6
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case PinkCyclops:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.X = 0
		c.kpos.X = 6 * 16
		c.pos.Y = 15 * 16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case RedCyclops:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.X = 3 * 16
		c.kpos.X = 6 * 16
		c.pos.Y = 15 * 16
		c.kpos.Y = 15 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case BlueRabbit:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.pos.X = 0
		c.kpos.X = 4 * 16
		c.n = 4
		c.kn = 8
		c.seesaw = false
		c.pos.Y = 13 * 16
		c.kpos.Y = 13 * 16
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case PinkFlower:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.view = Intro
		c.introFinished = false
		c.hasIntro = true
		c.pos = pixel.V(6*16, 12*16)
		c.ipos = pixel.V(0, 12*16)
		c.kpos = pixel.V(9*16, 12*16)
		c.iwidth = pixel.V(16, 16)
		c.in = 6
		c.kn = 12
		c.seesaw = true
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Fireball:
		c := new(basicAnimation)
		c.whatAmI = t
		*c = *be
		c.view = Intro
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
		c.sprite = pixel.NewSprite(CharacterImage, CharacterImage.Bounds())
		c.intervall = 2e8
		c.lastUpdate = time.Now().UnixNano()
		return c
	case PowerItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(15*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case BombItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(15*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case PunchItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(15*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case HeartItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(17*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case RollerbladeItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(17*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case SkullItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(19*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case WallghostItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(19*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case BombghostItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(19*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case LifeItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(21*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Exit:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(21*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case KickItem:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(23*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 2
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Bomb:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(9*16, 17*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 3
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Stub:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(22*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Brushwood:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(21*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Greenwall:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(25*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Greywall:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(26*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Brownwall:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(26*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Darkbrownwall:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(25*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Evergreen:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(25*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Tree:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(26*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Palmtree:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(27*16, 3*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Perl:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(26*16, 7*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Snowrock:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(27*16, 2*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Greenrock:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(25*16, 6*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case House:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(27*16, 6*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Christmastree:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(27*16, 4*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 32)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Perl1:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(23*16, 5*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Perl2:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(24*16, 5*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Perl3:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(25*16, 5*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Perl4:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(26*16, 5*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	case Littlesnowrock:
		c := new(basicAnimation)
		c.whatAmI = t
		c.view = Stay
		c.intervall = 1e8
		c.visible = true
		c.count = 1
		c.delta = 1
		c.pos = pixel.V(24*16, 6*16)
		c.kpos = pixel.V(0, 4*16)
		c.n = 1
		c.kn = 7
		c.width = pixel.V(16, 16)
		c.kwidth = pixel.V(32, 32)
		c.sprite = pixel.NewSprite(ItemImage, ItemImage.Bounds())
		c.seesaw = true
		c.lastUpdate = time.Now().UnixNano()
		return c
	default:
		panic("Unknown Animation")
	}

	// This line is never reached.
	return &basicAnimation{}
}
func NewExplosion(l, r, u, d uint8) Animation {
	b := new(explosionAnimation)
	b.lPower = l
	b.rPower = r
	b.uPower = u
	b.dPower = d
	b.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(16*(r+l+1)), float64(16*(u+d+1))))
	b.sprite = pixel.NewSprite(b.canvas, b.canvas.Bounds())
	b.batch = pixel.NewBatch(&pixel.TrianglesData{}, ItemImage)
	b.delta = 1
	b.count = 0
	b.drawCanvas()
	b.lastUpdate = time.Now().UnixNano()
	b.intervall = 2e7
	b.n = 4
	b.visible = false
	return b
}

func (c *explosionAnimation) getSpriteCenter(n uint8) pixel.Vec {
	return pixel.V(2*16+float64(n)*5*16, 8*16)
}
func (c *explosionAnimation) Die() {
}
func (c *explosionAnimation) ToCenter() pixel.Vec {
	return c.canvas.Bounds().Center().Sub(pixel.V(float64(c.lPower)*16+8, float64(c.dPower)*16+8))
}
func (c *explosionAnimation) ToBaseline() pixel.Vec {
	return c.ToCenter().Add(pixel.V(0, 8))
}
func (c *explosionAnimation) GetSize() (v pixel.Vec) {
	return c.canvas.Bounds().Size()
}
func (c *explosionAnimation) GetSprite() *pixel.Sprite {
	// GetSprite() liefert den aktuell zu zeichnenden Sprite.
	return c.sprite
}
func (c *explosionAnimation) IntroFinished() bool  { return true }
func (c *explosionAnimation) IsVisible() bool      { return c.visible }
func (c *explosionAnimation) SetView(uint8)        {}
func (c *explosionAnimation) SetIntervall(i int64) { c.intervall = i }
func (c *explosionAnimation) SetVisible(b bool)    { c.visible = b }
func (c *explosionAnimation) Show() {
	c.visible = true
	c.lastUpdate = time.Now().UnixNano()
}
func (c *explosionAnimation) drawCanvas() {
	c.batch.Clear()
	c.canvas.Clear(color.Transparent)

	if !c.visible {
		return
	}

	w := 16 * float64(c.lPower+c.rPower+1)
	h := 16 * float64(c.uPower+c.dPower+1)

	// Zeichne linke Flammenspitze
	s := pixel.NewSprite(ItemImage, pixel.R(float64(c.count)*5*16, 8*16, float64(c.count)*5*16+16, 9*16))
	s.Draw(c.batch, pixel.IM.Moved(pixel.V(8, float64(c.dPower)*16+8)))

	// Zeichne untere Flammenspitze
	s.Set(ItemImage, pixel.R(2*16+float64(c.count)*5*16, 6*16, 2*16+float64(c.count)*5*16+16, 7*16))
	s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(c.lPower)*16+8, 8)))

	// Zeichne obere Flammenspitze
	s.Set(ItemImage, pixel.R(2*16+float64(c.count)*5*16, 10*16, 3*16+float64(c.count)*5*16, 16*11))
	s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(c.lPower)*16+8, h-8)))

	// Zeichne rechte Flammenspitze
	s.Set(ItemImage, pixel.R(4*16+float64(c.count)*5*16, 8*16, 5*16+float64(c.count)*5*16, 9*16))
	s.Draw(c.batch, pixel.IM.Moved(pixel.V(w-8, float64(c.dPower)*16+8)))

	// Zeichne Mitte der Flamme
	s.Set(ItemImage, pixel.R(2*16+float64(c.count)*5*16, 8*16, 3*16+float64(c.count)*5*16, 9*16))
	s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(c.lPower*16+8), float64(c.dPower)*16+8)))

	// Zeichne linken Explosionsast
	for i := uint8(1); i < c.lPower; i++ {
		s.Set(ItemImage, pixel.R(1*16+float64(c.count)*5*16, 8*16, 2*16+float64(c.count)*5*16, 9*16))
		s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(i*16+8), float64(c.dPower)*16+8)))
	}

	// Zeichne rechten Explosionsast
	for i := uint8(1); i < c.rPower; i++ {
		s.Set(ItemImage, pixel.R(3*16+float64(c.count)*5*16, 8*16, 4*16+float64(c.count)*5*16, 9*16))
		s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64((c.lPower+i)*16+8), float64(c.dPower)*16+8)))
	}

	// Zeichne unteren Explosionsast
	for i := uint8(1); i < c.dPower; i++ {
		s.Set(ItemImage, pixel.R(2*16+float64(c.count)*5*16, 7*16, 3*16+float64(c.count)*5*16, 8*16))
		s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(c.lPower*16+8), float64(i*16+8))))
	}

	// Zeichne oberen Explosionsast
	for i := uint8(1); i < c.uPower; i++ {
		s.Set(ItemImage, pixel.R(2*16+float64(c.count)*5*16, 9*16, 3*16+float64(c.count)*5*16, 10*16))
		s.Draw(c.batch, pixel.IM.Moved(pixel.V(float64(c.lPower*16+8), float64((i+c.dPower)*16+8))))
	}

	c.batch.Draw(c.canvas)
}
func (c *explosionAnimation) Update() {
	if !c.visible {
		return
	}
	// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
	timenow := time.Now().UnixNano()
	if timenow-c.lastUpdate > c.intervall {
		c.lastUpdate = timenow
		if c.count == int8(c.n) { // rechts angekommen in der Bildfolge
			c.delta = -1
			c.count += c.delta
		} else if c.count == 0 && c.delta == -1 {
			c.visible = false
		} else {
			c.count += c.delta
		}
	}

	c.drawCanvas()
	c.sprite.Set(c.canvas, c.canvas.Bounds())
}

func (c *basicAnimation) Die() {
	if c.whatAmI == Bomb {
		c.visible = false
	} // Explosion einer Bombe benötigt eine spezielle Explosionsroutine, keine Standard-Todessequenz
	c.count = 1
	c.delta = 1
	c.view = Dead
}
func (c *basicAnimation) ToCenter() pixel.Vec {
	return c.GetSize().Scaled(0.5)
}
func (c *basicAnimation) ToBaseline() pixel.Vec {
	return pixel.V(0, c.GetSize().Y/2)
}
func (c *basicAnimation) GetSize() (v pixel.Vec) {
	switch c.view {
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
func (c *basicAnimation) getSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n uint8
	var width pixel.Vec
	if !c.visible {
		return pixel.R(0, 0, 1, 1)
	}
	switch c.view {
	case Dead:
		v = c.kpos
		n = c.kn
		width = c.kwidth
	case Intro:
		v = c.ipos
		n = c.in
		width = c.iwidth
	default:
		v = c.pos
		n = c.n
		width = c.width
	}
	// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
	if n > 1 {
		timenow := time.Now().UnixNano()
		if timenow-c.lastUpdate > c.intervall {
			c.lastUpdate = timenow
			if uint8(c.count) == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
				if c.view == Dead {
					c.visible = false
				} else if c.view == Intro {
					c.introFinished = true
					c.view = Stay
					c.count = 1
					v = c.pos
					n = c.n
					width = c.width
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
		v.X += width.X * float64(c.count-1)
	}

	return pixel.R(v.X, v.Y, v.X+width.X, v.Y+width.Y)
}
func (c *basicAnimation) IntroFinished() bool { return c.introFinished }
func (c *basicAnimation) IsVisible() bool {
	return c.visible
}
func (c *basicAnimation) SetView(view uint8) {
	// SetView() setzt den View neu.
	// Mögliche Eingabewerte sind Stay, Left, Right, Up, Down, Dead, Intro.
	if view == Intro {
		if c.hasIntro {
			c.introFinished = false
		} else {
			c.introFinished = true
			view = Stay
		}
	}
	c.view = view
	c.count = 1
}
func (c *basicAnimation) SetIntervall(i int64) { c.intervall = i }
func (c *basicAnimation) SetVisible(b bool)    { c.visible = b }
func (c *basicAnimation) Show() {
	c.visible = true
	c.lastUpdate = time.Now().UnixNano()
}
func (c *basicAnimation) Update() {
	r := c.getSpriteCoords()
	if !c.visible {
		c.sprite.Set(pixel.MakePictureData(r), r)
	} else if c.whatAmI >= Bomb {
		c.sprite.Set(ItemImage, r)
	} else {
		c.sprite.Set(CharacterImage, r)
	}
}

func (c *enhancedAnimation) getSpriteCoords() pixel.Rect {
	var v pixel.Vec
	var n uint8
	var width pixel.Vec

	if !c.visible {
		return pixel.R(0, 0, 1, 1)
	}

	switch c.view {
	case Stay:
		v = c.pos
		n = c.n
		width = c.width
	case Up:
		v = c.upos
		n = c.un
		width = c.width
	case Down:
		v = c.dpos
		n = c.dn
		width = c.width
	case Left:
		v = c.lpos
		n = c.ln
		width = c.width
	case Right:
		v = c.rpos
		n = c.rn
		width = c.width
	case Dead:
		v = c.kpos
		n = c.kn
		width = c.kwidth
	case Intro:
		v = c.ipos
		n = c.in
		width = c.iwidth
	}
	// Es wird geprüft, ob das nächste Sprite der Animation gezeigt werden muss, falls es eines gibt.
	if n > 1 {
		timenow := time.Now().UnixNano()
		if timenow-c.lastUpdate > c.intervall {
			c.lastUpdate = timenow
			if uint8(c.count) == n { // rechts angekommen in der Bildfolge --> Rückwärtsgang
				if c.view == Dead {
					c.visible = false
				} else if c.view == Intro {
					c.introFinished = true
					c.view = Stay
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
		v.X += width.X * float64(c.count-1)
	}

	return pixel.R(v.X, v.Y, v.X+width.X, v.Y+width.Y)
}
func (c *enhancedAnimation) SetView(view uint8) {
	if view == Stay {
		switch c.view {
		case Left:
			c.pos = c.lpos
			if c.ln > 1 {
				c.pos.X += c.width.X
			}
		case Right:
			c.pos = c.rpos
			if c.rn > 1 {
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
	if view == Intro {
		if c.hasIntro {
			c.introFinished = false
		} else {
			c.introFinished = true
			view = Stay
		}
	}
	if c.view != view {
		c.view = view
		c.count = 1
	}
}
func (c *enhancedAnimation) Update() {
	r := c.getSpriteCoords()
	if c.visible {
		c.sprite.Set(CharacterImage, r)
	} else {
		c.sprite.Set(pixel.MakePictureData(r), r)
	}
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
	CharacterImage = pixel.PictureDataFromImage(img)

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
	ItemImage = pixel.PictureDataFromImage(img)

	bm = new(enhancedAnimation)
	bm.visible = false
	bm.pos.X = 19
	bm.pos.Y = 19
	bm.width.X = 10
	bm.width.Y = 10
	bm.view = Down
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
	en.visible = false
	en.width.X = 10
	en.width.Y = 10
	en.view = Stay
	en.width.X = 16
	en.width.Y = 16
	en.count = 2
	en.delta = 1
	en.seesaw = true
	en.hasIntro = false
	en.introFinished = true
	en.pos.X = 304
	en.pos.Y = 368
	en.n = 1
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
	be.count = 1
	be.delta = 1
	be.seesaw = true
	be.visible = false
	be.hasIntro = false
	be.introFinished = true
	be.view = Stay
	be.width = pixel.V(16, 16)
	be.kwidth = pixel.V(16, 16)
	be.pos = pixel.V(304, 23*16)
	be.kpos = pixel.V(304+3*16, 23*16)
	be.n = 3
	be.kn = 7
}
