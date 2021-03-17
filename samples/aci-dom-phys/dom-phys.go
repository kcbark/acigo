package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kcbark/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s add|del|list|get args", os.Args[0])
	}

	a, errLogin := login(debug)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	execute(a, os.Args[1], os.Args[2:])

	// display existing

	aps, errList := a.PhysicalDomainList()
	if errList != nil {
		log.Printf("could not list physical domains: %v", errList)
		return
	}

	for _, t := range aps {
		name := t["name"]
		dn := t["dn"]
		log.Printf("FOUND physical domain: name=%s dn=%s\n", name, dn)
	}
}

func execute(a *aci.Client, cmd string, args []string) {
	switch cmd {
	case "add":
		if len(args) < 3 {
			log.Fatalf("usage: %s add phys-dom vlanpool-name vlanpool-mode", os.Args[0])
		}
		name := args[0]
		vlanpoolName := args[1]
		vlanpoolMode := args[2]
		errAdd := a.PhysicalDomainAdd(name, vlanpoolName, vlanpoolMode)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s", name)
	case "del":
		if len(args) < 1 {
			log.Fatalf("usage: %s del phys-dom", os.Args[0])
		}
		name := args[0]
		errDel := a.PhysicalDomainDel(name)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", name)
	case "list":
	case "get":
		if len(args) < 1 {
			log.Fatalf("usage: %s get phys-dom", os.Args[0])
		}
		name := args[0]
		pool, errGet := a.PhysicalDomainVlanPoolGet(name)
		if errGet != nil {
			log.Printf("FAILURE: get error: %v", errGet)
			return
		}
		log.Printf("SUCCESS: get: %s: => pool: %s", name, pool)
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
