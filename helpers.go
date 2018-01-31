package main

import (
	"fmt"
	"log"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/urfave/cli"
)

func verifyPath(c *cli.Context) {
	if c.NArg() > 0 {
		argPath := strings.Replace(c.Args().Get(0), nsCpcode, "", -1)
		nsPath = path.Clean(argPath)
		log.Println(nsPath)
	}

	if nsPath == "" {
		log.Println("Your path is pointing to root on Netstorage with CPCode: " + nsCpcode)
	}
}

func errorCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func executeNetstorageDirAction(dirPath, action string) {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	if dirPath != "" {
		nsPath = dirPath
	}

	location := path.Clean(path.Join("/", nsCpcode, nsPath))

	switch action {
	case "mkdir":
		r, b, e := ns.Mkdir(location)
		errorCheck(e)

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "list":
		r, b, e := ns.Dir(location)
		errorCheck(e)

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "remove":
		r, b, e := ns.Rmdir(location)
		errorCheck(e)

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "du":
		r, b, e := ns.Du(location)
		errorCheck(e)

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	default:
		r, b, e := ns.Dir(location)
		errorCheck(e)

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	}
}

func dirAction(action string, c *cli.Context) error {

	verifyPath(c)
	executeNetstorageDirAction(nsPath, action)

	return nil
}
