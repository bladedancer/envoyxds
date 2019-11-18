package xds

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/bladedancer/envoyxds/pkg/cache"
	"github.com/bladedancer/envoyxds/pkg/cache/redis"
	"github.com/sirupsen/logrus"
)

var log logrus.FieldLogger
var config *base.Config
var cacheCon cache.Cache

// Init Probably a nicer way to do this but for now it's good enough
func Init() {
	log = base.GetLog("xds")
	config = base.GetConfig()

	log.Infof("Connecting to Cache @ %s", config.CacheHost)
	cacheCon = redis.New([]string{config.CacheHost}, "", 0)
	err := cacheCon.Connect()
	if err != nil {
		log.Fatalf("Unable to connect: %s", err)
	}
}
