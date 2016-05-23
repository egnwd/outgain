package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/github"
)

type config struct {
	ClientID, ClientSecret, RedirectURI string
}

var c config
var endpoint = struct {
	user string
}{
	user: "https://api.github.com/user",
}

const githubScope = "user"

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	c.ClientID = os.Getenv("CLIENT_ID")
	c.ClientSecret = os.Getenv("CLIENT_SECRET")
	c.RedirectURI = os.Getenv("REDIRECT_URI")
}

// GetOAuthURL returns the authorize URL for the github api with the query
func GetOAuthURL(state string) string {
	values := make(url.Values)
	values.Add("client_id", c.ClientID)
	values.Add("redirect_uri", c.RedirectURI)
	values.Add("scope", githubScope)
	values.Add("state", state)

	u, _ := url.ParseRequestURI(github.Endpoint.AuthURL)
	u.RawQuery = values.Encode()

	return u.String()
}

// GetAccessToken exchanges a one-time code for a users access token for
// this session
func GetAccessToken(state, code string) (string, error) {
	values := make(url.Values)
	values.Add("client_id", c.ClientID)
	values.Add("client_secret", c.ClientSecret)
	values.Add("code", code)
	values.Add("state", state)

	u, _ := url.ParseRequestURI(github.Endpoint.TokenURL)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, u.String(), bytes.NewBufferString(values.Encode()))
	req.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	values, err = url.ParseQuery(string(bs))
	if err != nil {
		return "", err
	}

	fmt.Println("Got Access Token!")
	fmt.Printf("Access Token: %s\n", values.Get("access_token"))
	return values.Get("access_token"), nil
}

// GetUsername retrieves the github username of a user that's currently
// logged in
func GetUsername(accessToken string) (string, error) {
	values := make(url.Values)
	values.Add("access_token", accessToken)

	u, _ := url.ParseRequestURI(endpoint.user)
	u.RawQuery = values.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		Login string
	}
	json.Unmarshal(bs, &data)

	return data.Login, nil
}
