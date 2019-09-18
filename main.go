package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

// buildVersion should be populated at build time by build ldflags
var buildVersion string

func init() {
	Apps = make(map[Application]Metric)
}

func main() {
	start := time.Now()

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
		log.WithError(err).Fatal("unable to read in hosts")
	}

	var wg sync.WaitGroup
	for _, h := range hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()

			statusURL, err := HostStatusURL(flag.RootURL, h)
			if err != nil {
				log.WithError(err).Errorf("could not create URL for host '%s'", h)
			}
			host := Host{URL: statusURL}

			status, err := host.RequestHostStatus()
			if err != nil {
				log.WithError(err).Errorf("could not get status for host '%s'", h)
			}
			IncrementCounters(Apps, status)
		}(h)
	}
	wg.Wait()

	writeReport(os.Stdout)
	fmt.Printf("\nreport completed in %s\n", time.Now().Sub(start).Truncate(time.Millisecond))
}

func writeReport(w io.Writer) {
	for app, metrics := range Apps {
		var successRate float32
		if metrics.TotalSuccessCount == 0 || metrics.TotalRequestsCount == 0 {
			successRate = 0
		} else {
			successRate = float32(metrics.TotalSuccessCount) / float32(metrics.TotalRequestsCount)
		}
		if _, err := fmt.Fprintf(w, "%s,%s,%.2f\n", app.Name, app.Version, successRate); err != nil {
			log.WithError(err).Error("invalid printer format")
		}
	}
}
