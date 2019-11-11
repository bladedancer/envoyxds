package redis

import (
  rediscli "github.com/go-redis/redis"
    "time"
    "github.com/bladedancer/envoyxds/pkg/cache"
    "github.com/golang/protobuf/proto"
    "golang.org/x/net/context"
)

// DefaultEndpoint
const (
    DefaultEndpoint = "http://localhost:6379"
)



type redis struct {
    client *rediscli.Client
    endpoints  []string
    pathPrefix string
    timeout    time.Duration
    connected  bool
    }

// New returns an redis implementation of storage.Interface
func New(endpoints []string, prefix string, timeout time.Duration) cache.Cache {
    return &redis{endpoints: endpoints, pathPrefix: prefix, timeout: timeout}
}

// Endpoints gets the endpoints redis
func (r *redis) Endpoints() []string {
    if len(r.endpoints) == 0 {
        r.endpoints=[]string{DefaultEndpoint}
    }
    return r.endpoints
}
//Connect to the redis service
func (r *redis) Connect() error {
    r.client = rediscli.NewClient(&rediscli.Options{
        Addr:         r.endpoints[0],
        DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     10,
        PoolTimeout:  30 * time.Second,
    })
  return nil
}

// Close connection to etcd
func (r *redis) Close() error {

   // TODO
return r.client.Close()
}


// Put adds a new value to a key if the key does not already exist; else it updates the value passed
// to the specified key. 'ttl' is time-to-live in seconds (0 means forever).
func (r* redis) Set(ctx context.Context, key string, val proto.Message, ttl int64) error {
    data, err := proto.Marshal(val)
    if err != nil {
        return err
    }
    err = r.client.Set(key, data, 0).Err()
    if err != nil {

        log.Warnf("Error on Set %s", err)
    }
    return err
}

// Get unmarshals the protocol buffer message found at key into out, if found.
// If not found and ignoreNotFound is set, then out will be a zero object, otherwise
// error will be set to not found. A non-existing node or an empty response are both
// treated as not found.
func (r* redis) Get(ctx context.Context, key string, out proto.Message, ignoreNotFound bool) error {

  b, err:=r.client.Get("key").Bytes()
  proto.Unmarshal(b, out)
  return err
}

// Delete(ctx context.Context, key string, recurse bool, out proto.Message) error
// TODO: will need to add preconditions support
// if recurse then all the key having the same path under 'key' are going to be deleted
// if !recurse then only 'key' is going to be deleted
func (r *redis) Delete(ctx context.Context, key string, recurse bool, out proto.Message) error {
// TODO
return nil
}

