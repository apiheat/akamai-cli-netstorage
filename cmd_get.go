package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func cmdGet(c *cli.Context) error {
	return get(c)
}

func cmdDownload(c *cli.Context) error {
	return download(c)
}

func get(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)
	verifyPath(c)
	location := path.Clean(path.Join("/", nsCpcode, nsPath))

	resStat, bStat, errStat := ns.Stat(location)
	errorCheck(errStat)

	home, _ := homedir.Dir()
	pathLocal := path.Join(home, location)

	if resStat.StatusCode == 200 {
		var statObj StatNS
		xmlstr := strings.Replace(bStat, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &statObj)

		if c.String("destination") != "" {
			pathLocal = c.String("destination")
		}

		downloadDir := pathLocal
		if statObj.Files[0].Type == "file" {
			downloadDir = path.Dir(pathLocal)
		}

		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			os.MkdirAll(downloadDir, os.ModePerm)
		}

		switch statObj.Files[0].Type {
		case "dir":
			fmt.Printf("Going to download content of NETSTORAGE:%s directory to %s\n", location, pathLocal)
			// Download all files from given directory
			getDirFiles(ns, location, pathLocal)
		case "file":
			// Download given file
			getFile(ns, location, pathLocal, "")
		case "link":
			fmt.Printf("You are trying to download link NETSTORAGE:%s link\n", location)
			fmt.Printf("Please download original file or directory\n")
		}
	}
	return nil
}

func getFile(ns *netstorage.Netstorage, fileToGet, saveTo, prefix string) {
	fmt.Printf("%sDownloading NETSTORAGE:%s file to %s\n", prefix, fileToGet, saveTo)
	checkResponseCode(ns.Download(fileToGet, saveTo))
}

func getDirFiles(ns *netstorage.Netstorage, directory, saveTo string) {
	res, body, err := ns.Dir(directory)
	errorCheck(err)

	if res.StatusCode == 200 {
		var stat StatNS
		xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &stat)
		for i := range stat.Files {
			if stat.Files[i].Type == "file" {
				nsTargetPath := fmt.Sprintf("%s/%s", directory, stat.Files[i].Name)
				localDestinationPath := fmt.Sprintf("%s/%s", saveTo, stat.Files[i].Name)
				getFile(ns, nsTargetPath, localDestinationPath, "--> ")
			}
		}
		fmt.Println()
	}
}

func download(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	verifyPath(c)
	nsDestination := path.Clean(path.Join("/", nsCpcode, nsPath))
	fmt.Printf("Going to download content of NETSTORAGE:%s directory\n", nsDestination)

	res, body, err := ns.Dir(nsDestination)
	errorCheck(err)

	if res.StatusCode == 200 {
		var stat StatNS
		xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &stat)
		for i := range stat.Files {
			if stat.Files[i].Type == "file" {
				home, _ := homedir.Dir()
				pathLocal := path.Join(home, nsDestination)
				if c.String("to-directory") != "" {
					pathLocal = c.String("to-directory")
				}
				if _, err := os.Stat(pathLocal); os.IsNotExist(err) {
					os.MkdirAll(pathLocal, os.ModePerm)
				}

				nsTargetPath := fmt.Sprintf("%s/%s", nsDestination, stat.Files[i].Name)
				localDestinationPath := fmt.Sprintf("%s/%s", pathLocal, stat.Files[i].Name)
				fmt.Printf("\nDownloading from: %s to %s\n", nsTargetPath, localDestinationPath)
				checkResponseCode(ns.Download(nsTargetPath, localDestinationPath))
			}
		}
		fmt.Println()
	}
	return nil
}
