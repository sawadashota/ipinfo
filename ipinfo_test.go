package ipinfo_test

import (
	"testing"

	"github.com/sawadashota/ipinfo"
)

func TestIPInfo_IsEC2(t *testing.T) {
	type fields struct {
		IP       string
		Hostname string
		City     string
		Region   string
		Country  string
		Loc      string
		Org      string
	}
	cases := map[string]struct {
		fields fields
		want   bool
	}{
		"ec2 hostname": {
			fields: fields{
				Hostname: "ec2-0-0-0-0.ap-northeast-1.compute.amazonaws.com",
			},
			want: true,
		},
		"so-net hostname": {
			fields: fields{
				Hostname: "pb0000000.tubecm00.ap.so-net.ne.jp",
			},
			want: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			i := &ipinfo.IPInfo{
				IP:       c.fields.IP,
				Hostname: c.fields.Hostname,
				City:     c.fields.City,
				Region:   c.fields.Region,
				Country:  c.fields.Country,
				Loc:      c.fields.Loc,
				Org:      c.fields.Org,
			}
			if i.IsEC2() != c.want {
				t.Errorf("IPInfo.IsEC2() = %v, want %v", i.IsEC2(), c.want)
			}
		})
	}
}

func TestIPInfo_IsGCP(t *testing.T) {
	type fields struct {
		IP       string
		Hostname string
		City     string
		Region   string
		Country  string
		Loc      string
		Org      string
	}
	cases := map[string]struct {
		fields fields
		want   bool
	}{
		"GCE hostname": {
			fields: fields{
				Hostname: "0.0.0.0.bc.googleusercontent.com",
			},
			want: true,
		},
		"so-net hostname": {
			fields: fields{
				Hostname: "pb0000000.tubecm00.ap.so-net.ne.jp",
			},
			want: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			i := &ipinfo.IPInfo{
				IP:       c.fields.IP,
				Hostname: c.fields.Hostname,
				City:     c.fields.City,
				Region:   c.fields.Region,
				Country:  c.fields.Country,
				Loc:      c.fields.Loc,
				Org:      c.fields.Org,
			}
			if i.IsGCP() != c.want {
				t.Errorf("IPInfo.IsEC2() = %v, want %v", i.IsGCP(), c.want)
			}
		})
	}
}
