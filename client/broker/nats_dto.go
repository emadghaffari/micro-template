package broker

import (
	"context"
	"micro/config"
	zapLogger "micro/pkg/logger"
	"time"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Connect nats broker
func (n *nts) Connect(conf config.Config) error {
	n.logger = zapLogger.GetZapLogger(config.Confs.GetDebug())
	var err error
	once.Do(func() {
		var conn *nats.Conn
		opts := nats.Options{
			Name:           conf.Service.Name,
			Secure:         conf.Nats.Auth,
			User:           conf.Nats.Username,
			Password:       conf.Nats.Password,
			Servers:        conf.Nats.Endpoints,
			PingInterval:   time.Second * 60,
			AllowReconnect: conf.Nats.AllowReconnect,
			MaxReconnect:   conf.Nats.MaxReconnect,
			ReconnectWait:  time.Duration(conf.Nats.ReconnectWait) * time.Second,
			Timeout:        time.Duration(conf.Nats.Timeout) * time.Second,
			AsyncErrorCB:   n.errorReporter(n.logger),
		}

		// try to connect to nats message broker
		conn, err = opts.Connect()
		if err != nil {
			logger := zapLogger.GetZapLogger(config.Confs.GetDebug())
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
			return
		}

		nc, err = nats.NewEncodedConn(conn, conf.Nats.Encoder)
		if err != nil {
			return
		}

	})

	return err
}

// Conn get Connection
func (n *nts) Conn() *nats.EncodedConn {
	return nc
}

// Publish new message
func (n *nts) Publish(ctx context.Context, subject string, value interface{}) error {
	if err := nc.Publish(subject, &value); err != nil {
		logger := zapLogger.GetZapLogger(config.Confs.GetDebug())
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}

	return nil
}

// SendChan, send a untyped chan
func (n *nts) SendChan(subject string, ch chan interface{}) error {
	return nc.BindSendChan(subject, ch)
}

// SendByContext, send a request and get response with spesific context
func (n *nts) SendByContext(ctx context.Context, subject string, req interface{}, resp interface{}) error {
	if err := nc.RequestWithContext(ctx, subject, req, resp); err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return err
	}
	return nil
}

// RequestWithReply, send a request and get response of request
// Then call Flush
func (n *nts) RequestWithReply(subject string, req interface{}, resp string) error {
	if err := nc.PublishRequest(subject, resp, req); err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return err
	}

	if err := nc.Flush(); err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return err
	}

	return nil
}

// Subscribe, start to subscribe to a subject
func (n *nts) Subscribe(subject string, callBack func(resp *nats.Msg)) (*nats.Subscription, error) {
	sub, err := nc.Subscribe(subject, callBack)
	if err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return nil, err
	}

	return sub, nil
}

// RecvChan, BindRecvChan
func (n *nts) RecvChan(subject string, ch chan interface{}) (*nats.Subscription, error) {
	sub, err := nc.BindRecvChan(subject, ch)

	if err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return nil, err
	}

	return sub, nil
}

// RecvGroup, connect to subject in group mode
func (n *nts) RecvGroup(subject, queue string, callBack nats.Handler) (*nats.Subscription, error) {
	sub, err := nc.QueueSubscribe(subject, queue, callBack)
	if err != nil {
		n.logger.Info("msg", zap.Any("err", err.Error()))
		return nil, err
	}

	return sub, nil
}

// errorReporter, when nats has error
func (n *nts) errorReporter(log *zap.Logger) nats.ErrHandler {
	return func(_ *nats.Conn, sub *nats.Subscription, err error) {
		pendingMsgs, pendingBytes, _ := sub.Pending()
		droppedMsgs, _ := sub.Dropped()
		maxMsgs, maxBytes, _ := sub.PendingLimits()

		log.Error(err.Error(),
			zap.Any("subject", sub.Subject),
			zap.Any("queue", sub.Queue),
			zap.Any("pending_msgs", pendingMsgs),
			zap.Any("pending_bytes", pendingBytes),
			zap.Any("max_msgs_pending", maxMsgs),
			zap.Any("max_bytes_pending", maxBytes),
			zap.Any("dropped_msgs", droppedMsgs),
			zap.Any("message", "Error while consuming from nats"),
		)
	}
}
