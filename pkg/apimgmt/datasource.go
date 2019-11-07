package apimgmt

import (
	"fmt"
	"time"
)

// GetTenants Get the tenant data and a channel hat is updated on change
func GetTenants() ([]*Tenant, chan []*Tenant) {
	updateChan := make(chan []*Tenant)
	tenants := generateTenants(config.NumTenants)

	// Hack in a special case where the cache keeps growing
	if config.Pump > 0 {
		go pump(updateChan)
	}

	return tenants, updateChan
}

func pump(updateChan chan []*Tenant) {
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

func generateTenants(count int) []*Tenant {
	var tenants []*Tenant
	for i := 0; i < count; i++ {
		tenants = append(tenants, generateTenant(i))
	}
	return tenants
}

func generateTenant(id int) *Tenant {
	var proxies []*Proxy
	tls := true
	var port uint32
	for i := 0; i < config.NumRoutes; i++ {
		port = 80
		if tls {
			port = 443
		}
		proxies = append(
			proxies,
			&Proxy{
				Name: fmt.Sprintf("%d-google-%d", id, i),
				Frontend: &Frontend{
					BasePath: fmt.Sprintf("/route-%d", i),
				},
				Backend: &Backend{
					Host: "www.google.com",
					Port: port,
					Path: "/",
					TLS:  tls,
				},
			},
		)
		tls = !tls
	}
	return &Tenant{
		Name:    fmt.Sprintf("%d", id),
		Proxies: proxies,
	}
}
