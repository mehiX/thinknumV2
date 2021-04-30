package thinknum

import (
	"fmt"
	"sync"

	"github.com/mehiX/thinknumV2/internal/query"
)

// runresult Contains the data returned by a search
// The results are intended to be passed to the next step in the pipeline: persisting the results
type runresult struct {
	Data   query.RowsItems
	Error  error
	Search SearchDefinition
}

// RunAll Runs all the searches defined in the configuration file
func (c *Client) RunAll() <-chan runresult {

	resultsStream := make(chan runresult)

	go func() {
		defer close(resultsStream)

		searchesStream := make(chan SearchDefinition)
		// generate work
		go func() {
			defer close(searchesStream)

			// from slice to channel
			for _, s := range c.Searches {
				// skip disabled seaches
				if !s.Disabled {
					searchesStream <- s
				} else {
					fmt.Printf("Skip search: %s [disabled]\n", s.Name)
				}
			}
		}()

		workers := c.Workers

		var wg sync.WaitGroup
		wg.Add(workers)

		// start idle workers
		for i := 0; i < workers; i++ {
			go runner(c, searchesStream, &wg, resultsStream)
		}

		wg.Wait()
	}()

	return resultsStream

}

// runner A worker that sits idle waiting for work on the incoming channel
func runner(client *Client, searches <-chan SearchDefinition, wg *sync.WaitGroup, results chan<- runresult) {
	defer wg.Done()

	for s := range searches {
		ri, err := client.RunSearch(s)
		results <- runresult{ri, err, s}
	}
}
