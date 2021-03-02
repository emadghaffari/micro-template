package app

import (
	"fmt"
	"io"
	"log"
	"micro/config"
	zapLogger "micro/pkg/logger"
	pb "micro/proto"
	"micro/service"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	group "github.com/oklog/oklog/pkg/group"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var tracer opentracing.Tracer
var logger *zap.Logger

// StartApplication func
func StartApplication() {
	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init zap logger
	initLogger()

	// init configs
	if err := initConfigs(); err != nil {
		return
	}

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	closer, err := initJaeger()
	if err != nil {
		return
	}
	defer closer.Close()

	g := createService()
	initMetricsEndpoint(g)
	initCancelInterrupt(g)

	if err := g.Run(); err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.PanicLevel).Commit("Error in server run")
	}
}

func initLogger() {
	logger = zapLogger.GetZapLogger(false)
}

func initConfigs() error {
	// Current working directory
	dir, err := os.Getwd()
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.PanicLevel).Commit("init configs")
	}
	// read from file
	return config.LoadGlobalConfiguration(dir + "/config.yaml")
}

func initGRPCHandler(g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)
	// Add your GRPC options here

	lis, err := net.Listen("tcp", config.Global.GRPC.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.PanicLevel).Commit(err.Error())
	}

	g.Add(func() error {
		baseServer := grpc.NewServer(options...)

		// reflection for evans
		reflection.Register(baseServer)

		pb.RegisterMicroServer(baseServer, new(service.Micro))
		return baseServer.Serve(lis)
	}, func(error) {
		lis.Close()
	})
}

func initMetricsEndpoint(g *group.Group) {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", config.Global.HTTP.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "transport debug/HTTP during Listen err").Commit(err.Error())
	}
	g.Add(func() error {
		return http.Serve(debugListener, http.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}

func initCancelInterrupt(g *group.Group) {
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

func initJaeger() (io.Closer, error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: config.Global.Service.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.Global.Jaeger.LogSpans,
			LocalAgentHostPort: config.Global.Jaeger.HostPort,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	var closer io.Closer
	var err error
	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "during Listen jaeger err").Commit(err.Error())

		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
