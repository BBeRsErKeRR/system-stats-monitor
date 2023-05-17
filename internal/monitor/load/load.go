package load

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

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
	la, err := getLoad()
	if err != nil {
		return err
	}
	err = c.st.StoreLoadStat(ctx, *la)
	if err != nil {
		return err
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	var sumLoad1, sumLoad5, sumLoad15 float64
	lastLoadStats, err := c.st.GetLoadStats(ctx, period)
	if err != nil {
		return nil, err
	}
	for _, metric := range lastLoadStats {
		stat := metric.StatInfo.(storage.LoadStat)
		sumLoad1 += stat.Load1
		sumLoad5 += stat.Load5
		sumLoad15 += stat.Load15
	}
	totalLen := len(lastLoadStats)
	return storage.LoadStat{
		Load1:  sumLoad1 / float64(totalLen),
		Load5:  sumLoad5 / float64(totalLen),
		Load15: sumLoad15 / float64(totalLen),
	}, nil
}
