package xds

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/sirupsen/logrus"
)

var log logrus.FieldLogger
var config *base.Config

// Init Probably a nicer way to do this but for now it's good enough
func Init() {
	log = base.GetLog("xds")
	config = base.GetConfig()
}
