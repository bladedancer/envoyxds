package envoyxds

import "github.com/bladedancer/envoyxds/pkg/xdsconfig"

// GetTenantConfig Poll the data source to get teh tenant info.
func getTenants() []*xdsconfig.Tenant {
	return []*xdsconfig.Tenant{
		&xdsconfig.Tenant{
			Name: "test",
			Proxies: []*xdsconfig.Proxy{
				&xdsconfig.Proxy{
					Name: "google",
					Frontend: &xdsconfig.Frontend{
						BasePath: "/test",
					},
					Backend: &xdsconfig.Backend{
						Host: "www.google.com",
						Port: 443,
						Path: "/",
					},
				},
			},
		},
	}
}
