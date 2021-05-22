package grpc

import (
	"micro/app/router/middleware"
	"micro/client/jtrace"
	"micro/config"
	controller "micro/controller/grpc"
	zapLogger "micro/pkg/logger"
	pb "micro/proto"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (g micro) GetRouter() *grpc.Server {
	logger := zapLogger.GetZapLogger(config.Confs.GetDebug())
	options := g.defaultGRPCOptions(logger, jtrace.Tracer.GetTracer())

	baseServer := grpc.NewServer(options...)

	// reflection for evans
	reflection.Register(baseServer)

	pb.RegisterAuthServer(baseServer, new(controller.Micro))

	return baseServer
}

// defaultGRPCOptions
// add options for grpc connection
func (a *micro) defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
	options := []grpc.ServerOption{}

	// UnaryInterceptor and OpenTracingServerInterceptor for tracer
	options = append(options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		grpc_auth.UnaryServerInterceptor(middleware.M.JWT),
		grpc_prometheus.UnaryServerInterceptor,
	),
	))

	options = append(options, grpc.StreamInterceptor(
		grpc_auth.StreamServerInterceptor(middleware.M.JWT),
	))

	return options
}
