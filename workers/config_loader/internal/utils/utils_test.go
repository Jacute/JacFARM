package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandCIDR(t *testing.T) {
	testcases := []struct {
		name      string
		cidr      string
		addresses []string
	}{
		{
			name: "ok 10.10.0.0/30",
			cidr: "10.10.0.0/30",
			addresses: []string{
				"10.10.0.1",
				"10.10.0.2",
			},
		},
		{
			name: "ok 192.168.55.0/29",
			cidr: "192.168.55.0/29",
			addresses: []string{
				"192.168.55.1",
				"192.168.55.2",
				"192.168.55.3",
				"192.168.55.4",
				"192.168.55.5",
				"192.168.55.6",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			addresses, err := ExpandCIDR(tc.cidr)
			require.NoError(t, err)
			require.Equal(t, tc.addresses, addresses)
		})
	}
}

func TestExpandRange(t *testing.T) {
	testcases := []struct {
		name      string
		iprange   string
		addresses []string
	}{
		{
			name:    "ok 10.10.0.254-10.10.1.1",
			iprange: "10.10.0.254-10.10.1.1",
			addresses: []string{
				"10.10.0.254",
				"10.10.0.255",
				"10.10.1.0",
				"10.10.1.1",
			},
		},
		{
			name:    "ok 10.10.5.0-10.10.5.2",
			iprange: "10.10.5.0-10.10.5.2",
			addresses: []string{
				"10.10.5.0",
				"10.10.5.1",
				"10.10.5.2",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			addresses, err := ExpandRange(tc.iprange)
			require.NoError(t, err)
			require.Equal(t, tc.addresses, addresses)
		})
	}
}

func TestExpandIpFromN(t *testing.T) {
	testcases := []struct {
		name             string
		nStart, nEnd     int
		offsetX, offsetY int
		block            int
		ipTmpl           string
		addresses        []string
	}{
		{
			name:    "ok1",
			nStart:  0,
			nEnd:    6,
			offsetX: 32,
			offsetY: 1,
			block:   2,
			ipTmpl:  "10.{X}.{Y}.1",
			addresses: []string{
				"10.32.1.1",
				"10.32.2.1",
				"10.33.1.1",
				"10.33.2.1",
				"10.34.1.1",
				"10.34.2.1",
			},
		},
		{
			name:    "ok2",
			nStart:  0,
			nEnd:    6,
			offsetX: 25,
			offsetY: 0,
			block:   5,
			ipTmpl:  "10.{X}.{Y}.1",
			addresses: []string{
				"10.25.0.1",
				"10.25.1.1",
				"10.25.2.1",
				"10.25.3.1",
				"10.25.4.1",
				"10.26.0.1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			addresses := ExpandIpFromN(
				tc.nStart, tc.nEnd,
				tc.offsetX, tc.offsetY,
				tc.block, tc.ipTmpl,
			)
			require.Equal(t, tc.addresses, addresses)
		})
	}
}
