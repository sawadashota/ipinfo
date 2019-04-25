package ipinfo_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sawadashota/ipinfo"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//newTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

const (
	AWS = iota
	GCP
)

// client for testing normal response
// *http.Response: only give value when you want customize response
// use response duration time for testing timeout
func client(t *testing.T, host int, d time.Duration) *http.Client {
	t.Helper()

	var respJSON ipinfo.IPInfo
	switch host {
	case AWS:
		respJSON = ipinfo.IPInfo{
			IP:       "0.0.0.0",
			Hostname: "ec2-18.0.0.0.ap-northeast-1.compute.amazonaws.com",
		}
	case GCP:
		respJSON = ipinfo.IPInfo{
			IP:       "0.0.0.0",
			Hostname: "0.0.0.0.bc.googleusercontent.com",
		}
	default:
		respJSON = ipinfo.IPInfo{
			IP:       "0.0.0.0",
			Hostname: "pb0000000.tubecm00.ap.so-net.ne.jp",
		}
	}

	b, err := json.Marshal(&respJSON)
	if err != nil {
		t.Fatal(err)
	}

	return newTestClient(func(req *http.Request) *http.Response {
		time.Sleep(d)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})
}

func TestLookUp(t *testing.T) {

	type args struct {
		ip string
	}
	type clientArgs struct {
		provider         int
		responseDuration time.Duration
	}
	cases := map[string]struct {
		args           args
		clientArgs     clientArgs
		requestTimeout time.Duration
		want           *ipinfo.IPInfo
		wantErr        bool
	}{
		"AWS": {
			args: args{
				ip: "0.0.0.0",
			},
			clientArgs: clientArgs{
				provider: AWS,
			},
			want: &ipinfo.IPInfo{
				IP:       "0.0.0.0",
				Hostname: "ec2-18.0.0.0.ap-northeast-1.compute.amazonaws.com",
			},
			wantErr: false,
		},
		"timeout": {
			args: args{
				ip: "0.0.0.0",
			},
			clientArgs: clientArgs{
				provider:         AWS,
				responseDuration: 300 * time.Millisecond,
			},
			requestTimeout: 200 * time.Millisecond,
			want: &ipinfo.IPInfo{
				IP:       "0.0.0.0",
				Hostname: "ec2-18.0.0.0.ap-northeast-1.compute.amazonaws.com",
			},
			wantErr: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ipinfo.CustomHTTPClient(client(t, c.clientArgs.provider, c.clientArgs.responseDuration))

			ctx := context.Background()
			if c.requestTimeout > 0 {
				ctx, _ = context.WithTimeout(context.Background(), c.requestTimeout)
			}

			got, err := ipinfo.LookUp(ctx, c.args.ip)
			if (err != nil) != c.wantErr {
				t.Errorf("LookUp() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("LookUp() = %v, want %v", got, c.want)
			}
		})
	}
}
