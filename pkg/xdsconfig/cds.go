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
	clusterName := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	return &api.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(250 * time.Millisecond),
		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_EDS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		DnsRefreshRate:       ptypes.DurationProto(180 * time.Second),
		RespectDnsTtl:        true,
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		EdsClusterConfig: &api.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
					ApiConfigSource: &core.ApiConfigSource{
						ApiType: core.ApiConfigSource_GRPC,
						GrpcServices: []*core.GrpcService{
							&core.GrpcService{
								TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
									EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
										ClusterName: "service_xds",
									},
								},
							},
						},
					},
				},
			},
		},
		// LoadAssignment: &api.ClusterLoadAssignment{
		// 	ClusterName: clusterName,
		// 	Endpoints: []*endpoint.LocalityLbEndpoints{
		// 		&endpoint.LocalityLbEndpoints{
		// 			LbEndpoints: []*endpoint.LbEndpoint{
		// 				&endpoint.LbEndpoint{
		// 					HostIdentifier: &endpoint.LbEndpoint_EndpointName{
		// 						EndpointName: clusterName,
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		TlsContext: &auth.UpstreamTlsContext{
			Sni: proxy.Backend.Host,
		},
	}
}
