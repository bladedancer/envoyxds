package cmd

import (
	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/spf13/viper"
)

func syncConfigFromViper() xdsconfig.Config {
	return xdsconfig.Config{
		Port:           viper.GetUint32("port"),
		Path:           viper.GetString("path"),
		CertPath:       viper.GetString("certPath"),
		NumTenants:     viper.GetInt("tenants"),
		NumRoutes:      viper.GetInt("routes"),
		Domain:         viper.GetString("domain"),
		Pump:           viper.GetInt64("pump"),
		DNSRefreshRate: viper.GetInt64("dnsRefreshRate"),
		RespectDNSTTL:  viper.GetBool("respectDNSTTL"),
	}
}

func setupConfig() {
	config := syncConfigFromViper()
	log.Printf("Config: %+v", config)
	envoyxds.SetConfig(&config)
	xdsconfig.SetConfig(&config)
}
