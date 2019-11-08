package datasource

import (
	"github.com/bladedancer/envoyxds/pkg/apimgmt"
)

// TenantDatasource The configured tenant datasource
var TenantDatasource Datasource

// Datasource provides tenant information
type Datasource interface {
	GetTenants() ([]*apimgmt.Tenant, chan []*apimgmt.Tenant)
	GetTenant(name string) *apimgmt.Tenant
	UpsertTenant(tenant *apimgmt.Tenant)
}

func initDatasource() {
	TenantDatasource = MakeAggregateDatasource(
		MakeMemoryDatasource(),
		MakeFilesystemDatasource(),
	)
}
