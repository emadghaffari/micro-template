package jwt

import (
	"context"
	"micro/model"

	"go.uber.org/zap"
)

var (
	// JWT variable instance of intef
	Model  Micro = &micro{}
	logger *zap.Logger
)

// jwt meths interface
type Micro interface {
	Generate(ctx context.Context, model interface{}) (*model.JWT, error)
	GenerateJWT() (*model.JWT, error)
	genRefJWT(td *model.JWT) error
	store(ctx context.Context, model interface{}, td *model.JWT) error
	Get(ctx context.Context, token string, response interface{}) error
	Verify(tk string) (string, error)
}

// jwt struct
type micro struct{}
