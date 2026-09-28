package application

import (
	"io"
	"testing"
	"time"

	"github.com/jrsanchez/holidays-api/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// fakeRepository is a test double implementing ports.HolidayRepository.
type fakeRepository struct {
	holidays []domain.Holiday
}

func (f *fakeRepository) GetAll() []domain.Holiday {
	return f.holidays
}

func mustDate(s string) time.Time {
	t, err := time.Parse(domain.DateLayout, s)
	if err != nil {
		panic(err)
	}
	return t
}

func sampleHolidays() []domain.Holiday {
	return []domain.Holiday{
		{Date: mustDate("2026-01-01"), Title: "Año Nuevo", Type: "Civil", Inalienable: true, Extra: "Civil e Irrenunciable"},
		{Date: mustDate("2026-04-03"), Title: "Viernes Santo", Type: "Religioso", Inalienable: false, Extra: "Civil"},
		{Date: mustDate("2026-05-01"), Title: "Día del Trabajo", Type: "Civil", Inalienable: true, Extra: "Civil e Irrenunciable"},
		{Date: mustDate("2026-12-25"), Title: "Navidad", Type: "Religioso", Inalienable: true, Extra: "Civil e Irrenunciable"},
	}
}

func newTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger
}

func TestGetHolidays_NoFilter_ReturnsAll(t *testing.T) {
	repo := &fakeRepository{holidays: sampleHolidays()}
	svc := NewHolidayService(repo, newTestLogger())

	result := svc.GetHolidays(domain.HolidayFilter{})

	assert.Len(t, result, 4)
}

func TestGetHolidays_FilterByType_CaseInsensitive(t *testing.T) {
	repo := &fakeRepository{holidays: sampleHolidays()}
	svc := NewHolidayService(repo, newTestLogger())

	result := svc.GetHolidays(domain.HolidayFilter{Type: "civil"})

	assert.Len(t, result, 2)
	for _, h := range result {
		assert.Equal(t, "Civil", h.Type)
	}
}

func TestGetHolidays_FilterByDateRange(t *testing.T) {
	repo := &fakeRepository{holidays: sampleHolidays()}
	svc := NewHolidayService(repo, newTestLogger())

	start := mustDate("2026-02-01")
	end := mustDate("2026-11-30")

	result := svc.GetHolidays(domain.HolidayFilter{StartDate: &start, EndDate: &end})

	assert.Len(t, result, 2)
	assert.Equal(t, "Viernes Santo", result[0].Title)
	assert.Equal(t, "Día del Trabajo", result[1].Title)
}

func TestGetHolidays_FilterByTypeAndDateRange_NoMatch(t *testing.T) {
	repo := &fakeRepository{holidays: sampleHolidays()}
	svc := NewHolidayService(repo, newTestLogger())

	start := mustDate("2026-01-01")
	end := mustDate("2026-01-31")

	result := svc.GetHolidays(domain.HolidayFilter{Type: "Religioso", StartDate: &start, EndDate: &end})

	assert.Empty(t, result)
}

func TestGetHolidays_EmptyRepository(t *testing.T) {
	repo := &fakeRepository{holidays: []domain.Holiday{}}
	svc := NewHolidayService(repo, newTestLogger())

	result := svc.GetHolidays(domain.HolidayFilter{})

	assert.Empty(t, result)
}

