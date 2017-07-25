package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/dan-v/dosxvpn"
)

var (
	flagCli    = flag.Bool("cli", false, "Deploy using CLI. Must define DIGITALOCEAN_ACCESS_TOKEN")
	flagDelete = flag.Bool("delete", false, "Delete all dosxvpn instances")
	flagRegion = flag.String("region", "sfo2", "Region to deploy VPN (e.g. ams2,ams3,nyc1,nyc2,nyc3,sfo1,sfo2)")
)

const (
	digitalOceanClientId = "e731e4858af83d074073a9bb8507c5aed08611121b90ed6d602a6d4ce43d5c8c"
	port                 = 8999
)

func main() {
	flag.Parse()

	if *flagDelete {
		deleteAllInstances()
		return
	}

	if *flagCli {
		cliDeployment(*flagRegion)
		return
	}

	host := "http://localhost:" + strconv.Itoa(port)
	exec.Command("open", host).Start()
	handler := dosxvpn.Handler(digitalOceanClientId, host)
	err := http.ListenAndServe(":"+strconv.Itoa(port), handler)
	if err != nil {
		log.Fatal(err)
	}
}

func getCliToken() string {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Must have environment variable DIGITALOCEAN_ACCESS_TOKEN set")
	}
	return token
}

func deleteAllInstances() {
	token := getCliToken()
	instances, err := dosxvpn.RemoveAllDroplets(token)
	if err != nil {
		log.Fatal("Failed to remove instances.", err)
	}
	log.Println("Removed following instances: ", instances)
}

func cliDeployment(region string) {
	token := getCliToken()

	droplet, err := dosxvpn.Deploy(token, dosxvpn.DropletName("dosxvpn-"+randomString(6)+"-"+region), dosxvpn.DropletRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created DigitalOcean droplet", droplet.DropletID)

	log.Println("Waiting for SSH to start...")
	err = dosxvpn.WaitForSSH(droplet)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting VPN details...")
	vpnDetails, err := dosxvpn.GetVPNDetails(droplet, region)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Adding VPN to OSX...")
	err = dosxvpn.SetupVPN(vpnDetails)
	if err != nil {
		log.Println(err)
	}

	log.Println("##############################")
	log.Println("VPN IP:", droplet.IPv4Address)
	log.Println("##############################")
}

func randomString(n int) string {
	return strconv.Itoa(rand.Int())[:n]
}
