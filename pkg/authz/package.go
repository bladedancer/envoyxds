package authz

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/sirupsen/logrus"
    "github.com/bladedancer/envoyxds/pkg/cache"
    "github.com/bladedancer/envoyxds/pkg/cache/redis"
    "context"
)

var log logrus.FieldLogger
var config *base.Config

var c cache.Cache


// Init Probably a nicer way to do this but for now it's good enough
func Init() {
	log = base.GetLog("authz")
	config = base.GetConfig()
	log.Infof("Connecting to Cache @ %s", base.GetCacheConfig().Path)
    c=redis.New([]string{base.GetCacheConfig().Path},"",0)
    err:=c.Connect()
    if err!=nil {
        log.Fatalf("Unable to connect: %s", err)
    }
    c.Set(context.Background(), "APIKey", &ApiKeyMessage{}, 0)
}
