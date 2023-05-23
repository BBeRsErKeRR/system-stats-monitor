//go:build darwin
// +build darwin

package bpstalkers

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/monitor"
)

func getBps(ctx context.Context) (<-chan interface{}, <-chan error) {
	res := make(chan interface{})
	errC := make(chan error)
	defer close(errC)
	errC <- monitor.ErrNotImplemented
	return res, errC
}

func checkCall(ctx context.Context) error {
	return monitor.ErrNotImplemented
}
