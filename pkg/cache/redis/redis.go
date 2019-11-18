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

// New returns an redis implementation of cache.Cache Interface
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
    _, err := r.client.Ping().Result()
  return err
}

// Close connection to etcd
func (r *redis) Close() error {

    err:=r.client.Close()
    if err != nil {

        log.Warnf("Error on Set %s", err)
    }
    return err
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
func (r* redis) Get(ctx context.Context, key string, out proto.Message, ignoreNotFound bool) error {

  b, err:=r.client.Get(key).Bytes()
  if err!=nil {
      log.Warnf("Error on Get %s", err)
  }
  proto.Unmarshal(b, out)
  return err
}

// Delete Delete the key at a specific value
func (r *redis) Delete(ctx context.Context, key string, recurse bool, out proto.Message) error {
    r.client.Del(key).Err()
return r.client.Del(key).Err()
}

