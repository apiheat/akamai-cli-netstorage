package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/urfave/cli"
)

func cmdErase(c *cli.Context) error {
	return erase(c)
}

func erase(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	verifyPath(c)
	nsDestination := path.Clean(path.Join("/", nsCpcode, nsPath))
	fmt.Printf("Going to erase content of NETSTORAGE:%s\n", nsDestination)

	res, body, err := ns.Dir(nsDestination)
	errorCheck(err)

	if res.StatusCode == 200 {
		var stat StatNS
		xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &stat)
		for i := range stat.Files {
			nsTargetPath := fmt.Sprintf("%s/%s", nsDestination, stat.Files[i].Name)
			fmt.Printf("\nDeleting from: %s\n", nsTargetPath)
			if stat.Files[i].Type == "file" {
				checkResponseCode(ns.Delete(nsTargetPath))
			} else if stat.Files[i].Type == "dir" {
				checkResponseCode(ns.Rmdir(nsTargetPath))
			}
		}
		fmt.Println()
	}
	return nil
}
