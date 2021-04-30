package query

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TickerResponse struct {
	ResponseMetadata
	Items []TickerItem
}

type TickerItem struct {
	ID          string
	Sector      string
	Country     string
	Industry    string
	DisplayName string `json:"display_name"`
}

// TickerList Returns the list of tickers for the provided `datasetID`
func TickerList(hostname, version, token, datasetID string) ([]TickerItem, error) {

	if datasetID == "" {
		return nil, fmt.Errorf("dataset not specified when requesting the list of tickers")
	}

	URL := fmt.Sprintf("https://%s/connections/dataset/%s/tickers", hostname, datasetID)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	addRequestHeaders(req, token, version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickerResp TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return nil, err
	}

	log.Println("Count", tickerResp.Count)
	log.Println("Total", tickerResp.Total)

	return tickerResp.Items, nil
}
