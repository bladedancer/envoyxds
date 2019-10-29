package xdsconfig

import (
	"fmt"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

// GetRouteResources Get the route configuration data
func GetRouteResources(tenant *Tenant) []cache.Resource {
	var resources []cache.Resource

	for _, proxy := range tenant.Proxies {
		// Create the Routes
		config := makeRouteConfiguration(tenant.Name, proxy)
		resource := []cache.Resource{config}
		resources = append(resources, resource...)
	}

	return resources
}

// Create the envoy config for the tenant routes.
func makeRouteConfiguration(tenantName string, proxy *Proxy) *api.RouteConfiguration {
	vHosts := makeVHost(tenantName, "test", config.Domain, proxy)
	return &api.RouteConfiguration{
		Name:         "local_route",
		VirtualHosts: []*route.VirtualHost{vHosts},
	}
}

func makeVHost(tenantName string, env string, domain string, proxy *Proxy) *route.VirtualHost {
	id := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	vhost := route.VirtualHost{
		Name: id,
		Domains: []string{
			fmt.Sprintf("%s-%s.%s", env, tenantName, domain),
		},
		Routes: makeRoutes(tenantName, proxy),
	}
	return &vhost
}

func makeRoutes(tenantName string, proxy *Proxy) []*route.Route {
	id := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	r := &route.Route{
		Name: id,
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: proxy.Frontend.BasePath,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: id,
				},
				PrefixRewrite: proxy.Backend.Path,
				HostRewriteSpecifier: &route.RouteAction_HostRewrite{
					HostRewrite: proxy.Backend.Host,
				},
			},
		},
	}
	return []*route.Route{r}
}
