package storage

import (
	"context"
	"time"
)

type Metric struct {
	Date     time.Time
	StatInfo interface{}
}

func CreateMetric(data interface{}, date time.Time) Metric {
	return Metric{
		Date:     date,
		StatInfo: data,
	}
}

type Storage interface {
	Clear(context.Context, time.Time) error
	Connect(context.Context) error
	Close(context.Context) error

	StoreCPUTimeStat(context.Context, CPUTimeStat) error
	GetCPUTimeStats(context.Context, int64) ([]Metric, error)

	StoreLoadStat(context.Context, LoadStat) error
	GetLoadStats(context.Context, int64) ([]Metric, error)

	StoreNetworkStatesStat(context.Context, NetworkStatesStat) error
	GetNetworkStatesStats(context.Context, int64) ([]Metric, error)

	StoreNetworkStats(context.Context, []interface{}) error
	GetNetworkStats(context.Context, int64) ([]Metric, error)

	StoreUsageStat(context.Context, UsageStatItem) error
	GetUsageStats(context.Context, int64) ([]Metric, error)

	StorDiskIoStats(context.Context, []interface{}) error
	GetDiskIoStats(context.Context, int64) ([]Metric, error)

	StoreProtocolTalkersStat(context.Context, ProtocolTalkerItem) error
	GetProtocolTalkersStats(context.Context, int64) ([]Metric, error)

	StoreBpsTalkersStat(context.Context, BpsItem) error
	GetBpsTalkersStats(context.Context, int64) ([]Metric, error)
}
