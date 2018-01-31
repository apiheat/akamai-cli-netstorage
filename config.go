package main

import (
	"log"

	"github.com/go-ini/ini"
)

func config(configFile, configSection string) {
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

	return //nsHostname, nsKeyname, nsKey, nsCpcode, nsPath
}
