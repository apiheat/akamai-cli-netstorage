package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	netstorage "github.com/akamai/netstoragekit-golang"
	"github.com/go-ini/ini"
)

func userHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

type StatNS struct {
	XMLName   xml.Name `xml:"stat"`
	Directory string   `xml:"directory,attr"`
	Files     []FileNS `xml:"file"`
}

type FileNS struct {
	XMLName xml.Name `xml:"file"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
}

func main() {
	config := path.Join(userHome(), ".edgerc")
	cfg, err := ini.Load(config)
	if err != nil {
		fmt.Println("error:", err)
	}

	section, err := cfg.GetSection("netstorage")
	if err != nil {
		fmt.Println("error:", err)
	}

	nsHostname := section.Key("hostname").String()
	nsKeyname := section.Key("keyname").String()
	nsKey := section.Key("key").String()
	nsCpcode := section.Key("cpcode").String()
	nsPath := section.Key("path").String()

	ns := netstorage.NewNetstorage(nsHostname, nsKeyname, nsKey, true)

	nsDestination := fmt.Sprintf("/%s/%s", nsCpcode, nsPath)
	fmt.Printf("Going to download content of NETSTORAGE%s\n", nsDestination)

	res, body, err := ns.Dir(nsDestination)
	if err != nil {
		fmt.Println("error:", err)
	}

	if res.StatusCode == 200 {
		var stat StatNS
		xmlstr := strings.Replace(body, "ISO-8859-1", "UTF-8", -1)
		xml.Unmarshal([]byte(xmlstr), &stat)
		for i := range stat.Files {
			if stat.Files[i].Type == "file" {
				pathLocal := path.Join(userHome(), nsDestination)
				if _, err := os.Stat(pathLocal); os.IsNotExist(err) {
					os.MkdirAll(pathLocal, os.ModePerm)
				}

				nsTargetPath := fmt.Sprintf("%s/%s", nsDestination, stat.Files[i].Name)
				localDestinationPath := fmt.Sprintf("%s/%s", pathLocal, stat.Files[i].Name)
				fmt.Printf("\nDownloading from: %s to %s\n", nsTargetPath, localDestinationPath)
				f, body, err := ns.Download(nsTargetPath, localDestinationPath)
				if err != nil {
					fmt.Println("error:", err)
				}

				if f.StatusCode == 200 {
					fmt.Printf(body)
				}
			}
		}
		fmt.Println()
	}

	// Now time to upload
	targetDir := path.Join(userHome(), nsDestination) // Should be parameter
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		localPath := fmt.Sprintf("%s/%s", targetDir, f.Name())
		// nsTarget := fmt.Sprintf("%s/%s", nsDestination, f.Name())
		// Remove below and uncomment above
		nsTarget := fmt.Sprintf("%s/%s", fmt.Sprintf("/%s/%s", nsCpcode, "test"), f.Name())
		fmt.Printf("\nUploading from: %s to: %s\n", localPath, nsTarget)
		res, body, err := ns.Upload(localPath, nsTarget)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode == 200 {
			fmt.Printf(body)
		}
	}

	r, b, e := ns.Dir(fmt.Sprintf("/%s/%s", nsCpcode, "test"))
	if e != nil {
		log.Fatal(e)
	}

	if r.StatusCode == 200 {
		fmt.Printf(b)
	}

}
