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
		// Erg.: Die Breite des Titlebar wird geändert.
		Resize(uint16)

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
		// Erg.: Der Countdown wird gestoppt.
		StopCountdown()

		// Vor.: -
		// Erg.: Die Spielstandsanzeige wird aktualisiert.
		Update()
	}
)
