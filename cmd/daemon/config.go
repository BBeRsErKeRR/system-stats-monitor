package main

import (
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/config"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
	internalgrpc "github.com/BBeRsErKeRR/system-stats-monitor/internal/server/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/stats"
)

type Config struct {
	Logger *logger.Config `mapstructure:"logger"`
	App    *AppConf       `mapstructure:"app"`
}

type AppConf struct {
	GRPCServer  *internalgrpc.Config `mapstructure:"grpc_server"`
	StatsConfig *stats.Config        `mapstructure:"stats"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
