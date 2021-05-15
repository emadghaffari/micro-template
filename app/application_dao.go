package app

import (
	"io"

	group "github.com/oklog/oklog/pkg/group"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger *zap.Logger
	Base   Application = &App{}
)

// Application interface for start application
type Application interface {
	StartApplication()
	initLogger()
	initConfigs() error
	initGRPCHandler(g *group.Group)
	initHTTPEndpoint(g *group.Group)
	initCancelInterrupt(g *group.Group)
	initJaeger() (io.Closer, error)
	initConfigServer() error
	initMessageBroker() error
	initRedis() error
	createService() (g *group.Group)
	defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption
}

type App struct{}
