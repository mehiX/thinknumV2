/**
Reads a search configuration from standard input, together this time split parameters.
Generates a new array of searches split accordingly
**/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	thinknum "github.com/mehiX/thinknumV2"
)

const (
	dateFMT = "2006-01-02"
)

var (
	startDate = flag.String("from", "", "Start of the queried period. Format: "+dateFMT)
	endDate   = flag.String("to", "", "End of the queried period. Format: "+dateFMT)
	interval  = flag.Duration("interval", 7*24*time.Hour, "Interval to use for splitting the dates interval")
)

func init() {
	flag.Parse()
}

func main() {

	if err := validateFlags(); err != nil {
		fmt.Printf("Error: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	from, to, err := parseTime(dateFMT, *startDate, *endDate)
	if err != nil {
		fmt.Printf("Wrong time format. Error: %v\n", err)
		os.Exit(3)
	}

	srch, err := thinknum.ReadSearch(os.Stdin)
	if err != nil {
		fmt.Printf("Invalid json input. Error: %v\n", err)
		os.Exit(4)
	}

	splitSearches := srch.Split(from, to, *interval)

	// output the initial search
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)

	encoder.Encode(splitSearches)
}

func validateFlags() error {
	if *startDate == "" {
		return fmt.Errorf("invalid start date: %q", *startDate)
	}

	if *endDate == "" {
		return fmt.Errorf("invalid end date: %q", *endDate)
	}

	return nil
}

func parseTime(dateFormat, s1, s2 string) (time.Time, time.Time, error) {

	t1, err := time.Parse(dateFormat, s1)
	if err != nil {
		return t1, time.Time{}, err
	}

	t2, err := time.Parse(dateFormat, s2)
	if err != nil {
		return t2, time.Time{}, err
	}

	return t1, t2, nil
}
