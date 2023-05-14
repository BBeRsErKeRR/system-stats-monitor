package diskio

import (
	"context"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
)

type DiskIoStatItem struct {
	Device   string  `json:"device"`
	Tps      float64 `json:"tps"`
	KbReadS  float64 `json:"kb_read_s"`
	KbWriteS float64 `json:"kb_write_s"`
}

type DiskIoStat struct {
	Items []DiskIoStatItem
}

func NewDiskIoStatItem(device string, tps, kb_read_s, kb_write_s float64) DiskIoStatItem {
	return DiskIoStatItem{
		Device:   device,
		Tps:      tps,
		KbReadS:  kb_read_s,
		KbWriteS: kb_write_s,
	}
}

type DiskIoStatCollector struct {
	name string
	st   storage.Storage
}

func New(st storage.Storage) *DiskIoStatCollector {
	return &DiskIoStatCollector{
		name: "io",
		st:   st,
	}
}

func (c *DiskIoStatCollector) Grab(ctx context.Context) error {
	times, err := collectDiskIo(ctx)
	if err != nil {
		return err
	}
	return c.st.BulkStoreStats(ctx, c.name, times)
}

func (as *DiskIoStatCollector) GetStats(ctx context.Context, period int64) (interface{}, error) {
	statsItems, err := as.st.GetStats(ctx, as.name, period)
	if err != nil {
		return nil, err
	}
	buff := make(map[string]*DiskIoStatItem)
	buffLen := make(map[string]float64)
	for _, metric := range statsItems {
		stat := metric.StatInfo.(DiskIoStatItem)
		val, ok := buff[stat.Device]
		if !ok {
			buff[stat.Device] = &stat
			buffLen[stat.Device] = 1
		} else {
			buffLen[stat.Device] += 1
			val.Tps += stat.Tps
			val.KbReadS += stat.KbReadS
			val.KbWriteS += stat.KbWriteS
		}
	}
	ioStats := make([]DiskIoStatItem, 0, len(buff))
	for key, val := range buff {
		val.Tps = val.Tps / buffLen[key]
		val.KbReadS = val.KbReadS / buffLen[key]
		val.KbWriteS = val.KbWriteS / buffLen[key]
		ioStats = append(ioStats, *val)
	}
	return DiskIoStat{
		Items: ioStats,
	}, nil
}
