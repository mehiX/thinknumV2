package thinknum

import (
	"fmt"
	"sync"
)

// RunAll Runs all the searches defined in the configuration file
func runAllFor(c *client, resultsStream chan SearchResult) {

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

}

// runner A worker that sits idle waiting for work on the incoming channel
func runner(c *client, searches <-chan SearchDefinition, wg *sync.WaitGroup, results chan<- SearchResult) {
	defer wg.Done()

	for s := range searches {
		ri := c.RunSearch(s)
		results <- SearchResult{ri, s}
	}
}
