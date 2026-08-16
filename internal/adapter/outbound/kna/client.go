package kna

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL        = "https://apis.data.go.kr"
	defaultRequestTimeout = 60 * time.Second
	plantResourceBasePath = "/1400119/PlantResource"
)

// Client calls Korea National Arboretum APIs through the Public Data Portal.
type Client struct {
	baseURL    string
	httpClient *http.Client
	serviceKey string
}

// NewClient creates a client with a public data service key.
func NewClient(serviceKey string) (*Client, error) {
	if serviceKey == "" {
		return nil, errors.New("service key is required")
	}

	serviceKey, err := url.PathUnescape(serviceKey)
	if err != nil {
		return nil, fmt.Errorf("decode service key: %w", err)
	}

	return &Client{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultRequestTimeout},
		serviceKey: serviceKey,
	}, nil
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	response, err := c.httpClient.Do(request)
	if err == nil {
		return response, nil
	}

	var urlError *url.Error
	if errors.As(err, &urlError) {
		return nil, urlError.Err
	}
	return nil, err
}

func setQueryValue(query url.Values, key, value string) {
	if value != "" {
		query.Set(key, value)
	}
}
