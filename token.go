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
	secretFilename = ".auth"
	expiresFMT     = "20060102T150405Z"
	authURL        = "https://data.thinknum.com/api/authorize"
)

type auth struct {
	Token   string `json:"auth_token"`
	Expires string `json:"auth_expires"`
}

func Token(version, clientID, clientSecret string) (*auth, error) {
	token, err := tokenFromFile(secretFilename)
	if err == nil {
		if v, err := tokenIsValid(token); v && err == nil {
			// token from file is still valid, we can use it
			fmt.Println("Found cached valid token")
			return token, nil
		}
	}

	// token not present in local file or it is already expired
	token, err = tokenFromURL(version, clientID, clientSecret)
	if err == nil {
		// save token for later use
		log.Println("Got new token. Try to cache it.")
		if err := tokenToFile(secretFilename, token); err != nil {
			log.Printf("Error saving token to file: %s. Error: %v\n", secretFilename, err)
		} else {
			log.Println("Token successfully cached")
		}
	}

	return token, err
}

func tokenFromURL(version, clientID, clientSecret string) (*auth, error) {

	form := make(url.Values)
	form.Set("version", version)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	resp, err := http.PostForm(authURL, form)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth api responded with code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var a auth
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, err
	}

	return &a, nil

}

func tokenFromFile(fn string) (*auth, error) {

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

	var a auth
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, err
	}

	return &a, nil
}

func tokenToFile(fn string, t *auth) error {

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t)
}

func tokenIsValid(t *auth) (bool, error) {

	exp, err := time.Parse(expiresFMT, t.Expires)
	if err != nil {
		return false, err
	}

	return time.Now().Before(exp), nil
}
