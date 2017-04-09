package dosxvpn

const indexPageHTML = `<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link rel="stylesheet" href="/static/style.css">
		<link href="/static/favicon.png" rel="icon" type="image/png">
	</head>
	<body class="login">
	  <a href="/"><img class="logo registration-message" height="35" src="/static/logo.svg" alt="Digitalocean Logo"></a>
	  <p class="registration-message">One-Click OSX VPN</p>
	  <form class="vertical-form sign-in">
		  <p>This installer will create a fully configured VPN on a new 512MB droplet in your DigitalOcean account.</p>
		  <a href="{{.InstallLink}}" type="submit">Authenticate</a>
		  <br><br>
	  </form>
	  <div class="footer">
		  <p>
		  Don't have an account?
		  <a href="https://cloud.digitalocean.com/registrations/new">Sign Up</a>
		  </p>
	  </div>
	</body>
</html>`
const callbackHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link rel="stylesheet" href="/static/style.css">
		<script src="https://code.jquery.com/jquery-3.2.0.min.js"></script>
		<script src="/static/oauth_callback.js"></script>
	</head>
</html>
`
const regionPageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>One-Click VPN DigitalOcean</title>
		<link rel="stylesheet" href="/static/style.css">
		<script src="https://code.jquery.com/jquery-3.2.0.min.js"></script>
		<script src="/static/region.js"></script>
		<link href="/static/favicon.png" rel="icon" type="image/png">
	</head>
	<body class="login">
	  <a href="/"><img class="logo registration-message" height="35" src="/static/logo.svg" alt="Digitalocean Logo"></a>
	  <p class="registration-message">One-Click OSX VPN</p>
	  <form id="regions" class="vertical-form sign-in">
		<div class="region">Select Region:
			<select id="region-dropdown" name="region-dropdown">
				{{ range $key, $value := .Regions }}
					<option value="{{ $key }}">{{ $value }}</option>
				{{ end }}
			</select>
		</div>
		<button type="submit" class="go-btn">Setup VPN</button>
		<p><a class="advanced" href="#">Advanced</a></p>
		<div id="remove">
			<button type="submit" class="rem-btn">Remove VPN Droplets</button>
		</div>
	  </form>
	</body>
</html>
`
const uninstallPageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>One-Click VPN DigitalOcean</title>
		<link rel="stylesheet" href="/static/style.css">
		<script src="https://code.jquery.com/jquery-3.2.0.min.js"></script>
		<script src="/static/region.js"></script>
		<link href="/static/favicon.png" rel="icon" type="image/png">
	</head>
	<body class="login">
	  <a href="/"><img class="logo registration-message" height="35" src="/static/logo.svg" alt="Digitalocean Logo"></a>
	  <p class="registration-message">One-Click OSX VPN</p>
	  <form id="regions" class="vertical-form sign-in">
		<div class="region">Removed the following droplets:</div>
		{{ range $value := .RemovedDroplets }}
			<p>{{$value}}</p>
		{{ end }}
		<br>
		<a id="exit" type="submit" href="#">Exit</a>
	  </form>
	</body>
</html>
`
const progressPageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>One-Click VPN DigitalOcean</title>
		<link rel="stylesheet" href="/static/style.css">
		<script src="https://code.jquery.com/jquery-3.2.0.min.js"></script>
		<script type="text/javascript">
			window.installID = "{{.InstallID}}";
		</script>
		<script src="/static/progress_bar.js"></script>
		<link href="/static/favicon.png" rel="icon" type="image/png">
	</head>
	<body class="login">
	  <a href="/"><img class="logo registration-message" height="35" src="/static/logo.svg" alt="Digitalocean Logo"></a>
	  <p class="registration-message">One-Click OSX VPN</p>
	  <form class="vertical-form sign-in">
			<div id="progress-bar">
				<div id="current-progress"></div>
			</div>
			<p id="status-line">Initializing droplet&hellip;</p>
			<p id="initial-ip"></p>
			<p id="final-ip"></p>
			<a id="mobileconfig" type="submit" download="dosxvpn.mobileconfig" href="/static/dosxvpn.mobileconfig">Download VPN Configuration</a>
			<a id="exit" type="submit" href="#">Exit</a>
	  </form>
	</body>
</html>
`