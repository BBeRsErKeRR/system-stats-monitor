//go:build darwin
// +build darwin

package networkstates

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getNS(ctx context.Context) (*storage.NetworkStatesStat, error) { //nolint:revive
	return nil, monitor.ErrNotImplemented
}

func checkCall(ctx context.Context) error { //nolint:revive
	return monitor.ErrNotImplemented
}
