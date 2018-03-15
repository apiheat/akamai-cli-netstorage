package main

import (
	"encoding/xml"
	"os"
	"sort"

	"github.com/fatih/color"
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
	Mtime   string   `xml:"mtime,attr"`
	Size    string   `xml:"size,attr"`
	MD5     string   `xml:"md5,attr"`
}

// duNS output structure for du command
type duNS struct {
	XMLName   xml.Name `xml:"du"`
	Directory string   `xml:"directory,attr"`
	Info      duInfoNS `xml:"du-info"`
}

// duInfoNS output structure for DU info
type duInfoNS struct {
	XMLName xml.Name `xml:"du-info"`
	Files   string   `xml:"files,attr"`
	Bytes   string   `xml:"bytes,attr"`
}

var (
	version                                          string
	configSection, configFile, configCpcode, appName string
	nsHostname, nsKeyname, nsKey, nsCpcode, nsPath   string
	colorOn                                          bool
)

// VERSION
const (
	padding = 3
)

func main() {
	_, inCLI := os.LookupEnv("AKAMAI_CLI")

	appName := "akamai-netstorage"
	if inCLI {
		appName = "akamai netstorage"
	}

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "A CLI to interact with Akamai NetStorage"
	app.Version = version
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
			Name:        "cpcode",
			Value:       "",
			Usage:       "`CP CODE` to use",
			Destination: &configCpcode,
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       dir,
			Usage:       "Location of the credentials `FILE`",
			Destination: &configFile,
			EnvVar:      "AKAMAI_EDGERC",
		},
		cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable color output",
			Destination: &colorOn,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "put",
			Usage:     "Upload files from `DIRECTORY`",
			ArgsUsage: "--from /local/path [DIR]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from",
					Value: "",
					Usage: "Upload files from `DIRECTORY`",
				},
			},
			Action: cmdPut,
		},
		{
			Name:      "get",
			Aliases:   []string{"g"},
			Usage:     "Download from `OBJECT`",
			ArgsUsage: "--to /local/path [OBJECT]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "to",
					Value: "",
					Usage: "Download files to `DIRECTORY`",
				},
			},
			Action: cmdGet,
		},
		{
			Name:    "rm",
			Aliases: []string{"delete"},
			Usage:   "Delete `FILE`",
			Action:  cmdRm,
		},
		{
			Name:    "empty-directory",
			Aliases: []string{"e"},
			Usage:   "Erase all files from `DIRECTORY`, non empty directories inside target `DIRECTORY` will be ignored",
			Action:  cmdErase,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List `OBJECT` content in NetStorage",
			Action:  cmdList,
		},
		{
			Name:    "mkdir",
			Aliases: []string{"md"},
			Usage:   "Create `DIRECTORY` recursively",
			Action:  cmdMkdir,
		},
		{
			Name:  "rmdir",
			Usage: "Delete `DIRECTORY`",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "recursively",
					Usage: "Delete `DIRECTORY` recursively",
				},
			},
			Action: cmdRmdir,
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

		if c.String("cpcode") != "" {
			nsCpcode = c.String("cpcode")
		}

		if c.Bool("no-color") {
			color.NoColor = true
		}

		return nil
	}
	app.Run(os.Args)

}
