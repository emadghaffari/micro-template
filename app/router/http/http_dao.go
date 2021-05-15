package http

import "github.com/gin-gonic/gin"

var (
	Router Micro = &micro{}
)

// Micro service
type micro struct{}

type Micro interface {
	GetRouter() *gin.Engine
}
