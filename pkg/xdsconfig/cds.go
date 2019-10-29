package xdsconfig

import (
	"fmt"
	"time"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/golang/protobuf/ptypes"
)

// GetClusterResources Get test Cluster configuration data
func GetClusterResources(tenant *Tenant) []cache.Resource {
	var resources []cache.Resource

	for _, proxy := range tenant.Proxies {
		// Create the Routes
		config := makeCluster(tenant.Name, proxy)
		resource := []cache.Resource{config}
		resources = append(resources, resource...)
	}

	return resources
}

// Create a cluster for the proxy
func makeCluster(tenantName string, proxy *Proxy) *api.Cluster {
	log.Infof("creating cluster for proxy: %s", proxy.Name)
	address := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: proxy.Backend.Host,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: proxy.Backend.Port,
			},
		},
	}}
	return &api.Cluster{
		Name:                 fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name),
		ConnectTimeout:       ptypes.DurationProto(250 * time.Millisecond),
		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		Hosts:                []*core.Address{address},
		TlsContext: &auth.UpstreamTlsContext{
			Sni: proxy.Backend.Host,
		},
	}
}
