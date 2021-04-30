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
func NewClient(configFile string) (*Client, error) {

	cfg, err := FromJSON(configFile)
	if err != nil {
		return nil, err
	}

	token, err := Token(cfg.Version, cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config: *cfg,
		Token:  token.Token,
	}, nil
}

func (c *Client) Datasets(tickerID string) ([]query.DatasetItem, error) {
	return query.Datasets(c.Hostname, c.Version, c.Token, tickerID)
}

func (c *Client) Tickers(datasetID string) ([]query.TickerItem, error) {
	return query.TickerList(c.Hostname, c.Version, c.Token, datasetID)
}

func (c *Client) RunSearch(srch SearchDefinition) (query.RowsItems, error) {

	fmt.Printf("Running search: %s\n", srch.Name)

	dataset := query.DatasetItem{
		ID: srch.DatasetID,
	}

	return dataset.RunSearch(c.Hostname, c.Version, c.Token, c.PageSize, srch.Request)

}
