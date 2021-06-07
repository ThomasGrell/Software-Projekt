package main


import (	"./level"
			"fmt"
		)

func main (){
	var level1 level.Level
	level1 = level.NewLevel("./level/level1.txt")
	fmt.Println(level1.GetTilePos())
	fmt.Println(len(level1.GetTilePos()))
	fmt.Println(level1.GetLevelItems())
	fmt.Println(level1.GetLevelEnemies())
	fmt.Println(level1.GetBounds())
	fmt.Println(level1.GetArenaType())
	for {
		
	}
}
