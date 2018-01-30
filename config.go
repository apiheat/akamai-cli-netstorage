package main

import (
	"fmt"
	"log"
	"os/user"
	"path"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/go-ini/ini"
)

func config(configFile, configSection string) (nsHostname, nsKeyname, nsKey, nsCpcode, nsPath string) {
	if configFile == ".edgerc" {
		configFile = path.Join(userHome(), ".edgerc")
	}

	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Fatal("error:", err)
	}

	section, err := cfg.GetSection(configSection)
	if err != nil {
		log.Fatal("error:", err)
	}

	nsHostname = section.Key("hostname").String()
	nsKeyname = section.Key("keyname").String()
	nsKey = section.Key("key").String()
	nsCpcode = section.Key("cpcode").String()
	nsPath = section.Key("path").String()

	return nsHostname, nsKeyname, nsKey, nsCpcode, nsPath
}

func userHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func executeNetstorageDirAction(configFile, configSection, dirPath, action string) {
	nsHostname, nsKeyname, nsKey, nsCpcode, nsPath := config(configFile, configSection)

	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	if dirPath != "" {
		nsPath = dirPath
	}

	location := fmt.Sprintf("/%s", path.Join(nsCpcode, nsPath))

	switch action {
	case "mkdir":
		r, b, e := ns.Mkdir(location)
		if e != nil {
			log.Fatal(e)
		}

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "list":
		r, b, e := ns.Dir(location)
		if e != nil {
			log.Fatal(e)
		}

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "remove":
		r, b, e := ns.Rmdir(location)
		if e != nil {
			log.Fatal(e)
		}

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	case "du":
		r, b, e := ns.Du(location)
		if e != nil {
			log.Fatal(e)
		}

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	default:
		r, b, e := ns.Dir(location)
		if e != nil {
			log.Fatal(e)
		}

		if r.StatusCode == 200 {
			fmt.Printf(b)
		}
	}
}
