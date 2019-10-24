package envoyxds

import (
	"fmt"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

// GetRouteConfigurations Get test route configuration data
func GetRouteConfigurations(numTenants int, numRoutes int, domain string) *api.RouteConfiguration {
	var vHosts []*route.VirtualHost
	for i := 0; i < numTenants; i++ {
		vHosts = append(vHosts, getVirtualHost(i, numRoutes, domain))
	}

	return &api.RouteConfiguration{
		Name:         "local_route",
		VirtualHosts: vHosts,
	}
}

func getVirtualHost(id int, numRoutes int, domain string) *route.VirtualHost {
	var routes []*route.Route
	for i := 0; i < numRoutes; i++ {
		route := route.Route{
			Name: fmt.Sprintf("tenant-%d-route-%d", id, i),
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: fmt.Sprintf("/route-%d", i),
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{ // TODO
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: "web_service",
					},
					PrefixRewrite: "/",
					HostRewriteSpecifier: &route.RouteAction_HostRewrite{
						HostRewrite: "www.google.com",
					},
				},
			},
		}
		routes = append(routes, &route)
	}

	vhost := route.VirtualHost{
		Name: fmt.Sprintf("tenant-%d", id),
		Domains: []string{
			fmt.Sprintf("test-%d.%s", id, domain),
			fmt.Sprintf("prod-%d.%s", id, domain)},
		Routes: routes,
	}
	return &vhost
}
