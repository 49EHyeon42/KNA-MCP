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

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

func TestFngsPilbkInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/FungiService/fngsPilbkInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/FungiService/fngsPilbkInfo")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":     "test+/=",
			"reqFngsPilbkNo": "test-fungi-pictorial-book-number",
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
    <mshrmColorCdNm>mushroom color code name</mshrmColorCdNm>
    <crpphFomTpcdNm>carpophore form type code name</crpphFomTpcdNm>
    <familyKorNm>family Korean name</familyKorNm>
    <familyNm>family name</familyNm>
    <fngsEclgTpcdNm>fungi ecology type code name</fngsEclgTpcdNm>
    <fngsGnrlNm>fungi general name</fngsGnrlNm>
    <fngsPilbkNo>test-fungi-pictorial-book-number</fngsPilbkNo>
    <fngsPrpseTpcdNm>fungi purpose type code name</fngsPrpseTpcdNm>
    <fngsScnm>fungi scientific name</fngsScnm>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <grwEvrntDesc>growth environment description</grwEvrntDesc>
    <lastUpdtDtm> </lastUpdtDtm>
    <microShpe>microscopic shape</microShpe>
    <mshrmTpcdNm>mushroom type code name</mshrmTpcdNm>
    <occrrSsnNm>occurrence season name</occrrSsnNm>
    <rsrcActoClsscTpcdNm>resource classification type code name</rsrcActoClsscTpcdNm>
    <shpe>shape</shpe>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.FngsPilbkInfo(context.Background(), application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.FngsPilbkInfoResult{Item: &application.FngsPilbkInfoItem{
		MshrmColorCdNm:      "mushroom color code name",
		CrpphFomTpcdNm:      "carpophore form type code name",
		FamilyKorNm:         "family Korean name",
		FamilyNm:            "family name",
		FngsEclgTpcdNm:      "fungi ecology type code name",
		FngsGnrlNm:          "fungi general name",
		FngsPilbkNo:         "test-fungi-pictorial-book-number",
		FngsPrpseTpcdNm:     "fungi purpose type code name",
		FngsScnm:            "fungi scientific name",
		GenusKorNm:          "genus Korean name",
		GenusNm:             "genus name",
		GrwEvrntDesc:        "growth environment description",
		LastUpdtDtm:         " ",
		MicroShpe:           "microscopic shape",
		MshrmTpcdNm:         "mushroom type code name",
		OccrrSsnNm:          "occurrence season name",
		RsrcActoClsscTpcdNm: "resource classification type code name",
		Shpe:                "shape",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestFngsPilbkInfoReturnsEmptyItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.FngsPilbkInfo(context.Background(), application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "not-found"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestFngsPilbkInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.FngsPilbkInfo(context.Background(), application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"})
			var apiError *FngsPilbkInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *FngsPilbkInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestFngsPilbkInfoReturnsGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(response, `<OpenAPI_ServiceResponse><cmmMsgHeader><errMsg>SERVICE_KEY_IS_NULL</errMsg><returnAuthMsg>서비스 접근거부</returnAuthMsg><returnReasonCode>20</returnReasonCode></cmmMsgHeader></OpenAPI_ServiceResponse>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.FngsPilbkInfo(context.Background(), application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"})
	var apiError *FngsPilbkInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *FngsPilbkInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestFngsPilbkInfoReturnsResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "fngsPilbkInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "fngsPilbkInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "fngsPilbkInfo: response missing resultCode"},
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

			_, err = client.FngsPilbkInfo(context.Background(), application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "test-fungi-pictorial-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestFngsPilbkInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	searchResult, err := client.FngsPilbkSearch(ctx, application.FngsPilbkSearchQuery{PageNo: 1, NumOfRows: 2})
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult.Items) != 2 {
		t.Fatalf("fngsPilbkSearch items = %d, want 2", len(searchResult.Items))
	}

	for index, searchItem := range searchResult.Items {
		t.Run(fmt.Sprintf("search result %d", index+1), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.FngsPilbkInfo(ctx, application.FngsPilbkInfoQuery{ReqFngsPilbkNo: searchItem.FngsPilbkNo})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.FngsPilbkNo != searchItem.FngsPilbkNo || result.Item.FngsGnrlNm == "" || result.Item.FngsScnm == "" {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	t.Run("empty not-found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.FngsPilbkInfo(ctx, application.FngsPilbkInfoQuery{ReqFngsPilbkNo: "kna-mcp-no-result-20260818"})
		if err != nil {
			t.Fatal(err)
		}
		if result.Item != nil {
			t.Errorf("item = %#v, want nil", result.Item)
		}
	})
}
