//go:build windows
// +build windows

package protocoltalkers

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getTalkers(ctx context.Context) (<-chan storage.ProtocolTalkerItem, <-chan error) {
	return nil, monitor.ErrNotImplemented
}
