package xds

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	"github.com/bladedancer/envoyxds/pkg/xdsconfig"
)

//TODO - Simple Entry Point Needs to be refactored

//TennantRouter Router Struct for service
type TennantRouter struct {
}

//GetTennantRouter Return a new Router Server
func GetTennantRouter() *TennantRouter {
	return &TennantRouter{}
}

//Run Start the service
func (t *TennantRouter) Run() {

	http.HandleFunc("/shard",
		func(w http.ResponseWriter, r *http.Request) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			log.Infof("Received %s", buf.String())
			toks := strings.Split(buf.String(), ":")
			fmt.Fprintf(w, "%s", getSafeShard(toks[0], toks[1]))

		},
	)
	go http.ListenAndServe(":12001", nil)
	log.Info("Simple Service Started")
}

//getSafeShard -
func getSafeShard(host string, path string) string {
	tenant := extractTenant(host)
	shard := deploymentManager.GetShardName(tenant)
	s := deploymentManager.shards[shard]
	if s == nil {
		// Edge case for unknown tennat request
		// Should not occur, but will result in 404
		log.Warnf("shard should not be nil  %s", shard)
		return "back-0"
	}
	ten := getTenantFromShard(s, tenant)
	confirmOrMakeRoute(shard, ten, path)
	return shard

}
func confirmOrMakeRoute(shard string, t *xdsconfig.Tenant, path string) {
	found := false
	for _, p := range t.Proxies {
		for _, r := range p.Routes {
			if r.Match.GetPath() == path {
				found = true
				break
			}
		}

		if found {
			break
		}
	}
	if !found {

		path = strings.TrimPrefix(path, "/")
		prox := xdsconfig.MakeProxy(t.Name, apimgmt.MakeProxy(t.Name, path))
		t.Proxies = append(t.Proxies, prox)
		deploymentManager.OnChange <- ([]*xdsconfig.BackendShard{deploymentManager.shards[shard]})
		//TODO add some channel synchronization with a countdown latch
		time.Sleep(500 * time.Millisecond)

	}
}

// getTenantFromShard - Assumes tenant/shard affinity
func getTenantFromShard(shard *xdsconfig.BackendShard, tName string) *xdsconfig.Tenant {
	var res *xdsconfig.Tenant
	for _, tenant := range shard.Tenants {
		if tenant.Name == tName {
			res = tenant
		}
	}
	return res
}

//extractTenant = expecting something like this test-12.bladedancer.dynu.net, will return 12
func extractTenant(host string) string {
	tenant := ""
	s := strings.Split(host, "-")
	// re := regexp.MustCompile(`\-(.*?)\.`)
	// something like this ^^ might be better
	if len(s) > 1 {
		s = strings.Split(s[1], ".")
		if len(s) > 0 {
			tenant = s[0]
		}
	}
	return tenant
}
