package thinknum

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/mehiX/thinknumV2/internal/query"
)

func writeToFile(srchRes SearchResult, ftype string) SaveResult {

	var err error

	switch ftype {
	case "json":
		err = persistResultJSON(srchRes.Search.OutputFile, srchRes.RunResult)
	case "csv":
		err = persistResultCSV(srchRes.Search.OutputFile, srchRes.RunResult)
	default:
		err = errors.New("Type not supported")
	}

	return SaveResult{
		Search: srchRes.Search,
		Type:   ftype,
		Error:  err,
	}
}

func persistResultJSON(fn string, d query.RunResult) error {
	filename := fn + ".json"

	str, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, str, 0666)
}
