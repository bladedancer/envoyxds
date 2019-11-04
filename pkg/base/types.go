package base

import (
	"github.com/sirupsen/logrus"
)

// Config defines the configuration needed for Envoy XDS
type Config struct {
	Port           uint32
	Path           string
	CertPath       string
	NumTenants     int
	NumRoutes      int
	Domain         string
	Pump           int64
	DNSRefreshRate int64
	RespectDNSTTL  bool
	NumShards      int
}

var config *Config
var logger *logrus.Logger

// SetLog sets the logger for the package.
func SetLog(newLog *logrus.Logger) {
	logger = newLog
	return
}

// GetLog gets the logger.
func GetLog(name string) logrus.FieldLogger {
	return logger.WithField("package", name)
}

// SetConfig sets the config for the package.
func SetConfig(c *Config) {
	config = c
	return
}

// GetConfig returns the current config.
func GetConfig() *Config {
	return config
}
