package middleware

import (
	"context"
	"log"
	"net/http"
)

var (
	M Middleware = &middle{}
)

// Middleware interface
type Middleware interface {
	JWT(ctx context.Context) (context.Context, error)
	MiddlewareExample(next http.Handler) http.Handler
}

// middle struct
type middle struct{}

// JWT method
func (m *middle) JWT(ctx context.Context) (context.Context, error) {
	log.Println("Executing GRPC MiddlewareExample")
	return ctx, nil
}

func (m *middle) MiddlewareExample(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing HTTP MiddlewareExample")
		next.ServeHTTP(w, r)
	})
}
