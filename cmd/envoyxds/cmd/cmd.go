package cmd

import (
	"strings"

	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const pollInterval = "pollInterval"

// RootCmd configures the command params of the csa
var RootCmd = &cobra.Command{
	Use:     "csa",
	Short:   "The config sync agent synchronizes configuration between Axway SaaS and the remote cluster.",
	Version: "0.0.1",
	RunE:    run,
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().Int("port", 11000, "The XDS GRPC port.")
	RootCmd.Flags().String("logLevel", "info", "log level")
	RootCmd.Flags().String("logFormat", "json", "line or json")

	bindOrPanic("port", RootCmd.Flags().Lookup("port"))
	bindOrPanic("log.level", RootCmd.Flags().Lookup("logLevel"))
	bindOrPanic("log.format", RootCmd.Flags().Lookup("logFormat"))
}

func initConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix("csa")
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
	config := syncConfigFromViper()
	return envoyxds.Run(config)
}
