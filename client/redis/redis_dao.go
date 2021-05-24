package redis

import (
	"context"
	"micro/config"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Storage store = &rds{}
	logger  *zap.Logger
	once    sync.Once
)

// store interface is interface for store things into redis
type store interface {
	Connect(config config.Config) error
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key ...string) error
}

// rds struct for redis client
type rds struct {
	db *redis.Client
}
