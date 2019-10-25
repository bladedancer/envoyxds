package envoyxds

// Config defines the configuration needed for Envoy XDS
type Config struct {
	Port       int
	Path       string
	NumTenants int
	NumRoutes  int
	Domain     string
	Pump       int64
}
