package broker

import (
	"context"
	"micro/config"
	"sync"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	Nats NatsBroker = &nts{}
	nc   *nats.EncodedConn
	once sync.Once
)

// NatsBroker interface
type NatsBroker interface {
	Connect(conf config.Config) error
	Conn() *nats.EncodedConn
	Publish(ctx context.Context, subject string, value interface{}) error
	SendChan(subject string, ch chan interface{}) error
	SendByContext(ctx context.Context, subject string, req interface{}, resp interface{}) error
	RequestWithReply(subject string, req interface{}, resp string) error
	Subscribe(subject string, callBack func(resp *nats.Msg)) (*nats.Subscription, error)
	RecvChan(subject string, ch chan interface{}) (*nats.Subscription, error)
	RecvGroup(subject, queue string, callBack nats.Handler) (*nats.Subscription, error)
	errorReporter(log *zap.Logger) nats.ErrHandler
}

// nts struct for nats message broker
type nts struct {
	logger *zap.Logger
}
