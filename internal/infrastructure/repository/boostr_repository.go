// Package repository contains secondary (driven) adapters that implement
// the ports.HolidayRepository interface defined in the domain layer.
package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jrsanchez/holidays-api/internal/domain"
	"github.com/sirupsen/logrus"
)

// DefaultBoostrURL is the upstream endpoint that provides the raw Chilean
// holidays data used to populate this service's in-memory cache.
const DefaultBoostrURL = "https://api.boostr.cl/holidays.json"

// boostrResponse mirrors the JSON payload returned by api.boostr.cl.
type boostrResponse struct {
	Status string          `json:"status"`
	Data   []boostrHoliday `json:"data"`
}

// boostrHoliday mirrors a single holiday entry as returned by the upstream
// API. Note the upstream service does not provide a "phone" field.
type boostrHoliday struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Inalienable bool   `json:"inalienable"`
	Extra       string `json:"extra"`
}

// BoostrHolidayRepository is an in-memory implementation of
// ports.HolidayRepository. It fetches the full holiday list exactly once
// (at construction time, i.e. during application bootstrap) and serves all
// subsequent reads from memory, so the upstream service is never called
// again per incoming request.
type BoostrHolidayRepository struct {
	holidays []domain.Holiday
}

// NewBoostrHolidayRepository builds a BoostrHolidayRepository, performing a
// single HTTP call to the given URL (defaults to DefaultBoostrURL if empty)
// and parsing/caching the resulting holidays in memory.
//
// This function is intended to be called once during application startup
// (see cmd/api/main.go). If the upstream call fails, an error is returned so
// the caller can decide whether to abort startup or fall back to an empty
// cache.
func NewBoostrHolidayRepository(httpClient *http.Client, url string, logger *logrus.Logger) (*BoostrHolidayRepository, error) {
	if url == "" {
		url = DefaultBoostrURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	logger.WithField("url", url).Info("fetching holidays from upstream service (one-time bootstrap call)")

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("calling upstream holidays service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream holidays service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading upstream response body: %w", err)
	}

	var parsed boostrResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("parsing upstream response JSON: %w", err)
	}

	holidays := make([]domain.Holiday, 0, len(parsed.Data))
	for _, raw := range parsed.Data {
		date, err := time.Parse(domain.DateLayout, raw.Date)
		if err != nil {
			logger.WithFields(logrus.Fields{"date": raw.Date, "title": raw.Title}).
				Warn("skipping holiday with unparseable date")
			continue
		}
		holidays = append(holidays, domain.Holiday{
			Date:        date,
			Title:       raw.Title,
			Phone:       "", // not provided by the upstream API
			Type:        raw.Type,
			Inalienable: raw.Inalienable,
			Extra:       raw.Extra,
		})
	}

	logger.WithField("count", len(holidays)).Info("holidays successfully cached in memory")

	return &BoostrHolidayRepository{holidays: holidays}, nil
}

// GetAll returns every cached holiday. It never triggers a network call.
func (r *BoostrHolidayRepository) GetAll() []domain.Holiday {
	return r.holidays
}

