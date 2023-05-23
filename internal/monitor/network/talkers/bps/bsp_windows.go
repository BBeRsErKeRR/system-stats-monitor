//go:build windows
// +build windows

package bpstalkers

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func getBps(ctx context.Context) (<-chan storage.BpsItem, <-chan error) {
	errC := make(chan error)
	defer close(errC)
	errC <- monitor.ErrNotImplemented
	return nil, errC
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}
