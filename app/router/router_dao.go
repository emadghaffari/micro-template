package router

import (
	"net/http"

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
	GetRouter() error
	serverOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption
	dialOptions() []grpc.DialOption
	HandleFuncs(mux *http.ServeMux)
}
