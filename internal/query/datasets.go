package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DatasetResponse struct {
	ResponseMetadata
	Items []DatasetItem
}

type DatasetItem struct {
	State       string
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Summary     string
}

func (d DatasetItem) RunSearch(hostname, version, token string, pageSize int, srch Request) (RowsItems, error) {

	var items RowsItems

	f := func(params url.Values) (ResponseMetadata, error) {
		URL := fmt.Sprintf("https://%s/connections/dataset/%s/query/new", hostname, d.ID)

		req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(params.Encode()))
		if err != nil {
			return ResponseMetadata{}, err
		}

		addRequestHeadersPOST(req, token, version)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return ResponseMetadata{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(resp.Body)
			return ResponseMetadata{}, fmt.Errorf("code: %d, body: %s", resp.StatusCode, string(b))
		}

		var dsresp DatasetBasicQueryResponse
		if err := json.NewDecoder(resp.Body).Decode(&dsresp); err != nil {
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
		return RowsItems{}, err
	}
	frm["request"] = []string{string(paramsStr)}
	frm["limit"] = []string{strconv.Itoa(pageSize)}
	frm["start"] = []string{"0"}

	err = fetchAll(f, frm)

	return items, err
}

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

type DatasetBasicQueryResponse struct {
	ResponseMetadata
	Items RowsItems
}

type RowsItems struct {
	RowItemsMetadata
	Fields []interface{}
	Rows   []interface{}
}

type RowItemsMetadata struct {
	Total int
	Pages int
}
