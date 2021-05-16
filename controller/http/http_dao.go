package http

import "github.com/gin-gonic/gin"

var (
	Controller Micro = &micro{}
)

type Micro interface {
	Health(c *gin.Context)
	Metrics(c *gin.Context)
}

// micro service
type micro struct{}
