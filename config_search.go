package thinknum

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/mehiX/thinknumV2/internal/query"
)

// SearchDefinition Defines the parameters of a search.
type SearchDefinition struct {
	Name string `json:"name"`
	// Set this to true to ignore this search definition
	Disabled bool `json:"disabled"`
	// Path to a file where the results will be written
	// Since multiple formats are supported, this parameter should not have a type suffix.
	// The suffix will be added when the file is created
	OutputFile string `json:"output"`
	// Supported types: `json`, `csv`. Anything else will simply be ignored
	OutputTypes []string `json:"output_types"`
	DatasetID   string   `json:"dataset"`
	// A request object as defined by the Thinknum API Docs
	Request query.Request `json:"request"`
}

type timespan struct {
	Start time.Time
	End   time.Time
}

// Clone Performs a deep copy of the struct
// This solves an issue where the Filters array always points to the same memory location.
// I need to explicitely allocate a new memory location for the new struct
func (s SearchDefinition) Clone() SearchDefinition {
	newS := new(SearchDefinition)
	*newS = s
	newS.Request = s.Request.Clone()

	return *newS
}

// Split Split the current search definition into smaller time frames.
// It returns an array of search definitions, each having the same citeria as the original definition, plus a constraint on start and end time.
// The `interval` parameter is of time.Duration, therefor the largest avaialable time unit is `h` (hour). So to specify a week you should translate that in hours: 7 * 24h
func (s SearchDefinition) Split(from, to time.Time, interval time.Duration) []SearchDefinition {

	spans := splitTime(from, to, interval)

	searches := make([]SearchDefinition, len(spans))

	for index, span := range spans {
		ns := s.Clone()

		ns.Request.Filters = append(ns.Request.Filters,
			query.Filter{
				Column: filterColDateName,
				Type:   ">=",
				Value:  []string{span.Start.Format("2006-01-02")},
			},
			query.Filter{
				Column: filterColDateName,
				Type:   "<",
				Value:  []string{span.End.Format("2006-01-02")},
			})

		// each filter should write a different file
		ns.OutputFile = fmt.Sprintf("%s_%03d", ns.OutputFile, index)
		searches[index] = ns
	}

	return searches
}

// ReadSearchDefinition Reads a JSON object representing a SearchDefinition from an io.Reader.
// Returns a SearchDefinition or a decoding error if any
func ReadSearchDefinition(in io.Reader) (SearchDefinition, error) {

	var srch SearchDefinition

	err := json.NewDecoder(in).Decode(&srch)

	return srch, err
}

func splitTime(from, to time.Time, interval time.Duration) []timespan {
	s := make([]timespan, 0)

	var f1 time.Time

	for f1 = from; f1.Add(interval).Before(to); f1 = f1.Add(interval) {
		s = append(s, timespan{f1, f1.Add(interval)})
	}

	return append(s, timespan{f1, to})

}
