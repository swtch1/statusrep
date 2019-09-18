package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	ErrUnmarshal = errors.New("unable to unmarshal status response body")
)

// Host is an external host, with a URL, which has a status endpoint which can be queried.
type Host struct {
	// URL is the full URL where a host status is queried.
	URL string
}

// HostStatus contains status information for a single host, at the time of querying.
type HostStatus struct {
	Application   string `json:"application"`
	Version       string `json:"Version"`
	RequestsCount uint   `json:"requests_count"`
	SuccessCount  uint   `json:"success_count"`
	ErrorCount    uint   `json:"error_count"`
}

// RequestHostStatus gets the HostStatus by making an outbound request to the host status URL.
func (h *Host) RequestHostStatus() (HostStatus, error) {
	var status HostStatus
	b, err := h.getStatus()
	if err != nil {
		return status, err
	}
	if err := json.Unmarshal(b, &status); err != nil {
		return status, errors.Wrapf(err, "%s: '%s'", ErrUnmarshal, h.URL)
	}
	return status, nil
}

func (h *Host) getStatus() ([]byte, error) {
	resp, err := http.Get(h.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get status for '%s'", h.URL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read response body for '%s'", h.URL)
	}
	return body, nil
}

// ReadAllHosts returns a slice of all hosts from a newline delimited reader.
// For example, the output from reading r will look like this:
// host0
// host1
// host2
//
// White space surrounding the hosts will be stripped and empty lines will be discarded.
func ReadAllHosts(r io.ReadCloser) ([]string, error) {
	defer r.Close()
	var hosts []string

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return hosts, errors.Wrap(err, "reading hosts failed")
	}

	for _, h := range strings.Split(buf.String(), "\n") {
		// don't allow space around hosts
		trHost := strings.TrimSpace(h)
		// don't add empty lines
		if len(trHost) > 0 {
			hosts = append(hosts, trHost)
		}
	}
	return hosts, nil
}

// HostStatusURL adds a host and the status endpoint to the path of a root url,
// creating the full URL where a host's status page is expected.
func HostStatusURL(rootURL, host string) (string, error) {
	u, err := url.Parse(rootURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, host, "status")
	return u.String(), nil
}
