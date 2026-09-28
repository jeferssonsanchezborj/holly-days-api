// Package domain contains the core business entities of the Holidays API.
// This package has zero dependencies on frameworks, databases or transport
// mechanisms, following the Hexagonal Architecture (Ports & Adapters) style.
package domain

import "time"

// Holiday represents a Chilean holiday as defined by the business domain.
// It is intentionally decoupled from any external API or transport (JSON/XML)
// representation.
type Holiday struct {
	// Date is the calendar date on which the holiday occurs.
	Date time.Time
	// Title is the descriptive name of the holiday.
	Title string
	// Phone is a contact phone number associated with the holiday information
	// source. The upstream data provider (api.boostr.cl) does not currently
	// expose this field; it is kept as part of the domain model to satisfy
	// the exposed contract and defaults to an empty string when unavailable.
	Phone string
	// Type classifies the holiday, e.g. "Civil" or "Religioso".
	Type string
	// Inalienable indicates whether the holiday right cannot be waived
	// (in Spanish: "irrenunciable").
	Inalienable bool
	// Extra contains additional free-form information about the holiday.
	Extra string
}

// HolidayFilter groups all the optional criteria that can be used to narrow
// down the list of holidays returned by the application.
type HolidayFilter struct {
	// Type filters holidays by their type (e.g. "Civil", "Religioso").
	// Comparison is case-insensitive. Empty value means "no filter".
	Type string
	// StartDate, when not nil, excludes holidays occurring before this date.
	StartDate *time.Time
	// EndDate, when not nil, excludes holidays occurring after this date.
	EndDate *time.Time
}

// DateLayout is the canonical date format used across the application
// (matches the format returned by the upstream api.boostr.cl service).
const DateLayout = "2006-01-02"

