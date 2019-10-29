package envoyxds

import (
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/sirupsen/logrus"
)

var log logrus.FieldLogger = logrus.WithField("package", "envoyxds")
var config *xdsconfig.Config

// SetLog sets the logger for the package.
func SetLog(newLog logrus.FieldLogger) {
	log = newLog
	return
}

// SetConfig sets the config for the package.
func SetConfig(c *xdsconfig.Config) {
	config = c
	return
}
