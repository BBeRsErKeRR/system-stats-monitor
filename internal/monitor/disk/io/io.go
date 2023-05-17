package diskio

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

func NewDiskIoStatItem(device string, tps, kbReadS, kbWriteS float64) storage.DiskIoStatItem {
	return storage.DiskIoStatItem{
		Device:   device,
		Tps:      tps,
		KbReadS:  kbReadS,
		KbWriteS: kbWriteS,
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
	ios, err := collectDiskIo(ctx)
	if err != nil {
		return err
	}
	err = c.st.StorDiskIoStats(ctx, ios)
	if err != nil {
		return err
	}
	c.logger.Info("successful collect data")
	return nil
}

func (c *StatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	statsItems, err := c.st.GetDiskIoStats(ctx, period)
	if err != nil {
		return nil, err
	}
	buff := make(map[string]*storage.DiskIoStatItem)
	buffLen := make(map[string]float64)
	for _, metric := range statsItems {
		stat := metric.StatInfo.(storage.DiskIoStatItem)
		val, ok := buff[stat.Device]
		if !ok {
			buff[stat.Device] = &stat
			buffLen[stat.Device] = 1
		} else {
			buffLen[stat.Device]++
			val.Tps += stat.Tps
			val.KbReadS += stat.KbReadS
			val.KbWriteS += stat.KbWriteS
		}
	}
	ioStats := make([]storage.DiskIoStatItem, 0, len(buff))
	for key, val := range buff {
		val.Tps /= buffLen[key]
		val.KbReadS /= buffLen[key]
		val.KbWriteS /= buffLen[key]
		ioStats = append(ioStats, *val)
	}
	return storage.DiskIoStat{
		Items: ioStats,
	}, nil
}
