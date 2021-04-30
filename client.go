package thinknum

import (
	"fmt"

	"github.com/mehiX/thinknumV2/internal/query"
)

type Client struct {
	Config
	Token string
}

// NewClient Returns a new client for the Thinknum API. It will contain a valid token based on the received credentials
func NewClientFromJSON(configFile string) (*Client, error) {

	cfg, err := ConfigFromJSON(configFile)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(cfg.Version, cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}

	return NewClient(cfg, token), nil
}

func NewClient(cfg *Config, token *AuthToken) *Client {
	return &Client{
		Config: *cfg,
		Token:  token.Token,
	}

}

// Datasets Get a list of available datasets
// If a tickerID is provided (is not empty) then it is used to filter the datasets
func (c *Client) Datasets(tickerID string) ([]query.DatasetItem, error) {
	return query.Datasets(c.Hostname, c.Version, c.Token, tickerID)
}

// Tickers Get the list of tickers for the provided `datasetID`
func (c *Client) Tickers(datasetID string) ([]query.TickerItem, error) {
	if datasetID == "" {
		return nil, fmt.Errorf("no dataset provided when querying for tickers")
	}

	return query.TickerList(c.Hostname, c.Version, c.Token, datasetID)
}

func (c *Client) RunSearch(srch SearchDefinition) (query.RowsItems, error) {

	fmt.Printf("Running search: %s\n", srch.Name)

	dataset := query.DatasetItem{
		ID: srch.DatasetID,
	}

	return dataset.RunSearch(c.Hostname, c.Version, c.Token, c.PageSize, srch.Request)

}
