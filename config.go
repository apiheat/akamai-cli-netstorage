package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-ini/ini"
	"github.com/urfave/cli"
)

func config(configFile, configSection string) {
	var section *ini.Section

	cfg, err := ini.Load(configFile)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Printf("'%s' does not exist. Please run '%s --config Your_Configuration_File'...\n", configFile, appName)
		color.Unset()
	} else {
		section, err = cfg.GetSection(configSection)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Printf("Section '%s' does not exist in %s. Please run '%s --section Your_Section_Name' ...\n", configSection, configFile, appName)
			color.Unset()

			cli.NewExitError(err.Error(), 1)

			return
		}

		color.Set(color.FgBlue)
		fmt.Printf("# Credentials used: [%s] from %s\n", configSection, configFile)
		color.Unset()

		nsHostname = section.Key("hostname").String()
		nsKeyname = section.Key("keyname").String()
		nsKey = section.Key("key").String()
		nsCpcode = section.Key("cpcode").String()
		nsPath = section.Key("path").String()

	}

	return
}
