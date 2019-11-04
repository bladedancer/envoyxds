package xds

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var frontend *xdsconfig.FrontendShard
var backends []*xdsconfig.BackendShard

func version() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // good enough for now
}

// updateShard Update the snapshot cache with the shard details.
func updateShard(snapshotCache cache.SnapshotCache, shard xdsconfig.Shard) error {
	xds := shard.GetXDS()
	log.Infof("Updating shard %s (%d:%d:%d)", shard.GetName(), len(xds.CDS), len(xds.RDS), len(xds.LDS))
	err := snapshotCache.SetSnapshot(shard.GetName(), cache.NewSnapshot(version(), nil, xds.CDS, xds.RDS, xds.LDS))
	if err != nil {
		log.Error(err)
	}

	return err
}

func updateBackends(snapshotCache cache.SnapshotCache) {
	for _, shard := range backends {
		updateShard(snapshotCache, shard)
	}
}

func watch(snapshotCache cache.SnapshotCache) {
	tenants, updateChan := apimgmt.GetTenants()

	backends = toShards(tenants...)
	frontend = xdsconfig.MakeFrontendShard("front")

	updateBackends(snapshotCache)
	updateShard(snapshotCache, frontend)

	go func() {
		for tenants := range updateChan {
			backends = toShards(tenants...)
			updateBackends(snapshotCache)
		}
	}()
}
