package kna

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
)

func TestScnmInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/KlniService2/scnmInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/KlniService2/scnmInfo")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":    "test+/=",
			"reqLchnScnmId": "test-lichen-scientific-name-id",
		}
		if len(query) != len(wantQuery) {
			t.Errorf("query key count = %d, want %d", len(query), len(wantQuery))
		}
		for key, want := range wantQuery {
			if got := query.Get(key); got != want {
				t.Errorf("query %s = %q, want %q", key, got, want)
			}
		}

		response.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(response, `<?xml version="1.0" encoding="UTF-8"?>
<response>
  <header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header>
  <body><item>
    <stpltScnmRltnCdNm>scientific name relation code name</stpltScnmRltnCdNm>
    <falmNm>family name</falmNm>
    <falnKorNm> </falnKorNm>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <lastUpdtDtm>last update date time</lastUpdtDtm>
    <lchnGnrlNm>lichen general name</lchnGnrlNm>
    <lchnGnrlNm2> </lchnGnrlNm2>
    <lchnScnm>lichen scientific name</lchnScnm>
    <lchnScnmId>test-lichen-scientific-name-id</lchnScnmId>
    <ordscLtrtrNm>original description literature name</ordscLtrtrNm>
    <rmrk> </rmrk>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqLchnScnmID: "test-lichen-scientific-name-id"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ScnmInfoResult{Item: &application.ScnmInfoItem{
		StpltScnmRltnCdNm: "scientific name relation code name",
		FalmNm:            "family name",
		FalnKorNm:         " ",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		LastUpdtDtm:       "last update date time",
		LchnGnrlNm:        "lichen general name",
		LchnGnrlNm2:       " ",
		LchnScnm:          "lichen scientific name",
		LchnScnmID:        "test-lichen-scientific-name-id",
		OrdscLtrtrNm:      "original description literature name",
		Rmrk:              " ",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestScnmInfoReturnsEmptyItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqLchnScnmID: "not-found"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestScnmInfoReturnsDocumentedAPIErrors(t *testing.T) {
	for _, test := range []struct {
		code    string
		message string
	}{
		{code: "02", message: "DB_ERROR"},
		{code: "03", message: "NODATA_ERROR"},
		{code: "05", message: "SERVICETIME_OUT"},
		{code: "10", message: "INVALID_REQUEST_PARAMETER_ERROR"},
		{code: "11", message: "NO_MANDATORY_REQUEST_PARAMETERS_ERROR"},
		{code: "21", message: "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR"},
		{code: "33", message: "UNSIGNED_CALL_ERROR"},
	} {
		t.Run(test.code, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(response, `<response><header><resultCode>%s</resultCode></header></response>`, test.code)
			}))
			defer server.Close()

			client, err := NewClient("test-key")
			if err != nil {
				t.Fatal(err)
			}
			client.baseURL = server.URL

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqLchnScnmID: "test-lichen-scientific-name-id"})
			var apiError *ScnmInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ScnmInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestScnmInfoReturnsGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(response, `<OpenAPI_ServiceResponse><cmmMsgHeader><errMsg>SERVICE_KEY_IS_NOT_REGISTERED_ERROR</errMsg><returnAuthMsg>등록되지 않은 서비스키</returnAuthMsg><returnReasonCode>30</returnReasonCode></cmmMsgHeader></OpenAPI_ServiceResponse>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqLchnScnmID: "test-lichen-scientific-name-id"})
	var apiError *ScnmInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ScnmInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestScnmInfoReturnsResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "scnmInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "scnmInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "scnmInfo: response missing resultCode"},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				response.WriteHeader(test.statusCode)
				_, _ = io.WriteString(response, test.body)
			}))
			defer server.Close()

			client, err := NewClient("test-key")
			if err != nil {
				t.Fatal(err)
			}
			client.baseURL = server.URL

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqLchnScnmID: "test-lichen-scientific-name-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestScnmInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, lchnScnmID := range []string{"LS2017000708", "LS2017000825"} {
		t.Run(lchnScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqLchnScnmID: lchnScnmID})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.LchnScnmID != lchnScnmID || result.Item.LchnScnm == "" {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	for _, lchnScnmID := range []string{"KNA-MCP-NO-RESULT", "ls2017000708"} {
		t.Run("empty "+lchnScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqLchnScnmID: lchnScnmID})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item != nil {
				t.Errorf("item = %#v, want nil", result.Item)
			}
		})
	}

	t.Run("invalid service key", func(t *testing.T) {
		client, err := NewClient("invalid-service-key")
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		_, err = client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqLchnScnmID: "LS2017000708"})
		var apiError *ScnmInfoError
		if !errors.As(err, &apiError) {
			t.Fatalf("error = %v, want *ScnmInfoError", err)
		}
		if apiError.Code != "30" {
			t.Errorf("error = %#v, want code 30", apiError)
		}
	})
}
