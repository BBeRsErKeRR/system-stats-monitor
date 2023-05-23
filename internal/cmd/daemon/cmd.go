package daemon

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/app"
	daemonconfig "github.com/BBeRsErKeRR/system-stats-monitor/internal/config/daemon"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	internalgrpc "github.com/BBeRsErKeRR/system-stats-monitor/internal/server/grpc"
	"github.com/spf13/cobra"
)

var CfgFile string

var RootCmd = &cobra.Command{
	Short: "Scheduler application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		config, err := daemonconfig.NewConfig(CfgFile)
		if err != nil {
			log.Println("Error create config: " + err.Error())
			return
		}

		logg, err := logger.New(config.Logger)
		if err != nil {
			log.Println("Error create app logger: " + err.Error())
			return
		}

		monitor := app.New(logg, config.App, config.StatsConfig)
		gServer := internalgrpc.NewServer(logg, monitor, config.GRPCServer)

		go func() {
			if err := gServer.Start(ctx); err != nil {
				logg.Error("failed to start grpc server: " + err.Error())
				cancel()
			}
		}()

		defer gServer.GracefulStop()

		<-ctx.Done()
	},
}
