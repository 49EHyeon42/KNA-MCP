package kna

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL            = "https://apis.data.go.kr"
	defaultRequestTimeout     = 60 * time.Second
	insectResourceBasePath    = "/1400119/InsectService"
	insectResourceSuccessCode = "00"
)

var insectResourceResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

// Client calls the Korea National Arboretum insect resource API through the Public Data Portal.
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
