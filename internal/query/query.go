package query

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ResponseMetadata struct {
	Count       int
	Total       int
	Status      int
	Summary     string
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// addRequestHeadersPOST Set the correct Content-Type and then add the authorization headers by calling `addRequestHeaders`
func addRequestHeadersPOST(r *http.Request, token, version string) {
	addRequestHeaders(r, token, version)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

// addRequestHeaders Add the necessary authorization headers
func addRequestHeaders(r *http.Request, token, version string) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	r.Header.Set("X-API-Version", version)
}

// fetchAll handle pagination. `processResp` is a function that contains the logic for sending the request, receiving the response and persisting the data (usually appending it to a slice)
func fetchAll(processResp func(url.Values) (ResponseMetadata, error), params url.Values) error {

	for {
		start, err := strconv.Atoi(params.Get("start"))
		if err != nil {
			return err
		}

		resp, err := processResp(params)
		if err != nil || resp.Total <= resp.Count+start {
			return err
		}
		params.Set("start", strconv.Itoa(start+resp.Count))
	}
}
