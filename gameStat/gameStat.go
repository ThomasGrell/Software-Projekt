package gameStat

import (
	"../arena"
	. "../constants"
	"../tiles"
	"../level"
	"fmt"
	"github.com/faiface/pixel"
	"math/rand"
	"time"
)

type gs struct {
	tileMatrix    [][][]tiles.Tile
	freePos       [][]uint8 // 0 frei, 1 undestTile, 2 destTile, 3 blocked
	anzPlayer     uint8
	width, height int
	ar            arena.Arena
	lv 			  level.Level
}

func NewRandomGameStat (width, height int, anzPlayer uint8) *gs {
	l := newBlankLevel(rand.Intn(3), width, height, anzPlayer)
	(*l).setRandomTilesAndItems(width*height/2)
	return l
}

func NewGameStat (lv level.Level, anzPlayer uint8) *gs {
	g := new(gs)
	g.lv = lv
	var w,h int
	w,h = lv.GetBounds()
	g.width = w
	g.height = h
	ar := arena.NewArena(2,13,11)
	g.ar = ar
	fmt.Println(lv.GetArenaType(), w, h)
	for layer := 0; layer < g.height; layer++ {
		g.tileMatrix = append(g.tileMatrix, make([][]tiles.Tile, g.width))
		g.freePos = append(g.freePos, make([]uint8, g.width))
	}
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if ar.IsTile(x, y) {
				g.freePos[y][x] = Undestroyable
			}
		}
	}
	g.anzPlayer = anzPlayer
	g.setTilesAndItems (lv.GetTilePos(),lv.GetLevelItems(),lv.GetTileType ())
	return g
}

func (g *gs) GetBounds () (int,int){
	return g.width,g.height
}


func (g *gs) setTilesAndItems (partPos [][2]int,itemList []int,tile int) {
	rand.Seed(time.Now().UnixNano())
	var nt tiles.Tile
	var ni tiles.Item
	var cParts [][2]int = make ([][2]int,len(partPos))
	var cItems []int = make ([]int,len(itemList))
	var index,x,y int
	copy(cParts,partPos)
	copy(cItems,itemList)
	for len(cParts)!=0 {
		index = rand.Intn(len(cParts))
		x = cParts[index][0]
		y = cParts[index][1]
		if len(cItems)!=0 {
			ni = tiles.NewItem(uint8(cItems[0]), g.ar.CoordToVec(x, y))
			g.tileMatrix[y][x] = append(g.tileMatrix[y][x], ni)
			cItems = cItems[1:]
		}
		nt = tiles.NewTile(uint8(tile), g.ar.CoordToVec(x, y))
		g.tileMatrix[y][x] = append(g.tileMatrix[y][x], nt)
		g.freePos[y][x] = Destroyable
		cParts = append(cParts[:index], cParts[index+1:]...)
	}
}

func (l *gs) A() arena.Arena {
	return l.ar
}



func (l *gs) DrawColumn(y int,win pixel.Target) {
	for x, rowSlice := range (*l).tileMatrix[y] {
		if len(rowSlice) > 1 {
			rowSlice[1].Draw(win)
		} else if len(rowSlice) == 1 {
			rowSlice[0].Draw(win)
		}
		for i, tileORitem := range rowSlice {
			if !tileORitem.Ani().IsVisible() {
				(*l).tileMatrix[y][x] = append(rowSlice[:i], rowSlice[i+1:]...)
			}
		}
	}
}


func (l *gs) IsTile(x, y int) bool {
	if x >= l.width || x < 0 || y >= l.height || y < 0 {
		return true
	}
	return (*l).freePos[y][x] == Undestroyable || (*l).freePos[y][x] == Destroyable
}

func (l *gs) IsDestroyableTile(x, y int) bool {
	return (*l).freePos[y][x] == Destroyable
}

func (l *gs) GetPosOfNextTile(x, y int, dir pixel.Vec) (b bool, xx, yy int) {
	if dir.X != 0 && dir.Y != 0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
		return false, -1, -1
	} else {
		for i := 1; i <= int(dir.Len()); i++ {
			if (*l).IsTile(x+i*int(dir.X)/int(dir.Len()), y+i*int(dir.Y)/int(dir.Len())) {
				return true, x + i*int(dir.X)/int(dir.Len()), y + i*int(dir.Y)/int(dir.Len())
			}
		}
	}
	return false, -1, -1
}

