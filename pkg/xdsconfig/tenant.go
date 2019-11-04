package xdsconfig

import (
	"fmt"

	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

// Tenant representation of a tenant in envoy.
type Tenant struct {
	Name        string
	Proxies     []*Proxy
	Domain      string
	Environment string
}

// MakeTenant create a tenant config.
func MakeTenant(name string) *Tenant {
	return &Tenant{
		Name:        name,
		Domain:      config.Domain,
		Environment: "test",
		Proxies:     []*Proxy{},
	}
}

// GetVirtualHost get the virtual host configuration for this tenant.
func (t *Tenant) GetVirtualHost() *route.VirtualHost {
	var routes []*route.Route
	for _, proxy := range t.Proxies {
		routes = append(routes, proxy.Route)
	}

	id := fmt.Sprintf("t_%s", t.Name)
	vhost := route.VirtualHost{
		Name: id,
		Domains: []string{
			fmt.Sprintf("%s-%s.%s", t.Environment, t.Name, t.Domain),
		},
		Routes: routes,
	}
	return &vhost
}
