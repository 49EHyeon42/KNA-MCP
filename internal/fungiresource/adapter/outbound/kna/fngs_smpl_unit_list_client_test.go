package kna

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
)

func TestFngsSmplUnitList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != fngsSmplUnitListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, fngsSmplUnitListPath)
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"pageNo":     "2",
			"numOfRows":  "1",
			"reqFngsId":  "test-fungi-id",
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
      <clarDtlDscrt>collection site detail</clarDtlDscrt><clarHaslvVal> </clarHaslvVal>
      <cllcrNm>collector name</cllcrNm><familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <fngsEclgTpcdNm>fungi ecology type code name</fngsEclgTpcdNm><fngsGnrlNm>fungi general name</fngsGnrlNm>
      <fngsId>fungi ID</fngsId><fngsScnm>fungi scientific name</fngsScnm>
      <fngsSmplKindCdNm>fungi sample kind code name</fngsSmplKindCdNm><fngsSmplNo>fungi sample number</fngsSmplNo>
      <genusKorNm>genus Korean name</genusKorNm><genusNm>genus name</genusNm>
      <hbttChrcrCont>habitat characteristic content</hbttChrcrCont><hstCont>host content</hstCont>
      <lastUpdtDtm>last update date time</lastUpdtDtm><smplCllcnDt>sample collection date</smplCllcnDt>
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

	got, err := client.FngsSmplUnitList(context.Background(), application.FngsSmplUnitListQuery{
		PageNo:    2,
		NumOfRows: 1,
		ReqFngsID: "test-fungi-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.FngsSmplUnitListResult{
		Items: []application.FngsSmplUnitListItem{{
			ClarDtlDscrt:     "collection site detail",
			ClarHaslvVal:     " ",
			CllcrNm:          "collector name",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			FngsEclgTpcdNm:   "fungi ecology type code name",
			FngsGnrlNm:       "fungi general name",
			FngsID:           "fungi ID",
			FngsScnm:         "fungi scientific name",
			FngsSmplKindCdNm: "fungi sample kind code name",
			FngsSmplNo:       "fungi sample number",
			GenusKorNm:       "genus Korean name",
			GenusNm:          "genus name",
			HbttChrcrCont:    "habitat characteristic content",
			HstCont:          "host content",
			LastUpdtDtm:      "last update date time",
			SmplCllcnDt:      "sample collection date",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestFngsSmplUnitListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.FngsSmplUnitList(context.Background(), application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "unknown-id"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestFngsSmplUnitListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.FngsSmplUnitList(context.Background(), application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "test-id"})
			var apiError *FngsSmplUnitListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *FngsSmplUnitListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestFngsSmplUnitListReturnsGatewayError(t *testing.T) {
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

	_, err = client.FngsSmplUnitList(context.Background(), application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "test-id"})
	var apiError *FngsSmplUnitListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *FngsSmplUnitListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestFngsSmplUnitListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "fngsSmplUnitList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "fngsSmplUnitList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "fngsSmplUnitList: response missing resultCode"},
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

			_, err = client.FngsSmplUnitList(context.Background(), application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqFngsID: "test-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestFngsSmplUnitListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	searchResult, err := client.FngsSmplSearch(ctx, application.FngsSmplSearchQuery{PageNo: 1, NumOfRows: 10})
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult.Items) < 2 {
		t.Fatalf("fngsSmplSearch items = %d, want at least 2", len(searchResult.Items))
	}

	for index, searchItem := range searchResult.Items[:2] {
		t.Run(fmt.Sprintf("search result %d", index+1), func(t *testing.T) {
			wantTotalCount, err := strconv.Atoi(searchItem.Cnt)
			if err != nil {
				t.Fatalf("cnt = %q: %v", searchItem.Cnt, err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			result, err := client.FngsSmplUnitList(ctx, application.FngsSmplUnitListQuery{
				PageNo:    1,
				NumOfRows: 2,
				ReqFngsID: searchItem.FngsID,
			})
			if err != nil {
				t.Fatal(err)
			}
			wantItems := min(wantTotalCount, 2)
			if len(result.Items) != wantItems || result.TotalCount != wantTotalCount {
				t.Fatalf("result = %#v, want %d items and totalCount %d", result, wantItems, wantTotalCount)
			}
			for _, item := range result.Items {
				if item.FngsID != searchItem.FngsID || item.FngsSmplNo == "" {
					t.Errorf("item = %#v, want fngsId %q and fngsSmplNo", item, searchItem.FngsID)
				}
			}
		})
	}

	var pagedFngsID string
	for _, item := range searchResult.Items {
		count, err := strconv.Atoi(item.Cnt)
		if err == nil && count >= 4 {
			pagedFngsID = item.FngsID
			break
		}
	}
	if pagedFngsID == "" {
		t.Fatal("fngsSmplSearch returned no item with at least 4 samples")
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		first, err := client.FngsSmplUnitList(ctx, application.FngsSmplUnitListQuery{PageNo: 1, NumOfRows: 2, ReqFngsID: pagedFngsID})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.FngsSmplUnitList(ctx, application.FngsSmplUnitListQuery{PageNo: 2, NumOfRows: 2, ReqFngsID: pagedFngsID})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 2 || len(second.Items) != 2 || first.TotalCount != second.TotalCount || first.Items[0].FngsSmplNo == second.Items[0].FngsSmplNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		result, err := client.FngsSmplUnitList(ctx, application.FngsSmplUnitListQuery{
			PageNo:    1,
			NumOfRows: 1,
			ReqFngsID: "F999999999",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
