// Package application contains the use cases (business logic) of the
// Holidays API. It depends only on the domain package and on the ports
// (interfaces) it needs, never on concrete infrastructure implementations.
package application

import (
	"strings"
	"time"

	"github.com/jrsanchez/holidays-api/internal/domain"
	"github.com/jrsanchez/holidays-api/internal/domain/ports"
	"github.com/sirupsen/logrus"
)

// holidayService is the concrete implementation of ports.HolidayService.
type holidayService struct {
	repo   ports.HolidayRepository
	logger *logrus.Logger
}

// NewHolidayService builds a new application service given a holiday
// repository (secondary port) and a logger.
func NewHolidayService(repo ports.HolidayRepository, logger *logrus.Logger) ports.HolidayService {
	return &holidayService{repo: repo, logger: logger}
}

// GetHolidays returns the holidays known by the repository that satisfy the
// given filter. Filtering by type is case-insensitive; filtering by date
// range is inclusive on both ends.
func (s *holidayService) GetHolidays(filter domain.HolidayFilter) []domain.Holiday {
	all := s.repo.GetAll()
	result := make([]domain.Holiday, 0, len(all))

	for _, h := range all {
		if !matchesType(h, filter) {
			continue
		}
		if !matchesDateRange(h, filter) {
			continue
		}
		result = append(result, h)
	}

	s.logger.WithFields(logrus.Fields{
		"filter_type":       filter.Type,
		"filter_start_date": safeDate(filter.StartDate),
		"filter_end_date":   safeDate(filter.EndDate),
		"total_available":   len(all),
		"total_matched":     len(result),
	}).Info("holidays filtered")

	return result
}

func matchesType(h domain.Holiday, filter domain.HolidayFilter) bool {
	if filter.Type == "" {
		return true
	}
	return strings.EqualFold(h.Type, filter.Type)
}

func matchesDateRange(h domain.Holiday, filter domain.HolidayFilter) bool {
	if filter.StartDate != nil && h.Date.Before(*filter.StartDate) {
		return false
	}
	if filter.EndDate != nil && h.Date.After(*filter.EndDate) {
		return false
	}
	return true
}

func safeDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(domain.DateLayout)
}

