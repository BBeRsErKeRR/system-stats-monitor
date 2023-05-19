package cliui

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

type UI struct {
	cpuInfoWidget               *widgets.Paragraph
	loadInfoWidget              *widgets.Paragraph
	networkStateInfoWidget      *widgets.List
	networkStatisticsInfoWidget *widgets.List
	diskUsageInfoWidget         *widgets.List
	diskIoInfoWidget            *widgets.List
	protocolTalkersInfoWidget   *widgets.List
	bpsTalkersInfoWidget        *widgets.List

	stats  *stats.Stats
	logger logger.Logger
}

func NewUI(stats *stats.Stats, logger logger.Logger) (*UI, error) {
	if err := termui.Init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize termui")
	}

	ui := &UI{
		stats:  stats,
		logger: logger,
	}

	ui.initWidgets()
	ui.updateWidgets()

	return ui, nil
}

func (u *UI) initWidgets() {
	u.cpuInfoWidget = widgets.NewParagraph()
	u.cpuInfoWidget.Title = "CPU Info"
	u.cpuInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.loadInfoWidget = widgets.NewParagraph()
	u.loadInfoWidget.Title = "Load Info"
	u.loadInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.networkStateInfoWidget = widgets.NewList()
	u.networkStateInfoWidget.Title = "Network State Info"
	u.networkStateInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.networkStatisticsInfoWidget = widgets.NewList()
	u.networkStatisticsInfoWidget.Title = "Network Statistics Info"
	u.networkStatisticsInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.diskUsageInfoWidget = widgets.NewList()
	u.diskUsageInfoWidget.Title = "Disk Usage Info"
	u.diskUsageInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.diskIoInfoWidget = widgets.NewList()
	u.diskIoInfoWidget.Title = "Disk IO Info"
	u.diskIoInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.protocolTalkersInfoWidget = widgets.NewList()
	u.protocolTalkersInfoWidget.Title = "Protocol Talkers Info"
	u.protocolTalkersInfoWidget.BorderStyle.Fg = termui.ColorGreen

	u.bpsTalkersInfoWidget = widgets.NewList()
	u.bpsTalkersInfoWidget.Title = "BPS Talkers Info"
	u.bpsTalkersInfoWidget.BorderStyle.Fg = termui.ColorGreen
}

func (u *UI) updateWidgets() {
	u.cpuInfoWidget.Text = fmt.Sprintf("User: %.2f | System: %.2f | Idle: %.2f",
		u.stats.CPUInfo.User,
		u.stats.CPUInfo.System,
		u.stats.CPUInfo.Idle,
	)
	u.loadInfoWidget.Text = fmt.Sprintf("Load1: %.2f | Load5: %.2f | Load15: %.2f",
		u.stats.LoadInfo.Load1,
		u.stats.LoadInfo.Load5,
		u.stats.LoadInfo.Load15,
	)

	u.networkStateInfoWidget.Rows = make([]string, 0, len(u.stats.NetworkStateInfo.Counters))
	for k, v := range u.stats.NetworkStateInfo.Counters {
		u.networkStateInfoWidget.Rows = append(u.networkStateInfoWidget.Rows, fmt.Sprintf("%s: %d", k, v))
	}

	u.networkStatisticsInfoWidget.Rows = make([]string, 0, len(u.stats.NetworkStatisticsInfo.Items))
	for _, item := range u.stats.NetworkStatisticsInfo.Items {
		data := fmt.Sprintf("%s: %v %v %v %v", item.Command, item.PID, item.User, item.Protocol, item.Port)
		u.networkStatisticsInfoWidget.Rows = append(u.networkStatisticsInfoWidget.Rows, data)
	}

	u.diskUsageInfoWidget.Rows = make([]string, 0, len(u.stats.DiskUsageInfo.Items))
	for idx, item := range u.stats.DiskUsageInfo.Items {
		data := fmt.Sprintf("[%v] %s -> %s : used(%vM %v%%) inode(%vM %v%%)",
			idx,
			item.Path,
			item.Fstype,
			item.Used,
			ConvertFloat(item.AvailablePercent),
			item.InodesUsed,
			ConvertFloat(item.InodesAvailablePercent),
		)
		u.diskUsageInfoWidget.Rows = append(u.diskUsageInfoWidget.Rows, data)
	}

	u.diskIoInfoWidget.Rows = make([]string, 0, len(u.stats.DiskIoInfo.Items))
	for idx, item := range u.stats.DiskIoInfo.Items {
		data := fmt.Sprintf("[%v] %s -> tps(%v) kB_read/s(%v) kB_wrtn/s(%v)",
			idx,
			item.Device,
			ConvertFloat(item.Tps),
			ConvertFloat(item.KbReadS),
			ConvertFloat(item.KbWriteS),
		)
		u.diskIoInfoWidget.Rows = append(u.diskIoInfoWidget.Rows, data)
	}

	u.protocolTalkersInfoWidget.Rows = make([]string, 0, len(u.stats.ProtocolTalkersInfo.Items))
	for _, item := range u.stats.ProtocolTalkersInfo.Items {
		data := fmt.Sprintf("%s: %v  %v%%",
			item.Protocol,
			ConvertFloat(item.SendBytes),
			ConvertFloat(item.BytesPercentage),
		)
		u.protocolTalkersInfoWidget.Rows = append(u.protocolTalkersInfoWidget.Rows, data)
	}

	u.bpsTalkersInfoWidget.Rows = make([]string, 0, len(u.stats.BpsTalkersInfo.Items))
	for idx, item := range u.stats.BpsTalkersInfo.Items {
		data := fmt.Sprintf("[%v] (%s) %s -> %s: %v  %v b/s",
			idx,
			item.Protocol,
			item.Source,
			item.Destination,
			ConvertFloat(item.Numbers),
			ConvertFloat(item.Bps),
		)
		u.bpsTalkersInfoWidget.Rows = append(u.bpsTalkersInfoWidget.Rows, data)
	}
}

