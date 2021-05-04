package query

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Request A request deinition as defined by the Thinknum API Docs
type Request struct {
	Filters     []Filter `json:"filters,omitempty"`
	Tickers     []string `json:"tickers,omitempty"`
	Pointintime bool     `json:"pointintime,omitempty"`
}

// Clone Returns a copy of the original Request. This makes sure that also the internal slices are properly duplicated to a new memory address
func (r Request) Clone() Request {
	var newR Request

	newR.Filters = make([]Filter, len(r.Filters))
	for i := range r.Filters {
		newR.Filters[i] = r.Filters[i]
	}

	newR.Tickers = make([]string, len(r.Tickers))
	for i := range r.Tickers {
		newR.Tickers[i] = r.Tickers[i]
	}

	newR.Pointintime = r.Pointintime

	return newR
}

// Filter A single filter used to filter data from a dataset
type Filter struct {
	Column string   `json:"column"`
	Type   string   `json:"type"`
	Value  []string `json:"value"`
}

// ResponseMetadata Metadata returned with the response
type ResponseMetadata struct {
	// Number of records returned. If less than the specified limit (page size) then no more records
	Count int `json:"count"`
	// Total number of records that can be retrieved for this search
	Total       int    `json:"total"`
	Status      int    `json:"status"`
	Summary     string `json:"summary"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// addRequestHeadersPOST Set the correct Content-Type and then add the authorization headers by calling `addRequestHeaders`
func addRequestHeadersPOST(r *http.Request, token, version string) {
	addRequestHeaders(r, token, version)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

// addRequestHeaders Add the necessary authorization headers
func addRequestHeaders(r *http.Request, token, version string) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	r.Header.Set("X-API-Version", version)
}

// fetchAll handles pagination. `processResp` is a function that contains the logic for sending the request, receiving the response and persisting the data (usually appending it to a slice)
// Fetch each new page by calling `processResp` and advance to the next page based on the response metadata
func fetchAll(processResp func(url.Values) (ResponseMetadata, error), params url.Values) error {

	for {
		start, err := strconv.Atoi(params.Get("start"))
		if err != nil {
			return err
		}

		resp, err := processResp(params)
		if err != nil || resp.Total <= resp.Count+start {
			return err
		}
		params.Set("start", strconv.Itoa(start+resp.Count))
	}
}

// RunResult Contains the data returned by a search
// The results are intended to be passed to the next step in the pipeline: persisting the results
type RunResult struct {
	Data  RowsItems
	Error error
}
