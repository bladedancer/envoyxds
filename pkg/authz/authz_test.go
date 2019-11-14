package authz_test


import (
    "errors"
    "fmt"
    "sync"
    "time"

    rcli "github.com/go-redis/redis"
    "testing"
    "github.com/bladedancer/envoyxds/pkg/cache/redis"
    "context"
    "github.com/bladedancer/envoyxds/pkg/cache"
    "github.com/bladedancer/envoyxds/pkg/authz"
    "github.com/gogo/protobuf/proto"
)

func TestAuth(t *testing.T) {
    basic:=authz.BasicAuthCtx{Pass: "This should be pass"}
 //   api:=authz.ApiKeyCtx{ApiKey:"This is the keu"}
 //   oauth:=authz.OAuthCtx{Oath:"Oauth for now"}
    ctx:=authz.AuthEnvelope{CtxType:authz.ChangeType_BASIC, Context: &basic}
    authz.GetContext(ctx)


    }
