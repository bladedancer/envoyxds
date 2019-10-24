package envoyxds

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var version int = 0

//Load file and add a snapshot, incrementing version each time
func loadSnapshot(numTenants int, numRoutes int, domain string, c cache.SnapshotCache) {
	routeConfig := GetRouteConfigurations(numTenants, numRoutes, domain)
	rts := []cache.Resource{
		routeConfig,
	}
	log.Infof("Num Tenants: %d, Num Routes %d, Domain: %s", numTenants, numRoutes, domain)
	err := c.SetSnapshot("shard-0", cache.NewSnapshot(string(version), nil, nil, rts, nil))
	version++
	if err != nil {
		log.Error(err)
	}
}

func watch(snapshotCache cache.SnapshotCache, conf Config) {
	loadSnapshot(conf.NumTenants, conf.NumRoutes, conf.Domain, snapshotCache)
	// TODO SOME DYNAMIC UPDATES...hook up to db or timer
}
