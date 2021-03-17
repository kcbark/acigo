package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kcbark/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s add|del|list args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	tenant := os.Args[2]

	aps, errList := a.VrfList(tenant)
	if errList != nil {
		log.Printf("could not list VRFs: %v", errList)
		return
	}

	for _, t := range aps {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]
		log.Printf("found VRF: name=%s dn=%s descr=%s\n", name, dn, descr)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 2 {
			log.Fatalf("usage: %s add tenant vrf [descr]", os.Args[0])
		}
		tenant := args[0]
		vrf := args[1]
		var descr string
		if len(args) > 2 {
			descr = args[2]
		}
		errAdd := a.VrfAdd(tenant, vrf, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s %s", tenant, vrf)
	case "del":
		if len(args) < 2 {
			log.Fatalf("usage: %s del tenant vrf", os.Args[0])
		}
		tenant := args[0]
		vrf := args[1]
		errDel := a.VrfDel(tenant, vrf)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s %s", tenant, vrf)
	case "list":
	default:
		log.Printf("unknown command: %s", cmd)
	}
}

func login(debug bool) (*aci.Client, error) {

	a, errNew := aci.New(aci.ClientOptions{Debug: debug})
	if errNew != nil {
		return nil, fmt.Errorf("login new client error: %v", errNew)
	}

	errLogin := a.Login()
	if errLogin != nil {
		return nil, fmt.Errorf("login error: %v", errLogin)
	}

	return a, nil
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}
}
