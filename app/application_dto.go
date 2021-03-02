package app

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	group "github.com/oklog/oklog/pkg/group"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func createService() (g *group.Group) {
	g = &group.Group{}

	// init GRPC Handlers
	initGRPCHandler(g)
	return g
}

func defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
	options := []grpc.ServerOption{}

	// UnaryInterceptor and OpenTracingServerInterceptor for tracer
	options = append(options, grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
	))
	return options
}
