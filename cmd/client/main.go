package main

import (
	"log"

	versioncmd "github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd/client"
)

func init() {
	client.RootCmd.PersistentFlags().StringVar(
		&client.CfgFile,
		"config",
		"./configs/config.toml",
		"Configuration file path",
	)
	client.RootCmd.AddCommand(versioncmd.VersionCmd)
}

func main() {
	if err := client.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
