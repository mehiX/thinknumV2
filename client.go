package thinknum

import (
	"fmt"

	"github.com/mehiX/thinknumV2/internal/query"
)

// Client A Thinknum client should implement the `Client` interface
type Client interface {
	Datasets(string) ([]query.DatasetItem, error)
	Tickers(string) ([]query.TickerItem, error)
	RunSearch(SearchDefinition) query.RunResult
	RunAll() <-chan SearchResult
}

// SearchResult Brings together the search definition and the search results
type SearchResult struct {
	query.RunResult
	Search SearchDefinition
}

type client struct {
	Config
	Token string
}

// NewClientFromJSON Returns a new client for the Thinknum API. It will contain a valid token based on the received credentials
func NewClientFromJSON(configFile string) (Client, error) {

	cfg, err := ConfigFromJSON(configFile)
	if err != nil {
		return nil, err
	}

	token, err := GetToken(cfg.ConfigAuth)
	if err != nil {
		return nil, err
	}

	return NewClient(cfg, token), nil
}

// NewClient Create a new client providing your own configuration and token
func NewClient(cfg *Config, token *AuthToken) Client {
	return &client{
		Config: *cfg,
		Token:  token.Token,
	}

}

// Datasets Get a list of available datasets
// If a tickerID is provided (is not empty) then it is used to filter the datasets
func (c *client) Datasets(tickerID string) ([]query.DatasetItem, error) {
	return query.Datasets(c.Hostname, c.Version, c.Token, tickerID)
}

// Tickers Get the list of tickers for the provided `datasetID`
func (c *client) Tickers(datasetID string) ([]query.TickerItem, error) {
	if datasetID == "" {
		return nil, fmt.Errorf("no dataset provided when querying for tickers")
	}
	return query.TickerList(c.Hostname, c.Version, c.Token, datasetID)
}

// RunSearch Perform a search based on the SearchDefinition supplied
// Return a RunResult
func (c *client) RunSearch(sd SearchDefinition) query.RunResult {

	fmt.Printf("Running search: %s\n", sd.Name)

	dataset := query.DatasetItem{
		ID: sd.DatasetID,
	}

	return dataset.RunSearch(c.Hostname, c.Version, c.Token, c.PageSize, sd.Request)

}

// RunAll Runs all the searches defined in the configuration file
func (c *client) RunAll() <-chan SearchResult {

	resultsStream := make(chan SearchResult)

	go func() {
		defer close(resultsStream)

		runAllFor(c, resultsStream)
	}()

	return resultsStream
}
