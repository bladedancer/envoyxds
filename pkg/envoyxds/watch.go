package envoyxds

import (
	"fmt"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache"
)

var version int = 0

//Load file and add a snapshot, incrementing version each time
func loadSnapshot(numTenants int, numRoutes int, domain string, c cache.SnapshotCache) {
	listenerConfig := GetListener(0, domain)
	lrs := []cache.Resource{
		listenerConfig,
	}

	routeConfig := GetRouteConfigurations(numTenants, numRoutes, domain)
	rts := []cache.Resource{
		routeConfig,
	}
	log.Infof("Num Tenants: %d, Num Routes %d, Domain: %s", numTenants, numRoutes, domain)
	err := c.SetSnapshot("shard-0", cache.NewSnapshot(fmt.Sprintf("%d", version), nil, nil, rts, lrs))
	version++
	if err != nil {
		log.Error(err)
	}
}

func watch(snapshotCache cache.SnapshotCache, conf Config) {
	loadSnapshot(conf.NumTenants, conf.NumRoutes, conf.Domain, snapshotCache)

	// Hack in a special case where the cache keeps growing
	if conf.Pump > 0 {
		log.Infof("Pumping new route every %d seconds", conf.Pump)
		go pump(snapshotCache, conf)
	}
}

func pump(snapshotCache cache.SnapshotCache, conf Config) {
	tick := time.NewTicker(time.Duration(conf.Pump) * time.Second)
	i := 1

	listenerConfig := GetListener(0, conf.Domain)
	lrs := []cache.Resource{
		listenerConfig,
	}

	go func() {
		for {
			select {
			case <-tick.C:
				log.Infof("Pump %d", conf.NumRoutes+i)
				routeConfig := GetRouteConfigurations(conf.NumTenants, conf.NumRoutes+i, conf.Domain)
				i++
				rts := []cache.Resource{
					routeConfig,
				}
				err := snapshotCache.SetSnapshot("shard-0", cache.NewSnapshot(fmt.Sprintf("%d", version), nil, nil, rts, lrs))
				version++
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()
}
