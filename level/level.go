package level

import (
	. "../constants"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type startstatus struct {
	partsPos     [][2]int
	listOfItems  []int // never empty, because exit ist everytime in this list
	listOfEnemys []int
	w, h         int
	walltype     int
	arenaType    int
}

func NewLevel(path string) *startstatus {
	sts := new(startstatus)
	sts.listOfItems = append(sts.listOfItems, Exit)
	sts.statusFromFile(path)
	return sts
}

func (sts *startstatus) GetTilePos() [][2]int {
	return sts.partsPos
}

func (sts *startstatus) GetLevelItems() []int {
	return sts.listOfItems
}

func (sts *startstatus) GetLevelEnemys() []int {
	return sts.listOfEnemys
}

func (sts *startstatus) GetBounds() (int, int) {
	return sts.w, sts.h
}

func (sts *startstatus) GetArenaType() int {
	return sts.arenaType
}

func (sts *startstatus) GetTileType() int {
	return sts.walltype
}

func (sts *startstatus) statusFromFile(path string) {
	f, fileerr := os.Open(path)
	if fileerr != nil {
		fmt.Println("Leveldatei mit pfad: ", path, " konnte nicht erzeugt werden. Fehlermeldung: ", fileerr)
		return
	}
	var b []byte = make([]byte, 1)
	var status [7]string = [7]string{"w", "h", "arena", "pos", "item", "enemy", "walltype"}
	var n, i int
	var save []byte
	var err error
	for err != io.EOF {
		// liest eine Zeile aus
		n, err = f.Read(b)
		for b[0] != byte('\n') && err != io.EOF {
			save = append(save, b[0])
			n, err = f.Read(b)
		}
		// eine Zeile ist ausgelesen
		//fmt.Println("Status: ",status[i],"; ausgelesen: ",string(save))
		if string(save) == "" {
			return
		}
		if save[0] == byte('-') {
			save = make([]byte, 0)
			i++
		} else {
			switch status[i] {
			case "w":
				sts.w, _ = strconv.Atoi(string(save))
				save = make([]byte, 0)
			case "h":
				sts.h, _ = strconv.Atoi(string(save))
				save = make([]byte, 0)
			case "arena":
				sts.arenaType, _ = strconv.Atoi(string(save))
				save = make([]byte, 0)
			case "pos":
				sSlice := strings.Split(string(save), ",")
				var pos [2]int
				pos[0], _ = strconv.Atoi(sSlice[0])
				pos[1], _ = strconv.Atoi(sSlice[1])
				sts.partsPos = append(sts.partsPos, pos)
				save = make([]byte, 0)
			case "item":
				sSlice := strings.Split(string(save), ":")
				itemType := parseConstants(sSlice[0])
				anzItems, _ := strconv.Atoi(sSlice[1])
				for i := 0; i < anzItems; i++ {
					sts.listOfItems = append(sts.listOfItems, itemType)
				}
				save = make([]byte, 0)
			case "enemy":
				sSlice := strings.Split(string(save), ":")
				enemyType := parseConstants(sSlice[0])
				anzE, _ := strconv.Atoi(sSlice[1])
				for i := 0; i < anzE; i++ {
					sts.listOfEnemys = append(sts.listOfEnemys, enemyType)
				}
				save = make([]byte, 0)
			case "walltype":
				sts.walltype = parseConstants(string(save))
				save = make([]byte, 0)
			}
		}
		if n == 0 && err == io.EOF {
			break
		}
	}

}

func parseConstants(str string) int {
	switch str {
	case "Bomb":
		return Bomb
	case "PowerItem":
		return PowerItem
	case "BombItem":
		return BombItem
	case "PunchItem":
		return PunchItem
	case "HeartItem":
		return HeartItem
	case "RollerbladeItem":
		return RollerbladeItem
	case "SkullItem":
		return SkullItem
	case "WallghostItem":
		return WallghostItem
	case "BombghostItem":
		return BombghostItem
	case "LifeItem":
		return LifeItem
	case "Exit":
		return Exit
	case "KickItem":
		return KickItem
	case "Balloon":
		return Balloon
	case "Teddy":
		return Teddy
	case "Ghost":
		return Ghost
	case "Drop":
		return Drop
	case "Pinky":
		return Pinky
	case "BluePopEye":
		return BluePopEye
	case "Jellyfish":
		return Jellyfish
	case "Snake":
		return Snake
	case "Spinner":
		return Spinner
	case "YellowPopEye":
		return YellowPopEye
	case "Snowy":
		return Snowy
	case "YellowBubble":
		return YellowBubble
	case "PinkPopEye":
		return PinkPopEye
	case "Fire":
		return Fire
	case "Crocodile":
		return Crocodile
	case "Coin":
		return Coin
	case "Puddle":
		return Puddle
	case "PinkDevil":
		return PinkDevil
	case "Penguin":
		return Penguin
	case "PinkCyclops":
		return PinkCyclops
	case "RedCyclops":
		return RedCyclops
	case "BlueRabbit":
		return BlueRabbit
	case "PinkFlower":
		return PinkFlower
	case "BlueCyclops":
		return BlueCyclops
	case "Fireball":
		return Fireball
	case "Dragon":
		return Dragon
	case "BlueDevil":
		return BlueDevil
	case "Stub":
		return Stub
	case "Brushwood":
		return Brushwood
	case "Greenwall":
		return Greenwall
	case "Greywall":
		return Greywall
	case "Brownwall":
		return Brownwall
	case "Darkbrownwall":
		return Darkbrownwall
	case "Evergreen":
		return Evergreen
	case "Tree":
		return Tree
	case "Palmtree":
		return Palmtree
	case "Perl":
		return Perl
	case "Snowrock":
		return Snowrock
	case "Greenrock":
		return Greenrock
	case "House":
		return House
	case "Christmastree":
		return Christmastree
	case "Perl1":
		return Perl1
	case "Perl2":
		return Perl2
	case "Perl3":
		return Perl3
	case "Perl4":
		return Perl4
	case "Littlesnowrock":
		return Littlesnowrock
	default:
		fmt.Println("Unbekanntes Format parseConstants hat nicht geklappt")
	}
	return -1
}
