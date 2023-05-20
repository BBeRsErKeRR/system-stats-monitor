//go:build windows
// +build windows

package cpu

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	command "github.com/BBeRsErKeRR/system-stats-monitor/pkg/command"
)

var ErrInvalidData = errors.New("incorrect output")

func parseData(output string) (*storage.CPUTimeStat, error) {
	var err error
	lines := strings.Split(output, "\r\n")
	if len(lines) < 3 {
		return nil, ErrInvalidData
	}

	fields := strings.Split(lines[2], ",")

	user, err := strconv.ParseFloat(strings.Trim(fields[2], "\""), 64)
	if err != nil {
		return nil, err
	}

	system, err := strconv.ParseFloat(strings.Trim(fields[1], "\""), 64)
	if err != nil {
		return nil, err
	}

	idle, err := strconv.ParseFloat(strings.Trim(fields[3], "\""), 64)
	if err != nil {
		return nil, err
	}
	ct := NewCPUTimeStat(
		user/ClocksPerSec,
		system/ClocksPerSec,
		idle/ClocksPerSec,
	)
	return &ct, nil
}

func getCPUTimes(ctx context.Context) (*storage.CPUTimeStat, error) {
	out, err := command.WithContext(ctx, "typeperf",
		`Processor Information(_Total)\% Privileged Time`,
		`Processor Information(_Total)\% User Time`,
		`Processor Information(_Total)\% Idle Time`,
		"-sc",
		"1",
	)
	if err != nil {
		return nil, err
	}
	return parseData(out)
}
