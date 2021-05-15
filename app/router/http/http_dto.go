package http

import (
	"micro/controller/http"

	"github.com/gin-gonic/gin"
)

func (a *micro) GetRouter() *gin.Engine {
	router := gin.Default()

	// metrics
	router.GET("/metrics", http.Controller.Metrics)

	// health check
	router.GET("/health", http.Controller.Health)

	return router
}
