package xds

import (
	"fmt"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
)

// toShards converts an API Mgmt tenant to an envoy config on a shard.
func toShards(tenants ...*apimgmt.Tenant) []*xdsconfig.BackendShard {
	shards := []*xdsconfig.BackendShard{}

	// Create the shards
	for i := 0; i < config.NumShards; i++ {
		shards = append(shards, xdsconfig.MakeBackendShard(fmt.Sprintf("back-%d", i)))
	}

	// Populate them
	for i, tenant := range tenants {
		shardIndex := i % len(shards)
		shards[shardIndex].Tenants = append(shards[shardIndex].Tenants, convertTenant(tenant))
	}

	return shards
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
