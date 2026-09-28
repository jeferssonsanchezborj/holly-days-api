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

	// Serve the OpenAPI specification and a lightweight Swagger UI page so
	// the contract can be explored directly from the running service.
	engine.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	engine.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})

	v1 := engine.Group("/api/v1")
	{
		v1.GET("/holidays", holidayHandler.GetHolidays)
	}

	return engine
}

// swaggerUIHTML embeds a minimal Swagger UI page pointed at the locally
// served OpenAPI document, avoiding extra runtime dependencies.
const swaggerUIHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>Holidays API - Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '/openapi.yaml',
          dom_id: '#swagger-ui',
        });
      };
    </script>
  </body>
</html>`