func (u *UI) Run(ctx context.Context, ticker *time.Ticker) error {
	scrollTicker := time.NewTicker(1 * time.Second)
	defer termui.Close()
	// calculate the size of the terminal
	sizeX, sizeY, err := term.GetSize(0)
	if err != nil {
		return err
	}
	u.logger.Info(fmt.Sprintf("%v, %v", sizeX, sizeY))

	grid := termui.NewGrid()
	grid.SetRect(0, 0, sizeX, sizeY)
	grid.Set(
		termui.NewRow(
			0.5/7,
			termui.NewCol(1.0/2, u.cpuInfoWidget),
			termui.NewCol(1.0/2, u.loadInfoWidget),
		),
		termui.NewRow(
			1.5/7,
			termui.NewCol(1.0/2, u.networkStateInfoWidget),
			termui.NewCol(1.0/2, u.networkStatisticsInfoWidget),
		),
		termui.NewRow(
			2.0/7,
			termui.NewCol(1.0, u.diskUsageInfoWidget),
		),
		termui.NewRow(
			1.5/7,
			termui.NewCol(1.0, u.diskIoInfoWidget),
		),
		termui.NewRow(
			1.0/7,
			termui.NewCol(1.0/2, u.protocolTalkersInfoWidget),
			termui.NewCol(1.0/2, u.bpsTalkersInfoWidget),
		),
	)

	for {
		select {
		case <-scrollTicker.C:
			scrollList(u.networkStatisticsInfoWidget)
			scrollList(u.diskUsageInfoWidget)
			scrollList(u.diskIoInfoWidget)
			scrollList(u.protocolTalkersInfoWidget)
			scrollList(u.bpsTalkersInfoWidget)
		case e := <-termui.PollEvents():
			switch e.ID {
			case "q", "<C-c>":
				u.logger.Info("canceled")
				return nil
			case "<Resize>":
				termui.Clear()
				sizeX, sizeY, err := term.GetSize(0)
				if err != nil {
					return err
				}
				grid.SetRect(0, 0, sizeX, sizeY)
				u.updateWidgets()
				termui.Render(grid)
			}
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			u.logger.Info("render new data")
			u.updateWidgets()
			termui.Render(grid)
		}
	}
}

func (u *UI) UpdateStats(stats *stats.Stats) {
	u.stats = stats
}

func scrollList(list *widgets.List) {
	lenRows := len(list.Rows)
	if len(list.Rows) > 0 {
		if list.SelectedRow == lenRows-1 {
			list.ScrollTop()
		} else {
			list.ScrollHalfPageDown()
		}
	}
}

func ConvertFloat(item float64) string {
	return strconv.FormatFloat(item, 'f', 2, 64)
}
