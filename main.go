package main

import (
	"./animations"
	"./arena"
	"./characters"
	. "./constants"
	"./tiles"
	"./sounds"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	//"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
	"./level"
	"./level1"
)

var bombs []tiles.Bombe
var turfNtreesArena arena.Arena
var lev level.Level
var lev1 level1.Level
var tempAniSlice [][]interface{} // [Animation][Matrix]
var monster []characters.Enemy
var whiteBomberman characters.Player




// Vor: ...
// Eff: Ist der Counddown der Bombe abgelaufen passiert folgendes:
//     		- Eine neue Explosionsanimation ist erstellt und an die Position der ehemaligen bombe gesetzt.
//      	- Es ertönt der Explosionssound.
//      Ist der Countdown nicht abgelaufen, passiert nichts.

func checkForExplosions() []int {

	//var bomben []tiles.Bombe
	var indexe []int

	for index, item := range bombs {
		if ((item).GetTimeStamp()).Before(time.Now()) {
			//bomben = append (bomben,item)
			indexe = append(indexe, index)

			x, y := turfNtreesArena.GetFieldCoord(item.GetPos())
			power := int(item.GetPower())
			l, r, u, d := power, power, power, power
			if l > x {
				l = x
			}
			if 12-r < x {
				r = 12 - x
			}
			if d > y {
				d = y
			}
			if 10-u < y {
				u = 10 - y
			}

			// Falls es Hindernisse gibt, die zerstörbar oder unzerstörbar sind
			
			bl,xl,yl := lev1.GetPosOfNextTile(x,y,pixel.V(float64(-l),0))
			if bl {
				fmt.Println("Links: ",xl,yl)
				if lev1.IsDestroyableTile(xl,yl) {
					l=x-xl
				} else {
					l=x-xl-1
				}
			}
			br,xr,yr := lev1.GetPosOfNextTile(x,y,pixel.V(float64(r),0))
			if br {
				fmt.Println("Rechts: ",xr,yr)
				if lev1.IsDestroyableTile(xr,yr) {
					r=xr-x
				} else {
					r=xr-x-1
				}
			}
			bd,xd,yd := lev1.GetPosOfNextTile(x,y,pixel.V(0,float64(-d)))
			if bd {
				fmt.Println("Unten: ",xd,yd)
				if lev1.IsDestroyableTile(xd,yd) {
					d=y-yd
				} else {
					d=y-yd-1
				}
			}
			bu,xu,yu := lev1.GetPosOfNextTile(x,y,pixel.V(0,float64(u)))
			if bu {
				fmt.Println("Oben: ",xu,yu)
				if lev1.IsDestroyableTile(xu,yu) {
					u=yu-y
				} else {
					u=yu-y-1
				}
			}
			

			// falls sich ein Monster oder Player im Explosionsradius befindet

			for i := 1; i <= int(l); i++ {
				b := false
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if x-i == xx && y == yy {
						l = i
						m.Ani().Die()
						b = true
						break
					}
				}
				if b {
					break
				}
			}
			for i := 1; i <= int(r); i++ {
				b := false
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if x+i == xx && y == yy {
						r = i
						m.Ani().Die()
						b = true
						break
					}
				}
				if b {
					break
				}
			}
			for i := 1; i <= int(u); i++ {
				b := false
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if y+i == yy && x == xx {
						u = i
						m.Ani().Die()
						b = true
						break
					}
				}
				if b {
					break
				}
			}
			for i := 1; i <= int(d); i++ {
				b := false
				for _, m := range monster {
					xx, yy := turfNtreesArena.GetFieldCoord(m.GetPos())
					if y-i == yy && x == xx {
						d = i
						m.Ani().Die()
						b = true
						break
					}
				}
				if b {
					break
				}
			}
			
			if xl+l==x {
				lev1.RemoveTile(xl,yl)
			}
			if xr-r==x {
				lev1.RemoveTile(xr,yr)
			}
			if yd+d==y {
				lev1.RemoveTile(xd,yd)
			}
			if yu-u==y {
				lev1.RemoveTile(xu,yu)
			}
			
			// Items, die im Expolsionsradius liegen werden zerstört, die Expolion wird aber nicht kleiner!
			
			lev1.RemoveItems(x,y,pixel.V(float64(-l),0))
			lev1.RemoveItems(x,y,pixel.V(float64(r),0))
			lev1.RemoveItems(x,y,pixel.V(0,float64(-d)))
			lev1.RemoveItems(x,y,pixel.V(0,float64(u)))
			
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

	return indexe
}



// Vor.:...
// Eff.: Explodierte Bomben sind aus dem slice bombs gelöscht

func removeExplodedBombs(indexe []int) {
	for i := len(indexe) - 1; i >= 0; i-- {
		index := indexe[i]
		if len(bombs) == 1 {
			bombs = bombs[:0]
		} else {
			bombs = append(bombs[:index], bombs[index+1:]...)
		}
	}
}

