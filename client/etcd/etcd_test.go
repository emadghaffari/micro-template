package etcd

import (
	"context"
	"fmt"
	"micro/config"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"
	clientv3test "go.etcd.io/etcd/tests/v3/integration/clientv3"
)

func TestConnect(t *testing.T) {

	integration.BeforeTest(t)
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1, SkipCreatingClient: true})
	defer clus.Terminate(t)

	tests := []struct {
		step string
		conf config.Config
		err  error
	}{
		{
			step: "A",
			conf: config.Config{
				Debug: true,
				ETCD: config.ETCD{
					Endpoints: nil,
					Username:  "ruser",
					Password:  "T0pS3cr3t",
				},
			},
			err: fmt.Errorf("etcdclient: no available endpoints"),
		},
		{
			step: "B",
			conf: config.Config{
				Debug: true,
				ETCD: config.ETCD{
					Endpoints: []string{clus.Members[0].GRPCAddr()},
					Username:  "ruser",
					Password:  "T0pS3cr3t",
				}},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			err := Storage.Connect(tc.conf)
			once = sync.Once{}
			if err != nil && tc.err != nil {
				assert.Contains(t, err.Error(), tc.err.Error())
			}

			if err != nil && tc.err == nil {
				assert.Equal(t, nil, err)
			}

			if clientv3test.IsClientTimeout(err) {
				assert.Equal(t, true, err)
			}
		})
	}
}

func TestGetClient(t *testing.T) {
	integration.BeforeTest(t)
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1, SkipCreatingClient: true})
	defer clus.Terminate(t)

	c := Storage.GetClient()
	if c == nil {
		assert.Equal(t, c, c)
	}
}

func TestWatchKey(t *testing.T) {
	integration.BeforeTest(t)
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1, SkipCreatingClient: true})
	defer clus.Terminate(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	Storage.WatchKey(ctx, "key", func(e *clientv3.Event) {
		assert.Equal(t, "\"DATA\"", string(e.Kv.Value))
	})

	if err := Storage.Put(ctx, "key", "DATA"); err != nil {
		assert.Error(t, err, "error in put DATA")
	}

	time.Sleep(time.Second)
}

func TestGetKey(t *testing.T) {
	integration.BeforeTest(t)
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1, SkipCreatingClient: true})
	defer clus.Terminate(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := Storage.Put(ctx, "key2", "DATA"); err != nil {
		assert.Error(t, err, "error in put DATA")
	}

	Storage.GetKey(ctx, "key2", func(kv *mvccpb.KeyValue) {
		assert.Equal(t, "\"DATA\"", string(kv.Value))
	})

	cancel()
	if err := Storage.GetKey(ctx, "key2", func(kv *mvccpb.KeyValue) {}); err != nil {
		assert.Error(t, err, "invalid error")
	}
}

func TestPut(t *testing.T) {
	invalidData := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := Storage.Put(ctx, "key2", invalidData); err != nil {
		assert.Error(t, err, "error in put DATA")
	}

	validData := "DATA"
	ctx2, c := context.WithTimeout(context.Background(), 30*time.Second)
	c()

	if err := Storage.Put(ctx2, "key2", validData); err != nil {
		assert.Error(t, err, "error in put DATA")
	}
}
