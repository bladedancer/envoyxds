package xdsconfig

import (
	"fmt"
	"time"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/golang/protobuf/ptypes"
)

// GetFrontendClusterResources Get Cluster configuration data for the frontend
func GetFrontendClusterResources() []cache.Resource {
	var resources []cache.Resource

	for i := 0; i < config.NumShards; i++ {
		config := makeRoutingCluster(fmt.Sprintf("back-%d", i)) // embedding "back" not ideal
		resource := []cache.Resource{config}
		resources = append(resources, resource...)
	}

	return resources
}

// GetClusterResources Get test Cluster configuration data
func GetClusterResources(tenants []*Tenant) []cache.Resource {
	var resources []cache.Resource

	for _, tenant := range tenants {
		for _, proxy := range tenant.Proxies {
			// Create the Routes
			config := makeCluster(tenant.Name, proxy)
			resource := []cache.Resource{config}
			resources = append(resources, resource...)
		}
	}

	return resources
}

// Create a cluster for the proxy
func makeCluster(tenantName string, proxy *Proxy) *api.Cluster {
	//log.Infof("creating cluster for proxy: %s", proxy.Name)
	address := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: proxy.Backend.Host,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: proxy.Backend.Port,
			},
		},
	}}

	clusterName := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	return &api.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(250 * time.Millisecond),
		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		RespectDnsTtl:        config.RespectDNSTTL,
		DnsRefreshRate:       ptypes.DurationProto(time.Duration(config.DNSRefreshRate) * time.Millisecond),
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		LoadAssignment: &api.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*endpoint.LocalityLbEndpoints{
				&endpoint.LocalityLbEndpoints{
					LbEndpoints: []*endpoint.LbEndpoint{
						&endpoint.LbEndpoint{
							HostIdentifier: &endpoint.LbEndpoint_Endpoint{
								Endpoint: &endpoint.Endpoint{
									Address: address,
								},
							},
						},
					},
				},
			},
		},
		TlsContext: &auth.UpstreamTlsContext{
			Sni: proxy.Backend.Host,
		},
	}
}

// Create a cluster for the proxy
func makeRoutingCluster(shardName string) *api.Cluster {
	//log.Infof("creating cluster for proxy: %s", proxy.Name)
	address := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: shardName, // TODO shardName === pod name === clustername, not ideal
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: 443, // TODO
			},
		},
	}}

	return &api.Cluster{
		Name:                 shardName,
		ConnectTimeout:       ptypes.DurationProto(250 * time.Millisecond),
		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		RespectDnsTtl:        config.RespectDNSTTL,
		DnsRefreshRate:       ptypes.DurationProto(time.Duration(config.DNSRefreshRate) * time.Millisecond),
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		LoadAssignment: &api.ClusterLoadAssignment{
			ClusterName: shardName,
			Endpoints: []*endpoint.LocalityLbEndpoints{
				&endpoint.LocalityLbEndpoints{
					LbEndpoints: []*endpoint.LbEndpoint{
						&endpoint.LbEndpoint{
							HostIdentifier: &endpoint.LbEndpoint_Endpoint{
								Endpoint: &endpoint.Endpoint{
									Address: address,
								},
							},
						},
					},
				},
			},
		},
		TlsContext: &auth.UpstreamTlsContext{}, // Probably doesn't need to be ssl
	}
}
