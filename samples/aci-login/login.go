package main

import (
	"log"
	"os"
	"time"

	"github.com/udhos/acigo/aci"
)

func main() {
	a, errNew := aci.New(aci.ClientOptions{Debug: true})
	if errNew != nil {
		log.Printf("login new client error: %v", errNew)
		os.Exit(1)
	}

	errLogin := a.Login()
	if errLogin != nil {
		log.Printf("login error: %v", errLogin)
		os.Exit(2)
	}

	log.Printf("login ok: refresh=%v", a.RefreshTimeout())

	max := 3
	for i := 0; i < max; i++ {
		time.Sleep(5 * time.Second)
		errRefresh := a.Refresh()
		if errRefresh != nil {
			log.Printf("refresh %d/%d error: %v", i, max, errRefresh)
			os.Exit(3)
		}
		log.Printf("refresh %d/%d ok", i, max)
	}

	log.Printf("login done")
}
