package cmd

import (
	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/spf13/viper"
)

func syncConfigFromViper() envoyxds.Config {
	return envoyxds.Config{
		Port:       viper.GetUint32("port"),
		Path:       viper.GetString("path"),
		NumTenants: viper.GetInt("tenants"),
		NumRoutes:  viper.GetInt("routes"),
		Domain:     viper.GetString("domain"),
		Pump:       viper.GetInt64("pump"),
	}
}
