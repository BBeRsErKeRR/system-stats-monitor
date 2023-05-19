package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	v1grpc "github.com/BBeRsErKeRR/system-stats-monitor/api/v1/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	cliui "github.com/BBeRsErKeRR/system-stats-monitor/internal/ui"
	pkgnet "github.com/BBeRsErKeRR/system-stats-monitor/pkg/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host             string        `mapstructure:"host"`
	Port             string        `mapstructure:"port"`
	ResponseDuration time.Duration `mapstructure:"response_duration"`
	WaitDuration     time.Duration `mapstructure:"wait_duration"`
	IsTermUIEnable   bool          `mapstructure:"termui_enable"`
}

type Client struct {
	logger           logger.Logger
	Addr             string
	conn             *grpc.ClientConn
	responseDuration time.Duration
	waitDuration     time.Duration
	isTermUIEnable   bool
}

func NewClient(logger logger.Logger, conf *Config) *Client {
	addr, err := pkgnet.GetAddress(conf.Host, conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		logger:           logger,
		Addr:             addr,
		responseDuration: conf.ResponseDuration,
		waitDuration:     conf.WaitDuration,
		isTermUIEnable:   conf.IsTermUIEnable,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	var err error
	c.conn, err = grpc.DialContext(ctx, c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	return err
}

func (c *Client) StartMonitoring(ctx context.Context, cancelFunc context.CancelFunc) error {
	responseTicker := time.NewTicker(c.responseDuration)
	client := v1grpc.NewSystemStatsMonitorServiceV1Client(c.conn)
	req := &v1grpc.StartMonitoringRequest{
		ResponseDuration: int64(c.responseDuration / time.Second),
		WaitDuration:     int64(c.waitDuration / time.Second),
	}
	stream, err := client.StartMonitoring(context.Background(), req)
	if err != nil {
		return err
	}
	stats := &stats.Stats{}

	var ui *cliui.UI
	if c.isTermUIEnable {
		ui, err = cliui.NewUI(stats, c.logger)
		if err != nil {
			return err
		}

		go func() {
			defer cancelFunc()
			if err := ui.Run(ctx, responseTicker); err != nil {
				c.logger.Error("failed to start ui: " + err.Error())
			}
		}()
	}

	for {
		c.logger.Info("wait new data")
		payload, errRecv := stream.Recv()
		if errors.Is(errRecv, io.EOF) {
			c.logger.Info("statistics collection completed")
			break
		} else if errRecv != nil {
			return errRecv
		}
		stat := v1grpc.ResolveResponse(payload)
		if c.isTermUIEnable {
			ui.UpdateStats(stat)
		} else {
			c.printStats(stat)
		}
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) printStats(data *stats.Stats) {
	fmt.Println("CPU:")
	fmt.Println("  'user mode time':", cliui.ConvertFloat(data.CPUInfo.User))
	fmt.Println("  'system mode time':", cliui.ConvertFloat(data.CPUInfo.System))
	fmt.Println("  'idle time':", cliui.ConvertFloat(data.CPUInfo.Idle))
	fmt.Println("\nLA:")
	fmt.Println("  '1 minute':", cliui.ConvertFloat(data.LoadInfo.Load1))
	fmt.Println("  '5 minutes':", cliui.ConvertFloat(data.LoadInfo.Load5))
	fmt.Println("  '15 minutes':", cliui.ConvertFloat(data.LoadInfo.Load15))
	fmt.Println("\nNetwork:")
	fmt.Println("  States:")
	for key, value := range data.NetworkStateInfo.Counters {
		fmt.Printf("    %s: %v\n", key, value)
	}
	fmt.Println("  Listen items:")
	for _, item := range data.NetworkStatisticsInfo.Items {
		fmt.Printf("    %s: '%v %v %v %v'\n", item.Command, item.PID, item.User, item.Protocol, item.Port)
	}
	fmt.Println("\nDisk:")
	fmt.Println("  Usage:")
	for _, item := range data.DiskUsageInfo.Items {
		fmt.Printf("    - '%s -> %s : used(%vM %v%%) inode(%vM %v%%)'\n",
			item.Path,
			item.Fstype,
			item.Used,
			cliui.ConvertFloat(item.AvailablePercent),
			item.InodesUsed,
			cliui.ConvertFloat(item.InodesAvailablePercent),
		)
	}
	fmt.Println("  IO:")
	for _, item := range data.DiskIoInfo.Items {
		fmt.Printf("    - '%s -> tps(%v) kB_read/s(%v) kB_wrtn/s(%v)'\n",
			item.Device,
			cliui.ConvertFloat(item.Tps),
			cliui.ConvertFloat(item.KbReadS),
			cliui.ConvertFloat(item.KbWriteS),
		)
	}
	fmt.Println("\nTalkers:")
	fmt.Println("  Protocol:")
	for _, item := range data.ProtocolTalkersInfo.Items {
		fmt.Printf("    %s: '%v  %v%%'\n",
			item.Protocol,
			cliui.ConvertFloat(item.SendBytes),
			cliui.ConvertFloat(item.BytesPercentage),
		)
	}
	fmt.Println("  Protocol:")
	for _, item := range data.BpsTalkersInfo.Items {
		fmt.Printf("    - '(%s) %s -> %s: %v  %v b/s'\n",
			item.Protocol,
			item.Source,
			item.Destination,
			cliui.ConvertFloat(item.Numbers),
			cliui.ConvertFloat(item.Bps),
		)
	}
	fmt.Println(strings.Repeat("#", 100))
}
