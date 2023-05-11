package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/app"
	versioncmd "github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	internalgrpc "github.com/BBeRsErKeRR/system-stats-monitor/internal/server/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Scheduler application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		config, err := NewConfig(cfgFile)
		if err != nil {
			log.Println("Error create config: " + err.Error())
			return
		}

		logg, err := logger.New(config.Logger)
		if err != nil {
			log.Println("Error create app logger: " + err.Error())
			return
		}

		st := app.GetStorage()
		err = st.Connect(ctx)
		if err != nil {
			logg.Error("Error create db connection: " + err.Error())
			return
		}
		stats := stats.New(config.App.StatsConfig, st, logg)
		monitor := app.New(logg, stats)
		grpc := internalgrpc.NewServer(logg, monitor, config.App.GRPCServer)

		go func() {
			if err := grpc.Start(ctx); err != nil {
				logg.Error("failed to start grpc server: " + err.Error())
				cancel()
			}
		}()

		go func() {
			cleanTicker := time.NewTicker(config.App.StatsConfig.CleanDuration)
			for {
				select {
				case <-cleanTicker.C:
					err := stats.Clean(ctx)
					if err != nil {
						logg.Error("failed to clear storage: " + err.Error())
						cancel()
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			if err := monitor.CollectData(ctx, config.App.StatsConfig.ScanDuration); err != nil {
				logg.Error("failed to collect metrics: " + err.Error())
				cancel()
			}
		}()

		defer st.Close(ctx)
		defer grpc.Stop()

		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/config.toml", "Configuration file path")
	rootCmd.AddCommand(versioncmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
