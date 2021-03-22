<?php
session_start();

require 'account_class.php';

if ($account->sessionLogin()) {
  echo "home";
	header('Location: home.php');
	exit;
} else {
  echo "auth";
  header('Location: authenticate.php');
  exit;
}

 ?>
