package main

import (
	"./animations"
	"./arena"
	"./characters"
	. "./constants"
	"./level1"
	"./sounds"
	"./tiles"
	"./titlebar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
	"time"
)


var bombs []tiles.Bombe
var tb titlebar.Titlebar
var arena1 arena.Arena
var lev1 level1.Level

var tempAniSlice [][]interface{} // [Animation][Matrix]
var monster []characters.Enemy
var whiteBomberman characters.Player

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

			x, y := arena1.GetFieldCoord(item.GetPos())
			power := int(item.GetPower())
			l, r, u, d := power, power, power, power

			// Explosion darf nicht über Spielfeldrand hinausragen
			if l > x {
				l = x
			}
			if arena1.GetWidth()-1-r < x {
				r = arena1.GetWidth() - 1 - x
			}
			if d > y {
				d = y
			}
			if arena1.GetHeight()-1-u < y {
				u = arena1.GetHeight() - 1 - y
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

		A:
			for i := 0; i <= l; i++ {
				for j, m := range monster {
					xx, yy := arena1.GetFieldCoord(m.GetPos())
					if x-i == xx && y == yy {
						l = i
						m.Ani().Die()
						monster[j], monster[len(monster)-1] = monster[len(monster)-1], monster[j]
						monster = monster[:len(monster)-1]
						whiteBomberman.AddPoints(m.GetPoints()+100)
						break A
					}
					bmX, bmY := arena1.GetFieldCoord(whiteBomberman.GetPos())
					if x-i == bmX && y == bmY {
						l = i
						whiteBomberman.DecLife()
						if whiteBomberman.GetLife() == 0 {whiteBomberman.Ani().Die()}
						break A
					}
				}
			}

		B:
			for i := 0; i <= r; i++ {
				for j, m := range monster {
					xx, yy := arena1.GetFieldCoord(m.GetPos())
					if x+i == xx && y == yy {
						r = i
						m.Ani().Die()
						monster[j], monster[len(monster)-1] = monster[len(monster)-1], monster[j]
						monster = monster[:len(monster)-1]
						whiteBomberman.AddPoints(m.GetPoints()+100)
						break B
					}
					bmX, bmY := arena1.GetFieldCoord(whiteBomberman.GetPos())
					if x+i == bmX && y == bmY {
						r = i
						whiteBomberman.DecLife()
						if whiteBomberman.GetLife() == 0 {whiteBomberman.Ani().Die()}
						break B
					}
				}
			}

		C:
			for i := 0; i <= u; i++ {
				for j, m := range monster {
					xx, yy := arena1.GetFieldCoord(m.GetPos())
					if y+i == yy && x == xx {
						u = i
						m.Ani().Die()
						monster[j], monster[len(monster)-1] = monster[len(monster)-1], monster[j]
						monster = monster[:len(monster)-1]
						whiteBomberman.AddPoints(m.GetPoints()+100)
						break C
					}
					bmX, bmY := arena1.GetFieldCoord(whiteBomberman.GetPos())
					if x == bmX && y+i == bmY {
						u = i
						whiteBomberman.DecLife()
						if whiteBomberman.GetLife() == 0 {whiteBomberman.Ani().Die()}
						break C
					}
				}
			}

		D:
			for i := 0; i <= d; i++ {
				for j, m := range monster {
					xx, yy := arena1.GetFieldCoord(m.GetPos())
					if y-i == yy && x == xx {
						d = i
						m.Ani().Die()
						monster[j], monster[len(monster)-1] = monster[len(monster)-1], monster[j]
						monster = monster[:len(monster)-1]
						whiteBomberman.AddPoints(m.GetPoints()+100)
						break D
					}
					bmX, bmY := arena1.GetFieldCoord(whiteBomberman.GetPos())
					if x == bmX && y-i == bmY {
						d = i
						whiteBomberman.DecLife()
						if whiteBomberman.GetLife() == 0 {whiteBomberman.Ani().Die()}
						break D
					}
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
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(float64(-i)*TileSize, 0)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= r; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(float64(i)*TileSize, 0)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= u; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(0, float64(i)*TileSize)))
				if b {
					bom.SetTimeStamp(time.Now())
				}
			}
			for i := 1; i <= d; i++ {
				b, bom := isThereABomb(item.GetPos().Add(pixel.V(0, float64(-i)*TileSize)))
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

func isThereABomb(v pixel.Vec) (bool, tiles.Bombe) {
	for _, item := range bombs {
		if item.GetPos() == v {
			return true, item
		}
	}
	return false, nil
}
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
// Vor.: /
// Erg.: Die neue Bewegungsrichtung des Monsters ist zurückgegeben.
//		 Kann es sich nicht bewegen, ist die neue Bewegungsrichtung die alte (dann zittert es, weil die Bewegung auf 1 Pixel eingeschränkt ist).
//		 Gibt es nur die Möglichkeit zurück zu laufen, läuft es zurück,
//		 gibt es nur die Möglichkeit weiter oder zurück zu laufen, läuft es weiter,
//		 gibt es mehr als zwei Möglichkeiten, wird eine zufällige, nicht rückwärtsgewandte, Richtung zurückgegeben.
//		 links:0,rechts:1,oben:2,unten:3
func dirChoice(m characters.Enemy) (dir uint8){
	var grDirBool [4]bool = getGrantedDirections(m) // Wahrheitswerte der erlaubten Richtungen
	//fmt.Println("grDirBool",grDirBool)
	var grDirUint8 = make([]uint8,0)	// Zahlenwerte der erlaubten Richtungen (constants)
	var n int	// Zählvariable um festzustellen, wie viele mögliche Richtungen es gibt
	for j := range grDirBool {
		if grDirBool[j] {
			n++
			grDirUint8 = append(grDirUint8, uint8(j+1))
		}
	}
	if n == 0 {
		dir = 0
	}else if n == 1{
		dir = grDirUint8[0]
	}else if n == 2 {
		for _,d := range grDirUint8{
			if d == m.GetDirection() {
				dir = d
				return
			}
		}
		choice := rand.Intn(len(grDirUint8))
		dir = grDirUint8[choice]
	}else{
		for i, d := range grDirUint8{
			if d == homeDir(m.GetDirection()) {	// verhindert das Zurücklaufen bei mehr als zwei erlaubten Wegen
				grDirUint8[i] = grDirUint8[len(grDirUint8)-1]
				grDirUint8 = grDirUint8[:len(grDirUint8)-1]
			}
		}
		choice := rand.Intn(len(grDirUint8))
		dir = grDirUint8[choice]
	}
	return
}
//func dirChoice(monster characters.Enemy) (dir uint8) {
//	grDir := getGrantedDirections(monster)
//	var grDirInt = make([]uint8, 0)
//	var goingBack uint8
//	var n uint8
//	for j := range grDir {
//		if grDir[j] {
//			n++
//			switch monster.GetDirection() {
//			case Left:
//				if j != Right {
//					grDirInt = append(grDirInt, uint8(j))
//				} else {
//					goingBack = 1
//				}
//			case Right:
//				if j != Left {
//					grDirInt = append(grDirInt, uint8(j))
//				} else {
//					goingBack = 0
//				}
//			case Up:
//				if j != Down {
//					grDirInt = append(grDirInt, uint8(j))
//				} else {
//					goingBack = 3
//				}
//			case Down:
//				if j != Up {
//					grDirInt = append(grDirInt, uint8(j))
//				} else {
//					goingBack = 2
//				}
//			}
//		}
//	}
//	if n > 2 {
//		choice := rand.Intn(len(grDirInt))
//		dir = grDirInt[choice]
//	} else if n == 1 {
//		dir = goingBack
//	} else {
//		dir = monster.GetDirection()
//	}
//	return
//}
//func getGrantedDirections(c characters.Character) [4]bool {
//	var b [4]bool
//	b[0] = true
//	b[1] = true
//	b[2] = true
//	b[3] = true
//	pb := c.GetPosBox()
//	xll,yll := arena1.GetFieldCoord(pb.Min)
//	xur,yur := arena1.GetFieldCoord(pb.Max)
//	ll := pixel.V(float64(xll*TileSize),float64(yll*TileSize)).Sub(pixel.V(0,2))
//	ur := pixel.V(float64(xur*TileSize),float64(yur*TileSize)).Sub(pixel.V(0,2))
//	x,y := arena1.GetFieldCoord(c.GetPos())
//	fmt.Println("Coords",x,y,"Laut PosBox ll",ll.X/TileSize,ll.Y/TileSize,"ur",ur.X/TileSize,ur.Y/TileSize )
//	if lev1.IsTile(int((ll.X-1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ll.X-1)/TileSize), int(ur.Y/TileSize)) || ll.X-1 < 0 {
//		b[2] = false
//	}
//	if int((ur.X+1)/TileSize) > arena1.GetWidth()-1 {
//		b[3] = false
//	} else if lev1.IsTile(int((ur.X+1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ur.X+1)/TileSize), int(ur.Y/TileSize)) {
//		b[3] = false
//	}
//	if int((ur.Y+1)/TileSize) > arena1.GetHeight()-1 {
//		b[0] = false
//	} else if lev1.IsTile(int(ll.X/TileSize), int((ur.Y+1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ur.Y+1)/TileSize)) {
//		b[0] = false
//	}
//	if lev1.IsTile(int(ll.X/TileSize), int((ll.Y-1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ll.Y-1)/TileSize)) || ll.Y-1 < 0 {
//		b[1] = false
//	}
//	return b
//}
func getGrantedDirections(c characters.Character) [4]bool {
	var b [4]bool
	b[0] = true
	b[1] = true
	b[2] = true
	b[3] = true
	pb := c.GetPosBox()
	ll := pb.Min.Sub(arena1.GetLowerLeft())
	ur := pb.Max.Sub(arena1.GetLowerLeft())
	if lev1.IsTile(int((ll.X-1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ll.X-1)/TileSize), int(ur.Y/TileSize)) || ll.X-1 < 0 {
		b[2] = false
	}
	if int((ur.X+1)/TileSize) > arena1.GetWidth()-1 {
		b[3] = false
	} else if lev1.IsTile(int((ur.X+1)/TileSize), int(ll.Y/TileSize)) || lev1.IsTile(int((ur.X+1)/TileSize), int(ur.Y/TileSize)) {
		b[3] = false
	}
	if int((ur.Y+1)/TileSize) > arena1.GetHeight()-1 {
		b[0] = false
	} else if lev1.IsTile(int(ll.X/TileSize), int((ur.Y+1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ur.Y+1)/TileSize)) {
		b[0] = false
	}
	if lev1.IsTile(int(ll.X/TileSize), int((ll.Y-1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize), int((ll.Y-1)/TileSize)) || ll.Y-1 < 0 {
		b[1] = false
	}
	return b
}

func moveCharacter(c characters.Character, dt float64, dir uint8) (moved bool) {

	// Ist der Character Unsichtbar? Dann ist nichts zu bewegen.
	if !c.Ani().IsVisible() || c.Ani().GetView() == Dead || c.Ani().GetView() == Intro {
		return false
	}

	dist := c.GetSpeed() * dt
	if dist >= TileSize {
		dist = TileSize - 0.1
	}
	pb := c.GetPosBox()
	ll := pb.Min.Sub(arena1.GetLowerLeft())
	ur := pb.Max.Sub(arena1.GetLowerLeft())

	switch dir {
	case Left:
		dist = -dist
		bl, xl, _ := lev1.GetPosOfNextTile(int(ll.X/TileSize), int(ll.Y/TileSize), pixel.V(-TileSize, 0))
		bu, xu, _ := lev1.GetPosOfNextTile(int(ll.X/TileSize), int(ur.Y/TileSize), pixel.V(-TileSize, 0))
		if bl || bu {
			if bl && (xl >= xu || xu == -1) {
				if ll.X+dist < float64((xl+1)*TileSize) {
					dist = float64((xl+1)*TileSize) - ll.X + 0.1
				}
			} else if bu && (xu >= xl || xl == -1) {
				if ll.X+dist < float64((xu+1)*TileSize) {
					dist = float64((xu+1)*TileSize) - ll.X + 0.1
				}
			}
		}
		c.Move(pixel.V(dist, 0))
	case Right:
		bl, xl, _ := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int(ll.Y/TileSize), pixel.V(TileSize, 0))
		bu, xu, _ := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int(ur.Y/TileSize), pixel.V(TileSize, 0))
		if bl || bu {
			if bl && (xl <= xu || xu == -1) {
				if ur.X+dist > float64((xl)*TileSize) {
					dist = float64((xl)*TileSize) - ur.X - 0.1
				}
			} else if bu && (xu <= xl || xl == -1) {
				if ur.X+dist > float64((xu)*TileSize) {
					dist = float64((xu)*TileSize) - ur.X - 0.1
				}
			}
		}
		c.Move(pixel.V(dist, 0))
	case Up:
		bl, _, yl := lev1.GetPosOfNextTile(int((ll.X)/TileSize), int((ur.Y)/TileSize), pixel.V(0, TileSize))
		br, _, yr := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int((ur.Y)/TileSize), pixel.V(0, TileSize))
		if bl || br {
			if bl && (yl <= yr || yr == -1) {
				if ur.Y+dist > float64((yl)*TileSize) {
					dist = float64((yl)*TileSize) - ur.Y - 0.1
				}
			} else if br && (yr <= yl || yl == -1) {
				if ur.Y+dist > float64((yr)*TileSize) {
					dist = float64((yr)*TileSize) - ur.Y - 0.1
				}
			}
		}
		c.Move(pixel.V(0, dist))
	case Down:
		dist = -dist
		bl, _, yl := lev1.GetPosOfNextTile(int((ll.X)/TileSize), int((ll.Y)/TileSize), pixel.V(0, -TileSize))
		br, _, yr := lev1.GetPosOfNextTile(int((ur.X)/TileSize), int((ll.Y)/TileSize), pixel.V(0, -TileSize))
		if bl || br {
			//fmt.Println(br, xr, yr)
			if bl && (yl >= yr || yr == -1) {
				if ll.Y+dist < float64((yl+1)*TileSize) {
					dist = float64((yl+1)*TileSize) - ll.Y + 0.1
				}
			} else if br && (yr >= yl || yl == -1) {
				if ll.Y+dist < float64((yr+1)*TileSize) {
					dist = float64((yr+1)*TileSize) - ll.Y + 0.1
				}
			}
			//fmt.Println(dist, ll.Y, float64((yr)*TileSize))
		}
		c.Move(pixel.V(0, dist))
	}
	c.Ani().SetView(dir)
	return
}

