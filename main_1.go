package main

/*******************************************************************

 _________________________________
< Gemeinsam implementiert von     >
< Thomas Grell, Sebastian Rösch   >
< und Rayk von Ende               >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

 *******************************************************************/

import (
	"./animations"
	"./characters"
	. "./constants"
	"./gameStat"
	"./level"
	"./sounds"
	txt "./text"
	"./tiles"
	"./titlebar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var bombs []tiles.Bombe
var continu = true
var tb titlebar.Titlebar
var lv gameStat.GameStat
var music sounds.Sound
var levelDef level.Level
var tempAniSlice [][]interface{} // [Animation][Matrix]
var monster []characters.Enemy
var wB characters.Player
var win *pixelgl.Window
var pitchWidth int
var pitchHeight int
var itemBatch *pixel.Batch
var nextLevel bool

var clearingNeeded = false

func loadPic(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func showIntro(win *pixelgl.Window) {
	var zoom float64

	music = sounds.NewSound(ThroughSpace)
	go music.PlaySound()

	pic, err := loadPic("graphics/bomberman.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	win.Clear(colornames.Darkblue)
	win.SetSmooth(true)

	// Startbild: Zoom in
	if win.Bounds().H() > win.Bounds().W() {
		zoom = win.Bounds().W() / pic.Bounds().Size().Len()
	} else {
		zoom = win.Bounds().H() / pic.Bounds().Size().Len()
	}
	for i := float64(0); i <= zoom; i = i + 0.01 {
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
		win.Update()
	}

	// Startbild: Rotate
	for i := float64(0); i <= 12.564; /*6.282*/ i = i + 0.3141 {
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, zoom).Rotated(pixel.ZV, i))
		time.Sleep(1e7)
		win.Update()
	}

	time.Sleep(time.Second)

	txt.Print("Press Space").Draw(win, pixel.IM.Scaled(pixel.ZV, zoom*6).Moved(pixel.V(0, -win.Bounds().H()/4)))

	for !win.Pressed(pixelgl.KeySpace) {
		win.Update()
		time.Sleep(100 * time.Millisecond)
	}
	music.FadeOut()
	fadeOut(win)

	win.SetSmooth(false)
}

func fadeOut(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	imd.SetColorMask(pixel.Alpha(0.05))
	imd.Push(pixel.V(-win.Bounds().W()/2, -win.Bounds().H()/2))
	imd.Push(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2))
	imd.Rectangle(0)
	for i := 0; i < 100; i++ {
		imd.Draw(win)
		win.Update()
		time.Sleep(20 * time.Millisecond)
	}
}
func togglePics(win *pixelgl.Window, sprite1, sprite2 *pixel.Sprite, zoomFactor float64) {
	for {
		sprite1.Draw(win, pixel.IM.Scaled(pixel.ZV, zoomFactor))
		win.Update()
		time.Sleep(2e8)
		sprite2.Draw(win, pixel.IM.Scaled(pixel.ZV, zoomFactor))
		win.Update()
		time.Sleep(2e8)
	}
}

