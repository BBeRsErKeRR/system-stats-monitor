package main

import (
	"log"

	versioncmd "github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/cmd/daemon"
)

func init() {
	daemon.RootCmd.PersistentFlags().StringVar(
		&daemon.CfgFile,
		"config",
		"./configs/config.toml",
		"Configuration file path",
	)
	daemon.RootCmd.AddCommand(versioncmd.VersionCmd)
}

func main() {
	if err := daemon.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
