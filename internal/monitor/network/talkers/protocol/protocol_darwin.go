//go:build darwin
// +build darwin

package protocoltalkers

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getTalkers(ctx context.Context) (<-chan storage.ProtocolTalkerItem, <-chan error) {
	errC := make(chan error)
	defer close(errC)
	errC <- monitor.ErrNotImplemented
	return nil, errC
}
