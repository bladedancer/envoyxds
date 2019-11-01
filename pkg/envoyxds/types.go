package envoyxds

import (
	"fmt"

	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
)

// Shard information
type Shard interface {
	GetName() string
	GetXDS() *xdsconfig.XDS
}

// RoutingShard the sharded tenant details
type RoutingShard struct {
	Name      string // This will map to the envoy node id
	Resources *xdsconfig.XDS
}

// MakeRoutingShard creates an empty shard.
func MakeRoutingShard(name string) *RoutingShard {
	routingResources := xdsconfig.MakeXDS()
	routingResources.LDS = append(routingResources.LDS, xdsconfig.GetListenerResources()...) // TODO separate listener config for front & back
	routingResources.CDS = append(routingResources.CDS, xdsconfig.GetFrontendClusterResources()...)

	return &RoutingShard{
		Name:      name,
		Resources: routingResources,
	}
}

// GetName returns the xDS resources for the shard.
func (rs *RoutingShard) GetName() string {
	return rs.Name
}

// GetXDS returns the xDS resources for the shard.
func (rs *RoutingShard) GetXDS() *xdsconfig.XDS {
	return rs.Resources
}

// TenantShard the sharded tenant details
type TenantShard struct {
	Name    string // This will map to the envoy node id
	Tenants map[string]*xdsconfig.XDS
}

// TenantCluster the details of all the tenant shards
type TenantCluster struct {
	Shards []*TenantShard
}

// MakeTenantShard creates an empty shard.
func MakeTenantShard(name string) *TenantShard {
	return &TenantShard{
		Name:    name,
		Tenants: make(map[string]*xdsconfig.XDS),
	}
}

// MakeTenantCluster creates an empty shard cluster.
func MakeTenantCluster(name string) *TenantCluster {
	// TODO make dynamic yada yada
	shards := make([]*TenantShard, config.NumShards)
	for i := 0; i < len(shards); i++ {
		shards[i] = MakeTenantShard(
			fmt.Sprintf("%s-%d", name, i),
		)
	}
	return &TenantCluster{
		Shards: shards,
	}
}

// AddTenant add a tenant to the cluster. The appropriate shard
// is selected automatically.
func (sc *TenantCluster) AddTenant(tenant *xdsconfig.Tenant) {
	// This is slow but just messing around. Check if the tenant
	// is already in the cluster, if they are update them in place,
	// if not then add them to the one with the fewest.
	found := false
	resources := xdsconfig.MakeXDS(tenant)

	for _, shard := range sc.Shards {
		if _, found = shard.Tenants[tenant.Name]; found {
			// Update the shard with latest
			shard.Tenants[tenant.Name] = resources
			break
		}
	}

	if !found {
		var shortest *TenantShard
		for _, shard := range sc.Shards {
			if shortest == nil {
				shortest = shard
			} else if len(shard.Tenants) < len(shortest.Tenants) {
				shortest = shard
			}
		}
		shortest.Tenants[tenant.Name] = resources
	}
}

// GetName returns the name of the shard
func (ts *TenantShard) GetName() string {
	return ts.Name
}

// GetXDS returns the xDS resources for the shard.
func (ts *TenantShard) GetXDS() *xdsconfig.XDS {
	// TODO ... better
	xds := xdsconfig.MakeXDS()
	for t, r := range ts.Tenants {
		log.Infof("Adding %s to %s", t, ts.Name)
		xds.Add(r)
	}
	return xds
}
