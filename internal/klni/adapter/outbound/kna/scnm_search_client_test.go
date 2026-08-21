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

func TestScnmSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/KlniService2/scnmSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/KlniService2/scnmSearch")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"pageNo":     "2",
			"numOfRows":  "1",
			"reqGnrlNm":  "general name",
			"reqScnm":    "scientific name",
			"dateFrom":   "20240101",
			"dateTo":     "20241231",
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
      <stpltScnmRltnCdNm>scientific name relation code name</stpltScnmRltnCdNm>
      <classKorNm> </classKorNm>
      <classNm>class name</classNm>
      <falmNm>family name</falmNm>
      <falnKorNm>family Korean name</falnKorNm>
      <genusKorNm>genus Korean name</genusKorNm>
      <genusNm>genus name</genusNm>
      <lastUpdtDtm>last update date time</lastUpdtDtm>
      <lchnGnrlNm>lichen general name</lchnGnrlNm>
      <lchnScnm>lichen scientific name</lchnScnm>
      <lchnScnmId>lichen scientific name ID</lchnScnmId>
      <ordKorNm>order Korean name</ordKorNm>
      <ordNm>order name</ordNm>
      <phylumKorNm>phylum Korean name</phylumKorNm>
      <phylumNm>phylum name</phylumNm>
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

	got, err := client.ScnmSearch(context.Background(), application.ScnmSearchQuery{
		PageNo:    2,
		NumOfRows: 1,
		ReqGnrlNm: "general name",
		ReqScnm:   "scientific name",
		DateFrom:  "20240101",
		DateTo:    "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ScnmSearchResult{
		Items: []application.ScnmSearchItem{{
			StpltScnmRltnCdNm: "scientific name relation code name",
			ClassKorNm:        " ",
			ClassNm:           "class name",
			FalmNm:            "family name",
			FalnKorNm:         "family Korean name",
			GenusKorNm:        "genus Korean name",
			GenusNm:           "genus name",
			LastUpdtDtm:       "last update date time",
			LchnGnrlNm:        "lichen general name",
			LchnScnm:          "lichen scientific name",
			LchnScnmID:        "lichen scientific name ID",
			OrdKorNm:          "order Korean name",
			OrdNm:             "order name",
			PhylumKorNm:       "phylum Korean name",
			PhylumNm:          "phylum name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestScnmSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, key := range []string{"reqGnrlNm", "reqScnm", "dateFrom", "dateTo"} {
			if _, exists := request.URL.Query()[key]; exists {
				t.Errorf("%s query was sent for an empty value", key)
			}
		}
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.ScnmSearch(context.Background(), application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestScnmSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.ScnmSearch(context.Background(), application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *ScnmSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ScnmSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestScnmSearchReturnsObservedDateErrors(t *testing.T) {
	for _, message := range []string{
		"ORA-01839: date not valid for month specified",
		"ORA-01843: not a valid month",
	} {
		t.Run(message, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(response, `<response><header><resultCode>99</resultCode><resultMsg>%s</resultMsg></header></response>`, message)
			}))
			defer server.Close()

			client, err := NewClient("test-key")
			if err != nil {
				t.Fatal(err)
			}
			client.baseURL = server.URL

			_, err = client.ScnmSearch(context.Background(), application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *ScnmSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ScnmSearchError", err)
			}
			if apiError.Code != "99" || apiError.Message != message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestScnmSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.ScnmSearch(context.Background(), application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
	var apiError *ScnmSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ScnmSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestScnmSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "scnmSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "scnmSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "scnmSearch: response missing resultCode"},
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

			_, err = client.ScnmSearch(context.Background(), application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestScnmSearchLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name  string
		query application.ScnmSearchQuery
	}{
		{name: "without optional conditions", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10}},
		{name: "general name only", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, ReqGnrlNm: "가는"}},
		{name: "uppercase scientific name", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, ReqScnm: "Cladonia"}},
		{name: "lowercase scientific name", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, ReqScnm: "cladonia"}},
		{name: "both name conditions", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, ReqGnrlNm: "가락붉은열매지의", ReqScnm: "Cladonia"}},
		{name: "dateFrom only", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, DateFrom: "20170216"}},
		{name: "dateTo only", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, DateTo: "20170216"}},
		{name: "both date conditions", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, DateFrom: "20170101", DateTo: "20171231"}},
		{name: "all optional conditions", query: application.ScnmSearchQuery{PageNo: 1, NumOfRows: 10, ReqGnrlNm: "가락붉은열매지의", ReqScnm: "Cladonia", DateFrom: "20170101", DateTo: "20171231"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmSearch(ctx, test.query)
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 || result.TotalCount < len(result.Items) {
				t.Fatalf("scnmSearch result = %#v, want nonempty result", result)
			}
			for _, item := range result.Items {
				if test.query.ReqGnrlNm != "" && !strings.Contains(item.LchnGnrlNm, test.query.ReqGnrlNm) {
					t.Errorf("lchnGnrlNm = %q, want containing %q", item.LchnGnrlNm, test.query.ReqGnrlNm)
				}
				if test.query.ReqScnm != "" && !strings.Contains(strings.ToLower(item.LchnScnm), strings.ToLower(test.query.ReqScnm)) {
					t.Errorf("lchnScnm = %q, want containing %q", item.LchnScnm, test.query.ReqScnm)
				}
				checkLastUpdtDtm(t, item.LastUpdtDtm, test.query.DateFrom, test.query.DateTo)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.ScnmSearch(ctx, application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.ScnmSearch(ctx, application.ScnmSearchQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || reflect.DeepEqual(first.Items[0], second.Items[0]) {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.ScnmSearch(ctx, application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1, ReqGnrlNm: "KNA-MCP-NO-RESULT"})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		_, err := client.ScnmSearch(ctx, application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1, DateFrom: "20230230"})
		var apiError *ScnmSearchError
		if !errors.As(err, &apiError) || apiError.Code != "99" {
			t.Errorf("error = %#v, want code 99", err)
		}
	})

	t.Run("invalid service key", func(t *testing.T) {
		client, err := NewClient("invalid-service-key")
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		_, err = client.ScnmSearch(ctx, application.ScnmSearchQuery{PageNo: 1, NumOfRows: 1})
		var apiError *ScnmSearchError
		if !errors.As(err, &apiError) {
			t.Fatalf("error = %v, want *ScnmSearchError", err)
		}
		if apiError.Code != "30" {
			t.Errorf("error = %#v, want code 30", apiError)
		}
	})
}

func checkLastUpdtDtm(t *testing.T, value, dateFrom, dateTo string) {
	t.Helper()
	if dateFrom == "" && dateTo == "" {
		return
	}
	got, err := time.Parse("2006/01/02", value)
	if err != nil {
		t.Errorf("lastUpdtDtm = %q: %v", value, err)
		return
	}
	if dateFrom != "" {
		from, _ := time.Parse("20060102", dateFrom)
		if got.Before(from) {
			t.Errorf("lastUpdtDtm = %q, want on or after %s", value, dateFrom)
		}
	}
	if dateTo != "" {
		to, _ := time.Parse("20060102", dateTo)
		if !got.Before(to) {
			t.Errorf("lastUpdtDtm = %q, want before %s", value, dateTo)
		}
	}
}
