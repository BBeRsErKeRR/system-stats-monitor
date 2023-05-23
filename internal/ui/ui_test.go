package cliui

import (
	"fmt"
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestUpdateWidgets(t *testing.T) {
	mockLogger, err := logger.New(&logger.Config{})
	require.NoError(t, err)
	ui := &UI{
		stats: &stats.Stats{
			CPUInfo: storage.CPUTimeStat{
				User:   10.0,
				System: 5.0,
				Idle:   85.0,
			},
			LoadInfo: storage.LoadStat{
				Load1:  2.5,
				Load5:  3.0,
				Load15: 2.8,
			},
			NetworkStateInfo: storage.NetworkStatesStat{
				Counters: map[string]int32{
					"ESTABLISHED": 10,
					"CLOSE_WAIT":  5,
				},
			},
			NetworkStatisticsInfo: []storage.NetworkStatsItem{
				{
					Command:  "httpd",
					PID:      1234,
					User:     1001,
					Protocol: "tcp",
					Port:     80,
				},
			},
			DiskUsageInfo: []storage.UsageStatItem{
				{
					Path:                   "/",
					Fstype:                 "ext4",
					Used:                   102400,
					AvailablePercent:       70.0,
					InodesUsed:             2048,
					InodesAvailablePercent: 60.0,
				},
			},
			DiskIoInfo: []storage.DiskIoStatItem{
				{
					Device:   "sda",
					Tps:      10.0,
					KbReadS:  1024.0,
					KbWriteS: 512.0,
				},
			},
			ProtocolTalkersInfo: []storage.ProtocolTalkerItem{
				{
					Protocol:        "tcp",
					SendBytes:       1024.0,
					BytesPercentage: 70.0,
				},
			},
			BpsTalkersInfo: []storage.BpsItem{
				{
					Source:      "192.168.1.1",
					Destination: "192.168.1.2",
					Protocol:    "tcp",
					Bps:         1024.0,
					Numbers:     100.0,
				},
			},
		},
		logger: mockLogger,
	}

	ui.initWidgets()
	require.NoError(t, err)
	ui.updateWidgets()

	expectedCPU := "User: 10.00 | System: 5.00 | Idle: 85.00"
	require.Equal(t, expectedCPU, ui.cpuInfoWidget.Text,
		fmt.Sprintf("cpu info widget text should be '%s', got '%s'", expectedCPU, ui.cpuInfoWidget.Text),
	)

	expectedLa := "Load1: 2.50 | Load5: 3.00 | Load15: 2.80"
	require.Equal(t, expectedLa, ui.loadInfoWidget.Text,
		fmt.Sprintf("cpu info widget text should be '%s', got '%s'", expectedLa, ui.loadInfoWidget.Text),
	)

	require.Equal(t, 2,
		len(ui.networkStateInfoWidget.Rows),
		fmt.Sprintf("network state info widget rows should have 2 items, got %d rows", len(ui.networkStateInfoWidget.Rows)),
	)

	require.Equal(t, 1,
		len(ui.diskUsageInfoWidget.Rows),
		fmt.Sprintf("disk usage info widget row should have 1 items, got %d rows", len(ui.diskUsageInfoWidget.Rows)),
	)

	require.Equal(t, 1,
		len(ui.diskIoInfoWidget.Rows),
		fmt.Sprintf("disk io info widget row should have 1 items, got %d rows", len(ui.diskIoInfoWidget.Rows)),
	)

	require.Equal(t, 1,
		len(ui.protocolTalkersInfoWidget.Rows),
		fmt.Sprintf("protocol talkers widget row should have 1 items, got %d rows", len(ui.protocolTalkersInfoWidget.Rows)),
	)

	require.Equal(t, 1,
		len(ui.bpsTalkersInfoWidget.Rows),
		fmt.Sprintf("bps talkers info widget row should have 1 items, got %d rows", len(ui.bpsTalkersInfoWidget.Rows)),
	)
}
