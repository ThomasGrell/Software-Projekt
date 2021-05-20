package sounds

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
Eff.: Ein neues Soundobjekt wird geliefert. Mit nr kann der Sound ausgewählt werden.
      Zur Verfügung stehende Konstanten sind in constants.go definiert.
NewSound(nr uint8) Sound

****************************************************************************************

*/

type Sound interface {
	// Vor.: Muss mittels des Befehls "go" nebenläufig gestartet werden.
	// Eff.: Startet das Abspielen des Sounds. Musik wird in einer Endlosschleife
	//       abgespielt. Soundeffekte werden nur einmal abgespielt.
	PlaySound()

	// Vor.: Die Soundwiedergabe wurde mit PlaySound() nebenläufig gestartet.
	// Eff.: Beendet die Wiedergabe aller Sounds.
	StopSound()
}
