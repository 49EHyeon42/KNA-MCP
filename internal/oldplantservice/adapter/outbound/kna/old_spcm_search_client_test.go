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

	"github.com/49EHyeon42/KNA-MCP/internal/oldplantservice/application"
)

func TestOldSpcmSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/OldPlantService/oldSpcmSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/OldPlantService/oldSpcmSearch")
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"st":         "2",
			"sw":         "test-search-word",
			"dateGbn":    "1",
			"dateFrom":   "20240101",
			"dateTo":     "20241231",
			"numOfRows":  "1",
			"pageNo":     "2",
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
      <cprtCtnt>test-copyright</cprtCtnt><famlKorNm>test-family-Korean-name</famlKorNm>
      <famlNm>test-family-name</famlNm><frstRgstnDtm>test-first-registration-date-time</frstRgstnDtm>
      <imgUrl>https://example.com/test-old-plant.jpg</imgUrl><lastUpdtDtm> </lastUpdtDtm>
      <plantGnrlNm>test-plant-general-name</plantGnrlNm><plantOldSmplNo>test-old-specimen-number</plantOldSmplNo>
      <plantSpecsScnm>test-plant-scientific-name</plantSpecsScnm>
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

	got, err := client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{
		St:        "2",
		Sw:        "test-search-word",
		DateGbn:   "1",
		DateFrom:  "20240101",
		DateTo:    "20241231",
		NumOfRows: 1,
		PageNo:    2,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.OldSpcmSearchResult{
		Items: []application.OldSpcmSearchItem{{
			CprtCtnt:       "test-copyright",
			FamlKorNm:      "test-family-Korean-name",
			FamlNm:         "test-family-name",
			FrstRgstnDtm:   "test-first-registration-date-time",
			ImgURL:         "https://example.com/test-old-plant.jpg",
			LastUpdtDtm:    " ",
			PlantGnrlNm:    "test-plant-general-name",
			PlantOldSmplNo: "test-old-specimen-number",
			PlantSpecsScnm: "test-plant-scientific-name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestOldSpcmSearchOmitsEmptyDatesAndReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, key := range []string{"dateGbn", "dateFrom", "dateTo"} {
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

	result, err := client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{St: "1", Sw: "none", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestOldSpcmSearchReturnsDocumentedAPIErrors(t *testing.T) {
	tests := []struct {
		code    string
		message string
	}{
		{code: "01", message: "APPLICATION_ERROR"},
		{code: "02", message: "DB_ERROR"},
		{code: "03", message: "NODATA_ERROR"},
		{code: "04", message: "HTTP_ERROR"},
		{code: "05", message: "SERVICETIME_OUT"},
		{code: "10", message: "INVALID_REQUEST_PARAMETER_ERROR"},
		{code: "11", message: "NO_MANDATORY_REQUEST_PARAMETERS_ERROR"},
		{code: "12", message: "NO_OPENAPI_SERVICE_ERROR"},
		{code: "20", message: "SERVICE_ACCESS_DENIED_ERROR"},
		{code: "21", message: "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR"},
		{code: "22", message: "LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR"},
		{code: "30", message: "SERVICE_KEY_IS_NOT_REGISTERED_ERROR"},
		{code: "31", message: "DEADLINE_HAS_EXPIRED_ERROR"},
		{code: "32", message: "UNREGISTERED_IP_ERROR"},
		{code: "33", message: "UNSIGNED_CALL_ERROR"},
		{code: "99", message: "UNKNOWN_ERROR"},
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

			_, err = client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			var apiError *OldSpcmSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *OldSpcmSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestOldSpcmSearchReturnsObservedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>99</resultCode><resultMsg>ORA-01839: date not valid for month specified&#xA;</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20240230", DateTo: "20241231"})
	var apiError *OldSpcmSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *OldSpcmSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusOK || apiError.Code != "99" || apiError.Message != "ORA-01839: date not valid for month specified\n" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestOldSpcmSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	var apiError *OldSpcmSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *OldSpcmSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestOldSpcmSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "oldSpcmSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "oldSpcmSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "oldSpcmSearch: response missing resultCode"},
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

			_, err = client.OldSpcmSearch(context.Background(), application.OldSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestOldSpcmSearchLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		st   string
		sw   string
		want string
	}{
		{name: "partial scientific name", st: "1", sw: "Abies", want: "Abies"},
		{name: "case insensitive partial scientific name", st: "1", sw: "abies", want: "Abies"},
		{name: "exact scientific name", st: "2", sw: "Abies holophylla Maxim.", want: "Abies holophylla Maxim."},
		{name: "case insensitive exact scientific name", st: "2", sw: "abies holophylla maxim.", want: "Abies holophylla Maxim."},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{St: test.st, Sw: test.sw, PageNo: 1, NumOfRows: 10})
			if err != nil {
				t.Fatal(err)
			}
			matched := false
			for _, item := range result.Items {
				if strings.Contains(item.PlantSpecsScnm, test.want) {
					matched = true
				}
			}
			if !matched {
				t.Errorf("items = %#v, want matching %q", result.Items, test.want)
			}
		})
	}

	t.Run("registration date", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{
			St: "1", Sw: "Abies", DateGbn: "1", DateFrom: "20190401", DateTo: "20190430", PageNo: 1, NumOfRows: 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 1 || result.TotalCount == 0 {
			t.Errorf("result = %#v, want dated result", result)
		}
	})

	t.Run("modification date without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{
			St: "1", Sw: "Abies", DateGbn: "2", DateFrom: "20190401", DateTo: "20200101", PageNo: 1, NumOfRows: 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty dated result", result)
		}
	})

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{St: "1", Sw: "Abies", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{St: "1", Sw: "Abies", PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.Items[0].PlantOldSmplNo == second.Items[0].PlantOldSmplNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.OldSpcmSearch(ctx, application.OldSpcmSearchQuery{St: "1", Sw: "kna-mcp-no-result-20260822", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
