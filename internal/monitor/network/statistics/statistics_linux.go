//go:build linux
// +build linux

package networkstatistics

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

func parseNetstatOut(output string) ([]interface{}, error) {
	lines := strings.Split(output, "\n")
	startLine := 2
	if strings.Contains(lines[0], "Not all processes could") {
		startLine = 4
	}
	result := make([]interface{}, 0, len(output)-startLine)
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
		item := storage.NetworkStatsItem{
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

func getNS(ctx context.Context) ([]interface{}, error) {
	var out []byte
	var err error
	if os.Geteuid() == 0 { //nolint:nestif
		out, err = command.WithContext(ctx, "netstat", "-lntupe")
		if err != nil {
			return nil, err
		}
	} else {
		out, err = command.WithContext(ctx, "sudo", "netstat", "-lntupe")
		if err != nil {
			out, err = command.WithContext(ctx, "netstat", "-lntupe")
			if err != nil {
				return nil, err
			}
		}
	}

	return parseNetstatOut(string(out))
}
