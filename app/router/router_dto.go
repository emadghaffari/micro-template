package router

import (
	"context"
	"fmt"
	"log"
	"micro/app/router/middleware"
	"micro/client/jtrace"
	"micro/config"
	"micro/controller"
	zapLogger "micro/pkg/logger"
	"micro/proto/pb"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (g micro) GetRouter() error {
	logger := zapLogger.GetZapLogger(config.Confs.GetDebug())

	// net listen for grpc port
	lis, err := net.Listen("tcp", config.Confs.Get().Service.GRPC.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
	}

	// start new server
	baseServer := grpc.NewServer(Router.serverOptions(logger, jtrace.Tracer.GetTracer())...)

	// reflection for evans
	reflection.Register(baseServer)

	// register{MICRO}service
	pb.RegisterMicroServer(baseServer, &controller.Micro{})
	go func() {
		log.Fatalln(baseServer.Serve(lis))
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(context.Background(), config.Confs.Get().Service.GRPC.Host+config.Confs.Get().Service.GRPC.Port, Router.dialOptions()...)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
		return err
	}

	// new server from http package
	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	// register handler
	if err := pb.RegisterMicroHandler(context.Background(), gwmux, conn); err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(fmt.Sprintf("Failed to register gateway: %s", err.Error()))
		return err
	}

	// handle methods
	mux.Handle("/", gwmux)
	Router.HandleFuncs(mux)

	gwServer := &http.Server{
		Addr:    config.Confs.Get().Service.HTTP.Port,
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on ", config.Confs.Get().Service.HTTP.Port)
	log.Fatalln(gwServer.ListenAndServe())

	return nil
}

// defaultGRPCOptions
// add options for grpc connection
// In order to enable tracing of both upstream and downstream requests of the gRPC service, the gRPC client must also be initialized with client-side opentracing interceptor
// The parent spans created by the gRPC middleware are injected to the go context
func (a *micro) serverOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
	options := []grpc.ServerOption{}

	// UnaryInterceptor and OpenTracingServerInterceptor for tracer
	options = append(options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		grpc_auth.UnaryServerInterceptor(middleware.M.JWT),
		grpc_prometheus.UnaryServerInterceptor,
	),
	))

	options = append(options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_auth.StreamServerInterceptor(middleware.M.JWT),
		otgrpc.OpenTracingStreamServerInterceptor(tracer, otgrpc.LogPayloads()),
	)))

	return options
}

// dialOptions, options for dial connections
func (a *micro) dialOptions() []grpc.DialOption {
	options := []grpc.DialOption{}

	options = append(options, grpc.WithBlock())
	options = append(options, grpc.WithInsecure())

	return options
}

// HandleFuncs method for handler your basci methods
func (a *micro) HandleFuncs(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", controller.M.Metrics)
	mux.Handle("/health", middleware.M.MiddlewareExample(http.HandlerFunc(controller.M.Health)))
}
