package app

import (
	"io"

	group "github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
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
	initPostgres() error
}

type App struct{}
