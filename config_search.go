package thinknum

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type SearchDefinition struct {
	Name        string       `json:"name"`
	Disabled    bool         `json:"disabled"`
	OutputFile  string       `json:"output"`
	OutputTypes []string     `json:"output_types"`
	DatasetID   string       `json:"dataset"`
	Request     QueryRequest `json:"request"`
}

type QueryRequest struct {
	Filters []Filter `json:"filters,omitempty"`
	Tickers []string `json:"tickers,omitempty"`
}

type Filter struct {
	Column string   `json:"column"`
	Type   string   `json:"type"`
	Value  []string `json:"value"`
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
	newS.Request.Filters = make([]Filter, len(s.Request.Filters))
	for i := range s.Request.Filters {
		newS.Request.Filters[i] = s.Request.Filters[i]
	}

	newS.Request.Tickers = make([]string, len(s.Request.Tickers))
	for i := range s.Request.Tickers {
		newS.Request.Tickers[i] = s.Request.Tickers[i]
	}

	return *newS
}

func ReadSearch(in io.ReadCloser) (SearchDefinition, error) {

	var srch SearchDefinition

	err := json.NewDecoder(in).Decode(&srch)

	return srch, err
}

func SearchApplyDatesFilter(srch SearchDefinition, from, to time.Time, interval time.Duration) []SearchDefinition {

	spans := splitTime(from, to, interval)

	searches := make([]SearchDefinition, len(spans))

	for index, span := range spans {
		ns := srch.Clone()

		ns.Request.Filters = append(ns.Request.Filters,
			Filter{
				Column: filterColDateName,
				Type:   ">=",
				Value:  []string{span.Start.Format("2006-01-02")},
			},
			Filter{
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

func splitTime(from, to time.Time, interval time.Duration) []Timespan {
	s := make([]Timespan, 0)

	var f1 time.Time

	for f1 = from; f1.Add(interval).Before(to); f1 = f1.Add(interval) {
		s = append(s, Timespan{f1, f1.Add(interval)})
	}

	return append(s, Timespan{f1, to})

}
