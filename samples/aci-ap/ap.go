package main

import (
	"log"
	"os"

	"github.com/kcbark/acigo/aci"
)

func main() {

	debug := os.Getenv("DEBUG") != ""

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s add|del|list tenant ap [description]", os.Args[0])
	}

	cmd := os.Args[1]
	tenant := os.Args[2]

	isList := cmd == "list"

	var name, descr string

	if !isList {

		if len(os.Args) < 4 {
			log.Fatalf("usage: %s add|del|list tenant ap [description]", os.Args[0])
		}

		name = os.Args[3]
		if len(os.Args) > 4 {
			descr = os.Args[4]
		}

	}

	a := login(debug)
	defer logout(a)

	// add/del ap

	execute(a, cmd, tenant, name, descr)

	// display existing

	aps, errList := a.ApplicationProfileList(tenant)
	if errList != nil {
		log.Printf("could not list application profiles: %v", errList)
		return
	}

	for _, t := range aps {
		name := t["name"]
		dn := t["dn"]
		descr := t["descr"]
		log.Printf("FOUND application profile: name=%s dn=%s descr=%s\n", name, dn, descr)
	}
}

func execute(a *aci.Client, cmd, tenant, name, descr string) {
	switch cmd {
	case "add":
		errAdd := a.ApplicationProfileAdd(tenant, name, descr)
		if errAdd != nil {
			log.Printf("FAILURE: add error: %v", errAdd)
			return
		}
		log.Printf("SUCCESS: add: %s", name)
	case "del":
		errDel := a.ApplicationProfileDel(tenant, name)
		if errDel != nil {
			log.Printf("FAILURE: del error: %v", errDel)
			return
		}
		log.Printf("SUCCESS: del: %s", name)
	case "list":
	default:
		log.Printf("unknown command: %s", cmd)
	}
}

func login(debug bool) *aci.Client {

	a, errNew := aci.New(aci.ClientOptions{Debug: debug})
	if errNew != nil {
		log.Printf("login new client error: %v", errNew)
		os.Exit(1)
	}

	errLogin := a.Login()
	if errLogin != nil {
		log.Printf("login error: %v", errLogin)
		os.Exit(1)
	}

	return a
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}

	log.Printf("logout: done")
}
