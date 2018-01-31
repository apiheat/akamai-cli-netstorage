package main

import (
	"github.com/go-ini/ini"
)

func config(configFile, configSection string) {
	cfg, err := ini.Load(configFile)
	errorCheck(err)

	section, err := cfg.GetSection(configSection)
	errorCheck(err)

	nsHostname = section.Key("hostname").String()
	nsKeyname = section.Key("keyname").String()
	nsKey = section.Key("key").String()
	nsCpcode = section.Key("cpcode").String()
	nsPath = section.Key("path").String()

	return //nsHostname, nsKeyname, nsKey, nsCpcode, nsPath
}
