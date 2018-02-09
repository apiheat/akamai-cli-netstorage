package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	netstorage "github.com/akamai/netstoragekit-golang"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func cmdUpload(c *cli.Context) error {
	return upload(c)
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
