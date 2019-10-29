package envoyxds

import (
	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/golang/protobuf/ptypes"
)

// GetClusterConfigurations Get test Cluster configuration data
func GetClusterConfigurations() []cache.Resource {

	res := []cache.Resource{
		makeCluster("web_service", "google.com", "www.google.com"),
	}

	return res
}

func makeCluster(clusterName string, remoteHost string, sni string) *api.Cluster {
	/*
		 "clusters": [

		  {
		   "name": "web_service",
		   "type": "LOGICAL_DNS",
		   "connect_timeout": "0.250s",
		   "hosts": [
			{
			 "socket_address": {
			  "address": "google.com",
			  "port_value": 443
			 }
			}
		   ],
		   "tls_context": {
			"sni": "www.google.com"
		   },
		   "dns_lookup_family": "V4_ONLY"
		  }
		 ]
	*/

	log.Infof(">>>>>>>>>>>>>>>>>>> creating cluster " + clusterName)

	h := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: remoteHost,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: uint32(443),
			},
		},
	}}

	return &api.Cluster{
		Name: clusterName,

		ConnectTimeout: ptypes.DurationProto(250 * time.Millisecond),

		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		Hosts:                []*core.Address{h},
		TlsContext: &auth.UpstreamTlsContext{
			Sni: sni,
		},
	}
}
