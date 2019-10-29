package envoyxds

import (
	"fmt"

	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var version int = 0

func updateSnapshot(snapshotCache cache.SnapshotCache) {
	tenants := getTenants()
	xds := xdsconfig.MakeXDS(tenants)
	err := snapshotCache.SetSnapshot("shard-0", cache.NewSnapshot(fmt.Sprintf("%d", version), nil, xds.CDS, xds.RDS, xds.LDS))
	if err != nil {
		log.Error(err)
	}
	version++
}

func watch(snapshotCache cache.SnapshotCache) {
	updateSnapshot(snapshotCache)
	// TODO: PUMP
}
