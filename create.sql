/*
Damit ein user beim localhost sich mit dem Befehl
psql -d ag_manager -U ag_admin
anmelden kann, muss in
/etc/postgresql/12/main/pg_hba.conf
statt "peer" überall "md5" stehen.
Dafür sind zunächst alle "peer" auf "trust" zu setzen.
Dann muss der Server neu gestartet werden mit:
sudo /etc/init.d/postgresql restart
Nun kann man sich ohne Passwort mit
psql -d ag_manager -U ag_admin
anmelden.
Nun muss das Passwort neu gesetzt werden:
ALTER USER ag_admin with password 'kq9Ba8kf61;6]f';
Mit \q psql wieder verlassen und wieder
/etc/postgresql/12/main/pg_hba.conf
öffnen. Nun alle "trust" durch "md5" ersetzen.
Nochmals den Server neu starten.
Ab jetzt gibt es keine Fehlermeldung mehr, wenn man sich mittels
psql -d ag_manager -U ag_admin
anmelden möchte, sondern es wird das Passwort abgefragt.
*/

CREATE USER ag_admin WITH PASSWORD 'kq9Ba8kf61;6]f';
CREATE USER ag_user WITH PASSWORD 'kd83kCd[dj0i';

CREATE DATABASE ag_manager WITH OWNER ag_admin;


CREATE TABLE IF NOT EXISTS users (
id		SERIAL PRIMARY KEY,
firstname		VARCHAR(30),
lastname	VARCHAR(30),
username	VARCHAR(30) UNIQUE NOT NULL,
password	VARCHAR(255) NOT NULL,
email		VARCHAR(100),
roll		VARCHAR(10),
registrationtime TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
enabled BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
session_id  VARCHAR(255) PRIMARY KEY,
account_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
logintime  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

GRANT ALL PRIVILEGES ON DATABASE ag_manager TO ag_admin;
GRANT USAGE ON SCHEMA public TO ag_admin;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO ag_admin;

GRANT CONNECT ON DATABASE ag_manager TO ag_user;
GRANT USAGE ON SCHEMA public TO ag_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO ag_user;
