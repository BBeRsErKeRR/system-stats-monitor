package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	v1grpc "github.com/BBeRsErKeRR/system-stats-monitor/api/v1/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	pkgnet "github.com/BBeRsErKeRR/system-stats-monitor/pkg/net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host             string        `mapstructure:"host"`
	Port             string        `mapstructure:"port"`
	ResponseDuration time.Duration `mapstructure:"response_duration"`
	WaitDuration     time.Duration `mapstructure:"wait_duration"`
}

type Client struct {
	logger           logger.Logger
	Addr             string
	conn             *grpc.ClientConn
	responseDuration time.Duration
	waitDuration     time.Duration
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
	}
}

func (c *Client) Connect(ctx context.Context) error {
	var err error
	c.conn, err = grpc.DialContext(ctx, c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	return err
}

func (c *Client) StartMonitoring(ctx context.Context) error {
	client := v1grpc.NewSystemStatsMonitorServiceV1Client(c.conn)
	req := &v1grpc.StartMonitoringRequest{
		ResponseDuration: int64(c.responseDuration / time.Second),
		WaitDuration:     int64(c.waitDuration / time.Second),
	}
	stream, err := client.StartMonitoring(context.Background(), req)
	if err != nil {
		return err
	}
	for {
		c.logger.Info("wait new data")
		data, errRecv := stream.Recv()
		if errors.Is(errRecv, io.EOF) {
			c.logger.Info("statistics collection completed")
			break
		} else if errRecv != nil {
			return errRecv
		}

		c.printStats(data)
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) printStats(data *v1grpc.StatsResponse) {
	fmt.Println("\nCPU:")
	fmt.Println("  user mode time:", data.GetCpuInfo().GetUser())
	fmt.Println("  system mode time:", data.GetCpuInfo().GetSystem())
	fmt.Println("  idle time:", data.GetCpuInfo().GetIdle())
	fmt.Println()
}
