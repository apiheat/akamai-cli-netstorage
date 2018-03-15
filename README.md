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
   `go build -ldflags="-s -w -X main.version={{.Version}}" -o akamai-netstorage`

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

## Main Usage command
```shell
# akamai netstorage
NAME:
   akamai-netstorage - A CLI to interact with Akamai NetStorage

USAGE:
   akamai-netstorage [global options] command [command options] [arguments...]

VERSION:
   X.X.X

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     du                  Show disk usage of `DIRECTORY`
     empty-directory, e  Erase all files from `DIRECTORY`, non empty directories inside target `DIRECTORY` will be ignored
     get, g              Download from `OBJECT`
     list, ls            List `OBJECT` content in NetStorage
     mkdir, md           Create `DIRECTORY` recursively
     put                 Upload files from `DIRECTORY`
     rm, delete          Delete `FILE`
     rmdir               Delete `DIRECTORY`
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/partamonov/.edgerc") [$AKAMAI_EDGERC]
   --cpcode CP CODE         CP CODE to use
   --no-color               Disable color output
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "netstorage") [$AKAMAI_EDGERC_NETSTORAGE_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

### DU command

If you will DU not directory, then the output will be `ls`
```shell
> akamai netstorage du /test
# Directory               Files   Size
/<CPCODE>/test            10      4.3 MB
```

```shell
> akamai netstorage du /test/icons.png
Type    Name        Mtime                           Size     MD5
File:   icons.png   2018-02-12 10:59:30 +0100 CET   3.4 kB   055531dc54f773403865a704e0bdae21
```

### Empty-directory command

Purpose of this command is to delete all files and empty directories inside target `DIRECTORY`. This is a safety measure to prevent deletion of unexpected content.
If you want true recursive behaviour than open support case for Akamai to enable Quick Deletion for required CP Code.
Use `akamai netstorage rmdir --recursively [PATH]` command.

```shell
> akamai netstorage empty-directory test/super/tttt
```

### Get command
```shell
> akamai-netstorage get [command options] --to /local/path [OBJECT]
```

Purpose of this command is to download `OBJECT` from NetStorage.
If given `OBJECT` is directory, then all files will be downloaded. Directories inside target path won't be touched.

If you have such structure:
- dir1
  - dir2
    - file2
  - dir3
  - file1

Then if you will run `akamai netstorage get dir1`, file1 will be downloaded to `~/<CPCODE>/dir1`.

To set local path for download, please use `--to` option

### List command

Purpose of this command is to list directory content or object information

```shell
> akamai netstorage ls
Directory: /225406
Type   Name      Mtime                            Size   MD5
DIR:   test     2016-07-28 11:41:38 +0200 CEST
DIR:   test2    2018-01-31 07:05:58 +0100 CET
```

### Mkdir command

Purpose of this command is to create directories in Netstorage

### Put command

Purpose of this command is to put `OBJECT` to Netstorage. You can upload all files from given directory, if you will specify directory as source.

```shell
> akamai netstorage put --from /tmp/www /CPCODE/test
```

### Rm command

Purpose of this command is to delete `File` in Netstorage

### Rmdir command

Purpose of this command is to delete `Directory` in Netstorage