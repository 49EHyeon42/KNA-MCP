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

	"github.com/49EHyeon42/KNA-MCP/internal/kini/application"
)

func TestScnmInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/KiniService/scnmInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/KiniService/scnmInfo")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":     "test+/=",
			"reqInsctScnmId": "test-insect-scientific-name-id",
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
    <superFalmNm>superfamily name</superFalmNm>
    <classKorNm>class Korean name</classKorNm>
    <classNm>class name</classNm>
    <falmKorNm>family Korean name</falmKorNm>
    <falmNm>family name</falmNm>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <insctGnrlNm>insect general name</insctGnrlNm>
    <insctGnrlNm2> </insctGnrlNm2>
    <insctScnmId>test-insect-scientific-name-id</insctScnmId>
    <insctSpecsScnm>&lt;em&gt;insect&lt;/em&gt; scientific name</insctSpecsScnm>
    <lastUpdtDtm>last update date time</lastUpdtDtm>
    <ordKorNm>order Korean name</ordKorNm>
    <ordNm>order name</ordNm>
    <stpltScnmRltnCdNm>standard insect scientific name relation code name</stpltScnmRltnCdNm>
    <subFalmKorNm>subfamily Korean name</subFalmKorNm>
    <subFalmNm>subfamily name</subFalmNm>
    <superFalmKorNm>superfamily Korean name</superFalmKorNm>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqInsctScnmID: "test-insect-scientific-name-id"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ScnmInfoResult{Item: &application.ScnmInfoItem{
		SuperFalmNm:       "superfamily name",
		ClassKorNm:        "class Korean name",
		ClassNm:           "class name",
		FalmKorNm:         "family Korean name",
		FalmNm:            "family name",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		InsctGnrlNm:       "insect general name",
		InsctGnrlNm2:      " ",
		InsctScnmID:       "test-insect-scientific-name-id",
		InsctSpecsScnm:    "<em>insect</em> scientific name",
		LastUpdtDtm:       "last update date time",
		OrdKorNm:          "order Korean name",
		OrdNm:             "order name",
		StpltScnmRltnCdNm: "standard insect scientific name relation code name",
		SubFalmKorNm:      "subfamily Korean name",
		SubFalmNm:         "subfamily name",
		SuperFalmKorNm:    "superfamily Korean name",
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

	result, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqInsctScnmID: "not-found"})
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

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqInsctScnmID: "test-insect-scientific-name-id"})
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

	_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqInsctScnmID: "test-insect-scientific-name-id"})
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

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqInsctScnmID: "test-insect-scientific-name-id"})
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

	for _, insctScnmID := range []string{"Z000026622", "ZREB001001"} {
		t.Run(insctScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqInsctScnmID: insctScnmID})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.InsctScnmID != insctScnmID || result.Item.InsctGnrlNm == "" || !strings.Contains(result.Item.InsctSpecsScnm, "<em>") {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	for _, insctScnmID := range []string{"KNA-MCP-NO-RESULT", "zreb001001"} {
		t.Run("empty "+insctScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqInsctScnmID: insctScnmID})
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

		_, err = client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqInsctScnmID: "Z000026622"})
		var apiError *ScnmInfoError
		if !errors.As(err, &apiError) {
			t.Fatalf("error = %v, want *ScnmInfoError", err)
		}
		if apiError.Code != "30" {
			t.Errorf("error = %#v, want code 30", apiError)
		}
	})
}
