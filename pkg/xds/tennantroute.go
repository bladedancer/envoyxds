package xds

import (
    "net/http"
    "fmt"
    "bytes"
//    "regexp"
    "strings"
)

//TODO - Simple Entry Point Needs to be refactored

//TennantRouter Router Struct for service
type TennantRouter struct{
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
            newStr := buf.String()
            log.Infof("Received %s", newStr)
            fmt.Fprintf(w, "%s", getSafeShard(newStr))
        },
    )
    go http.ListenAndServe(":12001", nil)
    log.Info("Simple Service Started")
}
//getSafeShard - expecting something like this test-12.bladedancer.dynu.net
func getSafeShard(tennant string) string {
    s:=strings.Split(tennant,"-")
    shard:="back-0"
    // re := regexp.MustCompile(`\-(.*?)\.`)
    // something like this ^^ might be better
    if len(s)>1 {
        s=strings.Split(s[1],".")
        if len(s) > 0 {
            log.Infof("Looking for shard for %s", s[0])
            shard=deploymentManager.GetShardName(s[0])
            log.Infof("Found a good shard %s", shard)
        }
    }
    return shard

}
