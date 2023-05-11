package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	grpcclient "github.com/BBeRsErKeRR/system-stats-monitor/internal/client/grpc"
	versioncmd "github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
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

		client := grpcclient.NewClient(logg, config.App.GRPCClient)
		err = client.Connect(ctx)
		if err != nil {
			logg.Error("Error create db connection: " + err.Error())
			return
		}
		defer client.Close()

		go func() {
			if err := client.StartMonitoring(ctx); err != nil {
				logg.Error("failed to start get data: " + err.Error())
				cancel()
			}
		}()

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
