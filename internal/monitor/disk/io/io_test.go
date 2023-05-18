package diskio

import (
	"context"
	"errors"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	mockstorage "github.com/BBeRsErKeRR/system-stats-monitor/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errTestData = errors.New("command failed")

func TestGrab(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		commandRunner func(ctx context.Context) ([]interface{}, error)
		expectedError error
	}{
		{
			name: "success",
			commandRunner: func(ctx context.Context) ([]interface{}, error) {
				res := []interface{}{
					&storage.DiskIoStatItem{},
				}
				return res, nil
			},
			expectedError: nil,
		},
		{
			name: "error",
			commandRunner: func(ctx context.Context) ([]interface{}, error) {
				return nil, errTestData
			},
			expectedError: errTestData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := mockstorage.New()
			mockLogger, err := logger.New(&logger.Config{})
			require.NoError(t, err)
			collector := New(st, mockLogger)
			commandRunner = tt.commandRunner
			st.On("BulkStoreStats", mock.Anything, mock.Anything).Return(nil, nil).Once()
			err = collector.Grab(ctx)
			if tt.expectedError != nil {
				assert.Error(t, err)
				require.ErrorIs(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
