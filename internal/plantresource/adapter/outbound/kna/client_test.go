package kna

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestClientDoDoesNotExposeServiceKey(t *testing.T) {
	wantError := errors.New("network unavailable")
	client := &Client{httpClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, wantError
	})}}
	request, err := http.NewRequest(http.MethodGet, "https://example.com?serviceKey=secret", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.do(request)
	if !errors.Is(err, wantError) {
		t.Fatalf("error = %v, want %v", err, wantError)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("error exposes service key: %v", err)
	}
}
