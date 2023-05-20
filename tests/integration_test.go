package integrationtest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	v1grpc "github.com/BBeRsErKeRR/system-stats-monitor/api/v1/grpc"
	grpcclient "github.com/BBeRsErKeRR/system-stats-monitor/internal/client/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd/daemon"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type MockClient struct {
	*grpcclient.Client
	out chan stats.Stats
}

func (m *MockClient) Process(stat *stats.Stats) {
	fmt.Println("catched")
	m.out <- *stat
}

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTest Suite")
}

var _ = Describe("Daemon", Ordered, func() {
	var resChan chan stats.Stats
	var ctxClient context.Context
	var clientCancel context.CancelFunc

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer GinkgoRecover()
		daemon.RootCmd.PersistentFlags().StringVar(
			&daemon.CfgFile,
			"config",
			"./testdata/config.toml",
			"Configuration file path",
		)
		daemon.RootCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}
		err := daemon.RootCmd.ExecuteContext(ctx)
		require.NoError(GinkgoT(), err)
	}()

	time.Sleep(2 * time.Second)

	BeforeAll(func() {
		ctxClient, clientCancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		resChan = make(chan stats.Stats)

		conn, err := grpc.DialContext(ctxClient, "0.0.0.0:9081",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(GinkgoT(), err)

		client := v1grpc.NewSystemStatsMonitorServiceV1Client(conn)

		req := &v1grpc.StartMonitoringRequest{
			ResponseDuration: 4,
			WaitDuration:     3,
		}
		stream, err := client.StartMonitoring(ctxClient, req)
		if err != nil {
			if !(status.Code(err) == codes.Unavailable && strings.Contains(err.Error(), "error reading server preface")) {
				require.NoError(GinkgoT(), err)
			}
		}

		go func() {
			defer GinkgoRecover()
			defer close(resChan)
			defer conn.Close()
			payload, errRecv := stream.Recv()
			if !errors.Is(errRecv, io.EOF) {
				require.NoError(GinkgoT(), errRecv)
			}
			stat := v1grpc.ResolveResponse(payload)
			resChan <- *stat
		}()
	})

	AfterAll(func() {
		clientCancel()
	})

	Describe("Collect stats", Ordered, func() {
		var statsResponse stats.Stats
		BeforeAll(func() {
			waitTicker := time.NewTicker(7 * time.Second)
			select {
			case out := <-resChan:
				statsResponse = out
			case <-waitTicker.C:
				GinkgoT().Error("time out")
			}
		})

		It("CPUInfo", func() {
			cpuStat := statsResponse.CPUInfo
			require.NotEmpty(GinkgoT(), cpuStat)
			require.Greater(GinkgoT(), cpuStat.User, float64(0))
			require.Greater(GinkgoT(), cpuStat.Idle, float64(0))
			require.Greater(GinkgoT(), cpuStat.System, float64(0))
		})

		It("LoadInfo", func() {
			la := statsResponse.LoadInfo
			require.NotEmpty(GinkgoT(), la)
			require.Greater(GinkgoT(), la.Load1, float64(0))
			require.Greater(GinkgoT(), la.Load5, float64(0))
			require.Greater(GinkgoT(), la.Load15, float64(0))
		})

		It("NetworkStateInfo", func() {
			nst := statsResponse.NetworkStateInfo
			require.NotEmpty(GinkgoT(), nst)
			require.Greater(GinkgoT(), len(nst.Counters), 1)
		})

		It("NetworkStatisticsInfo", func() {
			ns := statsResponse.NetworkStatisticsInfo
			require.NotEmpty(GinkgoT(), ns)
			require.Greater(GinkgoT(), len(ns.Items), 1)
		})

		It("DiskUsageInfo", func() {
			du := statsResponse.DiskUsageInfo
			require.NotEmpty(GinkgoT(), du)
		})

		It("DiskIoInfo", func() {
			dio := statsResponse.DiskIoInfo
			require.NotEmpty(GinkgoT(), dio)
		})

		It("ProtocolTalkersInfo", func() {
			pt := statsResponse.ProtocolTalkersInfo
			require.NotEmpty(GinkgoT(), pt)
		})

		It("BpsTalkersInfo", func() {
			bps := statsResponse.BpsTalkersInfo
			require.NotEmpty(GinkgoT(), bps)
		})
	})
})
