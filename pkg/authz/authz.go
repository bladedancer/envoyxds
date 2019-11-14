package authz

import (
	"context"
	"fmt"
	"net"

	"github.com/bladedancer/envoyxds/pkg/base"
	"github.com/gogo/googleapis/google/rpc"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
)

// AuthorizationServer empty struct because this isn't a fancy example
type AuthorizationServer struct{}

// Check - Passthrough at the moment- check redis and add a header
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	out := &ApiKeyMessage{}
	c.Get(context.Background(), "APIKey", out, true)
	log.Infof("Passthrough with key %s", out.Key)

	log.Infof("Request Headers: %v", req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	log.Infof("Query String: %v", req.GetAttributes().GetRequest().GetHttp().GetQuery())
	log.Infof("Context Extensions: %+v", req.GetAttributes().GetContextExtensions())

	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.OK),
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
	}, nil

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
