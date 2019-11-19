package authz

import (
	"context"
	"fmt"
	"net"

	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/bladedancer/envoyxds/pkg/cache"
    code "google.golang.org/genproto/googleapis/rpc/code"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	_type "github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/golang/protobuf/ptypes"
)

// AuthorizationServer empty struct because this isn't a fancy example
type AuthorizationServer struct{}

// Check - Passthrough at the moment- check redis and add a header
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	var resp *auth.CheckResponse
	var headers map[string]string
	var err error
	authorized := checkAgainstScheme(req)

	if authorized {
		headers, err = backendAuth(req)
	}

	if err != nil {
		log.Error(err)
		resp = &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(code.Code_UNKNOWN),
			},
			HttpResponse: &auth.CheckResponse_DeniedResponse{
				DeniedResponse: &auth.DeniedHttpResponse{
					Status: &_type.HttpStatus{
						Code: _type.StatusCode_Forbidden,
					},
					Body: err.Error(),
				},
			},
		}

	} else if !authorized {
		log.Info("Not Authorized")
		resp = &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(code.Code_UNAUTHENTICATED),
			},
			HttpResponse: &auth.CheckResponse_DeniedResponse{
				DeniedResponse: &auth.DeniedHttpResponse{
					Status: &_type.HttpStatus{
						Code: _type.StatusCode_Forbidden,
					},
				},
			},
		}
	} else {
		// Happy days
		log.Infof("Adding headers to response: %+v", headers)
		respHeaders := []*core.HeaderValueOption{}

		if headers != nil {
			for hdrName, hdrVal := range headers {
				respHeaders = append(respHeaders, &core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   hdrName,
						Value: hdrVal,
					},
				})
			}
		}
		resp = &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(code.Code_OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: respHeaders,
				},
			},
		}
	}

	return resp, nil
}

//Run - Run the authz grpc service
func Run() error {
	// create a TCP listener
	Init()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", base.GetCacheConfig().Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	grpcServer := grpc.NewServer()
	authServer := &AuthorizationServer{}
	auth.RegisterAuthorizationServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	return nil
}
func getAPIContext(envelope *cache.AuthEnvelope) *cache.ApiKeyMessage {
	msg := cache.ApiKeyMessage{}
	ptypes.UnmarshalAny(envelope.Context, &msg)
	return &msg
}

func checkAgainstScheme(req *auth.CheckRequest) bool {
	authorized := false

	log.Infof("Request Headers: %v", req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	log.Infof("Query String: %v", req.GetAttributes().GetRequest().GetHttp().GetQuery())
	log.Infof("Context Extensions: %+v", req.GetAttributes().GetContextExtensions())

	extensions := req.GetAttributes().GetContextExtensions()
	log.Infof("Auth Type=%s", extensions["auth_type"])
	switch extensions["auth_type"] {
	case "apiKey":
		log.Infof("Calling API Authorize")
		authorized = NewAPIKey(req).Authorize()
	case "basic":
		//TODOO
	case "oauth":
		//TODOO
	case "opa":
		//TODOO
	case "custom":
		//TODOO
	default:

	}
	return authorized
}

// backendAuth generates the backend authorization headers.
func backendAuth(req *auth.CheckRequest) (map[string]string, error) {
	extensions := req.GetAttributes().GetContextExtensions()
	if authType, ok := extensions["be_type"]; ok {
		switch authType {
		case "apiKey":
			return backendAuthAPIKey(extensions)
		default:
			return nil, fmt.Errorf("no implementation for %s backend auth", authType)
		}
	}

	// No BE auth
	return nil, nil
}

func backendAuthAPIKey(extensions map[string]string) (map[string]string, error) {
	var tenant string
	var proxy string
	var name string
	var in string
	var ok bool

	if tenant, ok = extensions["tenant"]; !ok {
		return nil, fmt.Errorf("no tenant specified")
	}
	if proxy, ok = extensions["proxy"]; !ok {
		return nil, fmt.Errorf("no proxy specified")
	}

	if name, ok = extensions["be_name"]; !ok {
		return nil, fmt.Errorf("no name specified for %s-%s", tenant, proxy)
	}

	if in, ok = extensions["be_in"]; !ok {
		return nil, fmt.Errorf("no in specified for %s-%s", tenant, proxy)
	}

	// Get the credential
	creds := &cache.Credential{}
	if err := c.Get(context.Background(), fmt.Sprintf("%s-%s", tenant, proxy), creds, true); err != nil {
		return nil, err
	}
	if creds == nil || creds.Credential == nil {
		return nil, fmt.Errorf("no credentials available for %s-%s", tenant, proxy)
	}

	var cred string
	if cred, ok = creds.Credential[name]; !ok {
		return nil, fmt.Errorf("no credential with name %s specified for %s-%s", name, tenant, proxy)
	}

	// Insert the credential where it's supposed to go.
	if in != "header" {
		return nil, fmt.Errorf("have no clue how we're going to do credential injection in %s", in)
	}

	headers := make(map[string]string)
	headers[name] = cred
	return headers, nil
}
