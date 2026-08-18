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

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

func TestChildPilbkSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/ChildService2/childPilbkSearch" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/ChildService2/childPilbkSearch")
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
      <biogyNm>biology name</biogyNm><childLvbngPilbkNo>child pictorial book number</childLvbngPilbkNo>
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <lvbngTpcdNm>living thing type code name</lvbngTpcdNm><lvngKrlngNm>living thing Korean name</lvngKrlngNm>
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

	got, err := client.ChildPilbkSearch(context.Background(), application.ChildPilbkSearchQuery{
		PageNo:       2,
		NumOfRows:    1,
		ReqSearchWrd: "test-search-word",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ChildPilbkSearchResult{
		Items: []application.ChildPilbkSearchItem{{
			BiogyNm:           "biology name",
			ChildLvbngPilbkNo: "child pictorial book number",
			FamilyKorNm:       "family Korean name",
			FamilyNm:          "family name",
			LvbngTpcdNm:       "living thing type code name",
			LvngKrlngNm:       "living thing Korean name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestChildPilbkSearchReturnsEmptyItems(t *testing.T) {
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

	result, err := client.ChildPilbkSearch(context.Background(), application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestChildPilbkSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.ChildPilbkSearch(context.Background(), application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *ChildPilbkSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ChildPilbkSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestChildPilbkSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.ChildPilbkSearch(context.Background(), application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
	var apiError *ChildPilbkSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ChildPilbkSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestChildPilbkSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "childPilbkSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "childPilbkSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "childPilbkSearch: response missing resultCode"},
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

			_, err = client.ChildPilbkSearch(context.Background(), application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestChildPilbkSearchLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name         string
		reqSearchWrd string
		lvngKrlngNm  string
		biogyNm      string
	}{
		{name: "without search word"},
		{name: "exact Korean name", reqSearchWrd: "과립여우갓버섯", lvngKrlngNm: "과립여우갓버섯"},
		{name: "partial Korean name", reqSearchWrd: "여우갓", lvngKrlngNm: "과립여우갓버섯"},
		{name: "uppercase scientific name", reqSearchWrd: "Leucocoprinus", biogyNm: "leucocoprinus"},
		{name: "lowercase scientific name", reqSearchWrd: "leucocoprinus", biogyNm: "leucocoprinus"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{
				PageNo:       1,
				NumOfRows:    10,
				ReqSearchWrd: test.reqSearchWrd,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("childPilbkSearch returned no items")
			}
			if test.lvngKrlngNm != "" || test.biogyNm != "" {
				matched := false
				for _, item := range result.Items {
					if test.lvngKrlngNm != "" && item.LvngKrlngNm == test.lvngKrlngNm {
						matched = true
					}
					if test.biogyNm != "" && strings.Contains(strings.ToLower(item.BiogyNm), test.biogyNm) {
						matched = true
					}
				}
				if !matched {
					t.Errorf("childPilbkSearch items = %#v, want matching result", result.Items)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 2})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{PageNo: 2, NumOfRows: 2})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 2 || len(second.Items) != 2 || first.TotalCount != second.TotalCount || first.TotalCount < 4 || first.Items[0].ChildLvbngPilbkNo == second.Items[0].ChildLvbngPilbkNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{
			PageNo:       1,
			NumOfRows:    1,
			ReqSearchWrd: "kna-mcp-no-result-20260818",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
