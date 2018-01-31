# Akamai CLI for Netstorage
*NOTE:* This tool is intended to be installed via the Akamai CLI package manager, which can be retrieved from the releases page of the [Akamai CLI](https://github.com/akamai/cli) tool.

### Local Install, if you choose not to use the akamai package manager
If you want to compile it from source, you will need Go 1.9 or later, and the [Glide](https://glide.sh) package manager installed:
1. Fetch the package:
   `go get https://github.com/partamonov/akamai-cli-netstorage`
1. Change to the package directory:
   `cd $GOPATH/src/github.com/partamonov/akamai-cli-netstorage`
1. Install dependencies using Glide:
   `glide install`
1. Compile the binary:
   `go build -ldflags="-s -w" -o akamai-netstorage`

### Credentials
In order to use this configuration, you need to:
* Set up your credential files as described in the [authorization](https://developer.akamai.com/introduction/Prov_Creds.html) and [credentials](https://developer.akamai.com/introduction/Conf_Client.html) sections of the getting started guide on developer.akamai.com.

Expects `netstorage` section in .edgerc

```
[netstorage]
hostname = XXXXXXXX
key = XXXXXXX
keyname = XXXXXXXX
cpcode = 9999999
path = /some/path or ""
```

## Overview
The Akamai NetStorage Kit is a set of go libraries that wraps Akamai's {OPEN} APIs to help simplify common netstorage tasks.

## Usage
```shell
# akamai netstorage
NAME:
   akamai-netstorage - A CLI to interact with Akamai NetStorage

USAGE:
   akamai-netstorage [global options] command [command options] [arguments...]

VERSION:
   0.0.4

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     download, d  Download files from `DIRECTORY`
     du           Show disk usage of `DIRECTORY`
     erase, e     Erase all files from `DIRECTORY`
     list, ls     List `DIRECTORY` content in NetStorage
     mkdir, md    Create `DIRECTORY` recursively
     rmdir, rm    Delete empty `DIRECTORY`
     upload, u    Upload files from `DIRECTORY`
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/partamonov/.edgerc") [$AKAMAI_EDGERC]
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "netstorage") [$AKAMAI_EDGERC_NETSTORAGE_SECTION]
   --help, -h               show help
   --version, -v            print the version
```
