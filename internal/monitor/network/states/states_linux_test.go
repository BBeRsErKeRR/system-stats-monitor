//go:build linux
// +build linux

package networkstates

import (
	"context"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestParseSSOut(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   *storage.NetworkStatesStat
	}{
		{
			name: "valid input",
			output: `State      Recv-Q Send-Q Local Address:Port               Peer Address:Port              
ESTAB      0      52          127.0.0.1:80                 127.0.0.1:57016               
ESTAB      0      0           127.0.0.1:3306               127.0.0.1:60194`,
			want: &storage.NetworkStatesStat{
				Counters: map[string]int32{
					"ESTAB": 2,
				},
			},
		},
		{
			name:   "empty input",
			output: "",
			want: &storage.NetworkStatesStat{
				Counters: map[string]int32{},
			},
		},
		{
			name: "invalid input with empty lines",
			output: `State      Recv-Q Send-Q Local Address:Port               Peer Address:Port              
          
            `,
			want: &storage.NetworkStatesStat{
				Counters: map[string]int32{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSSOut(tt.output)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetNS(t *testing.T) {
	ctx := context.Background()
	_, err := getNS(ctx)
	require.NoError(t, err)
}
