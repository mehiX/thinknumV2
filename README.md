# THINKNUM - A client for the Thinknum API

Download and install Golang [here](https://golang.org/dl/)

Clone the repository and move to the project folder.

Generate a new `config.json` by copying the provided template `config.tmpl`. Fill in the API credentials and the desired query parameters.

## Tools
- [thinknumclient](#ThinknumClient) - perform searches
- [splitsrch](#SplitSearch) - split a search specification in time frames

### ThinknumClient

Performs a search based on a config file.

If you specified an output folder in `config.json` make sure to create it first with something like `make -p ./out`

```bash
go get -d -v ./...
go build -ldflags="-X 'main.buildTime=$(date)' -X 'main.commitID=$(git rev-parse HEAD)'" ./cmd/thinknumclient

./thinknumclient -config config.json
```

### SplitSearch

For searches that return too many results it is useful to split the search specification in smaller time frames. These smaller searches can run in parallel. The results can then be concatenated to form the desired result.

The command takes in 3 parameters:
- `from` - the start date (YYYY-MM-DD)
- `to` - the end date (YYYY-MM-DD)
- `interval` - the length of the desired time frame. This is expressed as Golang duration and hence it has the limitation that the biggest time unit is the hour. For this reason we need to express a day as `24h` and a week as `168h`

```bash
go get -d -v ./...
go build ./cmd/splitsrch

# print out the help
./splitsrch -h

# test on an empty search specification
echo '{}' | ./splitsrch -from 2020-01-01 -to 2020-05-30 -interval $(( 30 * 24))h > outsrch.json

# use a saved search specification
cat insrch.json | ./splitsrch -from 2020-01-01 -to 2020-05-30 -interval $(( 30 * 24))h > outsrch.json
```

Where `insrch.json` can be something like:

```json
{
    "name": "Golang in EU",
    "disabled": false,
    "output": "out/golang_nl",
    "output_types": ["json", "csv"],
    "dataset": "job_listings",
    "filters": [
        {
            "column": "description",
            "type": "(...)",
            "value": ["Golang", "GOLANG"]
        },
        {
            "column": "country",
            "type": "=",
            "value": ["EU"]
        }
    ]
}
```

And the output will be:

```json
[
    {
        "name": "Golang in EU",
        "disabled": false,
        "output": "out/golang_nl",
        "output_types": [
            "json",
            "csv"
        ],
        "dataset": "job_listings",
        "filters": [
            {
                "column": "description",
                "type": "(...)",
                "value": [
                    "Golang",
                    "GOLANG"
                ]
            },
            {
                "column": "country",
                "type": "=",
                "value": [
                    "EU"
                ]
            },
            {
                "column": "as_of_date",
                "type": ">=",
                "value": [
                    "2020-01-01"
                ]
            },
            {
                "column": "as_of_date",
                "type": "<",
                "value": [
                    "2020-01-31"
                ]
            }
        ]
    },
    {
        "name": "Golang in EU",
        "disabled": false,
        "output": "out/golang_nl",
        "output_types": [
            "json",
            "csv"
        ],
        "dataset": "job_listings",
        "filters": [
            {
                "column": "description",
                "type": "(...)",
                "value": [
                    "Golang",
                    "GOLANG"
                ]
            },
            {
                "column": "country",
                "type": "=",
                "value": [
                    "EU"
                ]
            },
            {
                "column": "as_of_date",
                "type": ">=",
                "value": [
                    "2020-01-31"
                ]
            },
            {
                "column": "as_of_date",
                "type": "<",
                "value": [
                    "2020-03-01"
                ]
            }
        ]
    },
    {
        "name": "Golang in EU",
        "disabled": false,
        "output": "out/golang_nl",
        "output_types": [
            "json",
            "csv"
        ],
        "dataset": "job_listings",
        "filters": [
            {
                "column": "description",
                "type": "(...)",
                "value": [
                    "Golang",
                    "GOLANG"
                ]
            },
            {
                "column": "country",
                "type": "=",
                "value": [
                    "EU"
                ]
            },
            {
                "column": "as_of_date",
                "type": ">=",
                "value": [
                    "2020-03-01"
                ]
            },
            {
                "column": "as_of_date",
                "type": "<",
                "value": [
                    "2020-03-31"
                ]
            }
        ]
    },
    {
        "name": "Golang in EU",
        "disabled": false,
        "output": "out/golang_nl",
        "output_types": [
            "json",
            "csv"
        ],
        "dataset": "job_listings",
        "filters": [
            {
                "column": "description",
                "type": "(...)",
                "value": [
                    "Golang",
                    "GOLANG"
                ]
            },
            {
                "column": "country",
                "type": "=",
                "value": [
                    "EU"
                ]
            },
            {
                "column": "as_of_date",
                "type": ">=",
                "value": [
                    "2020-03-31"
                ]
            },
            {
                "column": "as_of_date",
                "type": "<",
                "value": [
                    "2020-04-30"
                ]
            }
        ]
    },
    {
        "name": "Golang in EU",
        "disabled": false,
        "output": "out/golang_nl",
        "output_types": [
            "json",
            "csv"
        ],
        "dataset": "job_listings",
        "filters": [
            {
                "column": "description",
                "type": "(...)",
                "value": [
                    "Golang",
                    "GOLANG"
                ]
            },
            {
                "column": "country",
                "type": "=",
                "value": [
                    "EU"
                ]
            },
            {
                "column": "as_of_date",
                "type": ">=",
                "value": [
                    "2020-04-30"
                ]
            },
            {
                "column": "as_of_date",
                "type": "<",
                "value": [
                    "2020-05-30"
                ]
            }
        ]
    }
]
```

The result can then be paste in the original client configuration.

