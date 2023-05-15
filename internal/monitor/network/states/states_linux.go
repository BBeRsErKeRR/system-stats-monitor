//go:build linux
// +build linux

package networkstates

import (
	"context"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseSSOut(output string) (*storage.NetworkStatesStat, error) {
	lines := strings.Split(output, "\n")
	nsCounters := make(map[string]int32)

	for _, line := range lines[1:] {
		values := strings.Fields(line)
		if len(values) < 1 {
			continue
		}
		name := values[0]
		_, ok := nsCounters[name]
		if ok {
			nsCounters[name]++
		} else {
			nsCounters[name] = 1
		}
	}

	ret := &storage.NetworkStatesStat{
		Counters: nsCounters,
	}
	return ret, nil
}

func getNS(ctx context.Context) (*storage.NetworkStatesStat, error) {
	out, err := command.WithContext(ctx, "ss", "-ta")
	if err != nil {
		return nil, err
	}
	return parseSSOut(string(out))
}
