//go:build linux
// +build linux

package cpu

import (
	"errors"
	"strconv"
	"strings"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	files "github.com/BBeRsErKeRR/system-stats-monitor/pkg/files"
	"github.com/tklauser/go-sysconf"
)

var (
	ErrorGetStat = errors.New("stat does not contain cpu info")
	ErrorGetCPU  = errors.New("not contain cpu")
	ClocksPerSec = float64(100)
)

func init() {
	clkTck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	// ignore errors
	if err == nil {
		ClocksPerSec = float64(clkTck)
	}
}

func parseStatLine(line string) (*storage.CPUTimeStat, error) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return nil, ErrorGetStat
	}

	if !strings.HasPrefix(fields[0], "cpu") {
		return nil, ErrorGetCPU
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

	ct := NewCPUTimeStat(
		user/ClocksPerSec,
		system/ClocksPerSec,
		idle/ClocksPerSec,
	)
	return &ct, nil
}

func getCPUTimes() (*storage.CPUTimeStat, error) {
	lines, err := files.ReadLinesOffsetN("/proc/stat", 0, 1)
	if err != nil {
		return nil, err
	}
	return parseStatLine(lines[0])
}
