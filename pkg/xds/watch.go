package xds

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var frontend *xdsconfig.FrontendShard
var deploymentManager = MakeDeploymentManager()

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

func watch(snapshotCache cache.SnapshotCache) {
	// Frontend is static for now
	frontend = xdsconfig.MakeFrontendShard("front")
	updateShard(snapshotCache, frontend)

	// Tenants and shard contents are dynamic so listen for
	// changes and update accordingly
	go func() {
		for shards := range deploymentManager.OnChange {
			for i := 0; i < len(shards); i++ {
				updateShard(snapshotCache, shards[i])
			}
		}
	}()

	tenants, updateChan := apimgmt.GetTenants()
	go func() {
		for tenants := range updateChan {
			deploymentManager.AddTenants(tenants...)
		}
	}()

	deploymentManager.AddTenants(tenants...)
}
