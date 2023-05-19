package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage_All(t *testing.T) {
	ctx := context.Background()
	st := New()
	now := time.Now()

	// Add some metrics to the MemStores.
	err := st.StoreCPUTimeStat(ctx, storage.CPUTimeStat{User: 1, System: 2})
	require.NoError(t, err)
	err = st.StoreLoadStat(ctx, storage.LoadStat{Load1: 0.15, Load5: 0.2})
	require.NoError(t, err)
	err = st.StoreNetworkStatesStat(ctx, storage.NetworkStatesStat{})
	require.NoError(t, err)
	err = st.StoreUsageStat(ctx, storage.UsageStatItem{Path: "/tmp", Used: 1024})
	require.NoError(t, err)
	err = st.StorDiskIoStats(ctx, []interface{}{
		storage.DiskIoStatItem{Device: "sda", KbReadS: 1024},
	})
	require.NoError(t, err)
	err = st.StoreProtocolTalkersStat(ctx, storage.ProtocolTalkerItem{Protocol: "tcp", SendBytes: 100})
	require.NoError(t, err)
	err = st.StoreBpsTalkersStat(ctx, storage.BpsItem{Protocol: "tcp", Bps: 1024})
	require.NoError(t, err)

	// Clear the metrics before now.
	err = st.Clear(ctx, now.Add(time.Minute))
	require.NoError(t, err)

	// Check that there are no metrics with a date before now.
	metrics, err := st.GetCPUTimeStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetLoadStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetNetworkStatesStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetUsageStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetDiskIoStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetProtocolTalkersStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)

	metrics, err = st.GetBpsTalkersStats(ctx, 60)
	require.NoError(t, err)
	require.Empty(t, metrics)
}

func TestStorage_ConnectAndClose(t *testing.T) {
	ctx := context.Background()
	st := New()

	err := st.Connect(ctx)
	require.NoError(t, err)

	err = st.Close(ctx)
	require.NoError(t, err)
}
