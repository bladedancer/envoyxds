package xdsconfig

import (
	"fmt"
	"net"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

// GetEndpointResources Get endpoint configuration data
func GetEndpointResources(tenants []*Tenant) []cache.Resource {
	var resources []cache.Resource

	for _, tenant := range tenants {
		for _, proxy := range tenant.Proxies {
			// Create the endpoints
			config := makeEndpoint(tenant.Name, proxy)
			if config != nil {
				resource := []cache.Resource{config}
				resources = append(resources, resource...)
			}
		}
	}

	return resources
}

// Create the endpoint
func makeEndpoint(tenantName string, proxy *Proxy) *api.ClusterLoadAssignment {
	ips, err := net.LookupIP(proxy.Backend.Host)
	if err != nil {
		log.Errorf("Unable to resolve %s", proxy.Backend.Host)
	}

	address := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: ips[0].String(), // TODO
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: proxy.Backend.Port,
			},
			ResolverName: "envoy.ip",
		},
	}}
	clusterName := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	return &api.ClusterLoadAssignment{
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
	}
}
