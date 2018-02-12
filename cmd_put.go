package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func cmdPut(c *cli.Context) error {
	return put(c)
}

func put(c *cli.Context) error {
	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	verifyPath(c)
	nsDestination := path.Join("/", nsCpcode, nsPath)

	home, _ := homedir.Dir()
	from := path.Clean(path.Join(home, nsDestination))
	if c.String("from") != "" {
		from = c.String("from")
	}

	fi, err := os.Stat(from)
	errorCheck(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		putDirFiles(ns, from, nsDestination)
	case mode.IsRegular():
		putFile(ns, from, nsDestination, "")
	}
	return nil
}

func putFile(ns *netstorage.Netstorage, fileToPut, nsLocation, prefix string) error {
	fileToPut = path.Clean(fileToPut)
	nsLocation = path.Clean(nsLocation)

	if strings.HasPrefix(path.Base(fileToPut), ".") {
		color.Set(color.FgYellow)
		fmt.Printf("\n%sNetStorage Error: dot(.) as first character is not allowed", prefix)
		fmt.Printf("\n%sSkipping file upload: %s \n", prefix, fileToPut)
		color.Unset()
	} else {
		fmt.Printf("%sUploading file from: %s to: %s\n", prefix, fileToPut, nsLocation)
		checkResponseCode(ns.Upload(fileToPut, nsLocation))
	}
	return nil
}

func putDirFiles(ns *netstorage.Netstorage, directory, nsLocation string) error {
	files, err := ioutil.ReadDir(directory)
	errorCheck(err)

	for _, f := range files {
		localFile := path.Clean(fmt.Sprintf("%s/%s", directory, f.Name()))
		nsFile := path.Clean(fmt.Sprintf("%s/%s", nsLocation, f.Name()))
		putFile(ns, localFile, nsFile, "--> ")
	}
	return nil
}
