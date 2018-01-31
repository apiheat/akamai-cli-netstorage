package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	netstorage "github.com/akamai/netstoragekit-golang"
	humanize "github.com/dustin/go-humanize"
	"github.com/urfave/cli"
)

func verifyPath(c *cli.Context) {
	if c.NArg() > 0 {
		argPath := strings.Replace(c.Args().Get(0), nsCpcode, "", -1)
		nsPath = path.Clean(argPath)
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
			printBody(b)
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

func printBody(body string) {
	var statDir StatNS
	xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
	xml.Unmarshal([]byte(xmlstr), &statDir)

	fmt.Printf("\nDirectory: %s\n", statDir.Directory)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, fmt.Sprint("Type\tName\tMtime\tSize\tMD5"))
	if len(statDir.Files) == 0 {
		fmt.Fprintln(w, fmt.Sprint("..\t..\t-----\t----\t---"))
	}
	for i := range statDir.Files {
		date64, _ := strconv.ParseInt(statDir.Files[i].Mtime, 10, 64)
		size64, _ := strconv.ParseUint(statDir.Files[i].Size, 10, 64)
		size := humanize.Bytes(size64)
		if statDir.Files[i].Type == "file" {
			fmt.Fprintln(w, fmt.Sprintf("File:\t%s\t%s\t%s\t%s", statDir.Files[i].Name, time.Unix(date64, 0), size, statDir.Files[i].MD5))
		}

		if statDir.Files[i].Type == "dir" {
			fmt.Fprintln(w, fmt.Sprintf("DIR:\t%s\t%s\t%s\t%s", statDir.Files[i].Name, time.Unix(date64, 0), "", ""))
		}
	}
	w.Flush()
}

func dirAction(action string, c *cli.Context) error {

	verifyPath(c)
	executeNetstorageDirAction(nsPath, action)

	return nil
}
