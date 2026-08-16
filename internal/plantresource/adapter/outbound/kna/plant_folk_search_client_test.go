package kna

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

func TestPlantFolkSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantFolkSearchPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantFolkSearchPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "1",
			"numOfRows":    "2",
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
      <flcstPlantIdntfDscrt>folk plant identification description</flcstPlantIdntfDscrt>
      <flpltId>folk plant ID</flpltId>
      <plantBrdgFomTpcdNm>plant breeding form type name</plantBrdgFomTpcdNm>
      <plantGnrlNm>plant general name</plantGnrlNm>
      <plantSpecsScnm>plant species scientific name</plantSpecsScnm>
      <ptnt>patent information</ptnt>
    </item></items>
    <numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>7</totalCount>
  </body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.PlantFolkSearch(context.Background(), application.PlantFolkSearchQuery{
		PageNo:       1,
		NumOfRows:    2,
		ReqSearchWrd: "test-search-word",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantFolkSearchResult{
		Items: []application.PlantFolkSearchItem{{
			FlcstPlantIdntfDscrt: "folk plant identification description",
			FlpltID:              "folk plant ID",
			PlantBrdgFomTpcdNm:   "plant breeding form type name",
			PlantGnrlNm:          "plant general name",
			PlantSpecsScnm:       "plant species scientific name",
			Ptnt:                 "patent information",
		}},
		NumOfRows:  2,
		PageNo:     1,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantFolkSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if _, exists := request.URL.Query()["reqSearchWrd"]; exists {
			t.Error("reqSearchWrd query was sent for an empty search word")
		}
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantFolkSearch(context.Background(), application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantFolkSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantFolkSearch(context.Background(), application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *PlantFolkSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantFolkSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantFolkSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantFolkSearch(context.Background(), application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 1})
	var apiError *PlantFolkSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantFolkSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantFolkSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantFolkSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantFolkSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantFolkSearch: response missing resultCode"},
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

			_, err = client.PlantFolkSearch(context.Background(), application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantFolkSearchLive(t *testing.T) {
	serviceKey := os.Getenv("DATA_GO_KR_SERVICE_KEY")
	if serviceKey == "" {
		t.Skip("DATA_GO_KR_SERVICE_KEY is not set")
	}

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name           string
		pageNo         int
		numOfRows      int
		reqSearchWrd   string
		plantGnrlNm    string
		plantSpecsScnm string
	}{
		{name: "without search word", pageNo: 1, numOfRows: 1},
		{name: "exact Korean name", pageNo: 1, numOfRows: 10, reqSearchWrd: "가는금불초", plantGnrlNm: "가는금불초"},
		{name: "partial Korean name", pageNo: 1, numOfRows: 10, reqSearchWrd: "금불초", plantGnrlNm: "가는금불초"},
		{name: "uppercase scientific name", pageNo: 1, numOfRows: 10, reqSearchWrd: "Inula", plantSpecsScnm: "inula"},
		{name: "lowercase scientific name", pageNo: 1, numOfRows: 10, reqSearchWrd: "inula", plantSpecsScnm: "inula"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantFolkSearch(ctx, application.PlantFolkSearchQuery{
				PageNo:       test.pageNo,
				NumOfRows:    test.numOfRows,
				ReqSearchWrd: test.reqSearchWrd,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantFolkSearch returned no items")
			}
			if test.plantGnrlNm != "" || test.plantSpecsScnm != "" {
				matched := false
				for _, item := range result.Items {
					if test.plantGnrlNm != "" && item.PlantGnrlNm == test.plantGnrlNm {
						matched = true
					}
					if test.plantSpecsScnm != "" && strings.Contains(strings.ToLower(item.PlantSpecsScnm), test.plantSpecsScnm) {
						matched = true
					}
				}
				if !matched {
					t.Errorf("plantFolkSearch items = %#v, want matching result", result.Items)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		first, err := client.PlantFolkSearch(ctx, application.PlantFolkSearchQuery{PageNo: 1, NumOfRows: 1, ReqSearchWrd: "금불초"})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantFolkSearch(ctx, application.PlantFolkSearchQuery{PageNo: 2, NumOfRows: 1, ReqSearchWrd: "금불초"})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != 2 || second.TotalCount != 2 || first.Items[0].FlpltID == second.Items[0].FlpltID {
			t.Errorf("first = %#v, second = %#v, want distinct pages with totalCount 2", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantFolkSearch(ctx, application.PlantFolkSearchQuery{
			PageNo:       1,
			NumOfRows:    1,
			ReqSearchWrd: "KNA-MCP-NO-RESULT",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
