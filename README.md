# Akamai CLI for Netstorage
*NOTE:* This tool is intended to be installed via the Akamai CLI package manager, which can be retrieved from the releases page of the [Akamai CLI](https://github.com/akamai/cli) tool.

### Local Install, if you choose not to use the akamai package manager
* Go 19.2
* go get https://github.com/partamonov/akamai-cli-netstorage
* cd $GOPATH/src/github.com/partamonov/akamai-cli-netstorage
* go build

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
   netstorage - Akamai CLI

USAGE:
   akamai-netstorage [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
  upload, u    Upload files from directory
  download, d  add a task to the list
  list, ls     list remote directory
  mkdir, md    Create directory
  rmdir, rm    Delete directory
  du           Delete directory
  help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --section NAME, -s NAME  NAME of section to use from .edgerc (default: "netstorage")
   --config FILE, -c FILE   Load configuration from FILE (default: ".edgerc")
   --help, -h               show help
   --version, -v            print the version
```
