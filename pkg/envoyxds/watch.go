package envoyxds

import (
	"fmt"
	"time"

	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

func version() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // good enough for now
}

func updateSnapshot(snapshotCache cache.SnapshotCache, tenants []*xdsconfig.Tenant) {
	xds := xdsconfig.MakeXDS(tenants)
	err := snapshotCache.SetSnapshot("shard-0", cache.NewSnapshot(version(), nil, xds.CDS, xds.RDS, xds.LDS))
	if err != nil {
		log.Error(err)
	}
}

func watch(snapshotCache cache.SnapshotCache) {
	tenants, updateChan := getTenants()
	updateSnapshot(snapshotCache, tenants)

	go func() {
		for tenants := range updateChan {
			updateSnapshot(snapshotCache, tenants)
		}
	}()
}
