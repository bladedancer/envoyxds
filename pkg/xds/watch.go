package xds

import (
	"context"
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	redis "github.com/bladedancer/envoyxds/pkg/cache"
	"github.com/bladedancer/envoyxds/pkg/datasource"
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

func updateCredentials(tenants ...*apimgmt.Tenant) {
	// Obviously not ideal to not detect credential updates separately from any other tenant update
	// but it'll do for now....seem to be writing that a lot.
	if tenants != nil {
		for _, tenant := range tenants {
			for _, proxy := range tenant.Proxies {
				if proxy.Backend.Credential != nil {
					cacheCon.Set(context.Background(), fmt.Sprintf("%s-%s", tenant.Name, proxy.Name), &redis.Credential{Credential: proxy.Backend.Credential}, 0)
				}
			}
		}
	}
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

	tenants, updateChan := datasource.TenantDatasource.GetTenants()
	go func() {
		for tenants := range updateChan {
			// There'd be an sync concern here in the real world.
			deploymentManager.AddTenants(tenants...)
			updateCredentials(tenants...)
		}
	}()

	// Update BE credentials in Redis
	updateCredentials(tenants...)

	// Update deployments
	deploymentManager.AddTenants(tenants...)
}
