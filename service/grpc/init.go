package grpc

var (
	Service Micro = &micro{}
)

// micro service
type micro struct{}

type Micro interface{}
