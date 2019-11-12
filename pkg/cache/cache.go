package cache

import (
    "golang.org/x/net/context"
    "github.com/golang/protobuf/proto"
)

// Cache - Interface for KV Storage
type Cache interface{

    // Concept Stolen from Kubernetes

    // Endpoints returns an array of endpoints for the storage
    Endpoints() []string

    // Connect to etcd using client v3 api
    Connect() error

    // Close connection to etcd
    Close() error

    // Put adds a new value to a key if the key does not already exist; else it updates the value passed
    // to the specified key. 'ttl' is time-to-live in seconds (0 means forever).
    Set(ctx context.Context, key string, val proto.Message, ttl int64) error

    // Get unmarshals the protocol buffer message found at key into out, if found.
    // If not found and ignoreNotFound is set, then out will be a zero object, otherwise
    // error will be set to not found. A non-existing node or an empty response are both
    // treated as not found.
    Get(ctx context.Context, key string, out proto.Message, ignoreNotFound bool) error

    // Delete(ctx context.Context, key string, recurse bool, out proto.Message) error
    // TODO: will need to add preconditions support
    // if recurse then all the key having the same path under 'key' are going to be deleted
    // if !recurse then only 'key' is going to be deleted
    Delete(ctx context.Context, key string, recurse bool, out proto.Message) error
/*
    // Update performs a guaranteed update, which means it will continue to retry until an update succeeds or the request is canceled.
    Update(ctx context.Context, key string, uf UpdateFunc, template proto.Message) error

    // List returns all the values that match the filter.
    List(ctx context.Context, key string, filter Filter, obj proto.Message, out *[]proto.Message) error

    // Watch begins watching the specified key.
    Watch(ctx context.Context, key string, resourceVersion int64, filter Filter) (WatchInterface, error)

    // WatchList begins watching the specified key's items.
    WatchList(ctx context.Context, key string, resourceVersion int64, filter Filter) (WatchInterface, error)

    // CompareAndSet atomically sets the value to the given updated value if the current value == the expected value
    CompareAndSet(ctx context.Context, key string, expect proto.Message, update proto.Message) error
*/
}

