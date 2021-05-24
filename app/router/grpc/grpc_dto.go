package grpc

import (
	"micro/app/router/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (g micro) GetRouter() *grpc.Server {
	// logger := zapLogger.GetZapLogger(config.Confs.GetDebug())
	// options := g.defaultGRPCOptions(logger, jtrace.Tracer.GetTracer())

	// baseServer := grpc.NewServer(options...)

	// // reflection for evans
	// reflection.Register(baseServer)

	// pb.RegisterAuthServer(baseServer, new(controller.Micro))

	// return baseServer
	return nil
}

// defaultGRPCOptions
// add options for grpc connection
// In order to enable tracing of both upstream and downstream requests of the gRPC service, the gRPC client must also be initialized with client-side opentracing interceptor
// The parent spans created by the gRPC middleware are injected to the go context
func (a *micro) defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
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
