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

Rollen? Admin, Leiter, SoS -> Views?

Wie wird ID festgelegt? Wie hängen ID und Kürzel bzw. Schülernummer zusammen? Benötigt man beides? Und noch Benutzername (können Benutzer ohne Benutzername existieren?)?


Situation:
Leiter bieten AGs an, SuS kommen am Schuljahresbeginn zum Vortreffen (keine Zweit- und Drittwünsche?), Leiter wählt SuS aus (was ist bei zu vielen SuS?).
Leiter & SuS sind an den Terminen anwesend, an denen die AGs stattfinden (was ist bei Krankheit? Des Leiters?). Sie können sich auf dem AG-Portal mit einer Session anmelden.
Leiter können dort AGs erstellen. Sie tragen die Anwesenheit dort ein.
SuS können AGs sehen. (Und WunschAG angeben?) Sie können auch die Eintragung ihrer Anwesenheit prüfen (ihre Eltern auch, wenn sie Benutzername und Passwort kennen)
