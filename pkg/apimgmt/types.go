package apimgmt

// Backend The backend service being proxied.
type Backend struct {
	Host string
	Port uint32
	Path string
	TLS  bool
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