func sun() {
	const zoomFactor = 3
	const typ = 2
	const pitchWidth = 13
	const pitchHeight = 11
	var winSizeX float64 = zoomFactor * ((3 + pitchWidth) * TileSize) // TileSize = 16
	var winSizeY float64 = zoomFactor * ((1+pitchHeight)*TileSize + 32)

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}

	s1 := sounds.NewSound(ThroughSpace)
	go s1.PlaySound()

	arena1 = arena.NewArena(typ, pitchWidth, pitchHeight)

	lev1 = level1.NewBlankLevel(arena1, 1)
	lev1.SetRandomTilesAndItems(40, 10)

	whiteBomberman = characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()
	whiteBomberman.IncPower()
	whiteBomberman.IncPower()

	tb = titlebar.New((3 + pitchWidth) * TileSize)
	tb.SetMatrix(pixel.IM.Moved(pixel.V((3+pitchWidth)*TileSize/2, (1+pitchHeight)*TileSize+16)))
	tb.SetLifePointers(whiteBomberman.GetLifePointer())
	tb.SetPointsPointer(whiteBomberman.GetPointsPointer())
	tb.SetPlayers(1)
	go tb.Manager()
	tb.SetSeconds(5 * 60)
	tb.StartCountdown()
	tb.Update()

	// 2 Enemys
	monster = append(monster, characters.NewEnemy(YellowPopEye))
	//monster = append(monster, characters.NewEnemy(Drop))

	rand.Seed(time.Now().UnixNano())

	// Bomberman is in lowleft Corner
	whiteBomberman.MoveTo(arena1.GetLowerLeft())

	///////////////////////// ToDo Enyemys should be a Part of Level //////////////////////////////////////////////
	xx, yy := arena1.GetFieldCoord(whiteBomberman.GetPos())

	for _, m := range monster {
		for {
			i := rand.Intn(arena1.GetWidth())
			j := rand.Intn(arena1.GetHeight())
			if !lev1.IsTile(i, j) && xx != i && yy != j {
				m.MoveTo(arena1.GetLowerLeft().Add(pixel.V(float64(i)*
					TileSize, float64(j)*TileSize)))
				m.Ani().SetVisible(true)
				break
			}
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	itemBatch := pixel.NewBatch(&pixel.TrianglesData{}, animations.ItemImage)
	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))
	win.Update()
	last := time.Now()
	dt := time.Since(last).Seconds()

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
		keypressed := false
		dt = time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyLeft) {
			moveCharacter(whiteBomberman, dt, Left)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyRight) {
			moveCharacter(whiteBomberman, dt, Right)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyUp) {
			moveCharacter(whiteBomberman, dt, Up)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyDown) {
			moveCharacter(whiteBomberman, dt, Down)
			keypressed = true
		}
		if !keypressed {
			whiteBomberman.Ani().SetView(Stay)
		}
		if win.JustPressed(pixelgl.KeyB) {
			pb := whiteBomberman.GetPosBox()
			loleft := arena1.GetLowerLeft()
			b, _ := isThereABomb(pixel.Vec{X: math.Round(pb.Center().X/TileSize) * TileSize, Y: math.Round(pb.Center().Y/TileSize) * TileSize})
			c := lev1.IsTile(int((pb.Min.X-loleft.X)/TileSize), int((pb.Min.Y-loleft.Y)/TileSize))
			if !b && !c {
				bombs = append(bombs, tiles.NewBomb(whiteBomberman))
			}
		}

		/////////////////////////////////////Moving Enemys ///////////////////////////////////////////////////////////

		for _, m := range monster {
			var d uint8 = dirChoice(m)
			m.SetDirection(d)
			pos1 := math.Round(10*(m.GetPos().X+m.GetPos().Y)) / 10 // Auf eine Nachkommastelle runden.
			moveCharacter(m, dt, m.GetDirection())
			//fmt.Println("d:",d,"m.GetDir",m.GetDirection(),"moved",b)
			pos2 := math.Round(10*(m.GetPos().X+m.GetPos().Y)) / 10
			if pos1 == pos2 { // monster konnte sich nicht bewegen --> neue Richtung probieren.
				// Dadurch zittert es in der Falle bzw. biegt in Ecken ab oder läuft zurück.
				m.SetDirection(uint8(rand.Intn(4) + 1))
			}
		}

		/*
			for _,m := range(monster) {
				if !m.IsFollowing() {
					dirEn := rand.Intn(4)
					switch dirEn {
						case 0:									// l
							if !arena1.IsTile(xx-1,yy) {
								m.Move(pixel.V(-stepSize,0))
								m.Ani().SetView(Left)
							}
						case 1:									// r
							if !arena1.IsTile(xx+1,yy) {
								m.Move(pixel.V(stepSize,0))
								m.Ani().SetView(Right)
							}
						case 2:									// up
							if !arena1.IsTile(xx,yy+1) {
								m.Move(pixel.V(0,stepSize))
								m.Ani().SetView(Up)
							}
						case 3:
							if	!arena1.IsTile(xx,yy-1) {
								m.Move(pixel.V(0,-stepSize))
								m.Ani().SetView(Down)
							}
					}
				}
			}
		*/

		/////////////////////////////////////////////////////////////////////////////////////////////////////////////7

		arena1.GetCanvas().Draw(win, *(arena1.GetMatrix()))

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

		showExplosions(win)
		tempAniSlice = clearExplosions(tempAniSlice)

		whiteBomberman.Draw(win)
		for _, m := range monster {
			m.Draw(win)
		}

		tb.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
