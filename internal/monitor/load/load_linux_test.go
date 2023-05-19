//go:build linux
// +build linux

package load

import (
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestParseStatLine(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       *storage.LoadStat
		wantErrMsg string
	}{
		{
			name:  "valid input",
			input: "0.01 0.05 0.15",
			want: &storage.LoadStat{
				Load1:  0.01,
				Load5:  0.05,
				Load15: 0.15,
			},
		},
		{
			name:       "invalid input",
			input:      "0.01 abc 0.15",
			wantErrMsg: "strconv.ParseFloat: parsing \"abc\": invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStatLine(tt.input)
			if tt.wantErrMsg != "" {
				require.EqualError(t, err, tt.wantErrMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
