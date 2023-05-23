package cpu

import (
	"context"
	"errors"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
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
	mockLogger, err := logger.New(&logger.Config{})
	require.NoError(t, err)
	collector := New(mockLogger)
	commandRunner = CommandRunner
	data, err := collector.Grab(ctx)
	require.NoError(t, err)
	switch v := data.(type) {
	case *storage.CPUTimeStat:
		require.Equal(t, v.User, 10.0)
		require.Equal(t, v.System, 20.0)
		require.Equal(t, v.Idle, 30.0)
	default:
		t.Error("bad returned type")
	}

}

func TestStatCollector_Check(t *testing.T) {
	ctx := context.Background()
	mockLogger, err := logger.New(&logger.Config{})
	require.NoError(t, err)
	collector := New(mockLogger)
	err = collector.CheckCall(ctx)
	skipIfNotImplementedErr(t, err)
	require.NoError(t, err)
}
