One click personal VPN server on [DigitalOcean](https://digitalocean.com) with automated OSX VPN setup. The deployed VPN server features auto upgrading of both the OS and VPN software so you don't need to worry about managing a server.

## Features
* Personal IPSec VPN ([strongSwan](https://www.strongswan.org/)) deployed on DigitalOcean.
* Graphical deployment with automated OSX VPN setup.
* No separate software required - uses native OSX built in VPN.
* Set it and forget it. Automated OS and VPN software updates.
* Downloadable config file that can be used to setup VPN on other computers

## Installation

### Binary
The easiest way is to download a pre-built binary from the [GitHub Releases](https://github.com/dan-v/dosxvpn/releases) page. This is a packaged OSX app.

### Source
1. Fetch the project with `go get`:

  ```sh
  go get github.com/dan-v/dosxvpn
  cd $GOPATH/src/github.com/dan-v/dosxvpn
  ```

2. Install dependencies (using [Glide](https://github.com/Masterminds/glide) for dependency management)

  ```sh
  glide install
  ```
  
3. Run make to build (will need to install [platypus cli](http://www.sveinbjorn.org/platypus)). CLI and OSX app can then be found under build/osx/x86-64.

  ```sh
  make
  ```

## CLI Usage

```bash
go install github.com/dan-v/dosxvpn/cmd/...
DIGITALOCEAN_ACCESS_TOKEN=... dosxvpn -cli
```

Prints output like:
```
2017/04/05 15:58:57 Created DigitalOcean droplet 44882920
2017/04/05 15:58:57 Waiting for SSH to start...
2017/04/05 15:59:32 Getting VPN details...
2017/04/05 15:59:51 Adding VPN to OSX...
2017/04/05 15:59:55 ##############################
2017/04/05 15:59:55 VPN IP: 10.10.10.10
2017/04/05 15:59:55 ##############################
```

## How it works
A web server is started on application launch and directs you to your web browser. It uses client OAuth authentication to request access to your DigitalOcean account (this permission is revoked after deployment). Once authenticated, a 512MB droplet is deployed running CoreOS that is configured to auto update on new releases. The OS is configured to launch a container ([dosxvpn/strongswan](https://hub.docker.com/r/dosxvpn/strongswan/)) on boot running [strongSwan](https://www.strongswan.org/). 

## FAQ
1. <b>Should I use dosxvpn?</b> That's up to you. Use at your own risk.
2. <b>Are you going to support other VPS providers?</b> Possibly.
3. <b>Will this make me completely anonymous?</b> No, absolutely not. All of your traffic is going through a VPS which could be traced back to your account. You can also be tracked still with [browser fingerprinting](https://panopticlick.eff.org/), etc. Your [IP address may still leak](https://ipleak.net/) due to WebRTC, Flash, etc.
4. <b>How much does this cost?</b> This spins up a 512MB DigitalOcean droplet that costs $5 a month.
5. <b>How do I uninstall this thing?</b> Go to System Preferences->Network, click on dosxvpn-* and click the '-' button in the bottom left to delete the VPN. Don't forget to also remove the droplet that is deployed in your DigitalOcean account.

# Powered by
* [Golang](https://golang.org/)
* [jbowens/dochaincore](https://github.com/jbowens/dochaincore) - Deployment code was borrowed from this project
* [vimagick/strongswan](https://github.com/vimagick/dockerfiles/tree/master/strongswan) - Using forked version of this docker image for VPN
* [platypus](http://www.sveinbjorn.org/platypus) - Used to generate OSX app 