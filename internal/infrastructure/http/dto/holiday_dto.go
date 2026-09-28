// Package dto contains the Data Transfer Objects used by the HTTP adapter.
// They are responsible for the wire representation (JSON and XML) of the
// domain entities, keeping serialization concerns out of the domain layer.
package dto

import (
	"encoding/xml"

	"github.com/jrsanchez/holidays-api/internal/domain"
)

// HolidayDTO is the wire representation of a domain.Holiday. It exposes the
// fields required by the challenge: Fecha (date), título (title), teléfono
// (phone), tipo (type), inalienable and extra.
type HolidayDTO struct {
	XMLName     xml.Name `json:"-" xml:"holiday"`
	Date        string   `json:"date" xml:"date" example:"2026-01-01"`
	Title       string   `json:"title" xml:"title" example:"Año Nuevo"`
	Phone       string   `json:"phone" xml:"phone" example:""`
	Type        string   `json:"type" xml:"type" example:"Civil"`
	Inalienable bool     `json:"inalienable" xml:"inalienable" example:"true"`
	Extra       string   `json:"extra" xml:"extra" example:"Civil e Irrenunciable"`
}

// HolidaysResponse is the top level wrapper returned by the /holidays
// endpoint. XML requires a single root element, so this wrapper is used for
// both JSON and XML representations for consistency.
type HolidaysResponse struct {
	XMLName xml.Name     `json:"-" xml:"holidays"`
	Total   int          `json:"total" xml:"total,attr"`
	Data    []HolidayDTO `json:"data" xml:"holiday"`
}

// FromDomain converts a slice of domain.Holiday into a HolidaysResponse
// ready to be serialized as JSON or XML.
func FromDomain(holidays []domain.Holiday) HolidaysResponse {
	items := make([]HolidayDTO, 0, len(holidays))
	for _, h := range holidays {
		items = append(items, HolidayDTO{
			Date:        h.Date.Format(domain.DateLayout),
			Title:       h.Title,
			Phone:       h.Phone,
			Type:        h.Type,
			Inalienable: h.Inalienable,
			Extra:       h.Extra,
		})
	}
	return HolidaysResponse{Total: len(items), Data: items}
}

// ErrorResponse is the standard error payload returned by the API.
type ErrorResponse struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Message string   `json:"message" xml:"message"`
}

