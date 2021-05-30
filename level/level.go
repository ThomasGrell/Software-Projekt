package level

import (	. "../constants"
			"os"
			"io"
			"fmt"
			"strconv"
			"strings"
		)

type startstatus struct {
	partsPos [][2]int
	listOfItems []int // never empty, because exit ist everytime in this list
	listOfEnemys []int
	w,h int
	walltype int
	arenaType int
}

func NewLevel(path string) (*startstatus) {
	sts := new(startstatus)
	sts.listOfItems = append (sts.listOfItems,Exit)
	sts.statusFromFile(path)
	return sts
}

func (sts *startstatus) GetTilePos () [][2]int {
	return sts.partsPos
}

func (sts *startstatus) GetLevelItems () []int {
	return sts.listOfItems
}

func (sts *startstatus) GetLevelEnemys () []int {
	return sts.listOfEnemys
}

func (sts *startstatus) GetBounds () (int,int) {
	return sts.w,sts.h
}

func (sts *startstatus) GetArenaType () int	{
	return sts.arenaType
}

func (sts *startstatus) GetTileType () int {
	return sts.walltype
}

func (sts *startstatus) statusFromFile (path string) {
	f,fileerr := os.Open(path)
	if fileerr!=nil {
		fmt.Println("Leveldatei mit pfad: ",path," konnte nicht erzeugt werden. Fehlermeldung: ",fileerr)
		return
	}
	var b []byte = make([]byte,1)
	var status [7]string = [7]string{"w","h","arena","pos","item","enemy","walltype"}
	var n,i int
	var save []byte
	var err error
	for err!=io.EOF {
		// liest eine Zeile aus
		n,err = f.Read(b)
		for b[0]!=byte('\n') && err!=io.EOF {
			save = append(save,b[0])
			n,err = f.Read(b)
		}
		// eine Zeile ist ausgelesen
		//fmt.Println("Status: ",status[i],"; ausgelesen: ",string(save))
		if string(save)=="" {
			return
		}
		if save[0] == byte('-') {
			save = make([]byte,0)
			i++
		} else {
			switch status[i] {
				case "w":
					sts.w,_= strconv.Atoi(string(save))
					save = make([]byte,0)
				case "h":
					sts.h,_= strconv.Atoi(string(save))
					save = make([]byte,0)
				case "arena":
					sts.arenaType,_ = strconv.Atoi(string(save))
					save = make([]byte,0)
				case "pos":
					sSlice := strings.Split(string(save),",")
					var pos [2]int
					pos[0],_=strconv.Atoi(sSlice[0])
					pos[1],_=strconv.Atoi(sSlice[1])
					sts.partsPos = append(sts.partsPos,pos)
					save = make([]byte,0)
				case "item":
					sSlice := strings.Split(string(save),":")
					itemType := parseConstants(sSlice[0])
					anzItems,_ := strconv.Atoi(sSlice[1])
					for i:=0; i<anzItems; i++ {
						sts.listOfItems = append(sts.listOfItems,itemType)
					}
					save = make([]byte,0)
				case "enemy":
					sSlice := strings.Split(string(save),":")
					enemyType := parseConstants(sSlice[0])
					anzE,_ := strconv.Atoi(sSlice[1])
					for i:=0; i<anzE; i++ {
						sts.listOfEnemys = append(sts.listOfEnemys,enemyType)
					}
					save = make([]byte,0)
				case "walltype":
					sts.walltype = parseConstants(string(save))
					save = make([]byte,0)
			}
		}
		if n==0 && err ==  io.EOF {
			break
		}
	}
	
	
}

func parseConstants(str string) int {
	switch str {
		case "Bomb":
			return 100
		case "PowerItem":
			return 101
		case "BombItem":
			return 102
		case "PunchItem":
			return 103
		case "HeartItem":
			return 104
		case "RollerbladeItem":
			return 105
		case "SkullItem":
			return 106
		case "WallghostItem":
			return 107
		case "BombghostItem":
			return 108
		case "LifeItem":
			return 109
		case "Exit":
			return 110
		case "KickItem":
			return 111
		case "Balloon":
			return 9
		case "Teddy":
			return 10
		case "Ghost":
			return 11
		case "Drop":
			return 12
		case "Pinky":
			return 13
		case "BluePopEye":
			return 14
		case "Jellyfish":
			return 15
		case "Snake":
			return 16
		case "Spinner":
			return 17
		case "YellowPopEye":
			return 18
		case "Snowy":
			return 19
		case "YellowBubble":
			return 20
		case "PinkPopEye":
			return 21
		case "Fire":
			return 22
		case "Crocodile":
			return 23
		case "Coin":
			return 24
		case "Puddle":
			return 25
		case "PinkDevil":
			return 26
		case "Penguin":
			return 27
		case "PinkCyclops":
			return 28
		case "RedCyclops":
			return 29
		case "BlueRabbit":
			return 30
		case "PinkFlower":
			return 31
		case "BlueCyclops":
			return 32
		case "Fireball":
			return 33
		case "Dragon":
			return 34
		case "BlueDevil":
			return 35
		case "Stub":
			return 120
		case "Brushwood":
			return 121
		case "Greenwall":
			return 122
		case "Greywall":
			return 123
		case "Brownwall":
			return 124
		case "Darkbrownwall":
			return 125
		case "Evergreen":
			return 126
		case "Tree":
			return 127
		case "Palmtree":
			return 128
		case "Perl":
			return 129
		case "Snowrock":
			return 130
		case "Greenrock":
			return 131
		case "House":
			return 132
		case "Christmastree":
			return 133
		case "Perl1":
			return 134
		case "Perl2":
			return 135
		case "Perl3":
			return 136
		case "Perl4":
			return 137
		case "Littlesnowrock":
			return 138
		default:
			fmt.Println("Unbekanntes Format parseConstants hat nicht geklappt")
	}
	return -1
}
