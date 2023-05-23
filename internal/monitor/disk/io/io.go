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
	logger logger.Logger
}

func New(logger logger.Logger) *StatCollector {
	return &StatCollector{
		logger: logger,
	}
}

var commandRunner = collectDiskIo

func (c *StatCollector) Grab(ctx context.Context) (interface{}, error) {
	c.logger.Info("start collect data")
	ios, err := commandRunner(ctx)
	if err != nil {
		return nil, err
	}
	c.logger.Info("successful collect data")
	return ios, nil
}

func (s *StatCollector) CheckCall(ctx context.Context) error {
	return checkCall(ctx)
}
