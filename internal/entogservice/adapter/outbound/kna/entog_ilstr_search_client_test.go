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

func TestEntogIlstrSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/EntogService/entogIlstrSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/EntogService/entogIlstrSearch")
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"st":         "2",
			"sw":         "test-search-word",
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
      <btnc>test-scientific-name</btnc><cprtCtnt>test-copyright</cprtCtnt>
      <detailYn>test-detail-availability</detailYn><entogOfnmKrlngNm>test-Korean-name</entogOfnmKrlngNm>
      <entogOfnmScnmId>test-scientific-name-ID</entogOfnmScnmId><entogPilbkNo>test-pictorial-book-number</entogPilbkNo>
      <familyKorNm>test-family-Korean-name</familyKorNm><familyNm>test-family-name</familyNm>
      <genusKorNm>test-genus-Korean-name</genusKorNm><genusNm>test-genus-name</genusNm>
      <imgUrl>test-image-URL</imgUrl><ordKorNm>test-order-Korean-name</ordKorNm><ordNm>test-order-name</ordNm>
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

	got, err := client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{
		St:        "2",
		Sw:        "test-search-word",
		NumOfRows: 1,
		PageNo:    2,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.EntogIlstrSearchResult{
		Items: []application.EntogIlstrSearchItem{{
			Btnc:             "test-scientific-name",
			CprtCtnt:         "test-copyright",
			DetailYn:         "test-detail-availability",
			EntogOfnmKrlngNm: "test-Korean-name",
			EntogOfnmScnmID:  "test-scientific-name-ID",
			EntogPilbkNo:     "test-pictorial-book-number",
			FamilyKorNm:      "test-family-Korean-name",
			FamilyNm:         "test-family-name",
			GenusKorNm:       "test-genus-Korean-name",
			GenusNm:          "test-genus-name",
			ImgURL:           "test-image-URL",
			OrdKorNm:         "test-order-Korean-name",
			OrdNm:            "test-order-name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestEntogIlstrSearchPreservesObservedValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items><item><genusKorNm> </genusKorNm><imgUrl>NONE</imgUrl></item></items><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>1</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].GenusKorNm != " " || result.Items[0].ImgURL != "NONE" {
		t.Errorf("result = %#v, want preserved single space and NONE", result)
	}
}

func TestEntogIlstrSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{St: "1", Sw: "none", PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestEntogIlstrSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			var apiError *EntogIlstrSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *EntogIlstrSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestEntogIlstrSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
	var apiError *EntogIlstrSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *EntogIlstrSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestEntogIlstrSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "entogIlstrSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "entogIlstrSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "entogIlstrSearch: response missing resultCode"},
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

			_, err = client.EntogIlstrSearch(context.Background(), application.EntogIlstrSearchQuery{St: "1", Sw: "test-search-word", PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestEntogIlstrSearchLive(t *testing.T) {
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
		{name: "partial Korean name", st: "1", sw: "가시발", want: "가시발톡토기"},
		{name: "exact Korean name", st: "3", sw: "가시발톡토기", want: "가시발톡토기"},
		{name: "partial scientific name", st: "2", sw: "Coecobrya", want: "Coecobrya"},
		{name: "case insensitive partial scientific name", st: "2", sw: "coecobrya", want: "Coecobrya"},
		{name: "exact scientific name", st: "4", sw: "Coecobrya dubiosa (Yosii, 1956)", want: "Coecobrya dubiosa (Yosii, 1956)"},
		{name: "case insensitive exact scientific name", st: "4", sw: "coecobrya dubiosa (yosii, 1956)", want: "Coecobrya dubiosa (Yosii, 1956)"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.EntogIlstrSearch(ctx, application.EntogIlstrSearchQuery{St: test.st, Sw: test.sw, PageNo: 1, NumOfRows: 10})
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

		first, err := client.EntogIlstrSearch(ctx, application.EntogIlstrSearchQuery{St: "1", Sw: "톡토기", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.EntogIlstrSearch(ctx, application.EntogIlstrSearchQuery{St: "1", Sw: "톡토기", PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.Items[0].EntogPilbkNo == second.Items[0].EntogPilbkNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.EntogIlstrSearch(ctx, application.EntogIlstrSearchQuery{St: "2", Sw: "kna-mcp-no-result-20260822", PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
