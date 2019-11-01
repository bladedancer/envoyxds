package xdsconfig

import (
	"fmt"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

// GetFrontendRouteResources Get the frontend route configuration data
func GetFrontendRouteResources() []cache.Resource {
	// At the moment this uses clusterheader to decide and has no 
	// business logic. The listener should have a filter configured
	// to decide the correct shard and then route based on that.
	config := &api.RouteConfiguration{
		Name:         "local_route",
		VirtualHosts: []*route.VirtualHost{
			&route.VirtualHost{
				Name: "front",
				Domains: []string{"*"},
				Routes: []*route.Route{
					&route.Route{
						Name: "front",
						Match: &route.RouteMatch{
							PathSpecifier: &route.RouteMatch_Prefix{
								Prefix: "/",
							},
						},
						Action: &route.Route_Route{
							Route: &route.RouteAction{
								ClusterSpecifier: &route.RouteAction_ClusterHeader{
									ClusterHeader: "x-shard", // TODO
								},
							},
						},
					},
				},
			},
		},
	}
	resources := []cache.Resource{config}
	return resources
}

// GetRouteResources Get the route configuration data
func GetRouteResources(tenants []*Tenant) []cache.Resource {
	// Create the Routes
	config := makeRouteConfiguration(tenants)
	resources := []cache.Resource{config}
	return resources
}

// Create the envoy config for the tenant routes.
func makeRouteConfiguration(tenants []*Tenant) *api.RouteConfiguration {
	var vhosts []*route.VirtualHost

	for _, t := range tenants {
		vhosts = append(vhosts, makeVHost(t.Name, "test", config.Domain, t.Proxies))
	}
	return &api.RouteConfiguration{
		Name:         "local_route",
		VirtualHosts: vhosts,
	}
}

func makeVHost(tenantName string, env string, domain string, proxies []*Proxy) *route.VirtualHost {
	var routes []*route.Route
	for _, p := range proxies {
		routes = append(routes, makeRoutes(tenantName, p)...)
	}

	id := fmt.Sprintf("t_%s", tenantName)
	vhost := route.VirtualHost{
		Name: id,
		Domains: []string{
			fmt.Sprintf("%s-%s.%s", env, tenantName, domain),
		},
		Routes: routes,
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
