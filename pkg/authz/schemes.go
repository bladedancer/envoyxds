package authz

import (
    auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
    rpcstatus "google.golang.org/genproto/googleapis/rpc/status"

    core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

    "golang.org/x/net/context"
    "github.com/golang/protobuf/ptypes"
)

func checkAgainstScheme(req *auth.CheckRequest) *auth.CheckResponse{

    out:=&ApiKeyMessage{}
    c.Get(context.Background(), "APIKey", out, true)
    log.Infof("Passthrough with key %s", out.Key)


    log.Infof("Request Headers: %v", req.GetAttributes().GetRequest().GetHttp().GetHeaders())
    log.Infof("Query String: %v", req.GetAttributes().GetRequest().GetHttp().GetQuery())
    log.Infof("Context Extensions: %+v", req.GetAttributes().GetContextExtensions())

    extensions:=req.GetAttributes().GetContextExtensions()
    var ret *auth.CheckResponse
    switch extensions["auth_type"] {
    case "apiKey":
        //TODOO
    case "basic":
        //TODOO
    case "oauth":
        //TODOO
    case "opa":
        //TODOO
    case "custom":
        //TODOO
    default:
        ret = &auth.CheckResponse{
            Status: &rpcstatus.Status{
                Code: int32(0),
            },
            HttpResponse: &auth.CheckResponse_OkResponse{
                OkResponse: &auth.OkHttpResponse{
                    Headers: []*core.HeaderValueOption{
                        {
                            Header: &core.HeaderValue{
                                Key:   "x-custom-header-from-authz",
                                Value: out.Key,
                            },
                        },
                    },
                },
            },
        }
    }
    return ret
    }

// GetBasicContext - move and refactor
func GetBasicContext(ctx AuthEnvelope, b *BasicAuthCtx)  {
    switch ctx.CtxType {
    case ChangeType_BASIC:
        ptypes.UnmarshalAny(ctx.Context,b)
    }
}

