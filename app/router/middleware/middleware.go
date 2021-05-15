package middleware

import "context"

var (
	M Middleware = &middle{}
)

// Middleware interface
type Middleware interface {
	JWT(ctx context.Context) (context.Context, error)
}

// middle struct
type middle struct{}

// JWT method
func (m *middle) JWT(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
