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
	endDate   = flag.String("to", time.Now().Format(dateFMT), "End of the queried period. Format: "+dateFMT)
	interval  = flag.Duration("interval", 7*24*time.Hour, "Interval to use for splitting the dates interval")
)

const usage = `
  USAGE:

  Splits a search definition into smaller time intervals. 
  
  Useful when the initial search would generate too many results, that would generate timeouts.
  Splitting a search definition can also make querying faster since the sliced definitions can run in parallel.
  Each search definition slice will write its own output file. These files can then be appended together in one file.

  Reads from standard input a JSON object of the form:
  
		{
			"name": "anything",
			"disabled": false,
			"output": "somefile",
			"output_types": ["json", "csv"],
			"dataset": "id of the dataset to query",
			"request": {
				"filters": [
					{
						"column": "",
						"type": "",
						"value": [""]
					}
				]
			}
		}

  Prints out to standard output an array of similar object, with time bound filters.
  Here an example where an interval of 30 days (-interval 720h) was used:

	[
		{
			"name": "anything",
			"disabled": false,
			"output": "somefile",
			"output_types": ["json", "csv"],
			"dataset": "id of the dataset to query",
			"request": {
				"filters": [
					{
						"column": "",
						"type": "",
						"value": [""]
					},
					{
						"column": "as_of_date",
						"type":">=",
						"value":["2021-03-02"]},
					{
						"column":"as_of_date",
						"type":"<",
						"value":["2021-04-01"]
					}
				]
			}
		},
		{
			"name": "anything",
			"disabled": false,
			"output": "somefile",
			"output_types": ["json", "csv"],
			"dataset": "id of the dataset to query",
			"request": {
				"filters": [
					{
						"column": "",
						"type": "",
						"value": [""]
					},
					{
						"column": "as_of_date",
						"type":">=",
						"value":["2021-04-01"]},
					{
						"column":"as_of_date",
						"type":"<",
						"value":["2021-05-01"]
					}
				]
			}
		}
	]
		
`

func main() {

	flag.Parse()

	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
	}

	if err := validateFlags(); err != nil {
		flag.Usage()
		fmt.Printf("\nError: %v\n", err)
		os.Exit(2)
	}

	from, to, err := parseTime(dateFMT, *startDate, *endDate)
	if err != nil {
		fmt.Printf("Wrong time format. Error: %v\n", err)
		os.Exit(3)
	}

	srch, err := thinknum.ReadSearchDefinition(os.Stdin)
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
		return fmt.Errorf("no start date provided")
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
		return time.Time{}, t2, err
	}

	return t1, t2, nil
}
