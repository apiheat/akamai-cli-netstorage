package main

import (
	"encoding/xml"
	"os"
	"sort"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

// StatNS output structure for stat command
type StatNS struct {
	XMLName   xml.Name `xml:"stat"`
	Directory string   `xml:"directory,attr"`
	Files     []FileNS `xml:"file"`
}

// FileNS output structure for file stat
type FileNS struct {
	XMLName xml.Name `xml:"file"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
}

var (
	configSection, configFile                      string
	nsHostname, nsKeyname, nsKey, nsCpcode, nsPath string
)

const (
	VERSION = "0.0.4"
)

func main() {
	_, inCLI := os.LookupEnv("AKAMAI_CLI")

	appName := "akamai-purge"
	if inCLI {
		appName = "akamai purge"
	}

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "A CLI to interact with Akamai NetStorage"
	app.Version = VERSION
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Petr Artamonov",
		},
		{
			Name: "Rafal Pieniazek",
		},
	}

	dir, _ := homedir.Dir()
	dir += string(os.PathSeparator) + ".edgerc"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "section, s",
			Value:       "netstorage",
			Usage:       "`NAME` of section to use from credentials file",
			Destination: &configSection,
			EnvVar:      "AKAMAI_EDGERC_NETSTORAGE_SECTION",
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       dir,
			Usage:       "Location of the credentials `FILE`",
			Destination: &configFile,
			EnvVar:      "AKAMAI_EDGERC",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "upload",
			Aliases:   []string{"u"},
			Usage:     "Upload files from `DIRECTORY`",
			ArgsUsage: "--from-directory /local/path [DIR]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from-directory",
					Value: "",
					Usage: "Upload files from `DIRECTORY`",
				},
			},
			Action: cmdUpload,
		},
		{
			Name:      "download",
			Aliases:   []string{"d"},
			Usage:     "Download files from `DIRECTORY`",
			ArgsUsage: "--to-directory /local/path [DIR]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "to-directory",
					Value: "",
					Usage: "Download files to `DIRECTORY`",
				},
			},
			Action: cmdDownload,
		},
		{
			Name:    "erase",
			Aliases: []string{"e"},
			Usage:   "Erase all files from `DIRECTORY`",
			Action:  cmdErase,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List `DIRECTORY` content in NetStorage",
			Action:  cmdList,
		},
		{
			Name:    "mkdir",
			Aliases: []string{"md"},
			Usage:   "Create `DIRECTORY` recursively",
			Action:  cmdMkdir,
		},
		{
			Name:    "rmdir",
			Aliases: []string{"rm"},
			Usage:   "Delete empty `DIRECTORY`",
			Action:  cmdRm,
		},
		{
			Name:   "du",
			Usage:  "Show disk usage of `DIRECTORY`",
			Action: cmdDu,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(c *cli.Context) error {
		config(configFile, configSection)

		return nil
	}
	app.Run(os.Args)

}
