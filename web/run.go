package web

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"time"
)

const (
	ListenPort           = "8999"
	DigitalOceanClientID = "e731e4858af83d074073a9bb8507c5aed08611121b90ed6d602a6d4ce43d5c8c"
)

func Run(cleanup bool) {
	// see https://github.com/sveinbjornt/Platypus/issues/57
	if cleanup {
		go cleanupWorkaround()
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	host := "http://localhost:" + ListenPort
	exec.Command("sudo", "-u", user.Username, "open", host).Start()
	handler := webHandler(host, DigitalOceanClientID)
	err = http.ListenAndServe("127.0.0.1:"+ListenPort, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupWorkaround() {
	for {
		_, err := exec.Command("pgrep", "-f", "-a", "Contents/MacOS/dosxvpn").Output()
		if err != nil {
			os.Exit(1)
		}
		time.Sleep(time.Second * 5)
	}
}
