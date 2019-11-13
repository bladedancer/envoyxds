package cmd

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/spf13/viper"
)

func syncConfigFromViper() base.CacheConfig {
	return base.CacheConfig{
		Port:           viper.GetUint32("port"),
		Path:           viper.GetString("path"),
	}
}

func setupConfig() {
	config := syncConfigFromViper()
	log.Printf("Config: %+v", config)
	base.SetCacheConfig(&config)
}
