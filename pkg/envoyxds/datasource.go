package envoyxds

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
)

// GetTenantConfig Poll the data source to get teh tenant info.
func getTenants() ([]*xdsconfig.Tenant, chan []*xdsconfig.Tenant) {
	updateChan := make(chan []*xdsconfig.Tenant)
	tenants := generateTenants(config.NumTenants)

	// Hack in a special case where the cache keeps growing
	if config.Pump > 0 {
		go pump(updateChan)
	}

	return tenants, updateChan
}

func pump(updateChan chan []*xdsconfig.Tenant) {
	log.Infof("Pumping new route every %d seconds", config.Pump)
	tick := time.NewTicker(time.Duration(config.Pump) * time.Second)
	i := 1

	go func() {
		for {
			select {
			case <-tick.C:
				log.Infof("Pump %d", config.NumTenants+i)
				tenants := generateTenants(config.NumTenants + i)
				updateChan <- tenants
			}
			i++
		}
	}()
}

func generateTenants(count int) []*xdsconfig.Tenant {
	var tenants []*xdsconfig.Tenant
	for i := 0; i < count; i++ {
		tenants = append(tenants, generateTenant(i))
	}
	return tenants
}

func generateTenant(id int) *xdsconfig.Tenant {
	var proxies []*xdsconfig.Proxy

	for i := 0; i < config.NumRoutes; i++ {
		proxies = append(
			proxies,
			&xdsconfig.Proxy{
				Name: fmt.Sprintf("%d-google-%d", id, i),
				Frontend: &xdsconfig.Frontend{
					BasePath: fmt.Sprintf("/route-%d", i),
				},
				Backend: &xdsconfig.Backend{
					Host: "www.google.com",
					Port: 443,
					Path: "/",
				},
			},
		)
	}
	return &xdsconfig.Tenant{
		Name:    fmt.Sprintf("%d", id),
		Proxies: proxies,
	}
}
