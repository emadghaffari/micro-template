package etcd

import (
	"fmt"
	"micro/config"
	"micro/pkg/logger"
	"sync"
	"time"

	client "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var (
	// Storage for etcd storage
	Storage Store = &etcd{}
)

// Store interface
type Store interface {
	Connect() error
	GetClient() *client.Client
	// WatchKey(ctx context.Context, key string, options ...client.OpOption)
	WatchKey(ctx context.Context, key string, callBack func(*client.Event), options ...client.OpOption)
	Put(ctx context.Context, key string, value interface{}) error
}

type etcd struct {
	cli  *client.Client
	once sync.Once
}

func (e *etcd) Connect() error {
	var err error
	e.once.Do(func() {
		e.cli, err = client.New(client.Config{
			Endpoints:   config.Global.ETCD.Endpoints,
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

func (e *etcd) GetClient() *client.Client {
	return e.cli
}

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

func (e *etcd) Put(ctx context.Context, key string, value interface{}) error {
	_, err := e.cli.Put(ctx, "sample_key", "sample_value")

	if err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("Config server error: %s", err))).
			Append(zap.Any("key", key)).
			Append(zap.Any("value", value)).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		return err
	}

	return nil
}
