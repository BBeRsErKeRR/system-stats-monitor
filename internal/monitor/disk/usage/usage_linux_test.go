//go:build linux
// +build linux

package diskusage

import (
	"context"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestParseDfOut(t *testing.T) {
	cmdK := `Filesystem     1K-blocks    Used Available Use% Mounted on
/dev/sda1       10321208 3414928   6374172  35% /
devtmpfs          499660       0    499660   0% /dev
tmpfs             506668       0    506668   0% /dev/shm
C:\Program Files\Docker\Docker\resources 306064700 239919728  66144972  79% /Docker/host`

	cmdI := `Filesystem     Inodes IUsed IFree IUse% Mounted on
/dev/sda1        655360 13675 641685    3% /
devtmpfs         124915   583 124332    1% /dev
tmpfs            127852     1 127851    1% /dev/shm
C:\Program Files\Docker\Docker\resources      999 -999001  1000000     - /Docker/host`

	expected := []storage.UsageStatItem{
		{
			Path: "/", Fstype: "/dev/sda1", Used: 3334, AvailablePercent: 65.11499524981868,
			InodesUsed: 13, InodesAvailablePercent: 97.91336059570312,
		},
		{
			Path: "/dev", Fstype: "devtmpfs", Used: 0, AvailablePercent: 100,
			InodesUsed: 0, InodesAvailablePercent: 99.5332826321899,
		},
		{
			Path: "/dev/shm", Fstype: "tmpfs", Used: 0, AvailablePercent: 100,
			InodesUsed: 0, InodesAvailablePercent: 99.999217845634,
		},
		{
			Path: "/Docker/host", Fstype: "C:\\Program Files\\Docker\\Docker\\resources",
			Used: 234296, AvailablePercent: 21.611434445069946,
			InodesUsed: 0, InodesAvailablePercent: 100,
		},
	}

	result, err := parseDfOut(cmdK, cmdI)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestGetDU(t *testing.T) {
	ctx := context.Background()
	result, err := getDU(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	for _, item := range result {
		require.NotEmpty(t, item.Path, item)
		require.NotEmpty(t, item.Fstype)
		require.GreaterOrEqual(t, item.Used, int64(0))
		require.GreaterOrEqual(t, item.AvailablePercent, float64(0))
		require.LessOrEqual(t, item.AvailablePercent, float64(100))
		require.GreaterOrEqual(t, item.InodesUsed, int64(0))
		require.GreaterOrEqual(t, int64(item.InodesAvailablePercent), int64(0))
		require.LessOrEqual(t, int64(item.InodesAvailablePercent), int64(100))
	}
}
