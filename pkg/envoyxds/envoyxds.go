package envoyxds

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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log logrus.FieldLogger = logrus.WithField("package", "sync")

// SetLog sets the logger for the package.
func SetLog(newLog logrus.FieldLogger) {
	log = newLog
	return
}

// Run entry point for Envoy XDS command line.
func Run(conf Config) error {
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	server := xds.NewServer(snapshotCache, nil)
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Error(err)
	}

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	api.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	api.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	go watch(snapshotCache)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	log.Infof("Listening on %d", conf.Port)
	sig := <-gracefulStop
	log.Debugf("Got signal: %s", sig)
	grpcServer.GracefulStop()
	log.Info("Shutdown")
	return nil
}
