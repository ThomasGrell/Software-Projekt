# Software-Project Bombermen

## Aktuelle Verzeichnisstruktur:

- graphics: Spritesheets und andere Grafiken
- characters: Interface mit den Bomberman- und Enemyeigenschaften
- animations: Interface für grafische Darstellung und Animation von Objekten der Spielwelt
- sounds: Interface für Geräusche und Musik
- soundeffects: Soundeffekte und Musik im OGG-Format mit Samplerate von 41 kHz
- constants: Definition von allgemeinen Konstanten zur Identifizierung von Enemies, Bombermen, Items etc.
- arena: ADO arena



bomberman_rayk.go:  Testfile für Grafik, Sound etc.
main.go: Steuerung eines Bomberman mit Pfeiltasten vor dem Hintergrund von Level2

# Diskussion
### Die Feuerkraft erhöht sich ohne Einsammeln von Items?
### Bereits von Beginn an können beliebig viele Bomben gelegt werden 
### gelegte Bomben bleiben beim Zurücksetzen erhalten
### Teddy läuft in die Ecke und bewegt sich nicht mehr
### Alles zurücksetzen nach Lebensverlust stört manchmal
### Game over Anzeige passt sich nicht in Fenster ein (braucht Intro-Fenster, bekommt Spielfenster).
