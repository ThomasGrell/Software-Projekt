package level

type Level interface {

	/*
	 * Vor.: -
	 * Erg.: Slice mit den Koordinaten der zerstörbaren Wände/Teile ist geliefert.
	 * 		 Dabei steht die x-Koordinate an Index 0 und die y-Koordinate an Index 1
	 * 		 im jeweiligen Feld [2]int
	 */
	GetTilePos() [][2]int

	/*
	 * Vor.: -
	 * Erg.: Slice aus Itemtypen (als int, vgl. constants), die im Level vorkommen sollen, ist geliefert.
	 * 		 Dabei kann der gleiche Itemtyp mehrmals auftreten. Jeder im Slice vorkommende Itemtyp,
	 * 		 steht für ein im Spiel vorkommendes Item.
	 */
	GetLevelItems() []int

	/*
	 * Vor.: -
	 * Erg.: Slice aus Enemytypen (als int, vgl. constants), die im Level vorkommen sollen, ist geliefert.
	 * 		 Dabei kann der gleiche Enemytyp mehrmals auftreten. Jeder im Slice vorkommende Enemytyp,
	 * 		 steht für ein im Spiel vorkommendes Enemy.
	 */
	GetLevelEnemys() []int

	/*
	 * Vor.: -
	 * Erg.: Weite/Tilesize und Höhe/Tilesize (vgl. constants) der Spielfläche ist als Paar weite,höhe geliefert.
	 */
	GetBounds() (int, int)

	/*
	 * Vor.: -
	 * Erg.: Der Arenatyp (als int, vgl. constants) ist geliefert.
	 */
	GetArenaType() int

	/*
	 * Vor.: -
	 * Erg.: Die Musik (als uint8, vgl. constants.go) ist geliefert.
	 */
	GetMusic() uint8

	/*
	 * Vor.: -
	 * Erg.: Der Typ der Zerstörbaren Teile (als int, vgl. constants) ist geliefert.
	 */
	GetTileType() int
}
