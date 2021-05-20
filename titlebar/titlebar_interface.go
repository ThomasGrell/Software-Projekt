package titlebar

import (
	"github.com/faiface/pixel"
)

/*
 _________________________________
< Implementiert von Rayk von Ende >
 ---------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

****************************************************************************************

Konstruktor:

Vor.: -
Eff.: New(width float64) liefert eine Anzeige der Breite width. Dabei wird width auf ein
      Vielfaches von 8 abgerundet.
New(width float64) Titlebar

****************************************************************************************

Beispielcode:

package main

import (
	"./titlebar"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	wincfg := pixelgl.WindowConfig{
		Title:  "Bombermen 2021",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(wincfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 3))

	bar := titlebar.New(16 * 16)
	bar.SetPlayers(4)
	var playerOneLifes uint8 = 3
	var playerTwoLifes uint8 = 1
	var playerThreeLifes uint8 = 0
	var playerFourLifes uint8 = 5
	bar.SetLifePointers(&playerOneLifes, &playerTwoLifes, &playerThreeLifes, &playerFourLifes)
	var points uint32 = 1243
	bar.SetPointsPointer(&points)
	go bar.Manager()
	bar.SetSeconds(5 * 60)
	bar.StartCountdown()
	win.Update()
	for {
		bar.Draw(win)
		win.Update()
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}
	}
}

func main() {
	pixelgl.Run(run)
}

****************************************************************************************


*/

type (
	Titlebar interface {

		// Vor.: -
		// Erg.: Zeichnet die Spielstandsanzeige auf das angegebene Target (z.B. canvas oder win)
		Draw(target pixel.Target)

		// Vor.: -
		// Erg.: Liefert die verbleibende Spielzeit in Sekunden. Da intern ein Counter alle
		//       2 Sekunden reduziert wird, ist die Restzeit ein Vielfaches von 2.
		GetSeconds() uint16

		// Vor.: -
		// Erg.: Nebenläufige Funktion, welche den Countdown steuert und die Anzeige entsprechend
		//		 aktualisiert. Muss vor Aufruf von StartCountdown oder StopCountdown mittels
		//       der Anweisung "go" einmalig im Hauptprogramm gestartet werden.
		Manager()

		// Vor.: -
		// Erg.: Die übergebenen Pointeradressen auf die uint8-Werte der Spielerleben werden
		//       im Objekt gespeichert. Die Anzeige zeigt nun die verbleibenden Leben der
		//       Spieler an.
		SetLifePointers(lifePointers ...*uint8)

		// Vor.: -
		// Erg.: Setzt die interne Matrix für das Verschieben und Skalieren des
		//       Canvas-Elements.
		SetMatrix(matrix pixel.Matrix)

		// Vor.: -
		// Erg.: Legt die Anzahl der anzuzeigenden Spieler fest. Ist der übergebene Wert
		//		 kleiner als 1 oder größer als 4, dann geschieht nichts.
		SetPlayers(uint8)

		// Vor.: -
		// Erg.: Die übergebene Pointeradresse auf den uint32-Wert der Spielerpunkte werden
		//		 im Objekt gespeichert. Die Anzeige zeigt nun die Punktzahl des Spielers
		//	     stets aktuell an.	 .
		SetPointsPointer(pointPointer *uint32)

		// Vor.: -
		// Erg.: Legt die Spieldauer in Sekunden fest. Der übergebene Wert wird auf eine
		//       gerade Zahl abgerundet, da intern in Zwei-Sekunden-Schritten gezählt wird.
		//       Die Spieldauer wird nun in Form von 2 Balken links und rechts von der Uhr
		//       angezeigt.
		SetSeconds(seconds uint16)

		// Vor.: Die Funktion Manager() wurde vorher einmalig nebenläufig gestartet.
		// Erg.: Der Countdown wird gestartet. (Spielbeginn)
		StartCountdown()

		// Vor.: Die Funktion Manager() wurde vorher einmalig nebenläufig gestartet.
		// Erg.: Die verbleibende Spielzeit ist auf null gesetzt.
		StopCountdown()

		// Vor.: -
		// Erg.: Die Spielstandsanzeige wird aktualisiert.
		Update()
	}
)