func (l *gs) CollectItem(x, y int) (typ uint8, b bool) {
	if l.freePos[y][x] != Free {
		return 0, false
	}
	if len(l.tileMatrix[y][x]) == 1 {
		typ = l.tileMatrix[y][x][0].GetType()
		b = true
		l.tileMatrix[y][x] = l.tileMatrix[y][x][:0]
	} else {
		typ = 0
		b = false
	}
	return typ, b
}

func (l *gs) RemoveTile(x, y int) {
	if len((*l).tileMatrix[y][x]) == 2 {
		(*l).tileMatrix[y][x][1].Ani().Die()
		(*l).freePos[y][x] = Free
	} else if len((*l).tileMatrix[y][x]) == 1 {
		(*l).tileMatrix[y][x][0].Ani().Die()
		(*l).freePos[y][x] = Free
	}
}

func (l *gs) RemoveItems(x, y int, dir pixel.Vec) {
	if dir.X != 0 && dir.Y != 0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
	} else {
		for i := 1; i <= int(dir.Len()); i++ {
			if len((*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())]) == 1 {
				(*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())][0].Ani().Die()
				(*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())][0].Ani().Update()
			}
		}
	}
}


func newBlankLevel(typ, width, height int, anzPlayer uint8) *gs {
	ar := arena.NewArena(typ, width, height)
	l := new(gs)
	(*l).width = ar.GetWidth()
	(*l).height = ar.GetHeight()
	(*l).ar = ar
	for layer := 0; layer < (*l).height; layer++ {
		(*l).tileMatrix = append((*l).tileMatrix, make([][]tiles.Tile, (*l).width))
		(*l).freePos = append((*l).freePos, make([]uint8, (*l).width))
	}
	for y := 0; y < (*l).height; y++ {
		for x := 0; x < (*l).width; x++ {
			if ar.IsTile(x, y) {
				(*l).freePos[y][x] = Undestroyable
			}
		}
	}
	(*l).anzPlayer = anzPlayer
	l.freePos[0][0] = Blocked
	l.freePos[0][1] = Blocked
	l.freePos[1][0] = Blocked
	if anzPlayer > 1 {
		l.freePos[l.height-1][l.width-1] = Blocked
		l.freePos[l.height-1][l.width-2] = Blocked
		l.freePos[l.height-2][l.width-1] = Blocked
	} 
	if anzPlayer > 2 {
		l.freePos[l.height-1][0] = Blocked
		l.freePos[l.height-1][1] = Blocked
		l.freePos[l.height-2][0] = Blocked
	}
	if anzPlayer > 3 {
		l.freePos[0][l.width-1] = Blocked
		l.freePos[0][l.width-2] = Blocked
		l.freePos[1][l.width-1] = Blocked
	}
	return l
}

func (l *gs) setRandomTilesAndItems(numberTiles int) {
	rand.Seed(time.Now().UnixNano())
	width := (*l).width
	height := (*l).height
	var freeTiles [][2]int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if (*l).freePos[y][x] == Free {
				freeTiles = append(freeTiles, [2]int{x, y})
			}
		}
	}
	if len(freeTiles)<numberTiles {
		fmt.Println("Nicht genügend freie Plätze für die übergebene Anzahl Teile.")
		fmt.Println("Es werden nur ", len(freeTiles), " Tiele zufällig platziert.")
		numberTiles = len(freeTiles)
	}
	var index, x, y, i, t int
	var tile = Greenwall //120 + rand.Intn(19)
	var nt tiles.Tile
	var ni tiles.Item
	for i < int(numberTiles/2) {
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		t = 100 + rand.Intn(12)
		ni = tiles.NewItem(uint8(t), l.ar.CoordToVec(x, y))
		(*l).tileMatrix[y][x] = append((*l).tileMatrix[y][x], ni)
		nt = tiles.NewTile(uint8(tile), l.ar.CoordToVec(x, y))
		(*l).tileMatrix[y][x] = append((*l).tileMatrix[y][x], nt)
		(*l).freePos[y][x] = Destroyable
		freeTiles = append(freeTiles[:index], freeTiles[index+1:]...)
		i++
	}
	i=0
	for i < numberTiles-int(numberTiles*3/4) {
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		nt = tiles.NewTile(uint8(tile), l.ar.CoordToVec(x, y))
		(*l).tileMatrix[y][x] = append((*l).tileMatrix[y][x], nt)
		(*l).freePos[y][x] = Destroyable
		freeTiles = append(freeTiles[:index], freeTiles[index+1:]...)
		i++
	}
	
} 

