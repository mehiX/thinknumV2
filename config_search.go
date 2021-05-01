package thinknum

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/mehiX/thinknumV2/internal/query"
)

type SearchDefinition struct {
	Name        string        `json:"name"`
	Disabled    bool          `json:"disabled"`
	OutputFile  string        `json:"output"`
	OutputTypes []string      `json:"output_types"`
	DatasetID   string        `json:"dataset"`
	Request     query.Request `json:"request"`
}

type Timespan struct {
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

func ReadSearchDefinition(in io.Reader) (SearchDefinition, error) {

	var srch SearchDefinition

	err := json.NewDecoder(in).Decode(&srch)

	return srch, err
}

func splitTime(from, to time.Time, interval time.Duration) []Timespan {
	s := make([]Timespan, 0)

	var f1 time.Time

	for f1 = from; f1.Add(interval).Before(to); f1 = f1.Add(interval) {
		s = append(s, Timespan{f1, f1.Add(interval)})
	}

	return append(s, Timespan{f1, to})

}
