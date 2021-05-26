package middleware

import (
	"context"
	"fmt"
	"log"
	"micro/config"
	"net/http"
	"reflect"

	"google.golang.org/grpc"
)

var (
	M Middleware = &middle{}
)

// Middleware interface
type Middleware interface {
	JWT(ctx context.Context) (context.Context, error)
	MiddlewareExample(next http.Handler) http.Handler
	assignMiddleware(ctx context.Context, req interface{}, middlewares []string) error
}

// middle struct
type middle struct{}

// JWT method
func (m *middle) JWT(ctx context.Context) (context.Context, error) {
	log.Println("Executing GRPC MiddlewareExample")
	return ctx, nil
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//skip auth when ListReleases requested
	fmt.Println("=====================")
	// loop for all routes we have in config file
	for _, r := range config.Confs.Get().Service.Router {
		// if method name != proto rpc name, then go to next method
		if r.Method != info.FullMethod {
			continue
		}
		if err := M.assignMiddleware(ctx, req, r.Middlewares); err != nil {
			return nil, err
		}
	}
	fmt.Println("=====================")

	h, err := handler(ctx, req)

	return h, err
}

func (m *middle) assignMiddleware(ctx context.Context, req interface{}, middlewares []string) error {
	for _, m := range middlewares {
		method := reflect.ValueOf(&middle{}).MethodByName(m)
		if !method.IsValid() {
			continue
		}
		// FIXME
		// fix the return from method call,
		// we need to check errors for middlewares
		method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
	}

	return nil
}

func (m *middle) CheckSome(ctx context.Context, req interface{}) error {
	fmt.Println(ctx, req)
	fmt.Println("CheckSome")
}

func (m *middle) MiddlewareExample(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing HTTP MiddlewareExample")
		next.ServeHTTP(w, r)
	})
}
