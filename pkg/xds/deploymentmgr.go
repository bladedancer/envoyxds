package xds

import (
	"fmt"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
)

// Deployment deployment of a tenant on a shard
type Deployment struct {
	shardName string
	apiTenant *apimgmt.Tenant
	xdsTenant *xdsconfig.Tenant
}

// DeploymentManager maps tenants to clusters
type DeploymentManager struct {
	byTenantName map[string]*Deployment
	nextShard    int
	shards       map[string]*xdsconfig.BackendShard
	OnChange     chan []*xdsconfig.BackendShard
}

// MakeDeploymentManager create a DeplyomentManager.
func MakeDeploymentManager() *DeploymentManager {
	return &DeploymentManager{
		byTenantName: make(map[string]*Deployment),
		nextShard:    0,
		shards:       make(map[string]*xdsconfig.BackendShard),
		OnChange:     make(chan []*xdsconfig.BackendShard),
	}
}

// GetShardName gets the shard that the tenant is deployed on
func (dm *DeploymentManager) GetShardName(tenantName string) string {
    dep:=dm.byTenantName[tenantName]
    shard:=""
    if dep!= nil{
        shard=dep.shardName
    }
    return shard
}

// GetShard return the named shard.
func (dm *DeploymentManager) GetShard(name string) *xdsconfig.BackendShard {
	return dm.shards[name]
}

// GetShards return the shards for the deployed tenants.
func (dm *DeploymentManager) GetShards() []*xdsconfig.BackendShard {
	backendShards := make([]*xdsconfig.BackendShard, 0, len(dm.shards))
	for _, value := range dm.shards {
		backendShards = append(backendShards, value)
	}
	return backendShards
}

// updateShard updates the backend configuration for the specified shard.
func (dm *DeploymentManager) updateShard(shardName string) {
	if _, exists := dm.shards[shardName]; !exists {
		dm.shards[shardName] = xdsconfig.MakeBackendShard(shardName)
	}

	tenants := []*xdsconfig.Tenant{}
	for name := range dm.byTenantName {
		if dm.byTenantName[name].shardName == shardName {
			tenants = append(tenants, dm.byTenantName[name].xdsTenant)
		}
	}
	dm.shards[shardName].Tenants = tenants
}

// RemoveTenant removes the tenant from the deployment.
func (dm *DeploymentManager) RemoveTenant(tenantName string) {
	if deployment, exists := dm.byTenantName[tenantName]; exists {
		delete(dm.byTenantName, tenantName)
		dm.updateShard(deployment.shardName)
		dm.OnChange <- ([]*xdsconfig.BackendShard{dm.shards[deployment.shardName]})
	}
}

// AddTenants adds teh tenants to the deployment.
func (dm *DeploymentManager) AddTenants(tenants ...*apimgmt.Tenant) {
	dirtyShards := make(map[string]bool)

	for _, tenant := range tenants {
		var shardName string
		if curdep, exists := dm.byTenantName[tenant.Name]; exists {
			shardName = curdep.shardName
		} else {
			shardName = dm.getNextShard()
		}
		dirtyShards[shardName] = true

		dm.byTenantName[tenant.Name] = &Deployment{
			shardName: shardName,
			apiTenant: tenant,
			xdsTenant: convertTenant(tenant),
		}
	}

	// Update the dirty shards
	for shardName := range dirtyShards {
		dm.updateShard(shardName)
		dm.OnChange <- ([]*xdsconfig.BackendShard{dm.shards[shardName]})
	}
}

// getNextShard Get the name of the next shard to assign a tenant too.
func (dm *DeploymentManager) getNextShard() string {
	shardName := fmt.Sprintf("back-%d", dm.nextShard)
	dm.nextShard = (dm.nextShard + 1) % config.NumShards
	return shardName
}

// convertTenant an api tenant to an envoy tenant
func convertTenant(apiTenant *apimgmt.Tenant) *xdsconfig.Tenant {
	tenant := xdsconfig.MakeTenant(apiTenant.Name)

	// Proxies
	for _, proxy := range apiTenant.Proxies {
		addProxy(tenant, proxy)
	}

	return tenant
}

// addProxy Add an xds proxy to the tenant.
func addProxy(tenant *xdsconfig.Tenant, apiProxy *apimgmt.Proxy) *xdsconfig.Proxy {
	p := xdsconfig.MakeProxy(tenant.Name, apiProxy)
	tenant.Proxies = append(tenant.Proxies, p)
	return p
}
