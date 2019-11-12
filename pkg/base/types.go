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
	DatabaseURL    string
	DatabasePoll   int
}

type CacheConfig struct {
    Port           uint32
    Path           string
}
var config *Config
var cachecfg *CacheConfig
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
// SetCacheConfig sets the config for the package.
func SetCacheConfig(c *CacheConfig) {
    cachecfg = c
    return
}

// GetCacheConfig returns the current config.
func GetCacheConfig() *CacheConfig {
    return cachecfg
}
