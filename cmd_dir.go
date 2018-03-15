package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	netstorage "github.com/akamai/netstoragekit-golang"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	return dirAction("list", c)
}

func cmdMkdir(c *cli.Context) error {
	return dirAction("mkdir", c)
}

func cmdDu(c *cli.Context) error {
	return dirAction("du", c)
}

func cmdRmdir(c *cli.Context) error {
	return dirAction("rmdir", c)
}

func dirAction(action string, c *cli.Context) error {

	verifyPath(c)
	recursiveAction := false
	if c.Bool("recursively") {
		recursiveAction = true
	}
	executeNetstorageDirAction(nsPath, action, recursiveAction)
	return nil
}
func executeNetstorageDirAction(dirPath, action string, recursive bool) {
	verifyCreds()
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	if dirPath != "" {
		nsPath = dirPath
	}

	location := path.Clean(path.Join("/", nsCpcode, nsPath))

	switch action {
	case "mkdir":
		// Always recursive
		checkResponseCode(ns.Mkdir(location))
	case "list":
		// We need to check if given object is dir or file
		resSt, bSt, eSt := ns.Stat(location)
		errorCheck(eSt)

		if resSt.StatusCode == 200 {
			var statObj StatNS
			xmlstr := strings.Replace(bSt, "ISO-8859-1", "UTF-8", -1)
			xml.Unmarshal([]byte(xmlstr), &statObj)

			if statObj.Files[0].Type == "dir" {
				r, b, e := ns.Dir(location)
				errorCheck(e)

				if r.StatusCode == 200 {
					printBody(b)
				}
			} else {
				printStat(statObj.Files[0])
			}
		}
	case "rmdir":
		if recursive {
			checkResponseCode(ns.QuickDelete(location))
		} else {
			checkResponseCode(ns.Rmdir(location))
		}
	case "du":
		// We need to check if given object is dir or file
		resSt, bSt, eSt := ns.Stat(location)
		errorCheck(eSt)

		if resSt.StatusCode == 200 {
			var statObj StatNS
			xmlstr := strings.Replace(bSt, "ISO-8859-1", "UTF-8", -1)
			xml.Unmarshal([]byte(xmlstr), &statObj)

			if statObj.Files[0].Type == "dir" {
				duDir(location, ns)
			} else {
				printStat(statObj.Files[0])
			}
		}
	default:
		checkResponseCode(ns.Dir(location))
	}
}

func duDir(location string, ns *netstorage.Netstorage) {
	resDU, bDU, eDU := ns.Du(location)
	errorCheck(eDU)

	if resDU.StatusCode == 200 {
		var duObj duNS
		xmlstr := strings.Replace(bDU, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &duObj)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		fmt.Fprintln(w, fmt.Sprint("# Directory\tFiles\tSize"))
		size64, _ := strconv.ParseUint(duObj.Info.Bytes, 10, 64)
		size := humanize.Bytes(size64)
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", duObj.Directory, duObj.Info.Files, size))
		w.Flush()
	} else {
		color.Set(color.FgRed)
		fmt.Printf("Something went wrong...\n Response code: %v\n Message: %s\n", resDU.StatusCode, strings.Replace(bDU, "\"", "", -1))
		color.Unset()
	}
}
