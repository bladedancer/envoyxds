package envoyxds

import (
	"bytes"
	"io/ioutil"
	"strings"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/fsnotify/fsnotify"
	"github.com/gogo/protobuf/jsonpb"
)

var version int = 0

//extract tenant based on filename this is brute force method assumes "fmt: of %s+ tennant id+ json
func extractTenant(fn string) string {
	tmp := strings.Split(fn, "-")
	res := strings.Split(tmp[1], ".")
	return res[0]
}

//Watch for fs event changes on the config filesystem
// and load new snapshot for the given tenant identified in the fs event
// Note: this is not the way we will want to implement, just a way to demonstrate
func watchCfg(c cache.SnapshotCache, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("modified file:", event.Name)
					t := extractTenant(event.Name)
					loadSnapshot(t, event.Name, c)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// GetRoute Load the RouteConfiguration proto from Json file and return
func GetRoute(path string) *api.RouteConfiguration {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
	}
	route := api.RouteConfiguration{}
	rdr := bytes.NewReader(raw)

	jsonpb.Unmarshal(rdr, &route)
	return &route

}

//Load file and add a snapshot, incrementing version each time
func loadSnapshot(t string, fileLoc string, c cache.SnapshotCache) {
	route := GetRoute(fileLoc)
	rts := []cache.Resource{
		route,
	}
	log.Info(route)
	err := c.SetSnapshot(t, cache.NewSnapshot(string(version), nil, nil, rts, nil))
	version++
	if err != nil {
		log.Error(err)
	}
}

func watch(snapshotCache cache.SnapshotCache, path string) {
	tenants := "TenantA,TenantB"
	tens := strings.Split(tenants, ",")
	for _, t := range tens {
		loadSnapshot(t, path+"/route-"+t+".json", snapshotCache)
	}
	go watchCfg(snapshotCache, path)
}
