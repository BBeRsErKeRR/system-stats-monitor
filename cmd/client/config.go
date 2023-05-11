package main

import (
	grpcclient "github.com/BBeRsErKeRR/system-stats-monitor/internal/client/grpc"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/config"
	"github.com/BBeRsErKeRR/system-stats-monitor/internal/logger"
)

type Config struct {
	Logger *logger.Config `mapstructure:"logger"`
	App    *AppConf       `mapstructure:"app"`
}

type AppConf struct {
	GRPCClient *grpcclient.Config `mapstructure:"grpc_client"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
