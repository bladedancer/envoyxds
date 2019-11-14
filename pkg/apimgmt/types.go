package apimgmt

// Backend The backend service being proxied.
type Backend struct {
	Host          string
	Port          uint32
	Path          string
	TLS           bool
	Authorization *Authorization
}

// Frontend The proxy frontend details.
type Frontend struct {
	BasePath      string
	Authorization []Authorization
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

// Authorization API Authorization interface
type Authorization interface {
	Type() AuthorizationType
}

// AuthorizationType The Authorization type
type AuthorizationType string

const (
	// AuthorizationTypePassthrough No authorization
	AuthorizationTypePassthrough AuthorizationType = "passthrough"

	// AuthorizationTypeHTTP HTTP authorization - basic or bearer
	AuthorizationTypeHTTP AuthorizationType = "http"
	// AuthorizationTypeAPIKey API Key authorization
	AuthorizationTypeAPIKey AuthorizationType = "apiKey"
	// TODO: OAuth/JWT
)

/* ----------------------------------
    Passthrough Authorization
---------------------------------- */

// PassthroughAuthorization No authorization.
type PassthroughAuthorization struct {
}

// Type Implement the Authorization
func (a PassthroughAuthorization) Type() AuthorizationType {
	return AuthorizationTypePassthrough
}

/* ----------------------------------
    API Key Authorization
---------------------------------- */

// APIKeyAuthorization API Key authorization scheme.
type APIKeyAuthorization struct {
	Name     string
	Location string
}

// Type Implement the Authorization
func (a APIKeyAuthorization) Type() AuthorizationType {
	return AuthorizationTypeAPIKey
}

/* ----------------------------------
    HTTP Authorization
---------------------------------- */

// HTTPAuthorization Http Basic/Bearer authorization scheme.
type HTTPAuthorization struct {
	Scheme string
}

// Type Implement the Authorization
func (a HTTPAuthorization) Type() AuthorizationType {
	return AuthorizationTypeHTTP
}
