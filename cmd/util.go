package cmd

import (
	"log"
	"os"
)

func digitalOceanTokenEnvSet() bool {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		return false
	}
	return true
}

func getCliToken() string {
	token := os.Getenv("DIGITALOCEAN_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Must have environment variable DIGITALOCEAN_ACCESS_TOKEN set")
	}
	return token
}
