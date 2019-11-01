package envoyxds

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var routingShard Shard
var tenantCluster *TenantCluster

func version() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // good enough for now
}

// updateShard Update the snapshot cache with the shard details.
func updateShard(snapshotCache cache.SnapshotCache, shard Shard) error {
	xds := shard.GetXDS()
	log.Infof("Updating shard %s", shard.GetName())
	err := snapshotCache.SetSnapshot(shard.GetName(), cache.NewSnapshot(version(), nil, xds.CDS, xds.RDS, xds.LDS))
	if err != nil {
		log.Error(err)
	}
	return err
}

func updateRoutingShard(snapshotCache cache.SnapshotCache) {
	// TODO Create config for front
	updateShard(snapshotCache, routingShard)
}

func updateTenantShards(snapshotCache cache.SnapshotCache, tenants []*xdsconfig.Tenant) {
	// TODO tenant removal and improve
	for _, tenant := range tenants {
		tenantCluster.AddTenant(tenant)
	}

	// TODO: Should track the updated shards and only snapshot those.
	for _, shard := range tenantCluster.Shards {
		updateShard(snapshotCache, shard)
	}
}

func watch(snapshotCache cache.SnapshotCache) {
	routingShard = MakeRoutingShard("front")  // The node id of the envoys in the frontend
	tenantCluster = MakeTenantCluster("back") // The base name of the envoys statefulset

	tenants, updateChan := getTenants()
	updateRoutingShard(snapshotCache)
	updateTenantShards(snapshotCache, tenants)

	go func() {
		for tenants := range updateChan {
			updateTenantShards(snapshotCache, tenants)
		}
	}()
}
