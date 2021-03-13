<?php
// Quelle: https://alexwebdevelop.com/user-authentication/
try {
  $pdo = new PDO ('pgsql:host=localhost; dbname=ag_manager',"ag_admin",'kq9Ba8kf61;6]f');
  $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
} catch (PDOException $e) {
  echo "Fehler: Verbindung mit der Datenbank schlug fehl.";
  echo "Fehlermeldung: " . htmlspecialchars ($e->getMessage ());
}


function sanitize($data) {
  $data = trim($data);
  $data = stripslashes($data);
  $data = htmlspecialchars($data);
  return $data;
}


class Account
{
    private   $id;            // ID des eingeloggten Users
    private   $name;          // Benutzername des eingeloggten Users
    private   $authenticated; // True, wenn sich der Benutzer authentifiziert hat

    // Fügt der Datenbank einen neuen Benutzer hinzu.
    public function addAccount(string $name, string $passwd): int
    {
      global $pdo;  // Objekt für Datenbankanbindung

      $name   = sanitize($name);
	    $passwd = sanitize($passwd);

	    if (!$this->isNameValid($name))
	    {
	       throw new Exception('Ungültiger Benutzername');
	    }

	    if (!$this->isPasswdValid($passwd))
	    {
		     throw new Exception('Ungültiges Passwort');
	    }

	    if (!is_null($this->getIdFromName($name)))
	    {
		     throw new Exception('Der Benutzename ist bereits vergeben');
	    }

	    $query = 'INSERT INTO users (username, password) VALUES (:name, :passwd)';

	    $hash = password_hash($passwd, PASSWORD_DEFAULT);

      $values = array(':name' => $name, ':passwd' => $hash);

      try
      {
        $res = $pdo->prepare($query);
        $res->execute($values);
      }
      catch (PDOException $e)
      {
        throw new Exception('Abfragefehler in der Datenbank');
      }

      return $pdo->lastInsertId();
    }

    // editAccount ändert die Daten für einen bestehenden Account
    // Es würd nur auf Gültigkeit hinsichtlich der Syntax geprüft.
    public function editAccount(int $id, string $name, string $passwd, bool $enabled, string $firstname, string $lastname, string $email, string $roll)
    {
    	global $pdo;
    	$name = sanitize($name);
    	$passwd = sanitize($passwd);
      $email = sanitize($email);
      $firstname = stripslashes(htmlspecialchars($firstname));
      $lastname = stripslashes(htmlspecialchars($lastname));

    	if (!$this->isIdValid($id))
    	{
    		throw new Exception('Ungültige Benutzer-ID');
    	}

    	if (!$this->isNameValid($name))
    	{
    		throw new Exception('Ungültiger Benutzername');
    	}

    	if (!$this->isPasswdValid($passwd))
    	{
    		throw new Exception('Ungültiges Passwort');
    	}

    	$idFromName = $this->getIdFromName($name);

    	if (!is_null($idFromName) && ($idFromName != $id))
    	{
    		throw new Exception('Der Benutzername ist schon vergeben');
    	}

    	$query = 'UPDATE users SET username = :na, password = :pwd, enabled = :en, firstname = :fn, lastname = :ln, email = :email, roll = :roll WHERE id = :id';

    	$hash = password_hash($passwd, PASSWORD_DEFAULT);

    	$values = array(':na' => $name, ':pwd' => $hash, ':en' => $enabled ? 'TRUE' : 'FALSE', ':id' => $id, ':fn' => $firstname, ':ln' => $lastname, ':email' => $email, ':roll'=>$roll);

    	try
    	{
    		$res = $pdo->prepare($query);
    		$res->execute($values);
    	}
    	catch (PDOException $e)
    	{
    	   throw new Exception('Datenbankfehler beim Ändern der Benutzerdaten');
    	}
    }

