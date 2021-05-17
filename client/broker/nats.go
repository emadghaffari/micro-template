package broker

import (
	"context"
	"micro/config"
	zapLogger "micro/pkg/logger"
	"strings"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	Nats  NatsBroker = &nts{}
	oncen sync.Once
)

// NatsBroker interface
type NatsBroker interface {
	Connect() error
	Conn() *nats.EncodedConn
	Publish(ctx context.Context, subject string, value interface{}) error
}

// nts struct for nats message broker
type nts struct {
	conn *nats.EncodedConn
}

// Connect nats broker
func (n *nts) Connect() error {
	var err error
	oncen.Do(func() {
		var conn *nats.Conn
		opts := nats.Options{
			Name:         config.Confs.Get().Service.Name,
			Secure:       config.Confs.Get().Nats.Auth,
			User:         config.Confs.Get().Nats.Username,
			Password:     config.Confs.Get().Nats.Password,
			MaxReconnect: 10,
			Url:          strings.Join(config.Confs.Get().Nats.Endpoints, ","),
			PingInterval: time.Minute * 10,
		}

		// try to connect to nats message broker
		conn, err = opts.Connect()
		if err != nil {
			logger := zapLogger.GetZapLogger(config.Confs.Debug())
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
			return
		}

		n.conn, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
		if err != nil {
			return
		}

	})

	return err
}

// Conn get Connection
func (n *nts) Conn() *nats.EncodedConn {
	return n.conn
}

// Publish new message
func (n *nts) Publish(ctx context.Context, subject string, value interface{}) error {
	if err := n.conn.Publish(subject, &value); err != nil {
		logger := zapLogger.GetZapLogger(config.Confs.Debug())
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}

	return nil
}
