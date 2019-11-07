package xds

import (
    "net/http"
    "fmt"
    "bytes"
    "strings"
    "github.com/bladedancer/envoyxds/pkg/xdsconfig"
    "github.com/bladedancer/envoyxds/pkg/apimgmt"
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
        func (w http.ResponseWriter, r *http.Request) {
            buf := new(bytes.Buffer)
            buf.ReadFrom(r.Body)
            log.Infof("Received %s", buf.String())
            toks:=strings.Split(buf.String(),":")
            fmt.Fprintf(w, "%s", getSafeShard(toks[0], toks[1]))

        },
    )
    go http.ListenAndServe(":12001", nil)
    log.Info("Simple Service Started")
}
//getSafeShard -
func getSafeShard(host string, path string) string {
    tenant:=extractTenant(host)
    shard:=deploymentManager.GetShardName(tenant)
    s:=deploymentManager.shards[shard]
    if s==nil {
        log.Warnf("shard should not be nil  %s", shard)
    }
    ten:=getTenantFromShard(s, tenant)
    confirmOrMakeRoute(shard, ten, path)
    return shard

}
func confirmOrMakeRoute(shard string, t *xdsconfig.Tenant, path string) {
    found:=false
    for _, p := range t.Proxies {
        if p.Route.Match.GetPath()==path {

            found=true
            break;
        }
    }
    if !found {

        path=strings.TrimPrefix(path,"/")
        prox:=xdsconfig.MakeProxy(t.Name, apimgmt.MakeProxy(t.Name, path))
        t.Proxies = append(t.Proxies, prox)
        deploymentManager.OnChange <- ([]*xdsconfig.BackendShard{deploymentManager.shards[shard]})

    }
}
// getTenantFromShard - Assumes tenant/shard affinity
func getTenantFromShard(shard *xdsconfig.BackendShard, tName string ) *xdsconfig.Tenant{
    var res *xdsconfig.Tenant
    for _, tenant := range shard.Tenants {
        if tenant.Name==tName {
            res=tenant
        }
    }
    return res
}
//extractTenant = expecting something like this test-12.bladedancer.dynu.net, will return 12
func extractTenant(host string) string {
    tenant:=""
    s := strings.Split(host, "-")
    // re := regexp.MustCompile(`\-(.*?)\.`)
    // something like this ^^ might be better
    if len(s) > 1 {
        s = strings.Split(s[1], ".")
        if len(s) > 0 {
           tenant=s[0]
        }
    }
    return tenant
}
