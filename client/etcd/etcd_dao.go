package etcd

import (
	"micro/config"
	"sync"

	"go.etcd.io/etcd/api/v3/mvccpb"
	client "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

var (
	// Storage for etcd storage
	Storage Store = &etcd{}
	once    sync.Once
)

// Store interface
type Store interface {
	Connect(conf config.Config) error
	GetClient() *client.Client
	GetKey(ctx context.Context, key string, callBack func(*mvccpb.KeyValue), options ...client.OpOption) error
	WatchKey(ctx context.Context, key string, callBack func(*client.Event), options ...client.OpOption)
	Put(ctx context.Context, key string, value interface{}, options ...client.OpOption) error
}

type etcd struct {
	cli *client.Client
}
