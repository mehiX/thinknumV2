package thinknum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	expiresFMT = "20060102T150405Z"
)

// AuthToken Token as recieved from the authentication endpoint
type AuthToken struct {
	// The token string to use in subsequent requests
	Token string `json:"auth_token"`
	// The expiration time for this token.
	// The token is valid as long as the expiration time is in the future.
	// Otherwise a new token should be requested
	Expires string `json:"auth_expires"`
}

// Cache Store the token data in a file on disk
// `fn` is the path to the cache file.
// The cache file should be treated as a secret and not checked into version control systems or shared with others.
func (t *AuthToken) Cache(fn string) error {

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t)
}

// IsExpired Checks if the token is expired
// Returns an error if the time value cannot be parsed from the string value
// Returns `true` if the expiry date is in the future and `false` otherwise
func (t *AuthToken) IsExpired() (bool, error) {

	exp, err := time.Parse(expiresFMT, t.Expires)
	if err != nil {
		return true, err
	}

	return time.Now().After(exp), nil
}

// GetToken Returns an authorization token that can be used for all subsequent requests in this session
// If a token is present in the local cache and it is still valid then it is returned
// If there is no cached token or the cached token is expired, then a new token is requested from the authentication server. The obtained token is cached locally before being returned
func GetToken(configAuth ConfigAuth) (*AuthToken, error) {
	token, err := LoadCachedToken(configAuth.TokenCachePath)
	if err == nil {
		if v, err := token.IsExpired(); !v && err == nil {
			// token from file is still valid, we can use it
			fmt.Println("Found cached valid token")
			return token, nil
		}
	}

	// token not present in local file or it is already expired
	token, err = RequestNewToken(configAuth)
	if err == nil {
		// save token for later use
		log.Println("Got new token. Try to cache it.")
		if err := token.Cache(configAuth.TokenCachePath); err != nil {
			log.Printf("Error saving token to file: %s. Error: %v\n", configAuth.TokenCachePath, err)
		} else {
			log.Println("Token successfully cached")
		}
	}

	return token, err
}

// RequestNewToken Makes a POST request to the authentication server.
// Upon success returns a valid token with and expiry date.
// In case of error it returns nil and the error the occurred.
func RequestNewToken(ca ConfigAuth) (*AuthToken, error) {

	data := make(url.Values)
	data.Set("version", ca.Version)
	data.Set("client_id", ca.ClientID)
	data.Set("client_secret", ca.ClientSecret)

	authURL := fmt.Sprintf("https://%s%s", ca.Hostname, ca.AuthEndpoint)
	resp, err := http.PostForm(authURL, data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth api responded with code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var a AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, err
	}

	return &a, nil

}

// LoadCachedToken Loads the token data from file
// `fn` is the path to the cache file
// Returns the token as read from the file or an error.
// The token is not validated and is returned as-is.
func LoadCachedToken(fn string) (*AuthToken, error) {

	info, err := os.Stat(fn)
	if err != nil {
		return nil, err
	}

	if info.IsDir() || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a file: %s", fn)
	}

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var a AuthToken
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, err
	}

	return &a, nil
}