    //deleteAccount löscht einen Benutzer aus der Tabelle users und aus der Tabelle sessions
    public function deleteAccount(int $id)
    {
    	global $pdo;
    	if (!$this->isIdValid($id))
    	{
    		throw new Exception('Unbekannte Benutzer-ID');
    	}
    	$query = 'DELETE FROM users WHERE id = :id';
    	$values = array(':id' => $id);
    	try
    	{
    		$res = $pdo->prepare($query);
    		$res->execute($values);
    	}
    	catch (PDOException $e)
    	{
    	   throw new Exception('Datenbankfehler beim Löschen des Benutzers');
    	}
      /* Unnötig, da Fremdschlüssel in Datenbank mittels ON DELETE CASCADE definiert wurde
    	$query = 'DELETE FROM sessions WHERE (account_id = :id)';
    	$values = array(':id' => $id);
    	try
    	{
    		$res = $pdo->prepare($query);
    		$res->execute($values);
    	}
    	catch (PDOException $e)
    	{
    	   throw new Exception('Datenbankfehler beim Löschen des gespeicherten Sessions.');
    	}
      */
    }

    public function login(string $name, string $passwd): bool
    {
    	global $pdo;
    	$name = sanitize($name);
    	$passwd = sanitize($passwd);
    	if (!$this->isNameValid($name))
    	{
    		return FALSE;
    	}
    	if (!$this->isPasswdValid($passwd))
    	{
    		return FALSE;
    	}
    	$query = 'SELECT * FROM users WHERE (username = :name) AND (enabled = TRUE)';
    	$values = array(':name' => $name);
    	try
    	{
    		$res = $pdo->prepare($query);
    		$res->execute($values);
    	}
    	catch (PDOException $e)
    	{
    	   echo "Datenbankfehler beim Login";
         exit;
      }
      $row = $res->fetch(PDO::FETCH_ASSOC);
      if (is_array($row))
      {
         if (password_verify($passwd, $row['password']))
         {
             $this->id = intval($row['id'], 10);
             $this->name = $name;
             $this->authenticated = TRUE;
             $this->registerLoginSession();
             return TRUE;
         }
      }
    	return FALSE;
    }

    // registerLoginSession speichert die SessionID und die userID mit Zeitstempel
    // in der Datenbank.
    private function registerLoginSession()
    {
    	global $pdo;

    	if (session_status() == PHP_SESSION_ACTIVE)
    	{
    		/* 	Use a REPLACE statement to:
    			- insert a new row with the session id, if it doesn't exist, or...
    			- update the row having the session id, if it does exist.
    		*/
    		$query = 'INSERT INTO sessions (session_id, account_id, logintime) VALUES (:sid, :accountId, now())
                  ON CONFLICT (session_id) DO UPDATE SET account_id = :accountId, logintime = now()';
    		$values = array(':sid' => session_id(), ':accountId' => $this->id);
    		try
    		{
    			$res = $pdo->prepare($query);
    			$res->execute($values);
    		}
    		catch (PDOException $e)
    		{
    		   echo "Datenbankfehler beim Setzen der Session während des Logins\n";
           echo $e->getMessage(), "\n";
           exit;
    		}
    	}
    }

    // sessionsLogin ermöglicht ein Login mittels einer gültigen SessionID,
    // welche nicht älter als 7 Tage ist, ohne dass Username und Passwort
    // eingegeben werden müssen.
    // Voraussetzung ist, dass serverseitig in der php.ini die Session-Cookie-Lifetime
    // hoch genug gewählt wurde. Diese muss in Sekunden angegeben werden:
    // session.cookie_lifetime = 604800
    // Ein Wert von 0 bedeutet, dass das Cookie beim Schließen des Browsers gelöscht wird.
    public function sessionLogin(): bool
    {
    	global $pdo;
    	if (session_status() == PHP_SESSION_ACTIVE)
    	{
    		/*
    			Query template to look for the current session ID on the account_sessions table.
    			The query also make sure the Session is not older than 7 days
    		*/

    		$query =
    		'SELECT * FROM sessions, users WHERE (sessions.session_id = :sid) '.
    		'AND (sessions.logintime >= (now() - INTERVAL \'7 days\')) AND (sessions.account_id = users.id) '.
    		'AND (users.enabled = TRUE)';
    		$values = array(':sid' => session_id());

        try
    		{
    			$res = $pdo->prepare($query);
    			$res->execute($values);
    		}
    		catch (PDOException $e)
    		{
    		   echo "Datenbankfehler beim Sessionlogin";
           exit;
    		}

    		$row = $res->fetch(PDO::FETCH_ASSOC);
    		if (is_array($row))
    		{
    			$this->id = intval($row['id'], 10);
    			$this->name = $row['username'];
    			$this->authenticated = TRUE;
    			return TRUE;
    		}
    	}

    	/* If we are here, the authentication failed */
    	return FALSE;
    }

