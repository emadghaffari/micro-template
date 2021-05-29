package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// CheckSome example middleware for grpc calls
func (m *middle) CheckSome(ctx context.Context, req interface{}) error {
	return nil
}

// Middleware3 example middleware for grpc calls
func (m *middle) Middleware3(ctx context.Context, req interface{}) error {
	return fmt.Errorf("error in check middleware")
}

// global middleware example for all routes
func (m *middle) JWT(ctx context.Context) (context.Context, error) {
	log.Println("Executing GRPC MiddlewareExample")
	return ctx, nil
}

// MiddlewareExample example middleware for extra http calls
func (m *middle) MiddlewareExample(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing HTTP MiddlewareExample")
		next.ServeHTTP(w, r)
	})
}
