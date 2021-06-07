# Package "characters"
Dieses Paket stellt die Interfaces "Player" und "Enemy" bereit. 
Diese umfassen sämtliche Eigenschaften der Spielcharaktere, 
z. B. die Kollisionsbox, Anzahl der Leben, Fähigkeiten etc. 
Der Zugriff auf diese Eigenschaften erfolgt ausschließlich 
über die bereitgestellten Methoden. Für die grafische Darstellung 
der Charaktere wird das Paket "animations" eingebunden. 
Beim Aufruf einer der Konstruktoren NewEnemy oder NewPlayer 
wird automatisch der Konstruktor NewAnimation ausgeführt. 
Die Methoden des Pakets "animations" sind dann über die Methode 
Ani des Charakters zugreifbar.

Dokumentation siehe characters_interface.go