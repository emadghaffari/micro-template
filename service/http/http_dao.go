package http

var (
	Service Micro = &micro{}
)

type Micro interface{}

// micro service
type micro struct{}
