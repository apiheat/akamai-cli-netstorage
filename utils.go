package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/urfave/cli"
)

func verifyCreds() {
	if nsHostname == "" && nsKeyname == "" && nsKey == "" {
		color.Set(color.FgRed)
		log.Fatal("Check you configuration file section. It might miss one of required fields")
		color.Unset()
	}
}

func verifyPath(c *cli.Context) {
	if c.NArg() > 0 {
		argPath := strings.Replace(c.Args().Get(0), nsCpcode, "", -1)
		nsPath = path.Clean(argPath)
	}

	if nsPath == "" {
		color.Set(color.FgYellow)
		log.Println("Your path is pointing to root on Netstorage with CPCode: " + nsCpcode)
		color.Unset()
	}
}

func errorCheck(e error) {
	if e != nil {
		color.Set(color.FgRed)
		log.Fatal(e)
		color.Unset()
	}
}

func checkResponseCode(response *http.Response, body string, err error) {
	errorCheck(err)
	if response.StatusCode == 200 {
		color.Set(color.FgGreen)
		fmt.Println(strings.TrimSuffix(strip.StripTags(body), "\n"))
	} else {
		color.Set(color.FgRed)
		fmt.Printf("Something went wrong...\n Response code: %v\n Message: %s\n", response.StatusCode, strings.Replace(body, "\"", "", -1))
	}
	color.Unset()
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

func printStat(obj FileNS) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("Type\tName\tMtime\tSize\tMD5"))
	date64, _ := strconv.ParseInt(obj.Mtime, 10, 64)
	size64, _ := strconv.ParseUint(obj.Size, 10, 64)
	size := humanize.Bytes(size64)
	fmt.Fprintln(w, fmt.Sprintf("File:\t%s\t%s\t%s\t%s", obj.Name, time.Unix(date64, 0), size, obj.MD5))
	w.Flush()
}
