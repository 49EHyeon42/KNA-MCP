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

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

func TestPlantWordList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantWordListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantWordListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "2",
			"numOfRows":    "1",
			"reqSearchWrd": "test-search-word",
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
  <body>
    <items><item>
      <englsWrdNm>English word name</englsWrdNm>
      <krnWrdNm>Korean word name</krnWrdNm>
      <prfcnWrdNm> </prfcnWrdNm>
      <wrddscrt>word description</wrddscrt>
    </item></items>
    <numOfRows>1</numOfRows><pageNo>2</pageNo><totalCount>7</totalCount>
  </body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.PlantWordList(context.Background(), application.PlantWordListQuery{
		PageNo:       2,
		NumOfRows:    1,
		ReqSearchWrd: "test-search-word",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantWordListResult{
		Items: []application.PlantWordListItem{{
			EnglsWrdNm: "English word name",
			KrnWrdNm:   "Korean word name",
			PrfcnWrdNm: " ",
			Wrddscrt:   "word description",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantWordListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if _, exists := request.URL.Query()["reqSearchWrd"]; exists {
			t.Error("reqSearchWrd query was sent for an empty value")
		}
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantWordList(context.Background(), application.PlantWordListQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantWordListReturnsDocumentedAPIErrors(t *testing.T) {
	tests := []struct {
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
	}

	for _, test := range tests {
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

			_, err = client.PlantWordList(context.Background(), application.PlantWordListQuery{PageNo: 1, NumOfRows: 1})
			var apiError *PlantWordListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantWordListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantWordListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantWordList(context.Background(), application.PlantWordListQuery{PageNo: 1, NumOfRows: 1})
	var apiError *PlantWordListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantWordListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantWordListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantWordList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantWordList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantWordList: response missing resultCode"},
	}

	for _, test := range tests {
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

			_, err = client.PlantWordList(context.Background(), application.PlantWordListQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantWordListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name         string
		reqSearchWrd string
		krnWrdNm     string
	}{
		{name: "without search word"},
		{name: "exact Korean word name", reqSearchWrd: "가는맥", krnWrdNm: "가는맥"},
		{name: "partial Korean word name", reqSearchWrd: "가는", krnWrdNm: "가는맥"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantWordList(ctx, application.PlantWordListQuery{PageNo: 1, NumOfRows: 10, ReqSearchWrd: test.reqSearchWrd})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantWordList returned no items")
			}
			if test.krnWrdNm != "" {
				matched := false
				for _, item := range result.Items {
					if item.KrnWrdNm == test.krnWrdNm {
						matched = true
					}
				}
				if !matched {
					t.Errorf("plantWordList items = %#v, want krnWrdNm %q", result.Items, test.krnWrdNm)
				}
			}
		})
	}

	for _, test := range []struct {
		name         string
		reqSearchWrd string
	}{
		{name: "English word name is not searched", reqSearchWrd: "versatile"},
		{name: "refined word name is not searched", reqSearchWrd: "세맥"},
		{name: "without result", reqSearchWrd: "KNA-MCP-NO-RESULT"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantWordList(ctx, application.PlantWordListQuery{PageNo: 1, NumOfRows: 1, ReqSearchWrd: test.reqSearchWrd})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) != 0 || result.TotalCount != 0 {
				t.Errorf("result = %#v, want empty result", result)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		first, err := client.PlantWordList(ctx, application.PlantWordListQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantWordList(ctx, application.PlantWordListQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || first.Items[0].KrnWrdNm == second.Items[0].KrnWrdNm {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})
}
