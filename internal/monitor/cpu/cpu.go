package cpu

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type CPUTimeStat struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

func NewCPUTimeStat(user, system, idle float64) CPUTimeStat {
	return CPUTimeStat{
		User:   user,
		System: system,
		Idle:   idle,
	}
}

type CPUStatCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *CPUStatCollector {
	return &CPUStatCollector{
		name: "cpu",
		st:   st,
	}
}

func (c *CPUStatCollector) Grab(ctx context.Context) error {
	times, err := GetCPUTimes(ctx)
	if err != nil {
		return err
	}
	return c.st.StoreStats(ctx, c.name, *times)
}

func (as *CPUStatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	var sumUser, sumSystem, sumIdle float64
	lastCPUTimes, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	for _, metric := range lastCPUTimes {
		stat := metric.StatInfo.(CPUTimeStat)
		sumUser += stat.User
		sumSystem += stat.System
		sumIdle += stat.Idle
	}
	totalLen := len(lastCPUTimes)
	return NewCPUTimeStat(sumUser/float64(totalLen), sumSystem/float64(totalLen), sumIdle/float64(totalLen)), nil
}
