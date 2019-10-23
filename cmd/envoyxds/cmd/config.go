package cmd

import (
	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/spf13/viper"
)

func syncConfigFromViper() envoyxds.Config {
	return envoyxds.Config{
		Port: viper.GetInt("port"),
		Path: viper.GetString("path"),
	}
}
