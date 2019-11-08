package datasource

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
)

// MemoryDatasource is a tenant datasource that creates the tenants in memory.
type MemoryDatasource struct {
	tenants    []*apimgmt.Tenant
	updateChan chan []*apimgmt.Tenant
}

// MakeMemoryDatasource Initialize the memory datasource
func MakeMemoryDatasource() *MemoryDatasource {
	ds := &MemoryDatasource{
		updateChan: make(chan []*apimgmt.Tenant),
		tenants:    generateTenants(config.NumTenants),
	}

	if config.Pump > 0 {
		go ds.pump()
	}

	return ds
}

// GetTenants Get the tenant data and a channel hat is updated on change
func (ds *MemoryDatasource) GetTenants() ([]*apimgmt.Tenant, chan []*apimgmt.Tenant) {
	return ds.tenants, ds.updateChan
}

// GetTenant finds the tenant with the specified name.
func (ds *MemoryDatasource) GetTenant(name string) *apimgmt.Tenant {
	var tenant *apimgmt.Tenant
	for i, t := range ds.tenants {
		if t.Name == name {
			tenant = ds.tenants[i]
			break
		}
	}
	return tenant
}

// UpsertTenant update existing tenant in the store or add if not already present.
func (ds *MemoryDatasource) UpsertTenant(tenant *apimgmt.Tenant) {
	// Could be more efficient
	found := false
	for i, t := range ds.tenants {
		if t.Name == tenant.Name {
			found = true
			ds.tenants[i] = tenant
			break
		}
	}

	if !found {
		ds.tenants = append(ds.tenants, tenant)
	}
	ds.updateChan <- []*apimgmt.Tenant{tenant}
}

// pump Grows the number of tenants over time.
func (ds *MemoryDatasource) pump() {
	log.Infof("Pumping new route every %d seconds", config.Pump)
	tick := time.NewTicker(time.Duration(config.Pump) * time.Second)
	i := 1

	go func() {
		for {
			select {
			case <-tick.C:
				log.Infof("Pump %d", config.NumTenants+i)
				tenants := generateTenants(config.NumTenants + i)
				ds.updateChan <- tenants
			}
			i++
		}
	}()
}

func generateTenants(count int) []*apimgmt.Tenant {
	var tenants []*apimgmt.Tenant
	for i := 0; i < count; i++ {
		tenants = append(tenants, generateTenant(i))
	}
	return tenants
}

func generateTenant(id int) *apimgmt.Tenant {
	var proxies []*apimgmt.Proxy
	tls := true
	var port uint32
	for i := 0; i < config.NumRoutes; i++ {
		port = 80
		if tls {
			port = 443
		}
		proxies = append(
			proxies,
			&apimgmt.Proxy{
				Name: fmt.Sprintf("%d-google-%d", id, i),
				Frontend: &apimgmt.Frontend{
					BasePath: fmt.Sprintf("/route-%d", i),
				},
				Backend: &apimgmt.Backend{
					Host: "www.google.com",
					Port: port,
					Path: "/",
					TLS:  tls,
				},
			},
		)
		tls = !tls
	}
	return &apimgmt.Tenant{
		Name:    fmt.Sprintf("%d", id),
		Proxies: proxies,
	}
}
