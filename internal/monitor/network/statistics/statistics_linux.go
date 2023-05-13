//go:build linux
// +build linux

package networkstatistics

import (
	"context"
	"strconv"
	"strings"

	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseNetstatOut(output string) ([]NetworkStatsItem, error) {
	lines := strings.Split(output, "\n")
	startLine := 1
	if strings.Contains(lines[0], "Not all processes could") {
		startLine = 4
	}
	result := make([]NetworkStatsItem, 0, len(output)-startLine)
	for _, line := range lines[startLine:] {
		values := strings.Fields(line)
		if len(values) < 1 {
			continue
		}
		pidWithCommand := strings.Split(values[8], "/")

		pid, err := strconv.ParseInt(pidWithCommand[0], 10, 32)
		if err != nil {
			return nil, err
		}
		user, err := strconv.ParseInt(values[6], 10, 32)
		if err != nil {
			return nil, err
		}
		portAddress := strings.Split(values[3], ":")
		port, err := strconv.ParseInt(portAddress[len(portAddress)-1], 10, 32)
		if err != nil {
			return nil, err
		}
		item := NetworkStatsItem{
			Command:  strings.Join(pidWithCommand[1:], "/"),
			PID:      int32(pid),
			User:     int32(user),
			Protocol: values[0],
			Port:     int32(port),
		}
		result = append(result, item)
	}
	return result, nil
}

func GetNS(ctx context.Context) ([]NetworkStatsItem, error) {
	out, err := command.CommandWithContext(ctx, "netstat", "-lntupe")
	if err != nil {
		return nil, err
	}
	return parseNetstatOut(string(out))
}
