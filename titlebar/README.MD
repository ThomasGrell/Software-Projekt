# Package "titlebar"
Das Paket dient der Darstellung einer Infoleiste für das Spielgeschehen dar. 
Diese zeigt die verbleibende Zeit, den Punktestand sowie die verbleibenden Leben
von bis zu 4 Bomberman an.

Nachdem mit dem Konstruktor New(width float64) eine Infoleiste erzeugt
wurde, müssen als nächstes die Anzahl der Spieler mit SetPlayers(n uint8) festgelegt 
und die Pointer des Punktestandes SetPointsPointer(points *uint32) sowie 
der verbleibenden Leben SetLifePointers(lifePointers ...*uint8) übergeben werden.
Mit SetSeconds(s uint16) kann die angezeigte Restzeit gesetzt werden. Die Voreinstellung
ist 60s. 

Da die Darstellung der Infoleiste nebenläufig erfolgt, muss nun mittels
"go" die Methode Manager() gestartet werden. Erst danach kann der Countdown
mittels StartCoundown() gestartet werden. Dieser Aufruf bewirkt, dass
intern die Infoleiste jede Sekunde neu auf eine Leinwand (canvas)
gezeichnet wird. Mittels Draw(target pixel.Target) wird die Infoleiste
auf ein gewünschtes Ziel, z.B. dem Fenster, kopiert. Durch die
Methode SetMatrix(m pixel.Matrix) kann eine Transformationsmatrix übergeben
werden, welche sich auf Draw auswirkt, um die Infoleiste z.B. beliebig verschieben
zu können. Standardmäßig ist hier die Identitätsmatrix pixel.IM gesetzt.

Für weitere Informationen siehe "titlebar_interface.go".
