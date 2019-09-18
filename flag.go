package main

import (
	"fmt"
	"github.com/integrii/flaggy"
)

var (
	defaultLogLevel = "info"
	defaultRootUrl  = "http://storage.googleapis.com/revsreinterview/hosts"
)

// Flag defines command line flags given at runtime.
type Flag struct {
	// Version is the application Version, as taken from the VERSION file.  Version is the exception
	// that is defined a build time rather than runtime.
	Version string
	// LogLevel determines at what level to write application logs.
	LogLevel string
	// HostsFile contains the list of hosts to query, one per line.
	HostsFile string
	// RootUrl is the base url for host status pages.
	RootUrl string
}

func (f *Flag) Parse() {
	flaggy.SetVersion(f.Version)
	flaggy.SetName("statusrep")
	flaggy.SetDescription("Generate reports for hosts with a status endpoint.")

	f.defineAllFlags()
	flaggy.Parse()
	f.setDefaults()
	f.enforceRequirements()
}

func (f *Flag) defineAllFlags() {
	flaggy.String(
		&f.LogLevel,
		"l",
		"log-level",
		fmt.Sprintf("Application log level. This should be one of debug, info, warn, error, fatal. (default: %s)", defaultLogLevel),
	)
	flaggy.String(
		&f.HostsFile,
		"f",
		"hosts-file",
		"File containing the list of servers to query, one per line.",
	)
	flaggy.String(
		&f.RootUrl,
		"r",
		"root-url",
		fmt.Sprintf("The root URL where host paths can be found.  This URL will be prepended to all queries. (default: %s)", defaultRootUrl),
	)
}

func (f *Flag) setDefaults() {
	if f.LogLevel == "" {
		f.LogLevel = defaultLogLevel
	}
	if f.RootUrl == "" {
		f.RootUrl = defaultRootUrl
	}
}

func (f *Flag) enforceRequirements() {
	if f.HostsFile == "" {
		flaggy.ShowHelpAndExit("hosts file is required.")
	}
}