    public function logout()
    {
    	global $pdo;
    	if (is_null($this->id))
    	{
    		return;
    	}
    	$this->id = NULL;
    	$this->name = NULL;
    	$this->authenticated = FALSE;
    	if (session_status() == PHP_SESSION_ACTIVE)
    	{
    		$query = 'DELETE FROM sessions WHERE (session_id = :sid)';
    		$values = array(':sid' => session_id());
    		try
    		{
    			$res = $pdo->prepare($query);
    			$res->execute($values);
    		}
    		catch (PDOException $e)
    		{
    		   throw new Exception('Denbankfehler beim Logout');
    		}
    	}
    }

    public function closeOtherSessions()
    {
    	global $pdo;
    	if (is_null($this->id))
    	{
    		return;
    	}
    	if (session_status() == PHP_SESSION_ACTIVE)
    	{
    		$query = 'DELETE FROM sessions WHERE (session_id != :sid) AND (account_id = :account_id)';
    		$values = array(':sid' => session_id(), ':account_id' => $this->id);
    		try
    		{
    			$res = $pdo->prepare($query);
    			$res->execute($values);
    		}
    		catch (PDOException $e)
    		{
    		   throw new Exception('Datenbankfehler beim Schließen aller weiteren Sessions');
    		}
    	}
    }

    // isAuthenticated gibt den Wert von authenticated des Objekts zurück
    public function isAuthenticated(): bool
    {
    	return $this->authenticated;
    }

    // isNameValid prüft, ob der Nutzername gültig ist.
    // Ggf. Sonderzeichen ausschließen?
    public function isNameValid(string $name): bool
    {
      $valid = TRUE;
      return $valid;
    }

    // isPasswdValid prüft das neue Passwort auf Gültigkeit
    public function isPasswdValid(string $passwd): bool
    {
	     $valid = TRUE;
       $len = mb_strlen($passwd);
       if ($len < 5)
       {
         $valid = FALSE;
       }
       return $valid;
    }

    // getIdFromName gibt die ID des Accounts zurück oder NULL, falls dieser nicht existiert
    public function getIdFromName(string $name): ?int
    {
       global $pdo;
       if (!$this->isNameValid($name))
       {
         throw new Exception('Ungültiger Benutzername');
       }
       $id = NULL;
	     $query = 'SELECT id FROM users WHERE (username = :name)';
	     $values = array(':name' => $name);
       try
       {
         $res = $pdo->prepare($query);
         $res->execute($values);
       }
       catch (PDOException $e)
       {
         throw new Exception('Abfragefehler der ID in der Datenbank');
       }
       $row = $res->fetch(PDO::FETCH_ASSOC);
       if (is_array($row))
       {
         $id = intval($row['id'], 10);
       }
       return $id;
    }

    // isIdValid prüft die formale Gültigkeit der ID
    public function isIdValid(int $id): bool
    {
      $valid = TRUE;
      return $valid;
    }

    // __construct ist der Konstruktor des PHP Daten Objektes
    public function __construct()
    {
      $this->id = NULL;
      $this->name = NULL;
      $this->authenticated = FALSE;
    }

    public function __destruct()
    {

    }
}

$account = new Account();

 ?>
