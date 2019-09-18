package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// buildVersion should be populated at build time by build ldflags
var buildVersion string

func main() {
	var flag Flag
	flag.Version = buildVersion
	flag.Parse()

	SetLogger(os.Stderr, flag.LogLevel, "text", false)

	f, err := os.Open(flag.HostsFile)
	if err != nil {
		log.WithError(err).Fatalf("unable to open file '%s'", flag.HostsFile)
	}

	hosts, err := ReadAllHosts(f)
	if err != nil {
		log.WithError(err).Fatalf("unable to open file '%s'", flag.HostsFile)
	}

	for _, host := range hosts {
		statusURL, err := HostStatusURL(flag.RootURL, host)
		if err != nil {
			log.WithError(err).Errorf("could not get status for host '%s'", host)
		}
		fmt.Println(statusURL)
	}
}
