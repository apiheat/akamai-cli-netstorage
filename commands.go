package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func cmdUpload(c *cli.Context) error {
	return upload(c)
}

func cmdDownload(c *cli.Context) error {
	return download(c)
}

func cmdErase(c *cli.Context) error {
	return erase(c)
}

func cmdDelete(c *cli.Context) error {
	return rm(c)
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
			if stat.Files[i].Type == "file" {
				nsTargetPath := fmt.Sprintf("%s/%s", nsDestination, stat.Files[i].Name)
				fmt.Printf("\nDeleting from: %s \n", nsTargetPath)
				checkResponseCode(ns.Delete(nsTargetPath))
			}
		}
		fmt.Println()
	}
	return nil
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

func upload(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	verifyPath(c)
	nsDestination := path.Join("/", nsCpcode, nsPath)

	home, _ := homedir.Dir()
	targetDir := path.Clean(path.Join(home, nsDestination))
	if c.String("from-directory") != "" {
		targetDir = c.String("from-directory")
	}

	fi, err := os.Stat(targetDir)
	errorCheck(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(targetDir)
		errorCheck(err)

		for _, f := range files {
			localPath := path.Clean(fmt.Sprintf("%s/%s", targetDir, f.Name()))
			nsTarget := path.Clean(fmt.Sprintf("%s/%s", nsDestination, f.Name()))
			fmt.Printf("\nUploading from: %s to: %s\n", localPath, nsTarget)
			res, body, err := ns.Upload(localPath, nsTarget)
			errorCheck(err)

			if res.StatusCode == 200 {
				fmt.Printf(body)
			}
		}
	case mode.IsRegular():
		localPath := path.Clean(targetDir)
		nsTarget := path.Clean(fmt.Sprintf("%s/%s", nsDestination, path.Base(targetDir)))
		fmt.Printf("\nUploading from: %s to: %s\n", localPath, nsTarget)
		checkResponseCode(ns.Upload(localPath, nsTarget))
	}
	return nil
}
