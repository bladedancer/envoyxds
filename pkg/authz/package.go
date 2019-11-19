package authz

import (
	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/sirupsen/logrus"
    "github.com/bladedancer/envoyxds/pkg/cache"
    "github.com/bladedancer/envoyxds/pkg/cache/redis"
    "context"
    "github.com/golang/protobuf/ptypes"
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

    c.Set(context.Background(), "APIKey", &ApiKeyMessage{Key:"This is the key for now"}, 0)
    //TODO Throw this away- just a apikey entry to validate Check rpc call is working
    pack,_:=ptypes.MarshalAny(&ApiKeyMessage{Key:"Gavin 1 API Key"})
    err=c.Set(context.Background(), "gavin-1-apikey", &AuthEnvelope{CtxType:ChangeType_API,Context:pack}, 0)
    if err!=nil {
        log.Warnf("Unable to save initial key %s", err)
    }
    ret:=&AuthEnvelope{}
    c.Get(context.Background(),"gavin-1-apikey", ret, true)
    log.Infof("Verified Key Entry %d", ret.CtxType)
}
