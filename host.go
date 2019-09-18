package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var ErrUnmarshal = errors.New("unable to unmarshal status response body")

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

// Host is an external host which has a status and can be queried.
type Host struct {
	// Url is the full URL where a host status is queried.
	Url string
	// Status is populated with status values after the request is made.
	Status HostStatus
}

// GetHostStatus makes an outbound request to the
func (h *Host) GetHostStatus() error {
	b, err := h.getStatus()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &h.Status); err != nil {
		return errors.Wrapf(err, "%s: '%s'", ErrUnmarshal, h.Url)
	}
	return nil
}

func (h *Host) getStatus() ([]byte, error) {
	resp, err := http.Get(h.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get status for '%s'", h.Url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read response body for '%s'", h.Url)
	}
	return body, nil
}

// HostStatus contains status information for a single host, at the time of querying.
type HostStatus struct {
	RequestsCount uint   `json:"requests_count"`
	Application   string `json:"application"`
	Version       string `json:"Version"`
	SuccessCount  uint   `json:"success_count"`
	ErrorCount    uint   `json:"error_count"`
}
