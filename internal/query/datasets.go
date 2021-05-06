package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DatasetResponse The json response when querying for the list of datasets
type DatasetResponse struct {
	ResponseMetadata
	Items []DatasetItem
}

// DatasetItem Represents the data returned for a single Dataset when querying for the list of all datasets
type DatasetItem struct {
	State       string
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Summary     string
}

// RunSearch Run a query in the current dataset based on the passed in `Request` definiton
// `pageSize` defines the limit on the records to be returned
func (d DatasetItem) RunSearch(hostname, version, token string, pageSize int, srch Request) RunResult {

	var items RowsItems

	f := func(params url.Values) (ResponseMetadata, error) {
		URL := fmt.Sprintf("https://%s/connections/dataset/%s/query/new", hostname, d.ID)

		req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(params.Encode()))
		if err != nil {
			return ResponseMetadata{}, err
		}

		addRequestHeadersPOST(req, token, version)

		var resp *http.Response
		var statusCode int
		var errCount, maxErrCount = 0, 3

		for statusCode != http.StatusOK {

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				// allow maxErrCount retries on error, after which abort
				log.Printf("Error: %v\n", err)
				if errCount >= maxErrCount {
					return ResponseMetadata{}, err
				}

				errCount++

				fmt.Printf("Retry (%d/%d)...\n", errCount, maxErrCount)

				continue
			}

			statusCode = resp.StatusCode

			// in case of timeout we can try again
			// https://docs.thinknum.com/docs/query-api#http-response-status-code
			// When you get 504 error, you can keep retrying until data is returned. Every retries will connect to existing queued query and does not start new query.
			if statusCode != http.StatusOK && statusCode != http.StatusGatewayTimeout {
				b, _ := ioutil.ReadAll(resp.Body)
				return ResponseMetadata{}, fmt.Errorf("code: %d, body: %s", resp.StatusCode, string(b))
			}

			if statusCode == http.StatusGatewayTimeout {
				fmt.Printf("%s => request timeout. Retrying...\n", d.DisplayName)
			}
		}

		defer resp.Body.Close()

		var dsresp datasetBasicQueryResponse
		if err = json.NewDecoder(resp.Body).Decode(&dsresp); err != nil {
			return ResponseMetadata{}, err
		}

		// these are the fields metadata so we only need to save them once
		if len(items.Fields) == 0 {
			items.Fields = append(items.Fields, dsresp.Items.Fields...)
		}
		items.Rows = append(items.Rows, dsresp.Items.Rows...)
		items.Pages++
		// The Total should be the same with each request so no problem re-writing the value
		if dsresp.Total > 0 {
			items.Total = dsresp.Total
		}

		return dsresp.ResponseMetadata, nil
	}

	frm := url.Values{}
	paramsStr, err := json.Marshal(srch)
	if err != nil {
		return RunResult{RowsItems{}, err}
	}
	frm["request"] = []string{string(paramsStr)}
	frm["limit"] = []string{strconv.Itoa(pageSize)}
	frm["start"] = []string{"0"}

	err = fetchAll(f, frm)

	return RunResult{items, err}
}

// Datasets Query the list of datasets
func Datasets(hostname, version, token string, tickerFilter string) ([]DatasetItem, error) {

	URL := fmt.Sprintf("https://%s/connections/datasets", hostname)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	if tickerFilter != "" {
		v := req.URL.Query()
		v.Set("ticker", tickerFilter)
		req.URL.RawQuery = v.Encode()
	}

	addRequestHeaders(req, token, version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dsResp DatasetResponse
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return nil, err
	}

	return dsResp.Items, nil

}

type datasetBasicQueryResponse struct {
	ResponseMetadata
	Items RowsItems
}

// RowsItems Query results as returned by the Thinknum API
type RowsItems struct {
	RowItemsMetadata
	// Metadata for the returned columns
	Fields []Field
	// Rows of effective data
	Rows []Row
}

// Field Metadata for one column of the results
type Field struct {
	DisplayName string   `json:"display_name"`
	Format      string   `json:"format"`
	Metric      bool     `json:"metric"`
	ID          string   `json:"id"`
	Length      int      `json:"length"`
	Summary     string   `json:"summary"`
	Type        string   `json:"type"`
	Options     []string `json:"options"`
}

// Row One row of data
type Row []interface{}

// RowItemsMetadata Metadata for returned rows
type RowItemsMetadata struct {
	Total int
	// The number of requests made to retrieve all the results.
	// The actual value depends on the pageSize ('limit') set through configuration
	Pages int
}
