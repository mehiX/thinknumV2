package thinknum

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html"
	"io/ioutil"
	"strings"

	"github.com/mehiX/thinknumV2/internal/query"
	"github.com/microcosm-cc/bluemonday"
)

func persistResultCSV(fn string, rd query.RunResult) error {

	filename := fn + ".csv"

	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	if err := w.Write(headerForCSV(rd.Data.Fields)); err != nil {
		return err
	}

	if err := w.WriteAll(prepareForCSV(rd.Data.Rows)); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf.Bytes(), 0666)
}

func headerForCSV(fields []query.Field) []string {
	header := make([]string, len(fields))

	for i := range fields {
		header[i] = fields[i].DisplayName
	}

	return header
}

func prepareForCSV(matrix []query.Row) [][]string {
	out := make([][]string, len(matrix))

	p := bluemonday.StrictPolicy()

	for index, r := range matrix {
		out[index] = make([]string, len(r))
		for j := range r {
			s := fmt.Sprintf("%v", matrix[index][j])
			s = html.UnescapeString(s)
			s = p.Sanitize(s)
			s = strings.ReplaceAll(s, "\n", " ")
			out[index][j] = s
		}
	}

	return out
}
