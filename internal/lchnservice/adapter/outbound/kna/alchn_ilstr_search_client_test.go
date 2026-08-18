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

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

func TestAlchnIlstrSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/LchnService/alchnIlstrSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/LchnService/alchnIlstrSearch")
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
      <btnc>lichen scientific name</btnc>
      <cprtCtnt>copyright</cprtCtnt><detailYn>Y</detailYn><engNm> </engNm>
      <familyKorNm> </familyKorNm><familyNm>family name</familyNm>
      <frstRgstnDtm>first registration date time</frstRgstnDtm><genusKorNm>genus Korean name</genusKorNm>
      <genusNm>genus name</genusNm><imgUrl>http://example.com/lichen.jpg</imgUrl><japNm> </japNm>
      <lastUpdtDtm>last update date time</lastUpdtDtm><lchnGnrlNm>lichen general name</lchnGnrlNm>
      <lchnInfrpNm> </lchnInfrpNm><lchnPilbkNo>test-lichen-pictorial-book-number</lchnPilbkNo>
      <lchnScnmId>lichen scientific name ID</lchnScnmId><lchnTtnm>lichen species epithet</lchnTtnm><prkNm> </prkNm>
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

	got, err := client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{
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

	want := application.AlchnIlstrSearchResult{
		Items: []application.AlchnIlstrSearchItem{{
			Btnc:         "lichen scientific name",
			CprtCtnt:     "copyright",
			DetailYn:     "Y",
			EngNm:        " ",
			FamilyKorNm:  " ",
			FamilyNm:     "family name",
			FrstRgstnDtm: "first registration date time",
			GenusKorNm:   "genus Korean name",
			GenusNm:      "genus name",
			ImgURL:       "http://example.com/lichen.jpg",
			JapNm:        " ",
			LastUpdtDtm:  "last update date time",
			LchnGnrlNm:   "lichen general name",
			LchnInfrpNm:  " ",
			LchnPilbkNo:  "test-lichen-pictorial-book-number",
			LchnScnmID:   "lichen scientific name ID",
			LchnTtnm:     "lichen species epithet",
			PrkNm:        " ",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestAlchnIlstrSearchOmitsEmptyDatesAndReturnsEmptyItems(t *testing.T) {
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

	result, err := client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{St: "2", Sw: "none", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestAlchnIlstrSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			var apiError *AlchnIlstrSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *AlchnIlstrSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestAlchnIlstrSearchReturnsObservedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>99</resultCode><resultMsg>ORA-01843: not a valid month</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1, DateGbn: "1", DateFrom: "20241301", DateTo: "20241231"})
	var apiError *AlchnIlstrSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *AlchnIlstrSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusOK || apiError.Code != "99" || apiError.Message != "ORA-01843: not a valid month" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestAlchnIlstrSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	var apiError *AlchnIlstrSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *AlchnIlstrSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestAlchnIlstrSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "alchnIlstrSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "alchnIlstrSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "alchnIlstrSearch: response missing resultCode"},
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

			_, err = client.AlchnIlstrSearch(context.Background(), application.AlchnIlstrSearchQuery{St: "2", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestAlchnIlstrSearchLive(t *testing.T) {
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
		{name: "partial Korean name", st: "1", sw: "가락", want: "가락붉은열매지의"},
		{name: "exact Korean name", st: "3", sw: "가락붉은열매지의", want: "가락붉은열매지의"},
		{name: "partial scientific name", st: "2", sw: "cladonia", want: "Cladonia"},
		{name: "exact scientific name", st: "4", sw: "Cladonia digitata (L.) Hoffm.", want: "Cladonia digitata (L.) Hoffm."},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{St: test.st, Sw: test.sw, PageNo: 1, NumOfRows: 10})
			if err != nil {
				t.Fatal(err)
			}
			matched := false
			for _, item := range result.Items {
				if item.LchnGnrlNm == test.want || strings.Contains(item.Btnc, test.want) {
					matched = true
				}
			}
			if !matched {
				t.Errorf("items = %#v, want matching %q", result.Items, test.want)
			}
		})
	}

	for _, dateGbn := range []string{"1", "2"} {
		t.Run("dateGbn "+dateGbn, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{
				St: "2", Sw: "Cladonia", DateGbn: dateGbn, DateFrom: "20100101", DateTo: "20261231", PageNo: 1, NumOfRows: 1,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) != 1 || result.TotalCount == 0 {
				t.Errorf("result = %#v, want dated result", result)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 1, NumOfRows: 2})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{St: "2", Sw: "Cladonia", PageNo: 2, NumOfRows: 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 2 || len(second.Items) != 2 || first.TotalCount != second.TotalCount || first.Items[0].LchnPilbkNo == second.Items[0].LchnPilbkNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.AlchnIlstrSearch(ctx, application.AlchnIlstrSearchQuery{St: "2", Sw: "kna-mcp-no-result-20260818", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
