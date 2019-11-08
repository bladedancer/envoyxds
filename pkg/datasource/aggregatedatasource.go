package datasource

import (
	"reflect"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
)

// AggregateDatasource is a tenant datasource that combines other datasource.
type AggregateDatasource struct {
	sources    []Datasource
	updateChan chan []*apimgmt.Tenant
}

// MakeAggregateDatasource Initialize the datasource
func MakeAggregateDatasource(sources ...Datasource) *AggregateDatasource {
	updateChan := make(chan []*apimgmt.Tenant)

	go func() {
		cases := make([]reflect.SelectCase, len(sources))
		for i, source := range sources {
			_, sourceChan := source.GetTenants()
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sourceChan)}
		}

		for {
			_, value, ok := reflect.Select(cases)
			if !ok {
				break
			}
			updateChan <- value.Interface().([]*apimgmt.Tenant)
		}
	}()

	ds := &AggregateDatasource{
		updateChan: updateChan,
		sources:    sources,
	}

	return ds
}

// GetTenants Get the tenant data and a channel hat is updated on change
func (ds *AggregateDatasource) GetTenants() ([]*apimgmt.Tenant, chan []*apimgmt.Tenant) {
	t := []*apimgmt.Tenant{}

	for _, s := range ds.sources {
		tenants, _ := s.GetTenants()
		t = append(t, tenants...)
	}
	return t, ds.updateChan
}

// GetTenant finds the tenant with the specified name.
func (ds *AggregateDatasource) GetTenant(name string) *apimgmt.Tenant {
	var tenant *apimgmt.Tenant
	for _, s := range ds.sources {
		tenant = s.GetTenant(name)
		if tenant != nil {
			break
		}
	}
	return tenant
}

// UpsertTenant update existing tenant in the store or add to the first soruce if not already present.
func (ds *AggregateDatasource) UpsertTenant(tenant *apimgmt.Tenant) {
	// Could be more efficient
	found := false
	for _, s := range ds.sources {
		tenant = s.GetTenant(tenant.Name)
		if tenant != nil {
			found = true
			s.UpsertTenant(tenant)
			break
		}
	}

	if !found {
		ds.sources[0].UpsertTenant(tenant)
	}
}
