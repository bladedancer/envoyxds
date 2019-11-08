package cmd

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/spf13/viper"
)

func syncConfigFromViper() base.Config {
	return base.Config{
		Port:           viper.GetUint32("port"),
		Path:           viper.GetString("path"),
		CertPath:       viper.GetString("certPath"),
		NumTenants:     viper.GetInt("tenants"),
		NumRoutes:      viper.GetInt("routes"),
		Domain:         viper.GetString("domain"),
		Pump:           viper.GetInt64("pump"),
		DNSRefreshRate: viper.GetInt64("dnsRefreshRate"),
		RespectDNSTTL:  viper.GetBool("respectDNSTTL"),
		NumShards:      viper.GetInt("shards"),
		DatabaseURL:    viper.GetString("databaseUrl"),
		DatabasePoll:   viper.GetInt("databasePoll"),
	}
}

func setupConfig() {
	config := syncConfigFromViper()
	log.Printf("Config: %+v", config)
	base.SetConfig(&config)
}
