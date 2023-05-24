//go:build linux
// +build linux

package bpstalkers

import (
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestParseTCPDumpOut(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
		err   error
	}{
		{
			input: "2023-05-19 19:23:26.743536 IP 10.0.0.1.42545 > 10.0.0.2.39510: tcp 100",
			want: storage.BpsItem{
				Source:      "10.0.0.2.39510",
				Destination: "10.0.0.1.42545",
				Numbers:     100,
				Protocol:    "tcp",
			},
			err: nil,
		},
		{
			input: "2023-05-19 20:00:01.184133 IP 10.0.0.1.42545 > 10.0.0.2.39510: UDP, length 31",
			want: storage.BpsItem{
				Source:      "10.0.0.2.39510",
				Destination: "10.0.0.1.42545",
				Numbers:     31,
				Protocol:    "UDP",
			},
			err: nil,
		},
		{
			input: "2023-05-19 20:00:05.935911 ARP, Request who-has 10.0.0.1.42545 (00:25:4d:0c:e7:b1) tell 10.0.0.2.39510, length 28", //nolint:lll
			want:  nil,
			err:   nil,
		},
		{
			input: "",
			want:  nil,
			err:   nil,
		},
	}

	for _, tt := range tests {
		got, err := parseTCPDumpOut(tt.input)
		require.Equal(t, tt.want, got)
		require.Equal(t, tt.err, err)
	}
}
