package main

import "sync"

var (
	// Apps stores all tracked applications.
	Apps map[Application]Metric
	// AppLock guards additions to Apps to ensure thread safety
	AppLock sync.Mutex
)

// Application represts a single version of a particular application.
type Application struct {
	Name    string
	Version string
}

// Metric tracks total metric counters for a version of an application.
type Metric struct {
	TotalRequestsCount uint
	TotalSuccessCount  uint
	TotalErrorCount    uint
}

func (m Metric) IncrementRequestCount(num uint) Metric {
	m.TotalRequestsCount += num
	return m
}

func (m Metric) IncrementSuccessCount(num uint) Metric {
	m.TotalSuccessCount += num
	return m
}

func (m Metric) IncrementErrorCount(num uint) Metric {
	m.TotalErrorCount += num
	return m
}

// IncrementCounters adds the status metrics to any applications defined in apps,
// or creates new entries in apps where necessary.
func IncrementCounters(apps map[Application]Metric, status HostStatus) {
	AppLock.Lock()
	defer AppLock.Unlock()

	var m Metric
	app := Application{Name: status.Application, Version: status.Version}
	_, ok := apps[app]
	if ok {
		m = apps[app]
	}
	m = m.IncrementRequestCount(status.RequestsCount)
	m = m.IncrementSuccessCount(status.SuccessCount)
	m = m.IncrementErrorCount(status.ErrorCount)
	apps[app] = m
}
