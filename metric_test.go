package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementingCounters(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		appName     string
		version     string
		status      []HostStatus
		expRequests uint
		expSuccess  uint
	}{
		{
			appName: "foo",
			version: "1.1",
			status: []HostStatus{
				{
					RequestsCount: 2,
					SuccessCount:  1,
				},
				{
					Version:       "1.1",
					RequestsCount: 2,
					SuccessCount:  1,
				}},
			expRequests: 4,
			expSuccess:  2,
		},
		{
			appName: "bar",
			version: "1.2",
			status: []HostStatus{
				{
					RequestsCount: 10,
					SuccessCount:  5,
				},
				{
					RequestsCount: 20,
					SuccessCount:  0,
				}},
			expRequests: 30,
			expSuccess:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.appName, func(t *testing.T) {
			apps := make(map[Application]Metric)
			app := Application{Name: tt.appName, Version: tt.version}
			for _, s := range tt.status {
				// just assign the name and ver here to reduce cruft in the tests struct
				s.Application = tt.appName
				s.Version = tt.version

				IncrementCounters(apps, s)
			}
			assert.Equal(tt.expRequests, apps[app].TotalRequestsCount)
			assert.Equal(tt.expSuccess, apps[app].TotalSuccessCount)
		})
	}

}
