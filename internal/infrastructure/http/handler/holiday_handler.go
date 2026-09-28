// Package handler contains the primary (driving) HTTP adapters. They
// translate incoming HTTP requests into calls to the application layer
// (ports.HolidayService) and render the results as JSON or XML depending on
// content negotiation.
package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jrsanchez/holidays-api/internal/domain"
	"github.com/jrsanchez/holidays-api/internal/domain/ports"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/dto"
	"github.com/sirupsen/logrus"
)

// HolidayHandler exposes HTTP endpoints backed by a ports.HolidayService.
type HolidayHandler struct {
	service ports.HolidayService
	logger  *logrus.Logger
}

// NewHolidayHandler builds a new HolidayHandler.
func NewHolidayHandler(service ports.HolidayService, logger *logrus.Logger) *HolidayHandler {
	return &HolidayHandler{service: service, logger: logger}
}

// GetHolidays handles GET /api/v1/holidays.
//
// Query parameters:
//   - type:       optional, filters by holiday type (e.g. "Civil", "Religioso"). Case-insensitive.
//   - start_date: optional, ISO date (YYYY-MM-DD), lower bound of the date range (inclusive).
//   - end_date:   optional, ISO date (YYYY-MM-DD), upper bound of the date range (inclusive).
//
// Content negotiation:
//   - Accept: application/xml  -> renders XML
//   - Accept: application/json or anything else (default) -> renders JSON
func (h *HolidayHandler) GetHolidays(c *gin.Context) {
	filter, err := parseFilter(c)
	if err != nil {
		h.renderError(c, http.StatusBadRequest, err.Error())
		return
	}

	holidays := h.service.GetHolidays(filter)
	response := dto.FromDomain(holidays)

	h.render(c, http.StatusOK, response)
}

// parseFilter extracts and validates the domain.HolidayFilter from the
// incoming request's query parameters.
func parseFilter(c *gin.Context) (domain.HolidayFilter, error) {
	filter := domain.HolidayFilter{
		Type: strings.TrimSpace(c.Query("type")),
	}

	if raw := strings.TrimSpace(c.Query("start_date")); raw != "" {
		t, err := time.Parse(domain.DateLayout, raw)
		if err != nil {
			return domain.HolidayFilter{}, invalidDateError("start_date", raw)
		}
		filter.StartDate = &t
	}

	if raw := strings.TrimSpace(c.Query("end_date")); raw != "" {
		t, err := time.Parse(domain.DateLayout, raw)
		if err != nil {
			return domain.HolidayFilter{}, invalidDateError("end_date", raw)
		}
		filter.EndDate = &t
	}

	return filter, nil
}

func invalidDateError(field, value string) error {
	return &invalidParamError{field: field, value: value}
}

// invalidParamError represents a validation error for a query parameter.
type invalidParamError struct {
	field string
	value string
}

func (e *invalidParamError) Error() string {
	return "invalid value for '" + e.field + "': '" + e.value + "', expected format YYYY-MM-DD"
}

// render writes the response body using either XML or JSON depending on the
// request's Accept header. JSON is the default when the header is missing
// or not explicitly "application/xml".
func (h *HolidayHandler) render(c *gin.Context, status int, payload interface{}) {
	accept := c.GetHeader("Accept")
	if strings.Contains(accept, "application/xml") {
		c.XML(status, payload)
		return
	}
	c.JSON(status, payload)
}

// renderError writes an error response honoring content negotiation.
func (h *HolidayHandler) renderError(c *gin.Context, status int, message string) {
	h.render(c, status, dto.ErrorResponse{Message: message})
}

