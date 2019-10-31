package xdsconfig

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

// Config defines the configuration needed for Envoy XDS
type Config struct {
	Port       uint32
	Path       string
	CertPath   string
	NumTenants int
	NumRoutes  int
	Domain     string
	Pump       int64
}

// Backend The backend service being proxied.
type Backend struct {
	Host string
	Port uint32
	Path string
}

// Frontend The proxy frontend details.
type Frontend struct {
	BasePath string
}

// Proxy The virtualized service.
type Proxy struct {
	Name     string
	Frontend *Frontend
	Backend  *Backend
}

// Tenant The tenant.
type Tenant struct {
	Name    string
	Proxies []*Proxy
}

// XDS The xds resources
type XDS struct {
	LDS []cache.Resource
	CDS []cache.Resource
	RDS []cache.Resource
	EDS []cache.Resource
}

// MakeXDS Helper for creating xds config.
func MakeXDS(tenants []*Tenant) *XDS {
	xds := &XDS{
		LDS: GetListenerResources(),
		CDS: GetClusterResources(tenants),
		RDS: GetRouteResources(tenants),
		EDS: GetEndpointResources(tenants),
	}
	return xds
}
