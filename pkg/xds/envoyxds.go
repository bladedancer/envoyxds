package xds

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run entry point for Envoy XDS command line.
func Run() error {
	callbacks := Calls{}

	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	server := xds.NewServer(snapshotCache, callbacks)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Port))
	if err != nil {
		log.Fatal(err)
	}

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	api.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	api.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	go watch(snapshotCache)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	log.Infof("Listening on %d", config.Port)
	sig := <-gracefulStop
	log.Debugf("Got signal: %s", sig)
	grpcServer.GracefulStop()
	log.Info("Shutdown")
	return nil
}
