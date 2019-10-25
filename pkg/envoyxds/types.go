package envoyxds

// Config defines the configuration needed for Envoy XDS
type Config struct {
	Port       uint32
	Path       string
	NumTenants int
	NumRoutes  int
	Domain     string
	Pump       int64
}
