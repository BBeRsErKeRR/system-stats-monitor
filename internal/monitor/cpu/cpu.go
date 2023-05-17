package cpu

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func NewCPUTimeStat(user, system, idle float64) storage.CPUTimeStat {
	return storage.CPUTimeStat{
		User:   user,
		System: system,
		Idle:   idle,
	}
}

type StatCollector struct {
	st     storage.Storage
	logger logger.Logger
}

func New(st storage.Storage, logger logger.Logger) *StatCollector {
	return &StatCollector{
		st:     st,
		logger: logger,
	}
}

func (c *StatCollector) Grab(ctx context.Context) error {
	c.logger.Info("start collect data")
	times, err := getCPUTimes()
	if err != nil {
		return err
	}
	err = c.st.StoreCPUTimeStat(ctx, *times)
	if err != nil {
		return err
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	var sumUser, sumSystem, sumIdle float64
	lastCPUTimes, err := c.st.GetCPUTimeStats(ctx, period)
	if err != nil {
		return nil, err
	}
	for _, metric := range lastCPUTimes {
		stat := metric.StatInfo.(storage.CPUTimeStat)
		sumUser += stat.User
		sumSystem += stat.System
		sumIdle += stat.Idle
	}
	totalLen := len(lastCPUTimes)
	return NewCPUTimeStat(sumUser/float64(totalLen), sumSystem/float64(totalLen), sumIdle/float64(totalLen)), nil
}
