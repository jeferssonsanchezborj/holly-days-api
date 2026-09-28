package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger
}

func TestNewBoostrHolidayRepository_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": [
				{"date": "2026-01-01", "title": "Año Nuevo", "type": "Civil", "inalienable": true, "extra": "Civil e Irrenunciable"},
				{"date": "2026-04-03", "title": "Viernes Santo", "type": "Religioso", "inalienable": false, "extra": "Civil"}
			]
		}`))
	}))
	defer server.Close()

	repo, err := NewBoostrHolidayRepository(server.Client(), server.URL, testLogger())

	require.NoError(t, err)
	holidays := repo.GetAll()
	assert.Len(t, holidays, 2)
	assert.Equal(t, "Año Nuevo", holidays[0].Title)
	assert.Equal(t, "Civil", holidays[0].Type)
	assert.True(t, holidays[0].Inalienable)
}

func TestNewBoostrHolidayRepository_SkipsInvalidDates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{
			"status": "success",
			"data": [
				{"date": "not-a-date", "title": "Broken", "type": "Civil", "inalienable": false, "extra": ""},
				{"date": "2026-05-01", "title": "Día del Trabajo", "type": "Civil", "inalienable": true, "extra": "Civil e Irrenunciable"}
			]
		}`))
	}))
	defer server.Close()

	repo, err := NewBoostrHolidayRepository(server.Client(), server.URL, testLogger())

	require.NoError(t, err)
	holidays := repo.GetAll()
	require.Len(t, holidays, 1)
	assert.Equal(t, "Día del Trabajo", holidays[0].Title)
}

func TestNewBoostrHolidayRepository_UpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := NewBoostrHolidayRepository(server.Client(), server.URL, testLogger())

	assert.Error(t, err)
}

func TestNewBoostrHolidayRepository_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	_, err := NewBoostrHolidayRepository(server.Client(), server.URL, testLogger())

	assert.Error(t, err)
}

