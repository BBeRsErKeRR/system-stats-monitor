//go:build linux
// +build linux

package cpu

import (
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStatLine(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		stats *storage.CPUTimeStat
		err   error
	}{
		{
			name: "valid line",
			line: "cpu 18876 2 8731 26088957 9116 0 1068 0 0 0",
			stats: &storage.CPUTimeStat{
				User:   188.76,
				System: 87.31,
				Idle:   260889.57,
			},
			err: nil,
		},
		{
			name:  "empty line",
			line:  "",
			stats: nil,
			err:   ErrorGetStat,
		},
		{
			name:  "invalid prefix",
			line:  "not_cpu 123456 789012 345678 901234 56789 123456 789012 345678 901234 56789",
			stats: nil,
			err:   ErrorGetCPU,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := parseStatLine(tt.line)
			assert.Equal(t, tt.err, err)
			if tt.stats != nil {
				require.Equal(t, tt.stats.User, stats.User)
				require.Equal(t, tt.stats.System, stats.System)
				require.Equal(t, tt.stats.Idle, stats.Idle)
			}
		})
	}
}

func TestGetCPUTimes(t *testing.T) {
	stats, err := getCPUTimes()
	assert.NotNil(t, stats)
	assert.NoError(t, err)
}
