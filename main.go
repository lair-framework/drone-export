package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/lair-framework/api-server/client"
)

const (
	version = "1.0.0"
	tool    = "drone-export"
	usage   = `
Downloads a lair project.

Usage:
  drone-export [options] <id>
Options:
  -v              show version and exit
  -h              show usage and exit
  -k              allow insecure SSL connections
`
)

func main() {
	showVersion := flag.Bool("v", false, "")
	insecureSSL := flag.Bool("k", false, "")
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()
	if *showVersion {
		log.Println(version)
		os.Exit(0)
	}
	lairURL := os.Getenv("LAIR_API_SERVER")
	if lairURL == "" {
		log.Fatal("Fatal: Missing LAIR_API_SERVER environment variable")
	}
	if len(flag.Args()) < 1 {
		log.Fatal("Fatal: Missing required argument")
	}
	lairPID := flag.Arg(0)
	u, err := url.Parse(lairURL)
	if err != nil {
		log.Fatalf("Fatal: Error parsing LAIR_API_SERVER URL. Error %s", err.Error())
	}
	if u.User == nil {
		log.Fatal("Fatal: Missing username and/or password")
	}
	user := u.User.Username()
	pass, _ := u.User.Password()
	if user == "" || pass == "" {
		log.Fatal("Fatal: Missing username and/or password")
	}
	c, err := client.New(&client.COptions{
		User:               user,
		Password:           pass,
		Host:               u.Host,
		Scheme:             u.Scheme,
		InsecureSkipVerify: *insecureSSL,
	})
	if err != nil {
		log.Fatalf("Fatal: Error setting up client: Error %s", err.Error())
	}
	project, err := c.ExportProject(lairPID)
	if err != nil {
		log.Fatalf("Fatal: Unable to import project. Error %s", err.Error())
	}
	data, err := json.Marshal(project)
	if err != nil {
		log.Fatalf("Fatal: Unable to parse JSON. Error %s", err.Error())
	}
	fmt.Println(string(data))
}
