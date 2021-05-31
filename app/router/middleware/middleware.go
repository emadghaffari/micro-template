package middleware

import (
	"context"
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
	MiddlewareUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	MiddlewareStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

// middle struct
type middle struct{}

func (m *middle) MiddlewareUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

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

	h, err := handler(ctx, req)

	return h, err
}

// MiddlewareStreamInterceptor in thos middleware you can check
// server streaming, client streaming in info
func (m *middle) MiddlewareStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return nil
}

func (m *middle) assignMiddleware(ctx context.Context, req interface{}, middlewares []string) error {

	// loop middlewares for every route
	for _, m := range middlewares {

		// get middleware methods by name
		method := reflect.ValueOf(&middle{}).MethodByName(m)
		if !method.IsValid() {
			continue
		}

		// check every middleware for method
		responses := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
		if err := responses[0].Interface(); err != nil {
			return err.(error)
		}
	}

	return nil
}
