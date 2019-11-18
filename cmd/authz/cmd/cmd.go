package cmd

import (
	"strings"

	"github.com/bladedancer/envoyxds/pkg/authz"
	"github.com/bladedancer/envoyxds/pkg/cache/redis"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// RootCmd configures the command params of the csa
var RootCmd = &cobra.Command{
	Use:     "authz",
	Short:   "The External Auth Filter for envoy.",
	Version: "0.0.1",
	RunE:    run,
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().Uint32("port", 20000, "The Authz GRPC port.")
	RootCmd.Flags().String("logLevel", "info", "log level")
	RootCmd.Flags().String("logFormat", "json", "line or json")
	RootCmd.Flags().String("path", "localhost:6379", "The path to the cache")

	bindOrPanic("port", RootCmd.Flags().Lookup("port"))
	bindOrPanic("path", RootCmd.Flags().Lookup("path"))
	bindOrPanic("log.level", RootCmd.Flags().Lookup("logLevel"))
	bindOrPanic("log.format", RootCmd.Flags().Lookup("logFormat"))
}

func initConfig() {
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix("authz")
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
	authz.Init()
	redis.Init()

	return authz.Run()
}
