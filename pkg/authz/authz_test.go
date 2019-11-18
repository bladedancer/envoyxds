package authz_test


import (
    "testing"
    "github.com/bladedancer/envoyxds/pkg/authz"
    "github.com/golang/protobuf/ptypes"
)

func TestAuth(t *testing.T) {
    tStr:="This should be pass"
    basic:=authz.BasicAuthCtx{Pass: tStr}
 //   api:=authz.ApiKeyCtx{ApiKey:"This is the keu"}
 //   oauth:=authz.OAuthCtx{Oath:"Oauth for now"}
    c,_:=ptypes.MarshalAny(&basic)
    ctx:=authz.AuthEnvelope{CtxType:authz.ChangeType_BASIC, Context:c}
    basic.User="This should get removed"
    ptypes.UnmarshalAny(ctx.Context,&basic)
    if basic.User != "" {
        t.Errorf("Expected %s Got %s","", basic.User)
    }
    if basic.Pass != tStr {
        t.Errorf("Expected %s Got %s",tStr, basic.Pass)
    }


    }
