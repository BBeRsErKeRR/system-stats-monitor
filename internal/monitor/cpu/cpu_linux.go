//go:build linux
// +build linux

package cpu

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/tklauser/go-sysconf"

	files "github.com/BBeRsErKeRR/system-stats-monitor/pkg/files"
)

var (
	ErrorGetStat = errors.New("stat does not contain cpu info")
	ErrorGetCpu  = errors.New("not contain cpu")
	ClocksPerSec = float64(100)
)

func init() {
	clkTck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	// ignore errors
	if err == nil {
		ClocksPerSec = float64(clkTck)
	}
}

func parseStatLine(line string) (*CPUTimeStat, error) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return nil, ErrorGetStat
	}

	if !strings.HasPrefix(fields[0], "cpu") {
		return nil, ErrorGetCpu
	}

	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}

	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}

	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}

	ct := &CPUTimeStat{
		User:   user / ClocksPerSec,
		System: system / ClocksPerSec,
		Idle:   idle / ClocksPerSec,
	}
	return ct, nil
}

func GetCpuTimes(ctx context.Context) (*CPUTimeStat, error) {
	lines, err := files.ReadLinesOffsetN("/proc/stat", 0, 1)
	if err != nil {
		return nil, err
	}
	return parseStatLine(lines[0])
}
