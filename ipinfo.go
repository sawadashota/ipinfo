// Package ipinfo look up IP's information via ipinfo.io
package ipinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// IPInfo is response body from ipinfo.io
type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
}

var httpClient = http.DefaultClient

// CustomHTTPClient changes from http.DefaultClient
func CustomHTTPClient(client *http.Client) {
	httpClient = client
}

// Look up given IP's information
func LookUp(ctx context.Context, ip string) (*IPInfo, error) {
	req, err := http.NewRequest("GET", "https://ipinfo.io/"+ip, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipinfo: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i IPInfo
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	if err := i.validate(); err != nil {
		return nil, err
	}

	return &i, nil
}

// validate ip's information is enough
func (i *IPInfo) validate() error {
	if i.IP == "" {
		return errors.New(`ipinfo: key "ip" should not be empty`)
	}
	if i.Hostname == "" {
		return errors.New(`ipinfo: key "hostname" should not be empty`)
	}
	return nil
}

// IsEC2 or not
func (i *IPInfo) IsEC2() bool {
	return strings.HasPrefix(i.Hostname, "ec2-")
}

// IsGCP or not
func (i *IPInfo) IsGCP() bool {
	return strings.HasSuffix(i.Hostname, ".bc.googleusercontent.com")
}

// IsGoogleBot or not
func (i *IPInfo) IsGoogleBot() bool {
	return strings.HasSuffix(i.Hostname, ".googlebot.com")
}
