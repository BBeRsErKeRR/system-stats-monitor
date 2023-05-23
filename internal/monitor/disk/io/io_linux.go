//go:build linux
// +build linux

package diskio

import (
	"context"
	"errors"
	"strconv"
	"strings"

	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

var ErrOutput = errors.New("bad output")

func parseSSOut(output string) ([]interface{}, error) {
	lines := strings.Split(output, "\n")
	length := len(output)
	if length < 3 {
		return []interface{}{}, ErrOutput
	}
	result := make([]interface{}, 0, length-3)
	for _, line := range lines[3:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		tps, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}

		kbReadS, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		kbWriteS, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, err
		}

		result = append(result, NewDiskIoStatItem(fields[0], tps, kbReadS, kbWriteS))
	}

	return result, nil
}

func collectDiskIo(ctx context.Context) ([]interface{}, error) {
	out, err := command.WithContext(ctx, "iostat", "-d", "-k")
	if err != nil {
		return nil, err
	}
	return parseSSOut(string(out))
}

func checkCall(ctx context.Context) error {
	_, err := command.WithContext(ctx, "iostat", "-d", "-k")
	return err
}
