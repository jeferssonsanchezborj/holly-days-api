// Package ports defines the boundaries (interfaces) between the domain
// core and the outside world, following the Hexagonal Architecture pattern.
//
// - HolidayRepository is a "driven" (secondary) port: it is implemented by
//   an infrastructure adapter that knows how to obtain holidays data.
// - HolidayService is a "driving" (primary) port: it is implemented by the
//   application layer and consumed by infrastructure adapters (e.g. HTTP
//   handlers).
package ports

import "github.com/jrsanchez/holidays-api/internal/domain"

// HolidayRepository is the secondary port used by the application layer to
// retrieve all known holidays. Implementations are responsible for the
// actual data source (an HTTP API, a database, an in-memory cache, etc.).
type HolidayRepository interface {
	// GetAll returns every holiday currently held by the repository.
	GetAll() []domain.Holiday
}

// HolidayService is the primary port exposed by the application layer.
// HTTP (or any other) adapters depend on this interface instead of a
// concrete implementation, keeping the core logic isolated and testable.
type HolidayService interface {
	// GetHolidays returns the holidays that match the given filter.
	GetHolidays(filter domain.HolidayFilter) []domain.Holiday
}

