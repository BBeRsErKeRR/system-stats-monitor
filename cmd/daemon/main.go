package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/app"
	versioncmd "github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	internalgrpc "github.com/BBeRsErKeRR/system-stats-monitor/internal/server/grpc"
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

		monitor := app.New(logg, config.App, config.StatsConfig)
		grpc := internalgrpc.NewServer(logg, monitor, config.GRPCServer)

		go func() {
			if err := grpc.Start(ctx); err != nil {
				logg.Error("failed to start grpc server: " + err.Error())
				cancel()
			}
		}()

		// defer st.Close(ctx)
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
