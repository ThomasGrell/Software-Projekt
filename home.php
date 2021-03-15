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
				<h1>AG-Manager</h1>
				<a href="profile.php">Profil</a>
				<a href="logout.php">Logout</a>
			</div>
		</nav>
		<div class="content">
			<h2>Deine AGs:</h2>
			<?php 
				echo "<p>Du bist für keine AG eingetragen.</p>";
			 ?>
		</div>
	</body>
</html>
