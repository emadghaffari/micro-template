package grpc

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	Router Micro = &micro{}
)

// Micro service
type micro struct{}

type Micro interface {
	GetRouter() *grpc.Server
	defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption
}
