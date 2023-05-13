//go:build linux
// +build linux

package networkstates

import (
	"context"
	"strings"

	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseSSOut(output string) (*NetworkStatesStat, error) {
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
			nsCounters[name] += 1
		} else {
			nsCounters[name] = 1
		}

	}

	ret := &NetworkStatesStat{
		Counters: nsCounters,
	}
	return ret, nil
}

func GetNS(ctx context.Context) (*NetworkStatesStat, error) {
	out, err := command.CommandWithContext(ctx, "ss", "-ta")
	if err != nil {
		return nil, err
	}
	return parseSSOut(string(out))
}
