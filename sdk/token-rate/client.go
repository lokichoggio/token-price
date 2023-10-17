package tokenrate

import (
	resty "github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

type ClientOption func(client *Client)

func NewClient(apiKey string, options ...ClientOption) *Client {
	c := &Client{
		client: resty.New().
			SetHeader("Content-Type", "application/json").
			SetHeader("X-CoinAPI-Key", apiKey),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithHeaders(h map[string]string) ClientOption {
	return func(client *Client) {
		client.client.SetHeaders(h)
	}
}
