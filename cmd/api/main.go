// Command api is the entry point of the Holidays API service.
//
// It follows Hexagonal Architecture: infrastructure adapters (HTTP,
// external API client) are wired here into the application's core, which
// remains completely decoupled from transport or delivery details.
//
// IMPORTANT: the external holidays data source (api.boostr.cl) is called
// exactly once, during this bootstrap phase, and cached in memory for the
// lifetime of the process. Subsequent requests to this service are served
// entirely from that in-memory cache.
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/jrsanchez/holidays-api/internal/application"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/handler"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/router"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/repository"
	"github.com/jrsanchez/holidays-api/pkg/logger"
)

func main() {
	log := logger.New()

	upstreamURL := os.Getenv("HOLIDAYS_UPSTREAM_URL")
	if upstreamURL == "" {
		upstreamURL = repository.DefaultBoostrURL
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	// Secondary adapter: performs the ONE-TIME call to the upstream service
	// and caches the result in memory for the whole application lifetime.
	holidayRepo, err := repository.NewBoostrHolidayRepository(httpClient, upstreamURL, log)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize holidays repository from upstream service")
	}

	// Application/use-case layer.
	holidayService := application.NewHolidayService(holidayRepo, log)

	// Primary adapter: HTTP handler + router.
	holidayHandler := handler.NewHolidayHandler(holidayService, log)
	engine := router.New(holidayHandler, log)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.WithField("port", port).Info("starting Holidays API server")
	if err := engine.Run(":" + port); err != nil {
		log.WithError(err).Fatal("server stopped unexpectedly")
	}
}

