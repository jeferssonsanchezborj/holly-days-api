// Package router wires together the HTTP routes exposed by the API,
// attaching middleware and handlers.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/handler"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/middleware"
	"github.com/sirupsen/logrus"
)

// New builds a fully configured Gin engine with all routes registered.
func New(holidayHandler *handler.HolidayHandler, logger *logrus.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(logger))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Serve the OpenAPI 3.0.3 specification plus a fully self-hosted Swagger
	// UI (assets vendored under docs/swagger-ui, no CDN/internet dependency)
	// so the API contract can be explored interactively at GET /docs.
	engine.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	engine.Static("/docs", "./docs/swagger-ui")

	v1 := engine.Group("/api/v1")
	{
		v1.GET("/holidays", holidayHandler.GetHolidays)
	}

	return engine
}


