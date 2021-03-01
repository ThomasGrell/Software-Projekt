<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
* {
  box-sizing:border-box;
  font-family: verdana;
  font-size: 16px;
}

h1 {
  font-size: 28px;
  font-variant: small-caps;
  background-color:#2196F3;
  border: 5px dashed blue;
  padding: 10px;
}

.left {
  background-color:#2196F3;
  padding:10px;
  float:left;
  /*width:40%; /* The width is 20%, by default */
  width: 200px;
  heigth: 50px;
}

.right {
  background-color:#4CAF50;
  padding:10px;
  float:left;
  /*width:60%; /* The width is 20%, by default */
  width: 400px;
  border-style: none;
  heigth: 50px;
  font-family: verdana;
  font-size: 16px;
}

.c1 {
  overflow: auto;
  margin-top: 10px;
  margin-bottom: 10px;
}

/* Use a media query to add a break point at 800px: */
@media screen and (max-width:800px) {
  .left, .right {
    width:100%; /* The width is 100%, when the viewport is 800px or smaller */
  }
}
</style>
</head>

<body style="  background-color: darkblue;">
<?php

  $user="";
  $pwd="";
  
  function test_input($data) {
    $data = trim($data);
    $data = stripslashes($data);
    $data = htmlspecialchars($data);
  return $data;
}

if ($_SERVER["REQUEST_METHOD"] == "POST") {
  $user  = test_input($_POST["user"]);
  $pwd = test_input($_POST["pwd"]);
}
?>

<h1>AG-Planer (Anmeldung)</h1>

<form method="post" action="<?php echo htmlspecialchars($_SERVER["PHP_SELF"]);?>" style="background-color: blue; padding: 10px; overflow: auto">
  <div class="c1">
    <div class="left">Benutzername:</div>   
    <input class = "right" type="text" name="user" value="<?php echo $user ?>"/>
  </div>

  <div class="c1">
    <div class="left">Passwort:</div>     
    <input class="right" type="password" name="pwd" value="<?php echo $pwd ?>"/>
  </div>
  
  <div class="c1">
    <input class="left" type="submit" value="anmelden">
  </div>
</form>

<p><?php if (!$user=="") {echo "Ihr Benutzername lautet: "; echo $user;}; ?></p>
<p><?php if (!$pwd=="") {echo "Ihr Passwort lautet: "; echo $pwd;}; ?></p>


</body>
</html> 