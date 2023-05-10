package cpu

import (
	"context"
	"strconv"
	"strings"
	"sync"
)

type CPUTimeStat struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

type AverageStat struct {
	sync.Mutex
	lastCPUTimes []CPUTimeStat
}

func (as *AverageStat) Add(ctx context.Context) error {
	times, err := GetCpuTimes(ctx)
	if err != nil {
		return err
	}
	as.Lock()
	defer as.Unlock()
	as.lastCPUTimes = append(as.lastCPUTimes, *times)
	return nil
}

func (as *AverageStat) Avg() (*CPUTimeStat, error) {
	var sumUser, sumSystem, sumIdle float64
	as.Lock()
	defer as.Unlock()
	for _, stat := range as.lastCPUTimes {
		sumUser += stat.User
		sumSystem += stat.System
		sumIdle += stat.Idle
	}
	totalLen := len(as.lastCPUTimes)
	return &CPUTimeStat{
		User:   sumUser / float64(totalLen),
		System: sumSystem / float64(totalLen),
		Idle:   sumIdle / float64(totalLen),
	}, nil
}

func (c CPUTimeStat) JsonString() string {
	v := []string{
		`"user":` + strconv.FormatFloat(c.User, 'f', 1, 64),
		`"system":` + strconv.FormatFloat(c.System, 'f', 1, 64),
		`"idle":` + strconv.FormatFloat(c.Idle, 'f', 1, 64),
	}

	return `{` + strings.Join(v, ",") + `}`
}
