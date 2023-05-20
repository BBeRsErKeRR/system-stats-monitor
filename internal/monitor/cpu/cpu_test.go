package cpu

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	mockstorage "github.com/BBeRsErKeRR/system-stats-monitor/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func skipIfNotImplementedErr(t *testing.T, err error) {
	t.Helper()
	if errors.Is(err, monitor.ErrNotImplemented) {
		t.Skip("not implemented")
	}
}

func CommandRunner(ctx context.Context) (*storage.CPUTimeStat, error) {
	return &storage.CPUTimeStat{User: 10.0, System: 20.0, Idle: 30.0}, nil
}

func TestStatCollector_Grab(t *testing.T) {
	ctx := context.Background()
	mockStorage := mockstorage.New()
	mockLogger, err := logger.New(&logger.Config{})
	require.NoError(t, err)
	mockStorage.On("StoreStats", mock.Anything, mock.Anything).Return(nil, nil).Once()
	collector := New(mockStorage, mockLogger)
	commandRunner = CommandRunner
	err = collector.Grab(ctx)
	require.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestStatCollector_GetStats(t *testing.T) {
	ctx := context.Background()
	mockStorage := mockstorage.New()
	mockLogger, err := logger.New(&logger.Config{})
	require.NoError(t, err)
	now := time.Now()
	data := []storage.Metric{
		{
			Date:     now.Add(-time.Second),
			StatInfo: storage.CPUTimeStat{User: 10.0, System: 20.0, Idle: 30.0},
		},
		{
			Date:     now,
			StatInfo: storage.CPUTimeStat{User: 20.0, System: 30.0, Idle: 40.0},
		},
	}
	mockStorage.On("GetStats", ctx, mock.Anything).Return(data, nil)
	collector := New(mockStorage, mockLogger)
	ret, err := collector.GetStats(ctx, int64(time.Minute.Seconds()))
	skipIfNotImplementedErr(t, err)
	require.NoError(t, err)
	expected := NewCPUTimeStat((10.0+20.0)/2, (20.0+30.0)/2, (30.0+40.0)/2)
	assert.Equal(t, expected, ret)
	mockStorage.AssertExpectations(t)
}
