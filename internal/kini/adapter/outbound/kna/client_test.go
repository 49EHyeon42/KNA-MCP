package kna

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func liveServiceKey() string {
	if os.Getenv("KNA_MCP_LIVE_TESTS") != "1" {
		return ""
	}
	return os.Getenv("DATA_GO_KR_SERVICE_KEY")
}

func requireLiveServiceKey(t *testing.T) string {
	t.Helper()
	serviceKey := liveServiceKey()
	if serviceKey == "" {
		t.Skip("live tests require KNA_MCP_LIVE_TESTS=1 and DATA_GO_KR_SERVICE_KEY")
	}
	return serviceKey
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

func TestNewClientUses60SecondTimeout(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("timeout = %s, want 60s", client.httpClient.Timeout)
	}
}

func TestLiveServiceKey(t *testing.T) {
	for _, test := range []struct {
		name       string
		liveTests  string
		serviceKey string
		want       string
	}{
		{name: "disabled", serviceKey: "test-key"},
		{name: "invalid flag", liveTests: "true", serviceKey: "test-key"},
		{name: "missing service key", liveTests: "1"},
		{name: "enabled", liveTests: "1", serviceKey: "test-key", want: "test-key"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("KNA_MCP_LIVE_TESTS", test.liveTests)
			t.Setenv("DATA_GO_KR_SERVICE_KEY", test.serviceKey)

			if got := liveServiceKey(); got != test.want {
				t.Errorf("liveServiceKey() = %q, want %q", got, test.want)
			}
		})
	}
}
