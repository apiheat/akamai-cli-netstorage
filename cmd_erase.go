package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/fatih/color"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/urfave/cli"
)

func cmdErase(c *cli.Context) error {
	return erase(c)
}

func erase(c *cli.Context) error {
	verifyCreds()
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
				// Check if directory is not empty, only one level down.
				// Keeping that as safe measure, if you want to delete recursive, please enable QuickDeletion for account
				// You need to raise support ticket
				r, b, e := ns.Rmdir(nsTargetPath)
				errorCheck(e)

				switch r.StatusCode {
				case 200:
					color.Set(color.FgGreen)
					fmt.Println(strings.TrimSuffix(strip.StripTags(b), "\n"))
				case 409:
					color.Set(color.FgYellow)
					fmt.Printf("Not empty directory %s will be skipped\n", nsTargetPath)
					fmt.Println("... if you want to be able to recursively delete, then open support case for Akamai and request:")
					fmt.Println("    NetStorage QuickDelete option, than use 'akamai netstorage rmdir --recursively [PATH]'")
				default:
					color.Set(color.FgRed)
					fmt.Printf("Something went wrong...\n Response code: %v\n Message: %s\n", r.StatusCode, strings.Replace(b, "\"", "", -1))
				}
				color.Unset()
			}
		}
		fmt.Println()
	}
	return nil
}
