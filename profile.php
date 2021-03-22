<?php
session_start();
require 'account_class.php';
// If the user is not logged in redirect to the login page...
if (!($account->sessionLogin())) {
	header('Location: authenticate.php');
	exit;
}
?>

<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>AG-Manager</title>
		<link href="style.css" rel="stylesheet" type="text/css">
	</head>
	<body class="loggedin">
		<nav class="navtop">
			<div>
				<h1><a href="home.php">AG-Manager</a></h1>
				<a href="profile.php">Profil</a>
				<a href="logout.php">Logout</a>
			</div>
		</nav>
		<div class="content">
			<h2>Dein Profil:</h2>
      <form method="post" action="profile.php">
        <label for="username">Benutzername</label>
        <input type="text" id="username" name="username" placeholder="Benutzername..." required>

        <label for="password">Passwort</label>
        <input type="password" id="password" name="password" placeholder="Passwort..." required>

        <label for="password">Passwort (Wiederholung)</label>
        <input type="password" id="password2" name="password2" placeholder="Passwort..." required>

        <label for="firstname">Vorname</label>
        <input type="text" id="firstname" name="firstname" placeholder="Vorname...">

        <label for="lastname">Nachname</label>
        <input type="text" id="lastname" name="lastname" placeholder="Nachname...">

        <label for="email">E-Mail</label>
        <input type="text" id="email" name="email" placeholder="mailadresse@irgendwo.de">

        <input type="submit" value="Speichern">
      </form>
		</div>
	</body>
</html>
