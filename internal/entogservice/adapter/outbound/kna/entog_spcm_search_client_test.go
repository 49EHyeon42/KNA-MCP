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

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

func TestEntogSpcmSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/EntogService/entogSpcmSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/EntogService/entogSpcmSearch")
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=", "st": "2", "sw": "test-search-word",
			"dateGbn": "1", "dateFrom": "20200101", "dateTo": "20201231",
			"numOfRows": "1", "pageNo": "2",
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
      <btnc>test-scientific-name</btnc><clctDyDesc>test-collection-date</clctDyDesc>
      <cprtCtnt>test-copyright</cprtCtnt><detailYn>test-detail-availability</detailYn>
      <entogOfnmKrlngNm>test-Korean-name</entogOfnmKrlngNm><entogSmplNo>test-specimen-number</entogSmplNo>
      <familyKorNm>test-family-Korean-name</familyKorNm><familyNm>test-family-name</familyNm>
      <frstRgstnDtm>test-first-registration-date</frstRgstnDtm><genusKorNm>test-genus-Korean-name</genusKorNm>
      <genusNm>test-genus-name</genusNm><imgUrl>test-image-URL</imgUrl>
      <lastUpdtDtm>test-last-update-date</lastUpdtDtm><ordKorNm>test-order-Korean-name</ordKorNm><ordNm>test-order-name</ordNm>
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

	got, err := client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{
		St: "2", Sw: "test-search-word", DateGbn: "1", DateFrom: "20200101", DateTo: "20201231", NumOfRows: 1, PageNo: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.EntogSpcmSearchResult{
		Items: []application.EntogSpcmSearchItem{{
			Btnc: "test-scientific-name", ClctDyDesc: "test-collection-date", CprtCtnt: "test-copyright",
			DetailYn: "test-detail-availability", EntogOfnmKrlngNm: "test-Korean-name", EntogSmplNo: "test-specimen-number",
			FamilyKorNm: "test-family-Korean-name", FamilyNm: "test-family-name", FrstRgstnDtm: "test-first-registration-date",
			GenusKorNm: "test-genus-Korean-name", GenusNm: "test-genus-name", ImgURL: "test-image-URL",
			LastUpdtDtm: "test-last-update-date", OrdKorNm: "test-order-Korean-name", OrdNm: "test-order-name",
		}},
		NumOfRows: 1, PageNo: 2, TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestEntogSpcmSearchOmitsEmptyOptionalQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, name := range []string{"dateGbn", "dateFrom", "dateTo"} {
			if request.URL.Query().Has(name) {
				t.Errorf("query contains optional %s", name)
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
	if _, err := client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestEntogSpcmSearchPreservesObservedValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items><item><frstRgstnDtm> </frstRgstnDtm><genusKorNm> </genusKorNm><imgUrl>NONE</imgUrl><lastUpdtDtm> </lastUpdtDtm></item></items><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>1</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].FrstRgstnDtm != " " || result.Items[0].GenusKorNm != " " || result.Items[0].ImgURL != "NONE" || result.Items[0].LastUpdtDtm != " " {
		t.Errorf("result = %#v, want preserved spaces and NONE", result)
	}
}

func TestEntogSpcmSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "none", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestEntogSpcmSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			var apiError *EntogSpcmSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *EntogSpcmSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestEntogSpcmSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	var apiError *EntogSpcmSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *EntogSpcmSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestEntogSpcmSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "entogSpcmSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "entogSpcmSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "entogSpcmSearch: response missing resultCode"},
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

			_, err = client.EntogSpcmSearch(context.Background(), application.EntogSpcmSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestEntogSpcmSearchLive(t *testing.T) {
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
		{name: "partial Korean name", st: "1", sw: "고려붓", want: "고려붓톡토기"},
		{name: "exact Korean name", st: "3", sw: "고려붓톡토기", want: "고려붓톡토기"},
		{name: "partial scientific name", st: "2", sw: "Homidia", want: "Homidia"},
		{name: "case insensitive partial scientific name", st: "2", sw: "homidia", want: "Homidia"},
		{name: "exact scientific name", st: "4", sw: "Homidia koreana Lee & Lee, 1981", want: "Homidia koreana"},
		{name: "case insensitive exact scientific name", st: "4", sw: "homidia koreana lee & lee, 1981", want: "Homidia koreana"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.EntogSpcmSearch(ctx, application.EntogSpcmSearchQuery{St: test.st, Sw: test.sw, PageNo: 1, NumOfRows: 10})
			if err != nil {
				t.Fatal(err)
			}
			matched := false
			for _, item := range result.Items {
				if strings.Contains(item.Btnc, test.want) || strings.Contains(item.EntogOfnmKrlngNm, test.want) {
					matched = true
				}
			}
			if !matched {
				t.Errorf("items = %#v, want matching %q", result.Items, test.want)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.EntogSpcmSearch(ctx, application.EntogSpcmSearchQuery{St: "2", Sw: "Homidia", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.EntogSpcmSearch(ctx, application.EntogSpcmSearchQuery{St: "2", Sw: "Homidia", PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.Items[0].EntogSmplNo == second.Items[0].EntogSmplNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.EntogSpcmSearch(ctx, application.EntogSpcmSearchQuery{St: "2", Sw: "kna-mcp-no-result-20260822", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})

	t.Run("dates are ignored", func(t *testing.T) {
		queries := []application.EntogSpcmSearchQuery{
			{St: "2", Sw: "Homidia", PageNo: 1, NumOfRows: 1},
			{St: "2", Sw: "Homidia", DateGbn: "1", PageNo: 1, NumOfRows: 1},
			{St: "2", Sw: "Homidia", DateFrom: "20200101", PageNo: 1, NumOfRows: 1},
			{St: "2", Sw: "Homidia", DateGbn: "1", DateFrom: "20200101", DateTo: "20201231", PageNo: 1, NumOfRows: 1},
			{St: "2", Sw: "Homidia", DateGbn: "2", DateFrom: "20200101", DateTo: "20201231", PageNo: 1, NumOfRows: 1},
			{St: "2", Sw: "Homidia", DateGbn: "1", DateFrom: "20201301", DateTo: "20201331", PageNo: 1, NumOfRows: 1},
		}

		var first application.EntogSpcmSearchResult
		for i, query := range queries {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			result, err := client.EntogSpcmSearch(ctx, query)
			cancel()
			if err != nil {
				t.Fatal(err)
			}
			if i == 0 {
				first = result
				continue
			}
			if result.TotalCount != first.TotalCount || len(result.Items) != 1 || len(first.Items) != 1 || result.Items[0].EntogSmplNo != first.Items[0].EntogSmplNo {
				t.Errorf("query = %#v, result = %#v, want same as %#v", query, result, first)
			}
		}
	})
}
