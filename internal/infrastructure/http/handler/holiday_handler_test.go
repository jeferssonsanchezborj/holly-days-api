package handler

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jrsanchez/holidays-api/internal/domain"
	"github.com/jrsanchez/holidays-api/internal/infrastructure/http/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeService is a test double implementing ports.HolidayService.
type fakeService struct {
	holidays []domain.Holiday
}

func (f *fakeService) GetHolidays(filter domain.HolidayFilter) []domain.Holiday {
	return f.holidays
}

func testLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger
}

func setupRouter(svc *fakeService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHolidayHandler(svc, testLogger())
	r.GET("/api/v1/holidays", h.GetHolidays)
	return r
}

func sampleHoliday() domain.Holiday {
	date, _ := time.Parse(domain.DateLayout, "2026-01-01")
	return domain.Holiday{
		Date: date, Title: "Año Nuevo", Phone: "", Type: "Civil",
		Inalienable: true, Extra: "Civil e Irrenunciable",
	}
}

func TestGetHolidays_JSON_Default(t *testing.T) {
	svc := &fakeService{holidays: []domain.Holiday{sampleHoliday()}}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/holidays", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var body dto.HolidaysResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 1, body.Total)
	assert.Equal(t, "Año Nuevo", body.Data[0].Title)
	assert.Equal(t, "Civil", body.Data[0].Type)
	assert.True(t, body.Data[0].Inalienable)
}

func TestGetHolidays_XML_WhenRequested(t *testing.T) {
	svc := &fakeService{holidays: []domain.Holiday{sampleHoliday()}}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/holidays", nil)
	req.Header.Set("Accept", "application/xml")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/xml")

	var body dto.HolidaysResponse
	require.NoError(t, xml.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "Año Nuevo", body.Data[0].Title)
}

func TestGetHolidays_InvalidStartDate_ReturnsBadRequest(t *testing.T) {
	svc := &fakeService{}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/holidays?start_date=not-a-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetHolidays_EmptyResult(t *testing.T) {
	svc := &fakeService{holidays: []domain.Holiday{}}
	router := setupRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/holidays?type=Inexistente", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body dto.HolidaysResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Total)
}



