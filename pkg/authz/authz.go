package authz

import (
    "context"
    "net"
    "google.golang.org/grpc"
    "fmt"
    "github.com/bladedancer/envoyxds/pkg/base"

    auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
    "github.com/golang/protobuf/ptypes"
)

// AuthorizationServer empty struct because this isn't a fancy example
type AuthorizationServer struct{}

// Check - Passthrough at the moment- check redis and add a header
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
    return checkAgainstScheme(req), nil
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
