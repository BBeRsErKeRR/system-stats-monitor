//go:build linux
// +build linux

package networkstatistics

import (
	"context"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestParseNetstatOut(t *testing.T) {
	output := "Active Internet connections (only servers)\n" +
		"Proto Recv-Q Send-Q Local Address           Foreign Address         State       User       Inode      PID/Program name\n" + //nolint:lll
		"tcp       0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      1000       15535      -\n" +
		"tcp       0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      1000       15535      1111/sshd\n" +
		"tcp6      0      0 :::80                   :::*                    LISTEN      1000       15535      2222/apache2\n" + //nolint:lll
		"tcp       0      0 0.0.0.0:22              0.0.0.0:*                           1000       15535      1111/sshd\n"

	expected := []storage.NetworkStatsItem{
		{
			User:     1000,
			Protocol: "tcp",
			Port:     22,
		},
		{
			Command:  "sshd",
			PID:      1111,
			User:     1000,
			Protocol: "tcp",
			Port:     22,
		},
		{
			Command:  "apache2",
			PID:      2222,
			User:     1000,
			Protocol: "tcp6",
			Port:     80,
		},
	}

	items, err := parseNetstatOut(output)
	require.NoError(t, err)
	require.Len(t, items, len(expected))

	for i := range items {
		require.Equal(t, expected[i], items[i])
	}
}

func TestGetNS(t *testing.T) {
	ctx := context.Background()

	items, err := getNS(ctx)
	require.NoError(t, err)

	// Make sure all items have valid PID and port numbers
	for _, item := range items {
		require.Greater(t, item.PID, int32(0))
		require.Greater(t, item.Port, int32(0))
	}
}
