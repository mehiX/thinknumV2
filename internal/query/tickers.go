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

func TickerList(hostname, version, token, datasetID string) ([]TickerItem, error) {

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
