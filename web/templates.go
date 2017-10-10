package web

const indexPageHTML = `<!doctype html>
<html lang="en">
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
	</head>
	<body>
		<div class="login">
			<p class="registration-message">d<img height="20" src="/static/images/logo.svg">sxvpn</p>
			<form class="vertical-form sign-in">
				<p>This installer will create a fully configured VPN on a new 512MB droplet in your DigitalOcean account.</p>
				<a href="{{.InstallLink}}" class="btn-auth" type="submit"><i class="fa fa-sign-in fa-lg"></i> Authenticate</a>
				<br><br>
			</form>
			<div class="footer">
				<p>
				Don't have an account?
				<a href="https://cloud.digitalocean.com/registrations/new">Sign Up</a>
				</p>
			</div>
		</div>
	</body>
</html>`
const callbackHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
		<script src="/static/js/callback.js"></script>
	</head>
</html>
`
const dashboardPageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
		<script src="/static/js/dashboard.js"></script>
	</head>
	<body class="login">
	  <p class="registration-message">d<img height="20" src="/static/images/logo.svg">sxvpn</p>
	  <form id="regions" class="vertical-form sign-in">
		<div class="region">Select Region:
			<select id="region-dropdown" name="region-dropdown">
				{{ range $key, $value := .Regions }}
					<option value="{{ $key }}">{{ $value }}</option>
				{{ end }}
			</select>
		</div>
		<button type="submit" class="go-btn">Deploy New VPN</button>
	  </form>
	  <br>
	  {{ if .VPNList }}
	  <form id="advanced" class="vertical-form sign-in">
		<div class="container">
			<div class="row">
				<div class="col"></div>
				<div class="col-6">Manage Existing VPNs</div>
				<div class="col"></div>
			</div>
		</div>
		<div>
		{{ range $index, $value := .VPNList }}
			<button id="{{ $value }}" onClick="removeVPN({{ $value }})" class="btn-del">Remove {{ $value }}</button>
		{{ end }}
		<div class="loader">Loading...</div>
		</div>
	  </form>
	  {{end}}
	</body>
</html>
`
const deletePageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
		<script src="/static/js/delete.js"></script>
	</head>
</html>
`
const progressPageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
		<script src="/static/js/progress_bar.js"></script>
	</head>
	<body class="login">
	  <p class="registration-message">d<img height="20" src="/static/images/logo.svg">sxvpn</p>
	  <form class="vertical-form sign-in">
			<div id="progress-bar">
				<div id="current-progress"></div>
			</div>
			<p id="status-line">Initializing droplet..</p>
			<p id="initial-ip"></p>
	  </form>
	</body>
</html>
`
const completePageHTML = `
<!DOCTYPE html>
<html>
	<head>
		<title>dosxvpn</title>
		<link href="/static/images/favicon.png" rel="icon" type="image/png">
		<link rel="stylesheet" href="/static/css/style.css">
		<link rel="stylesheet" href="/static/css/bootstrap/bootstrap.css">
		<link rel="stylesheet" href="/static/css/fontawesome/css/font-awesome.min.css">
		<script src="/static/js/jquery/jquery.js"></script>
		<script src="/static/js/complete.js"></script>
	</head>
	<body class="login">
	  <p class="registration-message">d<img height="20" src="/static/images/logo.svg">sxvpn</p>
	  <form class="vertical-form sign-in">
			<div class="container">
				<div class="row">
					<div class="col"><h4 class="text-center">Install Complete</h4></div>
				</div>
			</div>
			<p id="final-ip">VPN IP: <b>{{ .FinalIP }}</b></p>
			
			<div id="showpass" style="text-align: center; font-size: 14px;">Show Certificate Password</div>
			<div style="display:none; text-align: center; font-size: 11px;" id="password">{{ .Password }}</div>
			<div class="container">
			<div class="row">
					<div class="col">
						<a id="apple" type="submit" download="dosxvpn.apple.mobileconfig" href="/download?type=apple"><i class="fa fa-apple"></i> Download</a>
					</div>
					<div class="col">
						<a id="android" type="submit" download="dosxvpn.android.sswan" href="/download?type=android"><i class="fa fa-android"></i> Download</a>			
					</div>
				</div>
			</div>

			<div class="container">
				<div class="row">
					<div class="col">
						<a href="http://pi.hole/admin/index.php?login" type="submit" target="_blank"><i class="fa fa-ban"></i> Adblock Settings</a>						
					</div>
				</div>
			</div>

			<div class="container">
				<div class="row">
					<div class="col">
						<a id="back" type="submit" style="background: #949495; font-size: 11px;" href="/dashboard"><i class="fa fa-arrow-left"></i> Back</a>
					</div>
					<div class="col"></div>
					<div class="col">
						<a id="exit" type="submit" style="background: #949495; font-size: 11px;" href="#">Exit</a>
					</div>
				</div>
			</div>
	  </form>
	</body>
</html>
`
