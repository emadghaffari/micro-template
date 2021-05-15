package grpc

var (
	Router Micro = &micro{}
)

// Micro service
type micro struct{}

type Micro interface {
}
