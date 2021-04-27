package level


import(	"../arena" 
		"../tiles"
		"../characters"
		"math/rand"
		"time"
		"fmt"
		//. "../constants"
		"github.com/faiface/pixel"
		"github.com/faiface/pixel/pixelgl")


type lv struct {
	layer1 []characters.Enemy		// monster
	layer2 []characters.Player		// Bomberman and Battleman
	layer4 []tiles.Tile				// destroyable Tiles
	layer3 []tiles.Item				// Items
	loleft pixel.Vec
	width,height int
	ar arena.Arena
}

func NewBlankLevel (ar arena.Arena) *lv {
	l := new(lv)
	(*l).loleft = ar.GetLowerLeft()
	(*l).width = ar.GetWidth()
	(*l).height = ar.GetHeight()
	(*l).ar = ar
	return l
}

func (l *lv) SetMonster (m []characters.Enemy) {
	(*l).layer1 = m
}

func (l *lv) SetPLayer (p []characters.Player) {
	(*l).layer2 = p
}

func (l *lv) SetItems (it []tiles.Item) {
	(*l).layer3 = it
}

func (l *lv) SetTiles (t []tiles.Tile) {
	(*l).layer4 = t
}

func (l *lv) GetTiles () []tiles.Tile {
	return (*l).layer4
}

func (l *lv) SetRandomItems (number int) {//, ar arena.Arena) {
	ar := (*l).ar
	rand.Seed(time.Now().UnixNano())
	width := (*l).width 
	height := (*l).height
	var freeTiles [][2]int
	for x:=0; x<width; x++ {
		for y:=0; y<height; y++ {
			if ar.IsFreeTile(x,y) && x+y>1 && x+y<width+height-2  {freeTiles=append(freeTiles,[2]int{x,y})}
		}
	}
	if len(freeTiles)-10-len(l.layer4) < number {
		fmt.Println("Nicht genügend freie Plätze.")
		fmt.Println("Es werden nur ",(len(freeTiles)-10)/5," Tiele zufällig platziert.")
	}
	var index,x,y,i,t int
	var nt tiles.Item
	for i<number {
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		t = 100+rand.Intn(12)
		nt = tiles.NewItem(uint8(t),ar.GetLowerLeft(),x,y)
		/*for _,dTile := range l.layer4 {
			xx,yy := dTile.GetIndexPos()
			if xx==x && yy==y { nt.Ani().SetVisible(false) }
		}*/
		fmt.Println(x,y)
		fmt.Println(nt.Ani().IsVisible())
		(*l).layer3 = append((*l).layer3,nt)
		freeTiles = append(freeTiles[:index],freeTiles[index+1:]...)
		i++
	}
}

//////////////////////////////// ToDo Clear not visible Tiles //////////////////////////////////////77

func (l *lv) RemoveTile(x,y int) bool {
	for i := len( (*l).layer4 )-1; i>=0; i-- {
		xx,yy := (*l).layer4[i].GetIndexPos()
		if xx==x && yy==y {
			(*l).layer4[i].Ani().Die()
			//(*l).layer4 = append((*l).layer4[:i],(*l).layer4[i+1:]...)
			return true
		}
	}
	return false
} 

func (l *lv) IsTile (x,y int) bool {
	if (*l).ar.IsTile(x,y) { return true }
	for _,dTile := range (*l).layer4 {
		xx,yy := dTile.GetIndexPos()
		if xx==x && yy==y {
			return true
		}
	}
	return false
}

func (l *lv) SetRandomTiles (number int) {//, ar arena.Arena) {
	var index,x,y,i,t int
	ar := (*l).ar
	rand.Seed(time.Now().UnixNano())
	t = 120+rand.Intn(18)
	width := (*l).width 
	height := (*l).height
	var freeTiles [][2]int
	for x:=0; x<width; x++ {
		for y:=0; y<height; y++ {
			if ar.IsFreeTile(x,y) && x+y>2 && x+y<width+height-4 {freeTiles=append(freeTiles,[2]int{x,y})}
		}
	}
	if len(freeTiles)-10 < number {
		fmt.Println("Nicht genügend freie Plätze.")
		fmt.Println("Es werden nur ",len(freeTiles)/2," Tiele zufällig platziert.")
		number = len(freeTiles)/2
	}
	for i<number {
		fmt.Println(len(freeTiles))
		index = rand.Intn(len(freeTiles))
		x = freeTiles[index][0]
		y = freeTiles[index][1]
		if t>132 {t++}		// christmastree ist too big!!!
		//nt := tiles.NewTile(uint8(t),ar.GetLowerLeft().Add(pixel.V(float64(x)*TileSize+TileSize/2, float64(y)*TileSize+TileSize/2)))
		nt := tiles.NewTile(uint8(t),ar.GetLowerLeft(),x,y)
		(*l).layer4 = append((*l).layer4,nt)
		freeTiles = append(freeTiles[:index],freeTiles[index+1:]...)
		i++
	}
}

func (l *lv) DrawTiles (win *pixelgl.Window) {
	for _, dTile := range (*l).layer4 {
			dTile.Draw(win)
	}
}

func (l *lv) DrawItems (win *pixelgl.Window) {
	for _, dItem := range (*l).layer3 {
			dItem.Draw(win)
	}
}
