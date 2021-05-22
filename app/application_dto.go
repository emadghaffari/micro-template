package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"micro/app/router/grpc"
	router "micro/app/router/http"
	"micro/client/broker"
	"micro/client/etcd"
	"micro/client/jtrace"
	"micro/client/postgres"
	"micro/client/redis"
	"micro/config"
	zapLogger "micro/pkg/logger"
	"net"
	"os"
	"os/signal"
	"syscall"

	group "github.com/oklog/oklog/pkg/group"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// StartApplication func
func (a *App) StartApplication() {
	fmt.Println("\n\n--------------------------------")

	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init zap logger
	a.initLogger()

	// init configs
	if err := a.initConfigs(); err != nil {
		return
	}

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	closer, err := a.initJaeger()
	if err != nil {
		return
	}
	defer closer.Close()

	if !config.Confs.GetDebug() {
		if err := a.initConfigServer(); err != nil {
			fmt.Println(err.Error())
		}
		defer etcd.Storage.GetClient().Close()
	}

	if err := a.initRedis(); err != nil {
		return
	}

	if err := a.initPostgres(); err != nil {
		return
	}
	defer postgres.Storage.Get().Close()

	if err := a.initMessageBroker(); err != nil {
		return
	}
	defer broker.Nats.Conn().Close()

	// create service
	g := a.createService()
	fmt.Printf("--------------------------------\n\n")
	if err := g.Run(); err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit("server stopped")
	}
}

// init zap logger
func (a *App) initLogger() {
	defer fmt.Printf("zap logger is available \n")
	zapLogger.SetLogPath("logs")
	logger = zapLogger.GetZapLogger(config.Confs.GetDebug())
}

// init configs
func (a *App) initConfigs() error {

	// Current working directory
	dir, err := os.Getwd()
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit("init configs")
		return err
	}

	defer fmt.Printf("configs loaded from file successfully \n")
	// read from file
	return config.Confs.Load(dir + "/config.yaml")
}

// init grpc connection
func (a *App) initGRPCHandler(g *group.Group) {
	defer fmt.Printf("grpc connected port:%s \n", config.Confs.Get().Service.GRPC.Port)

	lis, err := net.Listen("tcp", config.Confs.Get().Service.GRPC.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
	}

	g.Add(func() error {
		return grpc.Router.GetRouter().Serve(lis)
	}, func(error) {
		lis.Close()
	})
}

// init HTTP Endpoint
// add rest endpoints
func (a *App) initHTTPEndpoint(g *group.Group) {
	defer fmt.Printf("metrics started port:%s \n", config.Confs.Get().Service.HTTP.Port)

	g.Add(func() error {
		if err := router.Router.GetRouter().Run(config.Confs.Get().Service.HTTP.Port); err != nil {
			zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "transport debug/HTTP during Listen err").Commit(err.Error())
			return err
		}
		return nil
	}, func(error) {})
}

// init cancle Interrupt
func (a *App) initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

// init jaeger tracer
func (a *App) initJaeger() (io.Closer, error) {
	closer, err := jtrace.Tracer.Connect()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Jaeger loaded successfully \n")
	return closer, nil
}

// in production you load envs from etcd storage
func (a *App) initConfigServer() error {
	defer fmt.Printf("etcd storage loaded successfully \n")
	if err := etcd.Storage.Connect(config.Confs.Get()); err != nil {
		return err
	}

	// get watch list for this service with service ID
	etcd.Storage.GetKey(context.Background(), config.Confs.Get().Service.ID, func(kv *mvccpb.KeyValue) {
		// set configs from storage to struct - if exists in Set method
		config.Confs.Set(string(kv.Key), kv.Value)

	}, clientv3.WithPrefix())

	// loop over watchList
	for _, key := range config.Confs.Get().ETCD.WatchList {

		// get configs for first time on app starts
		err := etcd.Storage.GetKey(context.Background(), key, func(kv *mvccpb.KeyValue) {
			// set configs from storage to struct - if exists in Set method
			config.Confs.Set(string(kv.Key), kv.Value)

		}, clientv3.WithPrefix())
		if err != nil {
			return err
		}

		// start to watch keys
		etcd.Storage.WatchKey(context.Background(), key, func(e *clientv3.Event) {

			// set configs from storage to struct - if exists in Set method
			config.Confs.Set(string(e.Kv.Key), e.Kv.Value)

		}, clientv3.WithPrefix())
	}

	return nil
}

// init message broker
func (a *App) initMessageBroker() error {
	if err := broker.Nats.Connect(config.Confs.Get()); err != nil {
		return err
	}

	fmt.Printf("nats message broker loaded successfully \n")
	return nil
}

// init Redis database
func (a *App) initRedis() error {
	if err := redis.Storage.Connect(config.Confs.Get()); err != nil {
		return err
	}

	fmt.Printf("redis database loaded successfully \n")
	return nil
}

// init postgres database
func (a *App) initPostgres() error {
	if err := postgres.Storage.Connect(config.Confs.Get()); err != nil {
		return err
	}

	fmt.Printf("postgres database loaded successfully \n")
	return nil
}

func (a *App) createService() (g *group.Group) {
	g = &group.Group{}

	// init GRPC Handlers
	Base.initGRPCHandler(g)

	// init http endpoints
	Base.initHTTPEndpoint(g)

	// init cancel
	Base.initCancelInterrupt(g)
	return g
}
