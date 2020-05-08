package corezoid

import (
	"errors"
	"net/http"
	"time"
)

var asyncApi = "https://api.corezoid.com/api/2"
var syncApi = "https://sync-api.corezoid.com/api/1"

type Client struct {
	contentType ContentType
	client      *http.Client
	apiKey      string
	apiSecret   string
}

func NewClient(
	contentType ContentType,
	httpClient *http.Client,
	apiKey,
	apiSecret string,
) (*Client, error) {
	if contentType == Unknown {
		return nil, errors.New("unknown content type")
	}
	if len(apiKey) == 0 {
		return nil, errors.New("empty api key")
	}
	if len(apiSecret) == 0 {
		return nil, errors.New("empty api secret")
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
		httpClient.Timeout = time.Minute
	}

	return &Client{
		contentType: contentType,
		client:      httpClient,
		apiKey:      apiKey,
		apiSecret:   apiSecret,
	}, nil
}

func (c *Client) SetClient(client *http.Client) {
	c.client = client
}

func (c *Client) SetTransport(transport http.RoundTripper) {
	c.client.Transport = transport
}

func (c *Client) SetTimeout(duration time.Duration) {
	c.client.Timeout = duration
}
