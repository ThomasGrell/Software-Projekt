package level1


import(	"../arena" 
		"../tiles"
		"fmt"
		"github.com/faiface/pixel"
		"math/rand"
		"time"
		"github.com/faiface/pixel/pixelgl"
		//"math"
		/*
		"../characters"
		. "../constants"
		*/
		)


type lv struct {
	tileMatrix [][][]tiles.Tile
	freePos [][]uint8				// 0 frei, 1 undestTile, 2 destTile, 3 blocked
	anzPlayer uint8
	loleft pixel.Vec
	width,height int
	ar arena.Arena
}

func NewBlankLevel (ar arena.Arena, anzPlayer uint8) *lv {
	l := new(lv)
	(*l).loleft = ar.GetLowerLeft()
	(*l).width = ar.GetWidth()
	(*l).height = ar.GetHeight()
	(*l).ar = ar
	for layer:= 0; layer<(*l).height; layer++ {
		(*l).tileMatrix = append( (*l).tileMatrix, make([][]tiles.Tile,(*l).width) )
		(*l).freePos = append( (*l).freePos, make([]uint8,(*l).width) )
	}
	for y:=0; y<(*l).height; y++ {
		for x:=0; x<(*l).width; x++ {
			if ar.IsTile(x,y) {
				(*l).freePos[y][x]=1
			}
		}
	}
	(*l).anzPlayer = anzPlayer
	(*l).freePos[0][0] = 3
	(*l).freePos[0][1] = 3
	(*l).freePos[1][0] = 3
	switch anzPlayer {
		case 1:
			//Defaulteinstellung: Startplatz (untenlinks) für einen Player im Spiel
		case 2: 
			(*l).freePos[(*l).height-1][(*l).width-1] = 3
			(*l).freePos[(*l).height-1][(*l).width-2] = 3
			(*l).freePos[(*l).height-2][(*l).width-1] = 3
			//Startplatz (untenlinks und obenrechts) für zwei Player im Spiel
		case 3: 
			(*l).freePos[(*l).height-1][(*l).width-1] = 3
			(*l).freePos[(*l).height-1][(*l).width-2] = 3
			(*l).freePos[(*l).height-2][(*l).width-1] = 3
			(*l).freePos[(*l).height-1][0] = 3
			(*l).freePos[(*l).height-1][1] = 3
			(*l).freePos[(*l).height-2][0] = 3
			//Startplatz (untenlinks, obenrechts und obenlinks) für drei Player im Spiel
		case 4: 
			(*l).freePos[(*l).height-1][(*l).width-1] = 3
			(*l).freePos[(*l).height-1][(*l).width-1-2] = 3
			(*l).freePos[(*l).height-2][(*l).width-1] = 3
			(*l).freePos[(*l).height-1][0] = 3
			(*l).freePos[(*l).height-1][1] = 3
			(*l).freePos[(*l).height-2][0] = 3
			(*l).freePos[0][(*l).width-1] = 3
			(*l).freePos[0][(*l).width-2] = 3
			(*l).freePos[1][(*l).width-1] = 3
			//Startplatz (alle Ecken) für vier Player im Spiel
		default: 
			fmt.Println("Keine oder mehr als vier Startplätze gibt es nicht.")
	}
	return l
}

func (l *lv) SetRandomTilesAndItems (numberTiles, numberItems int) {
	rand.Seed(time.Now().UnixNano())
	width := (*l).width 
	height := (*l).height
	var freeTiles [][2]int
	for x:=0; x<width; x++ {
		for y:=0; y<height; y++ {
			if (*l).freePos[y][x]==0 {freeTiles=append(freeTiles,[2]int{x,y})}
		}
	}
	if len(freeTiles) < numberItems {
		fmt.Println("Nicht genügend freie Plätze für die übergebene Anzahl Items.")
		fmt.Println("Es werden nur ",len(freeTiles)," Items zufällig platziert.")
		numberItems = len(freeTiles)
	}
	var index,x,y,i,t int
	var puffer []int
	var b bool
	var ni tiles.Item
	for i<numberItems {
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		for _,p := range puffer {
			if p==index {
				b= true
				break	
			}
		}
		if !b {
			puffer = append(puffer,index)
			t = 100+rand.Intn(12)
			ni = tiles.NewItem(uint8(t),(*l).loleft,x,y)
			(*l).tileMatrix[y][x] = append ( (*l).tileMatrix[y][x], ni )
			i++
		}
		b=false
	}
	if len(freeTiles) < numberTiles {
		fmt.Println("Nicht genügend freie Plätze für die übergebene Anzahl Teile.")
		fmt.Println("Es werden nur ",len(freeTiles)," Tiele zufällig platziert.")
		numberTiles = len(freeTiles)
	}
	t = 120+rand.Intn(19)
	var nt tiles.Tile
	for i<numberTiles {
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		nt = tiles.NewTile(uint8(t),(*l).loleft,x,y)
		if len( (*l).tileMatrix [y][x] ) == 1 {
			((*l).tileMatrix [y][x][0]).Ani().SetVisible(false)
		}
		(*l).tileMatrix [y][x] = append ( (*l).tileMatrix[y][x], nt )
		(*l).freePos[y][x]=2
		freeTiles = append(freeTiles[:index],freeTiles[index+1:]...)
		i++
	}
}

func (l *lv) DrawColumn (y int,win *pixelgl.Window) {
	for x,rowSlice := range (*l).tileMatrix[y] {
		for i,tileORitem := range rowSlice {
			if !tileORitem.Ani().IsVisible() && (i==1 || len(rowSlice)==2) {
				(*l).tileMatrix[y][x] = append(rowSlice[:i],rowSlice[i+1:]...)
			} else {
				tileORitem.Draw(win)
			}
		}
	}
}

func (l *lv) IsTile (x,y int) bool {
	return (*l).freePos[y][x]==1 || (*l).freePos[y][x]==2
}

func (l *lv) IsDestroyableTile (x,y int) bool {
	return (*l).freePos[y][x]==2
}

func (l *lv) GetPosOfNextTile (x,y int, dir pixel.Vec) (b bool, xx,yy int) {
	if dir.X!=0 && dir.Y!=0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
		return false,-1,-1
	} else {
		for i:=1; i<=int(dir.Len()); i++ {
			if (*l).IsTile(x+i*int(dir.X)/int(dir.Len()),y+i*int(dir.Y)/int(dir.Len()))  {
				return true,x+i*int(dir.X)/int(dir.Len()),y+i*int(dir.Y)/int(dir.Len())
			}
		}
	}
	return false,-1,-1
}

func (l *lv) RemoveTile (x,y int) {
	if len((*l).tileMatrix[y][x])==2 {
		(*l).tileMatrix[y][x][1].Ani().Die()
		(*l).freePos[y][x]=0
	} else if len((*l).tileMatrix[y][x])==1 {
		(*l).tileMatrix[y][x][0].Ani().Die()
		(*l).freePos[y][x]=0
	}
}

func (l *lv) RemoveItems (x,y int, dir pixel.Vec) {
	if dir.X!=0 && dir.Y!=0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
	} else {
		for i:=1; i<=int(dir.Len()); i++ {
			if len((*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())])==1 {
				(*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())][0].Ani().Die()
				(*l).tileMatrix[y+i*int(dir.Y)/int(dir.Len())][x+i*int(dir.X)/int(dir.Len())][0].Ani().Update()
			}
		}
	}
}


