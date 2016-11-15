package main

import (
	"net/http"

	"github.com/mrjones/oauth"
)

// Client is like an http.Client, but limited and using OAuth
type Client struct {
	consumer *oauth.Consumer
	token    *oauth.AccessToken
}

var provider = oauth.ServiceProvider{
	AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
	RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
	AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
}

// NewClient produces a new oauth 1.0 client for use with Twitter
func NewClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Client {
	consumer := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		provider,
	)
	token := &oauth.AccessToken{
		Token:          accessToken,
		Secret:         accessTokenSecret,
		AdditionalData: make(map[string]string),
	}
	return &Client{
		consumer,
		token,
	}
}

// Get performs HTTP GET with the oauth creds
func (a *Client) Get(endpoint string, params map[string]string) (*http.Response, error) {
	return a.consumer.Get(
		endpoint,
		params,
		a.token,
	)
}

// Post performs HTTP POST with the oauth creds
func (a *Client) Post(endpoint string, params map[string]string) (*http.Response, error) {
	return a.consumer.Post(
		endpoint,
		params,
		a.token,
	)
}