func victory(win *pixelgl.Window) {
	var pic2 pixel.Picture
	pic1, err := loadPic("graphics/Screenshots/victory1.png")
	if err != nil {
		panic(err)
	}
	pic2, err = loadPic("graphics/Screenshots/victory2.png")
	if err != nil {
		panic(err)
	}
	sprite1 := pixel.NewSprite(pic1, pic1.Bounds())
	sprite2 := pixel.NewSprite(pic2, pic2.Bounds())
	win.SetBounds(pixel.R(0, 0, MaxWinSizeX, MaxWinSizeY))
	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))
	win.Clear(colornames.Black)
	win.SetSmooth(true)
	win.Update()
	// victory pic: zoom in
	winSize := win.Bounds().Size()
	picSize := pic1.Bounds().Size()
	zoomFactor := winSize.Len() / picSize.Len()
	for i := float64(0); i <= zoomFactor; i = i + 0.01 {
		sprite1.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
		win.Update()
	}
	// victory pic: toggle
	go togglePics(win, sprite1, sprite2, zoomFactor)
	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		time.Sleep(1e5)
		win.Update()
	}
	win.SetSmooth(false)
}
func gameOver(win *pixelgl.Window) {
	music.StopSound()
	music = sounds.NewSound(JuhaniJunkalaEnd)
	go music.PlaySound()
	win.SetBounds(pixel.R(0, 0, MaxWinSizeX, MaxWinSizeY))
	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))
	win.Update()
	var picGoOn pixel.Picture
	picEnd, err := loadPic("graphics/Screenshots/gameOverEnd.png")
	if err != nil {
		panic(err)
	}
	picGoOn, err = loadPic("graphics/Screenshots/gameOverGoOn.png")
	if err != nil {
		panic(err)
	}
	spriteGoOn := pixel.NewSprite(picGoOn, picGoOn.Bounds())
	spriteEnd := pixel.NewSprite(picEnd, picEnd.Bounds())
	win.Clear(colornames.Black)
	win.SetSmooth(true)
	winSize := win.Bounds().Size()
	picSize := picGoOn.Bounds().Size()
	zoomFactor := winSize.Len() / picSize.Len()
	for i := float64(0); i <= zoomFactor; i = i + 0.01 {
		//spriteGoOn.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
		if win.Pressed(pixelgl.KeyDown) {
			spriteEnd.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
			continu = false
		} else if win.Pressed(pixelgl.KeyUp) {
			spriteGoOn.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
			continu = true
		} else if win.Pressed(pixelgl.KeyEnter) {
			win.SetSmooth(false)
			music.StopSound()
			return
		} else {
			if continu {
				spriteGoOn.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
			} else {
				spriteEnd.Draw(win, pixel.IM.Scaled(pixel.ZV, i))
			}
		}
		win.Update()
	}
	// game over picGoOn: toggle
	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		if win.Pressed(pixelgl.KeyDown) {
			spriteEnd.Draw(win, pixel.IM.Scaled(pixel.ZV, zoomFactor))
			continu = false
		} else if win.Pressed(pixelgl.KeyUp) {
			spriteGoOn.Draw(win, pixel.IM.Scaled(pixel.ZV, zoomFactor))
			continu = true
		} else if win.Pressed(pixelgl.KeyEnter) {
			win.SetSmooth(false)
			music.StopSound()
			return
		}
		time.Sleep(1e7)
		win.Update()
	}
	continu = false
	win.SetSmooth(false)
}

func clearMonsters() {
	remains := make([]characters.Enemy, 0)
	for _, m := range monster {
		if m.IsAlife() || !m.Ani().SequenceFinished() {
			remains = append(remains, m)
		}
	}
	monster = remains[:]
}

func killEnemy(m characters.Enemy) {
	m.DecLife()
	if !m.IsAlife() {
		m.Ani().Die()
	} else {
		m.Ani().SetView(Intro)
	}
	go sounds.NewSound(Falling1).PlaySound()
}

func killPlayer(bm characters.Player) {
	bm.DecLife()
	bm.Ani().Die()
	bm.SetWallghost(false)
	bm.SetBombghost(false)
	bm.SetRemote(false)
	for _, bom := range bombs {
		bom.SetTimeStamp(time.Now())
	}
	go sounds.NewSound(Falling10).PlaySound()
}

