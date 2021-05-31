package main

import (
	"./animations"
	"./characters"
	. "./constants"
	"./gameStat"
	"./level"
	"./sounds"
	"./tiles"
	"./titlebar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
	"time"
)

var bombs []tiles.Bombe
var tb titlebar.Titlebar
var lev1 gameStat.GameStat
var lv level.Level
var tempAniSlice [][]interface{} // [Animation][Matrix]
var monster []characters.Enemy
var whiteBomberman characters.Player
var win *pixelgl.Window
var pitchWidth int
var pitchHeight int
var itemBatch *pixel.Batch

var clearingNeeded bool = false

func clearMonsters() {
	remains := make([]characters.Enemy, 0)
	for _, m := range monster {
		if m.IsAlife() || !m.Ani().SequenceFinished() {
			remains = append(remains, m)
		}
	}
	monster = remains[:]
}

// Vor: ...
// Eff: Ist der Counddown der Bombe abgelaufen passiert folgendes:
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

			x, y := lev1.A().GetFieldCoord(item.GetPos())
			power := int(item.GetPower())
			l, r, u, d := power, power, power, power

			// Explosion darf nicht über Spielfeldrand hinausragen
			if l > x {
				l = x
			}
			if lev1.A().GetWidth()-1-r < x {
				r = lev1.A().GetWidth() - 1 - x
			}
			if d > y {
				d = y
			}
			if lev1.A().GetHeight()-1-u < y {
				u = lev1.A().GetHeight() - 1 - y
			}

			// Falls es Hindernisse gibt, die zerstörbar oder unzerstörbar sind
			bl, xl, yl := lev1.GetPosOfNextTile(x, y, pixel.V(float64(-l), 0))
			if bl {
				if lev1.IsDestroyableTile(xl, yl) {
					l = x - xl
				} else {
					l = x - xl - 1
				}
			}
			br, xr, yr := lev1.GetPosOfNextTile(x, y, pixel.V(float64(r), 0))
			if br {
				if lev1.IsDestroyableTile(xr, yr) {
					r = xr - x
				} else {
					r = xr - x - 1
				}
			}
			bd, xd, yd := lev1.GetPosOfNextTile(x, y, pixel.V(0, float64(-d)))
			if bd {
				if lev1.IsDestroyableTile(xd, yd) {
					d = y - yd
				} else {
					d = y - yd - 1
				}
			}
			bu, xu, yu := lev1.GetPosOfNextTile(x, y, pixel.V(0, float64(u)))
			if bu {
				if lev1.IsDestroyableTile(xu, yu) {
					u = yu - y
				} else {
					u = yu - y - 1
				}
			}

			// falls sich ein Monster oder Player im Explosionsradius befindet

			bmX, bmY := lev1.A().GetFieldCoord(whiteBomberman.GetPosBox().Center())
			for i := 0; i <= l; i++ {
				for _, m := range monster {
					if !m.Ani().SequenceFinished() {
						continue
					}
					if !m.IsAlife() {
						clearingNeeded = true
						continue
					}
					xx, yy := lev1.A().GetFieldCoord(m.GetPosBox().Center())
					if x-i == xx && y == yy {
						l = i
						m.DecLife()
						if !m.IsAlife() {
							m.Ani().Die()
						}
						whiteBomberman.AddPoints(m.GetPoints())
						break
					}
				}
				if x-i == bmX && y == bmY && whiteBomberman.Ani().SequenceFinished() {
					l = i
					whiteBomberman.DecLife()
					whiteBomberman.Ani().Die()
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
					xx, yy := lev1.A().GetFieldCoord(m.GetPosBox().Center())
					if x+i == xx && y == yy {
						r = i
						m.DecLife()
						if !m.IsAlife() {
							m.Ani().Die()
						}
						whiteBomberman.AddPoints(m.GetPoints())
						break
					}
				}
				if x+i == bmX && y == bmY && whiteBomberman.Ani().SequenceFinished() {
					r = i
					whiteBomberman.DecLife()
					whiteBomberman.Ani().Die()
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
					xx, yy := lev1.A().GetFieldCoord(m.GetPosBox().Center())
					if y+i == yy && x == xx {
						u = i
						m.DecLife()
						if !m.IsAlife() {
							m.Ani().Die()
						}
						whiteBomberman.AddPoints(m.GetPoints() + 100)
						break
					}
				}
				if x == bmX && y+i == bmY && whiteBomberman.Ani().SequenceFinished() {
					u = i
					whiteBomberman.DecLife()
					whiteBomberman.Ani().Die()
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
					xx, yy := lev1.A().GetFieldCoord(m.GetPosBox().Center())
					if y-i == yy && x == xx {
						d = i
						m.DecLife()
						if !m.IsAlife() {
							m.Ani().Die()
						}
						whiteBomberman.AddPoints(m.GetPoints() + 100)
						break
					}
				}
				if x == bmX && y-i == bmY && whiteBomberman.Ani().SequenceFinished() {
					d = i
					whiteBomberman.DecLife()
					whiteBomberman.Ani().Die()
					break
				}
			}

			if xl+l == x {
				lev1.RemoveTile(xl, yl)
			}
			if xr-r == x {
				lev1.RemoveTile(xr, yr)
			}
			if yd+d == y {
				lev1.RemoveTile(xd, yd)
			}
			if yu-u == y {
				lev1.RemoveTile(xu, yu)
			}

			// Items, die im Expolsionsradius liegen werden zerstört, die Expolion wird aber nicht kleiner!

			lev1.RemoveItems(x, y, pixel.V(float64(-l), 0))
			lev1.RemoveItems(x, y, pixel.V(float64(r), 0))
			lev1.RemoveItems(x, y, pixel.V(0, float64(-d)))
			lev1.RemoveItems(x, y, pixel.V(0, float64(u)))

			// falls weitere Bomben im Explosionsradius liegen, werden auch gleich explodieren

			for i := 1; i <= l; i++ {
				b, bom := isThereABomb(lev1.A().GetFieldCoord(item.GetPos().Add(pixel.V(float64(-i)*TileSize, 0))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= r; i++ {
				b, bom := isThereABomb(lev1.A().GetFieldCoord(item.GetPos().Add(pixel.V(float64(i)*TileSize, 0))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= u; i++ {
				b, bom := isThereABomb(lev1.A().GetFieldCoord(item.GetPos().Add(pixel.V(0, float64(i)*TileSize))))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= d; i++ {
				b, bom := isThereABomb(lev1.A().GetFieldCoord(item.GetPos().Add(pixel.V(0, float64(-i)*TileSize))))
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

func showExplosions(win *pixelgl.Window) {
	for _, a := range tempAniSlice {
		ani := (a[0]).(animations.Animation)
		ani.Update()
		ani.GetSprite().Draw(win, (a[1]).(pixel.Matrix))
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
		xx, yy := lev1.A().GetFieldCoord(item.GetPos())
		if xx == x && yy == y {
			return true, item
		}
	}
	return false, nil
}

/*
func isThereABomb(v pixel.Vec) (bool, tiles.Bombe) {
	for _, item := range bombs {
		if item.GetPos() == v {
			return true, item
		}
	}
	return false, nil
}
*/

/*
// Herkunftsrichtung
func homeDir(dir uint8) (hdir uint8) {
	switch dir {
	case Up: hdir = Down
	case Down: hdir = Up
	case Left: hdir = Right
	case Right: hdir = Left
	}
	return
}
*/

/*
// Vor.: /
// Erg.: Die neue Bewegungsrichtung des Monsters ist zurückgegeben.
//		 Kann es sich nicht bewegen, ist die neue Bewegungsrichtung die alte (dann zittert es, weil die Bewegung auf 1 Pixel eingeschränkt ist).
//		 Gibt es nur die Möglichkeit zurück zu laufen, läuft es zurück,
//		 gibt es nur die Möglichkeit weiter oder zurück zu laufen, läuft es weiter,
//		 gibt es mehr als zwei Möglichkeiten, wird eine zufällige, nicht rückwärtsgewandte, Richtung zurückgegeben.
//		 links:0,rechts:1,oben:2,unten:3
func dirChoice(m characters.Enemy) (dir uint8){
	var grDirBool [4]bool = getGrantedDirections(m) // Wahrheitswerte der erlaubten Richtungen, bei Index 0: Up, 1: Down, 2: Left, 3: Right
	var grDirUint8 = make([]uint8,0)	// Zahlenwerte der erlaubten Richtungen (constants): Up: 1, Down: 2, Left: 3, Right: 4
	var n int	// Zählvariable um festzustellen, wie viele mögliche Richtungen es gibt
	for j := range grDirBool {	// Schleife zum Zählen der erlaubten Richtungen und um sie in den Richtungsslice grDirUint8 zu schreiben
		if grDirBool[j] {
			n++
			grDirUint8 = append(grDirUint8, uint8(j+1))
		}
	}
	if n == 0 {	// keine erlaubte Richtung
		dir = 0	// Stay
	}else if n == 1{	// 1 erlaubte Richtung --> lauf sie
		dir = grDirUint8[0]
	}else if n == 2 {	//	2 erlaubte Richtungen
		for _,d := range grDirUint8{	// wenn es nur weiter oder zurück geht, lauf weiter!
			if d == m.GetDirection() {
				dir = d
				return
			}
		}
		choice := rand.Intn(len(grDirUint8)) //	wenn du abbiegen kannst, tu das oder lauf zurück
		dir = grDirUint8[choice]
	}else{	// drei oder vier erlaubte Richtungen
		for i, d := range grDirUint8{
			if d == homeDir(m.GetDirection()) {	// verhindert das Zurücklaufen bei mehr als zwei erlaubten Wegen
				grDirUint8[i] = grDirUint8[len(grDirUint8)-1]
				grDirUint8 = grDirUint8[:len(grDirUint8)-1]
			}
		}
		choice := rand.Intn(len(grDirUint8))	// wähle eine zufällige (außer zurück)
		dir = grDirUint8[choice]
	}
	return
}

*/
func getPossibleDirections(x int, y int, inclBombs bool, inclTiles bool) (possibleDir [4]uint8, n uint8) {
	var b bool = false
	var isT func(int, int) bool
	if inclTiles {
		isT = lev1.IsTile
	} else {
		isT = lev1.IsUndestroyableTile
	}

	if inclBombs {
		b, _ = isThereABomb(x-1, y)
	}
	if x != 0 && !isT(x-1, y) && !b {
		possibleDir[n] = Left
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x+1, y)
	}
	if x != lev1.A().GetWidth()-1 && !isT(x+1, y) && !b {
		possibleDir[n] = Right
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x, y-1)
	}
	if y != 0 && !isT(x, y-1) && !b {
		possibleDir[n] = Down
		n++
	}

	if inclBombs {
		b, _ = isThereABomb(x, y+1)
	}
	if y != lev1.A().GetHeight()-1 && !isT(x, y+1) && !b {
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
	nextPos := getNextPosition(c, dt)
	var newDirChoice bool = false
	chr := c.(characters.Character)
	if !chr.Ani().SequenceFinished() {
		return
	}

	// Blickt man in Bewegungsrichtung, so werden von der hinteren linken Ecke (Min) der PosBox die
	// ganzzahligen Koordinaten im Spielfeld berechnet.
	xnow, ynow := lev1.A().GetFieldCoord(transformVecBack(chr.GetDirection(), transformRect(chr.GetDirection(), chr.GetPosBox()).Min))

	// Aus den Koordinaten wird nun eine Spielfeldnummer berechnet.
	newFieldNo := xnow + ynow*lev1.A().GetWidth()

	if !chr.IsAlife() {
		return
	}

	// Koordinaten des Spielfeldes, in welchem sich die vordere rechte Ecke
	// der PosBox in Bezug zur Bewegungsrichtung des Characters befindet
	xv, yv := lev1.A().GetFieldCoord(transformVecBack(chr.GetDirection(), transformRect(chr.GetDirection(), chr.GetPosBox()).Max))

	// Abhängig davon, ob der Character durch zerstörbare Wände gehen kann oder nicht, wird die Prüffunktion isT
	// definiert.
	var isT func(int, int) bool
	if chr.IsWallghost() {
		isT = lev1.IsUndestroyableTile
	} else {
		isT = lev1.IsTile
	}

	// Versperren Wände den Weg? Falls ja, geht es in dieser Richtung nicht weiter.
	// Eine neue Richtung muss her, also wird newDirChoice auf true gesetzt.
	switch chr.GetDirection() {
	case Left:
		x1, y1 := lev1.A().GetFieldCoord(nextPos.Min)
		bombThere1, _ := isThereABomb(xv-1, yv)
		bombThere2, _ := isThereABomb(xnow-1, ynow)
		if isT(x1, y1) || x1 < 0 || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lev1.A().GetFieldCoord(pixel.Vec{nextPos.Min.X, nextPos.Max.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Right:
		x1, y1 := lev1.A().GetFieldCoord(nextPos.Max)
		bombThere1, _ := isThereABomb(xv+1, yv)
		bombThere2, _ := isThereABomb(xnow+1, ynow)
		if isT(x1, y1) || x1 > lev1.A().GetWidth() || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lev1.A().GetFieldCoord(pixel.Vec{nextPos.Max.X, nextPos.Min.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Up:
		x1, y1 := lev1.A().GetFieldCoord(nextPos.Max)
		bombThere1, _ := isThereABomb(xv, yv+1)
		bombThere2, _ := isThereABomb(xnow, ynow+1)
		if isT(x1, y1) || y1 > lev1.A().GetHeight() || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lev1.A().GetFieldCoord(pixel.Vec{nextPos.Min.X, nextPos.Max.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	case Down:
		x1, y1 := lev1.A().GetFieldCoord(nextPos.Min)
		bombThere1, _ := isThereABomb(xv, yv-1)
		bombThere2, _ := isThereABomb(xnow, ynow-1)
		if isT(x1, y1) || y1 < 0 || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
		x1, y1 = lev1.A().GetFieldCoord(pixel.Vec{nextPos.Max.X, nextPos.Min.Y})
		if isT(x1, y1) || (bombThere1 && bombThere2) {
			newDirChoice = true
		}
	}

	switch c.(type) {
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
			var possibleDirections [4]uint8
			var n uint8
			x, y := lev1.A().GetFieldCoord(chr.GetPosBox().Center())
			possibleDirections, n = getPossibleDirections(x, y, !chr.IsBombghost(), !chr.IsWallghost())
			if n == 0 { // keine erlaubte Richtung
				chr.SetDirection(Stay) // Stay
			} else if n == 1 { // 1 erlaubte Richtung --> lauf sie
				chr.SetDirection(possibleDirections[0])
			} else if n == 2 { //	2 erlaubte Richtungen
				// Wenn es nur vor oder zurück geht, dann lauf weiter
				if chr.GetDirection() != possibleDirections[0] && chr.GetDirection() != possibleDirections[1] {
					// Wenn du abbiegen kannst, tu das oder lauf zurück
					chr.SetDirection(possibleDirections[rand.Intn(2)])
				}
			} else { // drei oder vier erlaubte Richtungen
				// wähle eine zufällige
				chr.SetDirection(possibleDirections[rand.Intn(int(n))])
			}
		}

	case characters.Player:
		// Ist ein Hindernis im Weg?
		if newDirChoice {
			// Und tschüss. Keine Bewegung möglich.
			return
		}
		t, b := lev1.CollectItem(xv, yv)
		if b {
			switch t {
			case BombItem:
				whiteBomberman.IncMaxBombs()
			case LifeItem:
				whiteBomberman.IncLife()
			case PowerItem:
				whiteBomberman.IncPower()
			case RollerbladeItem:
				whiteBomberman.IncSpeed()
			case WallghostItem:
				whiteBomberman.SetWallghost(true)
			case BombghostItem:
				whiteBomberman.SetBombghost(true)
			case SkullItem:
				whiteBomberman.DecLife()
				whiteBomberman.Ani().Die()
			}
		}
	}
	if chr.GetDirection() != Stay {
		chr.Move(transformVecBack(chr.GetDirection(), pixel.Vec{Y: chr.GetSpeed() * dt}))
		chr.SetFieldNo(newFieldNo)
	}
	chr.Ani().SetView(chr.GetDirection())
}

func deathSequence () {
	for !whiteBomberman.Ani().SequenceFinished() {
		itemBatch.Clear()

		for i := 0; i < pitchHeight; i++ {
			lev1.DrawColumn(i, itemBatch)
		}

		itemBatch.Draw(win)

		showExplosions(win)
		tempAniSlice = clearExplosions(tempAniSlice)

		whiteBomberman.Draw(win)
		for _, m := range monster {
			m.Draw(win)
		}
		win.Update()
	}
	
}

func setMonster () {
	monster = monster[:0]
	
	// Enemys from level
	for _, enemyType := range lv.GetLevelEnemys() {
		monster = append(monster, characters.NewEnemy(uint8(enemyType)))
	}
	
	rand.Seed(time.Now().UnixNano())
	xx, yy := lev1.A().GetFieldCoord(whiteBomberman.GetPos())
	for _, m := range monster {
		for {
			i := rand.Intn(pitchWidth)
			j := rand.Intn(pitchHeight)
			if !lev1.IsTile(i, j) && xx != i && yy != j && i+j > 4{
				m.MoveTo(lev1.A().GetLowerLeft().Add(pixel.V(float64(i)*
					TileSize, float64(j)*TileSize)))
				m.Ani().SetVisible(true)
				break
			}
		}
	}
}

func sun() {
	lv = level.NewLevel("./level/level1.txt")
	const typ = 2
	pitchWidth, pitchHeight = lv.GetBounds()
	var zoomFactor float64 = 11/float64(pitchHeight)*3
	var winSizeX float64 = zoomFactor * ((3 + float64(pitchWidth)) * TileSize) // TileSize = 16
	var winSizeY float64 = zoomFactor * ((1+float64(pitchHeight))*TileSize + 32)
	var err error

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err = pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	s1 := sounds.NewSound(lv.GetMusic())
	go s1.PlaySound()


	lev1 = gameStat.NewGameStat(lv, 1)

	whiteBomberman = characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()

	tb = titlebar.New((3 + uint16(pitchWidth)) * TileSize)
	tb.SetMatrix(pixel.IM.Moved(pixel.V((3+float64(pitchWidth))*TileSize/2, (1+float64(pitchHeight))*TileSize+16)))
	tb.SetLifePointers(whiteBomberman.GetLifePointer())
	tb.SetPointsPointer(whiteBomberman.GetPointsPointer())
	tb.SetPlayers(1)
	go tb.Manager()
	tb.SetSeconds(lv.GetTime())
	tb.StartCountdown()
	tb.Update()
	
	setMonster()

	// Bomberman is in lowleft Corner
	whiteBomberman.MoveTo(lev1.A().GetLowerLeft())

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	itemBatch = pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))
	win.Update()
	last := time.Now()
	dt := time.Since(last).Seconds()

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		keypressed := false
		dt = time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyLeft) {
			whiteBomberman.SetDirection(Left)
			moveCharacter(whiteBomberman, dt)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyRight) {
			whiteBomberman.SetDirection(Right)
			moveCharacter(whiteBomberman, dt)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyUp) {
			whiteBomberman.SetDirection(Up)
			moveCharacter(whiteBomberman, dt)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyDown) {
			whiteBomberman.SetDirection(Down)
			moveCharacter(whiteBomberman, dt)
			keypressed = true
		}
		if !keypressed {
			if whiteBomberman.Ani().SequenceFinished() && whiteBomberman.IsAlife() {
				if !whiteBomberman.Ani().IsVisible() {
					whiteBomberman.Ani().SetVisible(true)
					whiteBomberman.Ani().SetView(Down)
				}
				whiteBomberman.Ani().SetView(Stay)
			}
		}
		if win.JustPressed(pixelgl.KeyB) && whiteBomberman.GetBombs() < whiteBomberman.GetMaxBombs() {
			x, y := lev1.A().GetFieldCoord(whiteBomberman.GetPosBox().Center())
			b, _ := isThereABomb(x, y)
			c := lev1.IsTile(x, y)
			if !b && !c && whiteBomberman.IsAlife() {
				bombs = append(bombs, tiles.NewBomb(whiteBomberman, lev1.A().CoordToVec(x, y)))
				whiteBomberman.IncBombs()
			}
		}

		/////////////////////////////////////Moving Enemys ///////////////////////////////////////////////////////////

		for _, m := range monster {
			if whiteBomberman.Ani().IsVisible() && m.IsAlife() && m.Ani().SequenceFinished() {
				if  whiteBomberman.Ani().SequenceFinished() && whiteBomberman.GetPosBox().Intersects(m.GetPosBox()) {
					whiteBomberman.DecLife()
					whiteBomberman.Ani().Die()
				}
			}
			moveCharacter(m, dt)

		}
		
		/*if !whiteBomberman.Ani().IsVisible() {
			lev1.Reset()
			whiteBomberman.MoveTo(lev1.A().GetLowerLeft())
			bombs = bombs[:0]
			tempAniSlice = tempAniSlice[:0]
		}
		*/
		/////////////////////////////////////////////////////////////////////////////////////////////////////////////7

		lev1.A().GetCanvas().Draw(win, *(lev1.A().GetMatrix()))

		checkForExplosions()
		bombs = removeExplodedBombs(bombs)

		itemBatch.Clear()

		for i := 0; i < pitchHeight; i++ {
			lev1.DrawColumn(i, itemBatch)
		}

		for _, item := range bombs {
			item.Draw(itemBatch)
		}

		itemBatch.Draw(win)


		whiteBomberman.Draw(win)
		for _, m := range monster {
			m.Draw(win)
		}

		tb.Draw(win)

		win.Update()

		/*if !whiteBomberman.IsAlife() {
			s3 := sounds.NewSound(Falling1)
			go s3.PlaySound()
		}
		*/
	}
}

func main() {
	pixelgl.Run(sun)
}
