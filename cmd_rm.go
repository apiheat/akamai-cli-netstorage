package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/urfave/cli"
)

func cmdRm(c *cli.Context) error {
	return rm(c)
}

func rm(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	verifyPath(c)
	// For now disable deletion in root of CPCode
	if nsPath == "" {
		log.Fatal("Path cannot be empty")
	}
	nsDestination := path.Clean(path.Join("/", nsCpcode, nsPath))
	fmt.Printf("Going to delete object in NETSTORAGE:%s\n", nsDestination)

	res, body, err := ns.Stat(nsDestination)
	errorCheck(err)

	if res.StatusCode == 200 {
		var stat StatNS
		xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &stat)

		if stat.Files[0].Type == "dir" {
			log.Fatal("For deleting directories please use 'rmdir' command")
		}
		if stat.Files[0].Type == "file" {
			fmt.Printf("\nDeleting from: %s \n", nsDestination)
			checkResponseCode(ns.Delete(nsDestination))
		}
		fmt.Println()
	}
	return nil
}
