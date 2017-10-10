One click personal VPN server with DNS ad blocking running on [DigitalOcean](https://digitalocean.com). The deployed VPN server includes automated updates of both the OS and software, so you don't need to worry about managing a server in the cloud.

![](/static/images/overview.gif?raw=true)

## Features
* Personal IPsec-based VPN ([strongSwan](https://strongswan.org/)).
* Ad blocking DNS setup by default ([Pi-hole](https://pi-hole.net/)).
* Download profiles for sharing VPN with iPhone and Android.
* No additional software required for OSX/iPhone - uses native VPN client.
* Simple Web or CLI installation methods.
* Automated OS and VPN software updates.

## Web Install (OSX) 
1. Download the latest pre-built app from the [GitHub Releases](https://github.com/dan-v/dosxvpn/releases) page.
2. Open the app and run through the web based installation wizard to setup the VPN.

## CLI Usage (OSX)
1. Download the latest pre-built cli from the [GitHub Releases](https://github.com/dan-v/dosxvpn/releases) page.
2. Make the binary executable: chmod +x dosxvpn
3. Create an API token (https://cloud.digitalocean.com/settings/api/tokens) and export it: export DIGITALOCEAN_ACCESS_TOKEN=efdddd442dc4b687361d801ddff999aaaf4bb17b689d59149b6bc3d5f9d0s0d0df9f9f9
4. ./dosxvpn

### CLI Examples
* Deploy a new VPN and configure for immediate use: ./dosxvpn deploy --region sfo2 --auto-configure
* List dosxvpn instances: ./dosxvpn ls
* Remove dosxvpn instance: ./dosxvpn rm --name dosxvpn-472-sfo2

## FAQ
1. <b>Should I use dosxvpn?</b> That's up to you. Use at your own risk.
2. <b>How is this different than [algo](https://github.com/trailofbits/algo)?</b> While both are IPSec VPNs, there are two primary differences. 1) Installation: dosxvpn has a simple streamlined web or CLI installation without any additional system dependencies. Algo's install process only supports CLI and has system dependencies on Python. 2) Updates: dosxvpn handles updates of the OS and VPN. This means any critical security updates or bug fixes will automatically be applied for you. Algo is a one shot deployment and there are no automatic updates. To get updates you would need to manage updates yourself or deploy a new VPN instance.
3. <b>How much does this cost?</b> This launches a 512MB DigitalOcean droplet that costs $5/month currently.
4. <b>What is the bandwidth limit?</b> The 512MB DigitalOcean droplet has a 1TB bandwidth limit. This does not appear to be strictly enforced.
5. <b>Are you going to support other VPS providers?</b> Not right now.
6. <b>Will this make me completely anonymous?</b> No, absolutely not. All of your traffic is going through a VPS which could be traced back to your account. You can also be tracked still with [browser fingerprinting](https://panopticlick.eff.org/), etc. Your [IP address may still leak](https://ipleak.net/) due to WebRTC, Flash, etc.
7. <b>How do I uninstall this thing on OSX?</b> Go to System Preferences->Network, click on dosxvpn-* and click the '-' button in the bottom left to delete the VPN. Don't forget to also remove the droplet that is deployed in your DigitalOcean account.

# Powered by
* [strongSwan](https://strongswan.org/) - IPsec-based VPN software
* [CoreOS](https://coreos.com/) - used for running containers and automatic OS updates capabilities
* [Pi-hole](https://pi-hole.net/) - used for DNS adblocking
* [Platypus](http://www.sveinbjorn.org/platypus) - used to build the native OSX app 

# Acknowledgements
* [trailofbits/algo](https://github.com/trailofbits/algo) - strongSwan configuration is borrowed from this project
* [jbowens/dochaincore](https://github.com/jbowens/dochaincore) - Deployment code is borrowed from this project
* [vimagick/strongswan](https://github.com/vimagick/dockerfiles/tree/master/strongswan) - Using a forked version of this docker image for VPN server

### Building yourself
1. Fetch the project with `go get`:

  ```sh
  go get github.com/dan-v/dosxvpn
  cd $GOPATH/src/github.com/dan-v/dosxvpn
  ```
  
2. Run make to build (will need to install [platypus cli](http://www.sveinbjorn.org/platypus)).

  ```sh
  make
  ```