func showExpolosions(win *pixelgl.Window) []int {
	var indexe []int
	for index, a := range tempAniSlice {
		ani := (a[0]).(animations.Animation)
		ani.Update()
		mtx := (a[1]).(pixel.Matrix)
		(ani.GetSprite()).Draw(win, mtx)
		if !ani.IsVisible() {
			indexe = append(indexe, index)
		}
	}
	return indexe
}

func clearExplosions(indexe []int) {
	for index := len(indexe) - 1; index >= 0; index-- {
		if len(tempAniSlice) != 0 {
			if len(tempAniSlice) == 1 {
				tempAniSlice = tempAniSlice[:0]
			} else {
				tempAniSlice = append(tempAniSlice[:index], tempAniSlice[index+1:]...)
			}
		}
	}
}

func isThereABomb(v pixel.Vec) (bool, tiles.Bombe) {
	for _, item := range bombs {
		if item.GetPos() == v {
			return true, item
		}
	}
	return false, nil
}


func getGrantedDirections (c characters.Character) [4]bool {
	var b [4]bool
	b[0]=true
	b[1]=true
	b[2]=true
	b[3]=true
	pb := c.GetPosBox()
	ll := pb.Min.Sub(turfNtreesArena.GetLowerLeft())
	ur := pb.Max.Sub(turfNtreesArena.GetLowerLeft())
	if lev1.IsTile(int((ll.X-1)/TileSize),int(ll.Y/TileSize))||lev1.IsTile(int((ll.X-1)/TileSize),int(ur.Y/TileSize))|| ll.X-1<0 {b[0]=false}
	if int((ur.X+1)/TileSize)>turfNtreesArena.GetWidth()-1 {
		b[1]=false
	} else if lev1.IsTile(int((ur.X+1)/TileSize),int(ll.Y/TileSize)) || lev1.IsTile(int((ur.X+1)/TileSize),int(ur.Y/TileSize)){
		b[1]=false
	}
	if int((ur.Y+1)/TileSize)>turfNtreesArena.GetHeight()-1 {
		b[2]=false
	} else if lev1.IsTile(int(ll.X/TileSize),int((ur.Y+1)/TileSize)) || lev1.IsTile(int(ur.X/TileSize),int((ur.Y+1)/TileSize)) {
		b[2]=false
	}
	if lev1.IsTile(int(ll.X/TileSize),int((ll.Y-1)/TileSize))|| lev1.IsTile(int(ur.X/TileSize),int((ll.Y-1)/TileSize))|| ll.Y-1<0 {b[3]=false}

	return b
}

