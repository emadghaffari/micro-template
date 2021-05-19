package etcd

import (
	"encoding/json"
	"fmt"
	"micro/config"
	"micro/pkg/logger"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	client "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
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

// connect method, connect to etcd db
func (e *etcd) Connect(conf config.Config) error {
	var err error
	once.Do(func() {
		e.cli, err = client.New(client.Config{
			Endpoints:   conf.ETCD.Endpoints,
			Username:    conf.ETCD.Username,
			Password:    conf.ETCD.Password,
			DialTimeout: 5 * time.Second,
		})
	})
	if err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config server error: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		return err
	}

	return nil
}

// get etcd client
func (e *etcd) GetClient() *client.Client {
	return e.cli
}

// watch on a key
func (e *etcd) WatchKey(ctx context.Context, key string, callBack func(*client.Event), options ...client.OpOption) {
	rch := e.cli.Watch(ctx, key, options...)

	go func(rch client.WatchChan) {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				callBack(ev)
			}
		}
	}(rch)
}

// get value of key
func (e *etcd) GetKey(ctx context.Context, key string, callBack func(*mvccpb.KeyValue), options ...client.OpOption) error {
	resp, err := e.cli.Get(ctx, key, options...)
	if err != nil {
		return err
	}

	for _, k := range resp.Kvs {
		callBack(k)
	}

	return nil
}

// put into etcd
func (e *etcd) Put(ctx context.Context, key string, value interface{}, options ...client.OpOption) error {
	bts, err := json.Marshal(value)
	if err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config server put error: %s", err))).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", value)).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		return err
	}

	if _, err := e.cli.Put(ctx, key, string(bts), options...); err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config server put error: %s", err))).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", value)).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		return err
	}

	return nil
}
