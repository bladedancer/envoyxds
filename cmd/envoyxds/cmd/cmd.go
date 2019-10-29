package cmd

import (
	"strings"

	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// RootCmd configures the command params of the csa
var RootCmd = &cobra.Command{
	Use:     "xds",
	Short:   "The XDS configures envoy.",
	Version: "0.0.1",
	RunE:    run,
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().Uint32("port", 10000, "The XDS GRPC port.")
	RootCmd.Flags().String("path", "/", "The path for the config.")
	RootCmd.Flags().String("certPath", "/certs", "The path for the listener certs.")
	RootCmd.Flags().String("logLevel", "info", "log level")
	RootCmd.Flags().String("logFormat", "json", "line or json")
	RootCmd.Flags().Int("tenants", 10, "The number of tenants.")
	RootCmd.Flags().Int("routes", 5, "The number of routes per tenant.")
	RootCmd.Flags().String("domain", "bladedancer.dynu.net", "The domain for the routes.")
	RootCmd.Flags().Int64("pump", 0, "If set this adds a new route every N seconds.")

	bindOrPanic("port", RootCmd.Flags().Lookup("port"))
	bindOrPanic("path", RootCmd.Flags().Lookup("path"))
	bindOrPanic("certPath", RootCmd.Flags().Lookup("certPath"))
	bindOrPanic("log.level", RootCmd.Flags().Lookup("logLevel"))
	bindOrPanic("log.format", RootCmd.Flags().Lookup("logFormat"))
	bindOrPanic("tenants", RootCmd.Flags().Lookup("tenants"))
	bindOrPanic("routes", RootCmd.Flags().Lookup("routes"))
	bindOrPanic("domain", RootCmd.Flags().Lookup("domain"))
}

func initConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix("xds")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func bindOrPanic(key string, flag *flag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if err := setupLogging(viper.GetString("log.level"), viper.GetString("log.format")); err != nil {
		return err
	}
	setupConfig()
	return envoyxds.Run()
}
