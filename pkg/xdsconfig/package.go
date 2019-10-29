package xdsconfig

import (
	"github.com/sirupsen/logrus"
)

var log logrus.FieldLogger = logrus.WithField("package", "xdsconfig")
var config *Config

// SetLog sets the logger for the package.
func SetLog(newLog logrus.FieldLogger) {
	log = newLog
	return
}

// SetConfig sets the config for the package.
func SetConfig(c *Config) {
	config = c
	return
}
