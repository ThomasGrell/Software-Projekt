# Ideen und Anmerkungen

Alles aus meiner sehr subjektiven Sicht.
Ich erwarte nicht, dass ich Recht bekomme. ;o)
Der Quelltext der Arena scheint mir aber aktuell
sehr unflexibel und schwer verständlich.

## Ist es sinnvoll, die Arena der Fenstergröße anzupassen?
Die Arena sollte eine feste Größe haben. 
Diese kann von Level zu Level unterschiedlich sein. 
Das Fenster zeigt ggf. einen bestimmten Ausschnitt, 
welcher der Position des Players folgt, 
oder der Zoom wird für eine bestimmte Fenstergröße angepasst. 

## Was sollte die Methode NewArena() können?
Breite und Höhe des inneren Spielfeldbereiches sollten vorgegeben werden können.
Aber nicht in Pixeln, sondern durch Anzahl der Felder in beide Richtungen.
Außerdem sollte das gewünschte Layout übergeben werden können.
Dieses könnte man über Konstanten in constants.go regeln.

## Konstanten in Methoden vermeiden
In mehreren Methoden werden bisher die Variablen "columns" und "rows"
konstant belegt. Eine Anpassung für andere Arenen erfordert viele manuelle Änderungen.
Gleiches gilt für locations[2][35].

## Speicherschwein Kollisionskarte
Die Kollisionskarte umfasst aktuell 15x17x4x8 = 8160 Byte.
Die Wahrheitskarte umfasst lediglich 255 Byte.
Aus den Werten y und x kann man schnell die Rechteckkoordinaten ermitteln.
Das kann ggf eine simple Methode erledigen.

## Arrays als Funktionsparameter
Statt in Funktionen komplette Arrays zu kopieren,
ist es überlegenswert, mit Slices zu arbeiten, da hier
nur Pointer übergeben werden. Die möglichen Werte des
Slices lassen sich dann mit range durchlaufen.
Das ist flexibler als starre Arrays.