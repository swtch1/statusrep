package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidJsonReturnsUnmarshalError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{"invalid}}`)
		assert.Nil(err)
	}))
	defer ts.Close()

	host := Host{URL: ts.URL}
	err := host.GetHostStatus()
	assert.True(strings.Contains(err.Error(), ErrUnmarshal.Error()))
}

func TestGettingHostStatusSetsRequestsCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		expRequestsCount uint
	}{
		{12345},
		{56789},
	}

	for _, tt := range tests {
		t.Run(string(tt.expRequestsCount), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"requests_count": %d}`, tt.expRequestsCount)
				assert.Nil(err)
			}))
			defer ts.Close()

			host := Host{URL: ts.URL}
			err := host.GetHostStatus()
			assert.Nil(err)
			assert.Equal(tt.expRequestsCount, host.Status.RequestsCount)
		})
	}
}

func TestGettingHostStatusSetsApplication(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		expApplication string
	}{
		{"app1"},
		{"app2"},
	}

	for _, tt := range tests {
		t.Run(tt.expApplication, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"application": "%s"}`, tt.expApplication)
				assert.Nil(err)
			}))
			defer ts.Close()

			host := Host{URL: ts.URL}
			err := host.GetHostStatus()
			assert.Nil(err)
			assert.Equal(tt.expApplication, host.Status.Application)
		})
	}
}

func TestGettingHostStatusSetsVersion(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		expVersion string
	}{
		{"ver1.2"},
		{"v5"},
	}

	for _, tt := range tests {
		t.Run(tt.expVersion, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"Version": "%s"}`, tt.expVersion)
				assert.Nil(err)
			}))
			defer ts.Close()

			host := Host{URL: ts.URL}
			err := host.GetHostStatus()
			assert.Nil(err)
			assert.Equal(tt.expVersion, host.Status.Version)
		})
	}
}

func TestGettingHostStatusSetsErrorCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		expErrorCount uint
	}{
		{12345},
		{56789},
	}

	for _, tt := range tests {
		t.Run(string(tt.expErrorCount), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"error_count": %d}`, tt.expErrorCount)
				assert.Nil(err)
			}))
			defer ts.Close()

			host := Host{URL: ts.URL}
			err := host.GetHostStatus()
			assert.Nil(err)
			assert.Equal(tt.expErrorCount, host.Status.ErrorCount)
		})
	}
}

func TestGettingHostStatusSetsSuccessCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		expSuccessCount uint
	}{
		{12345},
		{56789},
	}

	for _, tt := range tests {
		t.Run(string(tt.expSuccessCount), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"success_count": %d}`, tt.expSuccessCount)
				assert.Nil(err)
			}))
			defer ts.Close()

			host := Host{URL: ts.URL}
			err := host.GetHostStatus()
			assert.Nil(err)
			assert.Equal(tt.expSuccessCount, host.Status.SuccessCount)
		})
	}
}

type TestReadCloser struct{ data []byte }

func (rc TestReadCloser) Read(p []byte) (int, error) {
	r := bytes.NewReader(rc.data)
	n, err := r.Read(p)
	if err != nil {
		return 0, err
	}
	return n, io.EOF
}
func (rc TestReadCloser) Close() error { return nil }

func TestGettingHostsTrimsWhiteSpace(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		name     string
		hosts    string
		expSlice []string
	}{
		{"leading_space", " host1\n host2", []string{"host1", "host2"}},
		{"trailing_space", "host1  \nhost2       ", []string{"host1", "host2"}},
		{"extra_space", "  host1  \n   host2     ", []string{"host1", "host2"}},
		{"empty_lines", "  host1  \n \n \n  host2  \n ", []string{"host1", "host2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := []byte(tt.hosts)
			testRC := TestReadCloser{data: b}
			hosts, err := ReadAllHosts(testRC)
			assert.Nil(err)
			assert.Equal(tt.expSlice, hosts)
		})
	}
}

func TestHostURLCreation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		rootURL string
		host    string
		expURL  string
	}{
		{"http://some.root.com", "host1", "http://some.root.com/host1/status"},
		{"http://root.com", "host2", "http://root.com/host2/status"},
		{"http://root.com:80", "host3", "http://root.com:80/host3/status"},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			url, err := HostStatusURL(tt.rootURL, tt.host)
			assert.Nil(err)
			assert.Equal(tt.expURL, url)
		})
	}
}
