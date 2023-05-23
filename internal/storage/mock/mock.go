package mockstorage

import (
	"context"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/mock"
)

type Storage struct {
	mock.Mock
}

func (ms *Storage) GetStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	args := ms.Called(ctx, period)
	return args.Get(0).([]storage.Metric), args.Error(1)
}

func (ms *Storage) StoreStats(ctx context.Context, data interface{}) error {
	args := ms.Called(ctx, data)
	return args.Error(0)
}

func (ms *Storage) Clear(ctx context.Context, date time.Time) error {
	args := ms.Called(ctx, date)
	return args.Error(0)
}

func (ms *Storage) Connect(ctx context.Context) error { //nolint:revive
	return nil
}

func (ms *Storage) Close(ctx context.Context) error { //nolint:revive
	return nil
}

func (ms *Storage) StoreCPUTimeStat(ctx context.Context, data storage.CPUTimeStat) error {
	return ms.StoreStats(ctx, data)
}

func (ms *Storage) GetCPUTimeStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreLoadStat(ctx context.Context, data storage.LoadStat) error {
	return ms.StoreStats(ctx, data)
}

func (ms *Storage) GetLoadStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreNetworkStatesStat(ctx context.Context, data storage.NetworkStatesStat) error {
	return ms.StoreStats(ctx, data)
}

func (ms *Storage) GetNetworkStatesStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreNetworkStats(ctx context.Context, data []storage.NetworkStatsItem) error {
	args := ms.Called(ctx, data)
	return args.Error(0)
}

func (ms *Storage) GetNetworkStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreUsageStats(ctx context.Context, data []storage.UsageStatItem) error {
	args := ms.Called(ctx, data)
	return args.Error(0)
}

func (ms *Storage) GetUsageStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StorDiskIoStats(ctx context.Context, data []storage.NetworkStatsItem) error {
	args := ms.Called(ctx, data)
	return args.Error(0)
}

func (ms *Storage) GetDiskIoStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreProtocolTalkersStat(ctx context.Context, data storage.ProtocolTalkerItem) error {
	return ms.StoreStats(ctx, data)
}

func (ms *Storage) GetProtocolTalkersStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func (ms *Storage) StoreBpsTalkersStat(ctx context.Context, data storage.BpsItem) error {
	return ms.StoreStats(ctx, data)
}

func (ms *Storage) GetBpsTalkersStats(ctx context.Context, period int64) ([]storage.Metric, error) {
	return ms.GetStats(ctx, period)
}

func New() *Storage {
	return &Storage{}
}
