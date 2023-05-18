package diskio

import (
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestCollectDiskIo(t *testing.T) {

	tests := []struct {
		name               string
		output             string
		expectedStatsItems []storage.DiskIoStatItem
		expectedErr        error
	}{
		{
			name:               "empty output",
			output:             "",
			expectedStatsItems: []storage.DiskIoStatItem{},
			expectedErr:        ErrOutput,
		},
		{
			name: "valid output",
			output: `Linux 5.4.0-135-generic (SOME-HOST3)  05/18/23        _x86_64_        (16 CPU)

			Device            tps    kB_read/s    kB_wrtn/s    kB_dscd/s
			sda               1.00         0.00         4.00         0.00
			`,
			expectedStatsItems: []storage.DiskIoStatItem{
				{
					Device:   "sda",
					Tps:      1.0,
					KbReadS:  0.0,
					KbWriteS: 4.0,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskIoStats, err := parseSSOut(tt.output)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, diskIoStats, len(tt.expectedStatsItems))
				for i := range diskIoStats {
					actual, expected := diskIoStats[i].(storage.DiskIoStatItem), tt.expectedStatsItems[i]
					assert.Equal(t, expected.Device, actual.Device)
					assert.Equal(t, expected.Tps, actual.Tps)
					assert.Equal(t, expected.KbReadS, actual.KbReadS)
					assert.Equal(t, expected.KbWriteS, actual.KbWriteS)
				}
			}
		})
	}
}
