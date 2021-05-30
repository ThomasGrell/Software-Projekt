package gamePlay

import (
	. "../constants"
	"../tiles"
	"fmt"
	"github.com/faiface/pixel"
	"math/rand"
)

type ga struct {
	tileMatrix    	[][][]tiles.Tile		// slice of destroyable Tiles and Items
	freePos       	[][]uint8				// 0: Free; 1:Destroyable; 2=Undestroyable
	posPlayer 		[][2]int				// len(posPLayer)==anzPLayer im Spiel
	level 			level.Level				// Initialzustand
	bounds			[2]int					// width, height
}


func NewGame (lv level.Level) *ga {
	g := new(ga)
	
	return g
}


func (g *ga) IsTile(x, y int) bool {
	if x >= (g*).bounds[0] || x < 0 || y >= (*g).bounds[1] || y < 0 {
		return true
	}
	return (*g).freePos[y][x] == Undestroyable || (*g).freePos[y][x] == Destroyable
}

func (g *ga) IsDestroyableTile(x, y int) bool {
	return (*g).freePos[y][x] == Destroyable
}

func (g *ga) GetPosOfNextTile(x, y int, dir pixel.Vec) (b bool, xx, yy int) {
	if dir.X != 0 && dir.Y != 0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
		return false, -1, -1
	} else {
		for i := 1; i <= int(dir.Len()); i++ {
			if (*g).IsTile(x+i*int(dir.X)/int(dir.Len()), y+i*int(dir.Y)/int(dir.Len())) {
				return true, x + i*int(dir.X)/int(dir.Len()), y + i*int(dir.Y)/int(dir.Len())
			}
		}
	}
	return false, -1, -1
}

func (g *ga) CollectItem(x, y int) (typ uint8, b bool) {
	if l.freePos[y][x] != Free {
		return 0, false
	}
	if len(l.tileMatrix[y][x]) == 1 {
		typ = g.tileMatrix[y][x][0].GetType()
		b = true
		g.tileMatrix[y][x] = g.tileMatrix[y][x][:0]
	} else {
		typ = 0
		b = false
	}
	return typ, b
}

func (l *lv) RemoveTile(x, y int) {
	if len((*g).tileMatrix[y][x]) == 2 {
		if (*g).tileMatrix[y][x][1].Ani().IsVisible() {
			(*g).tileMatrix[y][x][1].Ani().Die()
			(*g).freePos[y][x] = Free
		} else {
			g.tileMatrix[y][x] = g.tileMatrix[y][x][:1]
		}
	} else if len((*l).tileMatrix[y][x]) == 1 {
		if (*g).tileMatrix[y][x][0].Ani().IsVisible() {
			(*g).tileMatrix[y][x][0].Ani().Die()
			(*g).freePos[y][x] = Free
		} else {
			g.tileMatrix[y][x] = g.tileMatrix[y][x][:0]
		}
	}
}

func (g *ga) RemoveItems(x, y int, dir pixel.Vec) {
	xx := int(dir.X)/int(dir.Len())
	yy := int(dir.Y)/int(dir.Len())
	if dir.X != 0 && dir.Y != 0 {
		fmt.Println("Kein Gültiger Vektor übergeben.")
	} else {
		for i := 1; i <= int(dir.Len()); i++ {
			if len((*g).tileMatrix[y+i*yy][x+i*xx]) == 1 {
				if (*g).tileMatrix[y+i*yy][x+i*xx][0].Ani().IsVisible() {
					(*g).tileMatrix[y+i*yy][x+i*xx][0].Ani().Die()
					(*g).tileMatrix[y+i*yy][x+i*xx][0].Ani().Update()
				} else {
					(*g).tileMatrix[y+i*yy][x+i*xx]=(*g).tileMatrix[y+i*yy][x+i*xx][:0]
				}
			}
		}
	}
}
