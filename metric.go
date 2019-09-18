package main

import "sync"

var (
	// Apps stores all tracked applications.
	Apps map[Application]Metric
	// AppLock guards additions to Apps to ensure thread safety
	AppLock sync.Mutex
)

type Application struct {
	Name    string
	Version string
}

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
