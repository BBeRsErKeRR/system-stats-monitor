package diskio

import (
	"context"
	"errors"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTestData = errors.New("command failed")

func TestGrab(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		commandRunner func(ctx context.Context) ([]storage.DiskIoStatItem, error)
		expected      []storage.DiskIoStatItem
		expectedError error
	}{
		{
			name: "success",
			commandRunner: func(ctx context.Context) ([]storage.DiskIoStatItem, error) {
				res := []storage.DiskIoStatItem{
					{},
				}
				return res, nil
			},
			expectedError: nil,
		},
		{
			name: "error",
			commandRunner: func(ctx context.Context) ([]storage.DiskIoStatItem, error) {
				return nil, errTestData
			},
			expectedError: errTestData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger, err := logger.New(&logger.Config{})
			require.NoError(t, err)
			collector := New(mockLogger)
			commandRunner = tt.commandRunner
			res, err := collector.Grab(ctx)
			if tt.expectedError != nil {
				assert.Error(t, err)
				require.ErrorIs(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			if len(tt.expected) > 0 {
				require.ElementsMatch(t, tt.expected, res)
			}
		})
	}
}