func sun() {
	const zoomFactor = 3
	const typ = 2
	const pitchWidth = 13
	const pitchHeight = 11
	var winSizeX float64 = zoomFactor * ((3 + pitchWidth) * TileSize) // TileSize = 16
	var winSizeY float64 = zoomFactor * ((1 + pitchHeight) * TileSize)
	var stepSize float64 = 1
	

	wincfg := pixelgl.WindowConfig{
		Title:  "Bomberman 2021",
		Bounds: pixel.R(0, 0, winSizeX, winSizeY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	
	
	turfNtreesArena = arena.NewArena(typ, pitchWidth, pitchHeight)
	
	lev = level.NewBlankLevel(turfNtreesArena)
	lev.SetRandomTiles(10)
	lev.SetRandomItems(10)
	lev1 = level1.NewBlankLevel(turfNtreesArena,1)
	lev1.SetRandomTilesAndItems (20,10)
	
	whiteBomberman = characters.NewPlayer(WhiteBomberman)
	whiteBomberman.Ani().Show()
	whiteBomberman.IncPower()
	whiteBomberman.IncPower()

	// 2 Enemys

	monster = append(monster, characters.NewEnemy(YellowPopEye))
	monster = append(monster, characters.NewEnemy(Drop))

	// not a real random number

	rand.Seed(42)

	// Put character at free space with at least two free neighbours in a row
	/*
A:
	for i := 2 * turfNtreesArena.GetWidth(); i < len(turfNtreesArena.GetPassability())-2*turfNtreesArena.GetWidth(); i++ { // Einschränkung des Wertebereichs von i um index out of range Probleme zu vermeiden
		if turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-1] && turfNtreesArena.GetPassability()[i-2] || // checke links, rechts, oben, unten
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+1] && turfNtreesArena.GetPassability()[i+2] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i+turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i+2*turfNtreesArena.GetWidth()] ||
			turfNtreesArena.GetPassability()[i] && turfNtreesArena.GetPassability()[i-turfNtreesArena.GetWidth()] &&
				turfNtreesArena.GetPassability()[i-2*turfNtreesArena.GetWidth()] {
			whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i%turfNtreesArena.GetWidth())*
				TileSize, float64(i/turfNtreesArena.GetWidth())*TileSize)))
			break A
		}
	}
*/

	// Bomberman is in lowleft Corner
	whiteBomberman.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(0, 0)))
	
///////////////////////// ToDo Enyemys should be a Part of Level //////////////////////////////////////////////
	xx, yy := turfNtreesArena.GetFieldCoord(whiteBomberman.GetPos())

	for _, m := range monster {
		for {
			i := rand.Intn(turfNtreesArena.GetWidth())
			j := rand.Intn(turfNtreesArena.GetHeight())
			if !turfNtreesArena.IsTile(i, j) && xx != i && yy != j {
				m.MoveTo(turfNtreesArena.GetLowerLeft().Add(pixel.V(float64(i)*
					TileSize, float64(j)*TileSize)))
				m.Ani().SetVisible(true)
				break
			}
		}
	}
	
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	win.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), zoomFactor))
	win.Update()

	for !win.Closed() && !win.Pressed(pixelgl.KeyEscape) {
//////////////////////////// ToDo Implement a similar Function inside Main /////////////////////////////////////
		//grDir := turfNtreesArena.GrantedDirections(whiteBomberman.GetPosBox()) // [4]bool left-right-up-down granted?
		grDir := getGrantedDirections(whiteBomberman)
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		keypressed := false

		if win.Pressed(pixelgl.KeyLeft) && grDir[0] {
			whiteBomberman.Move(pixel.V(-stepSize, 0))
			whiteBomberman.Ani().SetView(Left)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyRight) && grDir[1] {
			whiteBomberman.Move(pixel.V(stepSize, 0))
			whiteBomberman.Ani().SetView(Right)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyUp) && grDir[2] {
			whiteBomberman.Move(pixel.V(0, stepSize))
			whiteBomberman.Ani().SetView(Up)
			keypressed = true
		}
		if win.Pressed(pixelgl.KeyDown) && grDir[3] {
			whiteBomberman.Move(pixel.V(0, -stepSize))
			whiteBomberman.Ani().SetView(Down)
			keypressed = true
		}
		if !keypressed {
			whiteBomberman.Ani().SetView(Stay)
		}
		if win.JustPressed(pixelgl.KeyB) {
			pb := characters.Player(whiteBomberman).GetPosBox()
			loleft := turfNtreesArena.GetLowerLeft()
			b, _ := isThereABomb(pixel.Vec{math.Round(pb.Center().X/TileSize) * TileSize, math.Round(pb.Center().Y/TileSize) * TileSize})
			c := lev.IsTile(int((pb.Min.X-loleft.X)/TileSize),int((pb.Min.Y-loleft.Y)/TileSize))
			if !b && !c {
				var item tiles.Bombe
				item = tiles.NewBomb(characters.Player(whiteBomberman))
				bombs = append(bombs, item)
			}
		}

/////////////////////////////////////ToDO Moving Enemys ///////////////////////////////////////////////////////////

		/*
			for _,m := range(monster) {
				xx,yy := turfNtreesArena.GetFieldCoord(m.GetPos())
				x,y := turfNtreesArena.GetFieldCoord(whiteBomberman.GetPos())
				if x == xx && y == yy {
					whiteBomberman.DecLife()
				}
				if !m.IsFollowing() {
					dir := rand.Intn(4)
					switch dir {
						case 0:									// l
							if !turfNtreesArena.IsTile(xx-1,yy) {
								m.Move(pixel.V(-stepSize,0))
								m.Ani().SetView(Left)
							}
						case 1:									// r
							if !turfNtreesArena.IsTile(xx+1,yy) {
								m.Move(pixel.V(stepSize,0))
								m.Ani().SetView(Right)
							}
						case 2:									// up
							if !turfNtreesArena.IsTile(xx,yy+1) {
								m.Move(pixel.V(0,stepSize))
								m.Ani().SetView(Up)
							}
						case 3:
							if	!turfNtreesArena.IsTile(xx,yy-1) {
								m.Move(pixel.V(0,-stepSize))
								m.Ani().SetView(Down)
							}
					}
				}
			}
		*/
		
/////////////////////////////////////////////////////////////////////////////////////////////////////////////7

		turfNtreesArena.GetCanvas().Draw(win, *(turfNtreesArena.GetMatrix()))
		
		
		removeExplodedBombs(checkForExplosions())

		
		//lev.DrawItems(win)
		//lev.DrawTiles(win)
		for i:=0; i<pitchHeight; i++ {
			lev1.DrawColumn(i,win)
		}
		for _, item := range bombs {
			item.Draw(win)
		}
		
		
		

		clearExplosions(showExpolosions(win))
		
		

		whiteBomberman.Draw(win)

		for _, m := range monster {
			m.Draw(win)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(sun)
}
