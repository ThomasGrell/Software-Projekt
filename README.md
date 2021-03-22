# AG-Planer
## To-Do:
Wenn ihr mir eure Public ssh-keys sendet, kann ich diese hinterlegen.
Habe mit php und html schon mal eine Eingabemaske für den Login erstellt.
Sieht definitiv hübscher aus, als das Paket von Herrn Schäfer für gfx und go.

## Diskussion
Benutzer - eingeschrieben - AGs:
AG-Leiter und Schueler haben verschiedene Rollen - das kann hier vmtl über die Views gesteuert werden?
Kardinalitäten sind verschieden: eine AG kann nur von 1 Leiter geleitet werden, aber viele SuS können eingeschrieben sein.

Lehrer + Externe -> eine Entität "Leiter" mit Attribut "extern"?

# Rollen?
user:
editor:
admin:

-> Views?

# Wie wird ID festgelegt?
Automatisch vom System, daher der Datentyp SERIAL. (Siehe create.sql)

# Wie hängen ID und Kürzel bzw. Schülernummer zusammen?
Gar nicht. ID ist der eindeutige Primary Key für Benutzer. So vermeiden wir, dass wir unterschiedliche Schlüsseltypen für Lehrer und Schüler haben. Kürzel und Schülernummer sind aber für Datenimport und Synchronisation von Nutzen.

# Benötigt man beides?
Ja, siehe oben.

# Können Benutzer ohne Benutzername existieren?
Nein, daher ist der Benutzername als NOT NULL gekennzeichnet. Man könnte ihn auch als Schlüssel verwenden, da er UNIQUE ist. So wäre aber eine Änderung des Benutzernamens durch den Benutzer schwierig, da man alle Einträge in der Datenbank, die diesen als Fremdschlüssel verwenden, ebenfalls ändern müsste.

Situation:
Leiter bieten AGs an, SuS kommen am Schuljahresbeginn zum Vortreffen (keine Zweit- und Drittwünsche?), Leiter wählt SuS aus (was ist bei zu vielen SuS?).
Leiter & SuS sind an den Terminen anwesend, an denen die AGs stattfinden (was ist bei Krankheit? Des Leiters?). Sie können sich auf dem AG-Portal mit einer Session anmelden.
Leiter können dort AGs erstellen. Sie tragen die Anwesenheit dort ein.
SuS können AGs sehen. (Und WunschAG angeben?) Sie können auch die Eintragung ihrer Anwesenheit prüfen (ihre Eltern auch, wenn sie Benutzername und Passwort kennen)