// Vor: ...
// Eff: Ist der Countdown der Bombe abgelaufen passiert folgendes:
//     		- Eine neue Explosionsanimation ist erstellt und an die Position der ehemaligen bombe gesetzt.
//      	- Es ertönt der Explosionssound.
//      Ist der Countdown nicht abgelaufen, passiert nichts.
func checkForExplosions() {

	for _, item := range bombs {
		if ((item).GetTimeStamp()).Before(time.Now()) {
			//bomben = append (bomben,item)

			item.Ani().Die()
			b, owner := item.Owner()
			if b {
				owner.DecBombs()
			}

			x, y := lv.A().GetFieldCoord(item.GetPos())
			power := int(item.GetPower())
			l, r, u, d := power, power, power, power

			// Explosion darf nicht über Spielfeldrand hinausragen
			if l > x {
				l = x
			}
			if lv.A().GetWidth()-1-r < x {
				r = lv.A().GetWidth() - 1 - x
			}
			if d > y {
				d = y
			}
			if lv.A().GetHeight()-1-u < y {
				u = lv.A().GetHeight() - 1 - y
			}

			// Falls es Hindernisse gibt, die zerstörbar oder unzerstörbar sind
			bl, xl, yl := lv.GetPosOfNextTile(x, y, pixel.V(float64(-l), 0))
			if bl {
				if lv.IsDestroyableTile(xl, yl) {
					l = x - xl
				} else {
					l = x - xl - 1
				}
			}
			br, xr, yr := lv.GetPosOfNextTile(x, y, pixel.V(float64(r), 0))
			if br {
				if lv.IsDestroyableTile(xr, yr) {
					r = xr - x
				} else {
					r = xr - x - 1
				}
			}
			bd, xd, yd := lv.GetPosOfNextTile(x, y, pixel.V(0, float64(-d)))
			if bd {
				if lv.IsDestroyableTile(xd, yd) {
					d = y - yd
				} else {
					d = y - yd - 1
				}
			}
			bu, xu, yu := lv.GetPosOfNextTile(x, y, pixel.V(0, float64(u)))
			if bu {
				if lv.IsDestroyableTile(xu, yu) {
					u = yu - y
				} else {
					u = yu - y - 1
				}
			}

			// falls sich ein Monster oder Player im Explosionsradius befindet
			bmX, bmY := lv.A().GetFieldCoord(wB.GetPosBox().Center())
			for i := 0; i <= l; i++ {
				for _, m := range monster {
					if !m.Ani().SequenceFinished() {
						continue
					}
					if !m.IsAlife() {
						clearingNeeded = true
						continue
					}
					xx, yy := lv.A().GetFieldCoord(m.GetPosBox().Center())
					if x-i == xx && y == yy {
						l = i
						killEnemy(m)
						wB.AddPoints(m.GetPoints())
						break
					}
				}
				if x-i == bmX && y == bmY && wB.Ani().SequenceFinished() {
					l = i
					killPlayer(wB)
					break
				}
			}
			for i := 1; i <= r; i++ {
				for _, m := range monster {
					if !m.Ani().SequenceFinished() {
						continue
					}
					if !m.IsAlife() {
						clearingNeeded = true
						continue
					}
					xx, yy := lv.A().GetFieldCoord(m.GetPosBox().Center())
					if x+i == xx && y == yy {
						r = i
						killEnemy(m)
						wB.AddPoints(m.GetPoints())
						break
					}
				}
				if x+i == bmX && y == bmY && wB.Ani().SequenceFinished() {
					r = i
					killPlayer(wB)
					break
				}
			}
			for i := 1; i <= u; i++ {
				for _, m := range monster {
					if !m.Ani().SequenceFinished() {
						continue
					}
					if !m.IsAlife() {
						clearingNeeded = true
						continue
					}
					xx, yy := lv.A().GetFieldCoord(m.GetPosBox().Center())
					if y+i == yy && x == xx {
						u = i
						killEnemy(m)
						wB.AddPoints(m.GetPoints() + 100)
						break
					}
				}
				if x == bmX && y+i == bmY && wB.Ani().SequenceFinished() {
					u = i
					killPlayer(wB)
					break
				}
			}
			for i := 1; i <= d; i++ {
				for _, m := range monster {
					if !m.Ani().SequenceFinished() {
						continue
					}
					if !m.IsAlife() {
						clearingNeeded = true
						continue
					}
					xx, yy := lv.A().GetFieldCoord(m.GetPosBox().Center())
					if y-i == yy && x == xx {
						d = i
						killEnemy(m)
						wB.AddPoints(m.GetPoints() + 100)
						break
					}
				}
				if x == bmX && y-i == bmY && wB.Ani().SequenceFinished() {
					d = i
					killPlayer(wB)
					break
				}
			}

			if xl+l == x {
				lv.RemoveTile(xl, yl)
			}
			if xr-r == x {
				lv.RemoveTile(xr, yr)
			}
			if yd+d == y {
				lv.RemoveTile(xd, yd)
			}
			if yu-u == y {
				lv.RemoveTile(xu, yu)
			}

			// Items, die im Explosionsradius liegen werden zerstört, die Explosion wird aber nicht kleiner!

			lv.RemoveItems(x, y, pixel.V(float64(-l), 0))
			lv.RemoveItems(x, y, pixel.V(float64(r), 0))
			lv.RemoveItems(x, y, pixel.V(0, float64(-d)))
			lv.RemoveItems(x, y, pixel.V(0, float64(u)))

			// falls weitere Bomben im Explosionsradius liegen, werden auch gleich explodieren

			for i := 1; i <= l; i++ {
				b, bom := isThereABomb(lv.A().GetFieldCoord(item.GetPos().Add(pixel.V(float64(-i)*TileSize, 0))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= r; i++ {
				b, bom := isThereABomb(lv.A().GetFieldCoord(item.GetPos().Add(pixel.V(float64(i)*TileSize, 0))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= u; i++ {
				b, bom := isThereABomb(lv.A().GetFieldCoord(item.GetPos().Add(pixel.V(0, float64(i)*TileSize))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= d; i++ {
				b, bom := isThereABomb(lv.A().GetFieldCoord(item.GetPos().Add(pixel.V(0, float64(-i)*TileSize))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}

			ani := animations.NewExplosion(uint8(l), uint8(r), uint8(u), uint8(d))
			ani.Show()
			tempAni := make([]interface{}, 2)
			tempAni[0] = ani
			tempAni[1] = (item.GetMatrix()).Moved(ani.ToCenter())
			tempAniSlice = append(tempAniSlice, tempAni)
			s2 := sounds.NewSound(Deathflash)
			go s2.PlaySound()
			if clearingNeeded {
				clearMonsters()
			}
		}
	}

}

// Vor.:...
// Eff.: Nicht explodierte Bomben aus dem Slice existingBombs werden in den
//       Slice remainingBombs kopiert
func removeExplodedBombs(existingBombs []tiles.Bombe) (remainingBombs []tiles.Bombe) {
	j := 0
	for i, bomb := range existingBombs {
		if !bomb.IsVisible() {
			remainingBombs = append(remainingBombs, existingBombs[j:i]...)
			j = i + 1
		}
	}
	remainingBombs = append(remainingBombs, existingBombs[j:]...)
	return remainingBombs
}

func showExplosions(target pixel.Target) {
	for _, a := range tempAniSlice {
		ani := (a[0]).(animations.Animation)
		ani.Update()
		ani.GetSprite().Draw(target, (a[1]).(pixel.Matrix))
	}
}

func clearExplosions(existingExplosions [][]interface{}) (remainingExplosions [][]interface{}) {
	for _, exp := range existingExplosions {
		if exp[0].(animations.Animation).IsVisible() {
			remainingExplosions = append(remainingExplosions, exp)
		}
	}
	return remainingExplosions
}

func isThereABomb(x, y int) (bool, tiles.Bombe) {
	for _, item := range bombs {
		xx, yy := lv.A().GetFieldCoord(item.GetPos())
		if xx == x && yy == y {
			return true, item
		}
	}
	return false, nil
}

func getPossibleDirections(x int, y int, inclBombs bool, inclTiles bool, follow bool) (possibleDir [4]uint8, n uint8) {
	var b = false
	var hori, vert int8

	var isT func(int, int) bool
	if inclTiles {
		isT = lv.IsTile
	} else {
		isT = lv.IsUndestroyableTile
	}

	if follow {
		xb, yb := lv.A().GetFieldCoord(wB.GetPosBox().Center())
		if xb < x {
			hori = -1
		} else if xb > x {
			hori = 1
		}
		if yb < y {
			vert = -1
		} else if yb > y {
			vert = 1
		}
	}

	if inclBombs {
		b, _ = isThereABomb(x-1, y)
	}
	if x != 0 && !isT(x-1, y) && !b && (!follow || hori < 1) {
		possibleDir[n] = Left
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x+1, y)
	}
	if x != lv.A().GetWidth()-1 && !isT(x+1, y) && !b && (!follow || hori > -1) {
		possibleDir[n] = Right
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x, y-1)
	}
	if y != 0 && !isT(x, y-1) && !b && (!follow || vert < 1) {
		possibleDir[n] = Down
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x, y+1)
	}
	if y != lv.A().GetHeight()-1 && !isT(x, y+1) && !b && (!follow || vert > -1) {
		possibleDir[n] = Up
		n++
	}
	return
}

func getNextPosition(c interface{}, dt float64) pixel.Rect {
	dir := c.(characters.Character).GetDirection()
	box := transformRect(dir, c.(characters.Character).GetPosBox())
	return transformRectBack(dir, box.Moved(pixel.Vec{X: 0, Y: c.(characters.Character).GetSpeed()}.Scaled(dt)))
}

// Bewegungen in die 4 Richtungen sind formal identisch, müssen aber
// programmtechnisch unterschiedlich behandelt werden. Die Idee von
// transformRect und transformVec ist es nun, die Koordinaten so zu transformieren, dass
// man alle Berechnungen so ausführen kann, als ob der Character
// aufwärts läuft. Mit transformRectBack und transformVecBack transformiert man alles zurück.
func transformRect(dir uint8, box pixel.Rect) pixel.Rect {
	switch dir {
	case Left:
		return pixel.Rect{Min: pixel.Vec{X: box.Min.Y, Y: -box.Max.X}, Max: pixel.Vec{X: box.Max.Y, Y: -box.Min.X}}
	case Right:
		return pixel.Rect{Min: pixel.Vec{X: -box.Max.Y, Y: box.Min.X}, Max: pixel.Vec{X: -box.Min.Y, Y: box.Max.X}}
	case Down:
		return pixel.Rect{Min: box.Max.Scaled(-1), Max: box.Min.Scaled(-1)}
	default:
		return box
	}
}

func transformRectBack(dir uint8, box pixel.Rect) pixel.Rect {
	switch dir {
	case Right:
		return pixel.Rect{Min: pixel.Vec{X: box.Min.Y, Y: -box.Max.X}, Max: pixel.Vec{X: box.Max.Y, Y: -box.Min.X}}
	case Left:
		return pixel.Rect{Min: pixel.Vec{X: -box.Max.Y, Y: box.Min.X}, Max: pixel.Vec{X: -box.Min.Y, Y: box.Max.X}}
	case Down:
		return pixel.Rect{Min: box.Max.Scaled(-1), Max: box.Min.Scaled(-1)}
	default:
		return box
	}
}

func transformVec(dir uint8, v pixel.Vec) pixel.Vec {
	switch dir {
	case Left:
		return pixel.Vec{X: v.Y, Y: -v.X}
	case Right:
		return pixel.Vec{X: -v.Y, Y: v.X}
	case Down:
		return pixel.Vec{X: -v.X, Y: -v.Y}
	default:
		return v
	}
}

func transformVecBack(dir uint8, v pixel.Vec) pixel.Vec {
	switch dir {
	case Right:
		return pixel.Vec{X: v.Y, Y: -v.X}
	case Left:
		return pixel.Vec{X: -v.Y, Y: v.X}
	case Down:
		return pixel.Vec{X: -v.X, Y: -v.Y}
	default:
		return v
	}
}

func moveCharacter(c interface{}, dt float64) {
	var newDirChoice = false
	chr := c.(characters.Character)
	if !chr.Ani().SequenceFinished() || !chr.IsAlife() {
		return
	}

	nextPos := getNextPosition(c, dt)

	// Blickt man in Bewegungsrichtung, so werden von der hinteren linken Ecke (Min) der PosBox die
	// ganzzahligen Koordinaten im Spielfeld berechnet.
	xnow, ynow := lv.A().GetFieldCoord(transformVecBack(chr.GetDirection(), transformRect(chr.GetDirection(), chr.GetPosBox()).Min))

	// Aus den Koordinaten wird nun eine Spielfeldnummer berechnet.
	newFieldNo := xnow + ynow*lv.A().GetWidth()

	// Koordinaten des Spielfeldes, in welchem sich die vordere rechte Ecke
	// der PosBox in Bezug zur Bewegungsrichtung des Characters befindet
	xv, yv := lv.A().GetFieldCoord(transformVecBack(chr.GetDirection(), transformRect(chr.GetDirection(), chr.GetPosBox()).Max))

	// Abhängig davon, ob der Character durch zerstörbare Wände gehen kann oder nicht, wird die Prüffunktion isT
	// definiert.
	var isT func(int, int) bool
	if chr.IsWallghost() {
		isT = lv.IsUndestroyableTile
	} else {
		isT = lv.IsTile
	}

	var isB func(int, int) (bool, tiles.Bombe)
	if chr.IsBombghost() {
		isB = func(int, int) (bool, tiles.Bombe) { return false, nil }
	} else {
		isB = isThereABomb
	}

	// Versperren Wände oder Bomben den Weg? Falls ja, geht es in dieser Richtung nicht weiter.
	// Eine neue Richtung muss her, also wird newDirChoice auf true gesetzt.
	switch chr.GetDirection() {
	case Left:
		x1, y1 := lv.A().GetFieldCoord(nextPos.Min)
		bombThere1, _ := isB(xv-1, yv)
		bombThere2, _ := isB(xnow-1, ynow)
		if isT(x1, y1) || x1 < 0 || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lv.A().GetFieldCoord(pixel.Vec{X: nextPos.Min.X, Y: nextPos.Max.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Right:
		x1, y1 := lv.A().GetFieldCoord(nextPos.Max)
		bombThere1, _ := isB(xv+1, yv)
		bombThere2, _ := isB(xnow+1, ynow)
		if isT(x1, y1) || x1 > lv.A().GetWidth() || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lv.A().GetFieldCoord(pixel.Vec{X: nextPos.Max.X, Y: nextPos.Min.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Up:
		x1, y1 := lv.A().GetFieldCoord(nextPos.Max)
		bombThere1, _ := isB(xv, yv+1)
		bombThere2, _ := isB(xnow, ynow+1)
		if isT(x1, y1) || y1 > lv.A().GetHeight() || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lv.A().GetFieldCoord(pixel.Vec{X: nextPos.Min.X, Y: nextPos.Max.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Down:
		x1, y1 := lv.A().GetFieldCoord(nextPos.Min)
		bombThere1, _ := isB(xv, yv-1)
		bombThere2, _ := isB(xnow, ynow-1)
		if isT(x1, y1) || y1 < 0 || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lv.A().GetFieldCoord(pixel.Vec{X: nextPos.Max.X, Y: nextPos.Min.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	}

	switch chr := c.(type) {
	case characters.Enemy:

		// War kein Hindernis im Weg?
		if !newDirChoice {
			// Stehe ich auf einem neuen Feld?
			if newFieldNo != chr.GetFieldNo() {
				// Ein neues Feld wurde vollständig betreten.
				// Jetzt ist es Zeit zu überprüfen, ob die Bewegungsrichtung
				// geändert werden kann.
				newDirChoice = true
			}
		}

		// Ein Richtungswechsel steht ggf. an
		if newDirChoice {
			var psblDirs [4]uint8
			var n uint8
			x, y := lv.A().GetFieldCoord(chr.GetPosBox().Center())
			psblDirs, n = getPossibleDirections(x, y, !chr.IsBombghost(), !chr.IsWallghost(), chr.IsFollowing())
			if n == 0 { // keine erlaubte Richtung
				chr.SetDirection(Stay) // Stay
			} else if n == 1 { // 1 erlaubte Richtung --> lauf sie
				chr.SetDirection(psblDirs[0])
			} else if n == 2 { //	2 erlaubte Richtungen
				// Wenn es nur vor oder zurück geht, dann lauf weiter
				if chr.GetDirection() != psblDirs[0] && chr.GetDirection() != psblDirs[1] {
					// Wenn du abbiegen kannst, tu das oder lauf zurück
					chr.SetDirection(psblDirs[rand.Intn(2)])
				}
			} else { // drei oder vier erlaubte Richtungen
				// wähle eine zufällige
				isStraigtForwardPossible := func() bool {
					for _, val := range psblDirs {
						if val == chr.GetDirection() {
							return true
						}
					}
					return false
				}
				if isStraigtForwardPossible() && rand.Intn(int(chr.GetBehaviour())) != 0 {
					newDirChoice = false
					break
				}
				chr.SetDirection(psblDirs[rand.Intn(int(n))])
			}
		}

	case characters.Player:
		// Ist ein Hindernis im Weg?
		if newDirChoice {
			// Und tschüss. Keine Bewegung möglich.
			return
		}
		t, b := lv.CollectItem(xv, yv)
		if b {
			switch t {
			case BombItem:
				chr.IncMaxBombs()
			case LifeItem:
				chr.IncLife()
			case PowerItem:
				chr.IncPower()
			case RollerbladeItem:
				chr.IncSpeed()
			case WallghostItem:
				chr.SetWallghost(true)
			case BombghostItem:
				chr.SetBombghost(true)
			case SkullItem:
				chr.DecLife()
				chr.Ani().Die()
			case HeartItem:
				chr.SetRemote(true)
			case Exit:
				nextLevel = true
				continu = true
			}
		}
	}
	if chr.GetDirection() != Stay {
		chr.Move(transformVecBack(chr.GetDirection(), pixel.Vec{Y: chr.GetSpeed() * dt}))
		chr.SetFieldNo(newFieldNo)
	}
	chr.Ani().SetView(chr.GetDirection())
}

/*
func deathSequence() {
	for !wB.Ani().SequenceFinished() {
		itemBatch.Clear()

		for i := 0; i < pitchHeight; i++ {
			lv.DrawColumn(i, itemBatch)
		}

		itemBatch.Draw(win)

		showExplosions(win)
		tempAniSlice = clearExplosions(tempAniSlice)

		wB.Draw(win)
		for _, m := range monster {
			m.Draw(win)
		}
		win.Update()
	}
}
*/

func setMonster() {
	monster = monster[:0]
	// Enemies from level
	for _, enemyType := range levelDef.GetLevelEnemies() {
		monster = append(monster, characters.NewEnemy(uint8(enemyType)))
	}

	rand.Seed(time.Now().UnixNano())
	//xx, yy := lv.A().GetFieldCoord(wB.GetPos())
	for _, m := range monster {
		for {
			i := rand.Intn(pitchWidth)
			j := rand.Intn(pitchHeight)
			if !lv.IsTile(i, j) && i+j > 4 {
				m.MoveTo(lv.A().GetLowerLeft().Add(pixel.V(float64(i)*TileSize, float64(j)*TileSize)))
				m.Ani().Show()
				break
			}
		}
	}
}
func loadLevel(nr uint8) string {
	switch nr {
	case 1:
		return "./level/stufe_1_level_1.txt"
	case 2:
		return "./level/stufe_1_level_2.txt"
	case 3:
		return "./level/stufe_1_level_3.txt"
	case 4:
		return "./level/stufe_2_level_1.txt"
	case 5:
		return "./level/stufe_2_level_2.txt"
	case 6:
		return "./level/stufe_2_level_3.txt"
	case 7:
		return "./level/stufe_3_level_1.txt"
	case 8:
		return "./level/stufe_3_level_2.txt"
	case 9:
		return "./level/stufe_3_level_3.txt"
	case 10:
		return "./level/stufe_3_level_Boss.txt"
	}
	return "./level/stufe_3_level_Boss.txt"
}
func sun() {
	var levelCount uint8 = 1
	levelDef = level.NewLevel(loadLevel(levelCount))
	pitchWidth, pitchHeight = levelDef.GetBounds()
	var zoomFactor float64
	if float64((pitchHeight+1)*TileSize+32)/float64((pitchWidth+3)*TileSize) > float64(MaxWinSizeY)/MaxWinSizeX {
		zoomFactor = MaxWinSizeY / float64((pitchHeight+1)*TileSize+32)
	} else {
		zoomFactor = MaxWinSizeX / float64((pitchWidth+3)*TileSize)
	}
	var winSizeX = zoomFactor * float64(pitchWidth+3) * TileSize
	var winSizeY = zoomFactor * (float64(pitchHeight+1)*TileSize + 32)
	var err error

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, MaxWinSizeX, MaxWinSizeY),
		VSync:  true,
	}
	win, err = pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	showIntro(win) // INTROOOOOOOOOOOOOOOO

	win.Clear(colornames.Black)
	txt.Print("Let's Play").Draw(win, pixel.IM.Scaled(pixel.V(0, 0), 3))
	win.Update()
	time.Sleep(time.Second * 2)

	fadeOut(win)

	lv = gameStat.NewGameStat(levelDef, 1)

	wB = characters.NewPlayer(WhiteBomberman)
	wB.Ani().Show()
	tb = titlebar.New((3 + uint16(pitchWidth)) * TileSize)
	tb.SetMatrix(pixel.IM.Moved(pixel.V((3+float64(pitchWidth))*TileSize/2, (1+float64(pitchHeight))*TileSize+16)))
	tb.SetLifePointers(wB.GetLifePointer())
	tb.SetPointsPointer(wB.GetPointsPointer())
	tb.SetPlayers(1)
	go tb.Manager()

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	itemBatch = pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)

	// if player wants to continue:
	for continu {
		//println("continue! nextLevel:", nextLevel, "Levelcount:", levelCount)
		if nextLevel {
			music.FadeOut()
			levelCount++
			if levelCount == 11 {
				victory(win)
				break
			}
			nextLevel = false
			win.SetMatrix(pixel.IM)
			win.Clear(colornames.Black)
			txt.Print("Level "+strconv.Itoa(int(levelCount))).Draw(win, pixel.IM.Scaled(pixel.V(0, 0), 3).Moved(win.Bounds().Center()))
			for !win.Pressed(pixelgl.KeySpace) {
				win.Update()
				time.Sleep(time.Millisecond)
			}
		} else {
			levelCount = 1
			wB.Reset()
		}
		levelDef = level.NewLevel(loadLevel(levelCount))
		lv = gameStat.NewGameStat(levelDef, 1)
		pitchWidth, pitchHeight = levelDef.GetBounds()
		if float64((pitchHeight+1)*TileSize+32)/float64((pitchWidth+3)*TileSize) > float64(MaxWinSizeY)/MaxWinSizeX {
			zoomFactor = MaxWinSizeY / float64((pitchHeight+1)*TileSize+32)
		} else {
			zoomFactor = MaxWinSizeX / float64((pitchWidth+3)*TileSize)
		}
		winSizeX = zoomFactor * float64(pitchWidth+3) * TileSize
		winSizeY = zoomFactor * (float64(pitchHeight+1)*TileSize + 32)
		continu = false

		music.StopSound()
		music = sounds.NewSound(levelDef.GetMusic())
		go music.PlaySound()
		win.SetBounds(pixel.R(0, 0, winSizeX, winSizeY))
		win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))
		win.Update()
		//lv.Reset()
		setMonster()
		wB.MoveTo(lv.A().GetLowerLeft())
		wB.SetDirection(Stay)
		wB.Ani().SetView(Down)
		wB.Ani().SetView(Stay)
		wB.Ani().Show()
		//wB.Reset()
		//wB.SetMaxBombs(10)
		//wB.SetPower(10)
		tb.StopCountdown()
		tb.Resize((3 + uint16(pitchWidth)) * TileSize)
		tb.SetMatrix(pixel.IM.Moved(pixel.V((3+float64(pitchWidth))*TileSize/2, (1+float64(pitchHeight))*TileSize+16)))
		tb.SetSeconds(levelDef.GetTime())
		tb.Update()
		tb.StartCountdown()
		win.Update()
		last := time.Now()
		dt := time.Since(last).Seconds()

		for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
			if nextLevel {
				continu = true
				break
			}
			if tb.GetSeconds() == 0 && wB.Ani().GetView() != Dead {
				killPlayer(wB)
			}

			if wB.Ani().GetView() == Dead && wB.Ani().SequenceFinished() && wB.IsAlife() {
				lv.Reset()
				setMonster()
				tb.SetSeconds(levelDef.GetTime())
				tb.StartCountdown()
				wB.MoveTo(lv.A().GetLowerLeft())
				wB.SetDirection(Stay)
				wB.Ani().SetView(Down)
				wB.Ani().SetView(Stay)
				wB.Ani().Show()
			}

			keypressed := false
			dt = time.Since(last).Seconds()
			last = time.Now()

			if win.Pressed(pixelgl.KeyLeft) && wB.Ani().GetView() != Dead {
				wB.SetDirection(Left)
				moveCharacter(wB, dt)
				keypressed = true
			}
			if win.Pressed(pixelgl.KeyRight) && wB.Ani().GetView() != Dead {
				wB.SetDirection(Right)
				moveCharacter(wB, dt)
				keypressed = true
			}
			if win.Pressed(pixelgl.KeyUp) && wB.Ani().GetView() != Dead {
				wB.SetDirection(Up)
				moveCharacter(wB, dt)
				keypressed = true
			}
			if win.Pressed(pixelgl.KeyDown) && wB.Ani().GetView() != Dead {
				wB.SetDirection(Down)
				moveCharacter(wB, dt)
				keypressed = true
			}
			if win.Pressed(pixelgl.KeyX) && wB.Ani().GetView() != Dead && wB.HasRemote() {
				for _, bom := range bombs {
					bom.SetTimeStamp(time.Now())
				}
				keypressed = true
			}

			if !keypressed && wB.Ani().GetView() != Dead && wB.IsAlife() {
				wB.Ani().SetView(Stay)
			}

			if win.JustPressed(pixelgl.KeySpace) && wB.GetBombs() < wB.GetMaxBombs() {
				x, y := lv.A().GetFieldCoord(wB.GetPosBox().Center())
				b, _ := isThereABomb(x, y)
				c := lv.IsTile(x, y)
				if !b && !c && wB.IsAlife() {
					bombs = append(bombs, tiles.NewBomb(wB, lv.A().CoordToVec(x, y)))
					wB.IncBombs()
				}
			}

			/////////////////////////////////////Moving Enemies ///////////////////////////////////////////////////////////

			for _, m := range monster {
				if m.IsAlife() && m.Ani().SequenceFinished() {
					if wB.Ani().GetView() != Dead && wB.GetPosBox().Intersects(m.GetPosBox()) {
						killPlayer(wB)
					}
				}
				moveCharacter(m, dt)
			}

			/*if !wB.Ani().IsVisible() {
				lv.Reset()
				wB.MoveTo(lv.A().GetLowerLeft())
				bombs = bombs[:0]
				tempAniSlice = tempAniSlice[:0]
			}
			*/
			/////////////////////////////////////////////////////////////////////////////////////////////////////////////7

			lv.A().GetCanvas().Draw(win, *(lv.A().GetMatrix()))

			checkForExplosions()
			bombs = removeExplodedBombs(bombs)

			itemBatch.Clear()

			for i := 0; i < pitchHeight; i++ {
				lv.DrawColumn(i, itemBatch)
			}

			for _, item := range bombs {
				item.Draw(itemBatch)
			}

			itemBatch.Draw(win)

			showExplosions(win)

			tempAniSlice = clearExplosions(tempAniSlice)

			wB.Draw(win)

			for _, m := range monster {
				m.Draw(win)
			}

			tb.Draw(win)

			win.Update()

			if !wB.IsAlife() && wB.Ani().SequenceFinished() {
				break
			}

		}
		tb.StopCountdown()
		//if rand.Intn(2) == 1 {
		//	victory(win)
		//}else{
		if !nextLevel {
			gameOver(win)
		}
	}
}

func main() {
	pixelgl.Run(sun)
}
