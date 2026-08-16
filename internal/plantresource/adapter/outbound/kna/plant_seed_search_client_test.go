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

func TestPlantSeedSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantSeedSearchPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantSeedSearchPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "1",
			"numOfRows":    "2",
			"reqSearchWrd": "test-search-word",
			"dateFrom":     "test-date-from",
			"dateTo":       "test-date-to",
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
      <apgFamilyKorNm>apg family Korean name</apgFamilyKorNm><apgFamilyNm>apg family name</apgFamilyNm>
      <blprdEnmnt>bloom period end month</blprdEnmnt><blprdStmnt>bloom period start month</blprdStmnt>
      <clrngMthodCdNm>cleaning method</clrngMthodCdNm><familyKorNm>family Korean name</familyKorNm>
      <familyNm>family name</familyNm><fritCdNm>fruit type</fritCdNm>
      <frssnEnmnt>fruit season end month</frssnEnmnt><frssnStmnt>fruit season start month</frssnStmnt>
      <lastUpdtDtm>last update date time</lastUpdtDtm><plantGnrlNm>plant general name</plantGnrlNm>
      <plantSpecsScnm>plant species scientific name</plantSpecsScnm><rfrncLtrtrCont>reference literature</rfrncLtrtrCont>
      <seedCtsrfcDesc>seed surface description</seedCtsrfcDesc><seedCtsrfcTpcdNm>seed surface type</seedCtsrfcTpcdNm>
      <seedEmbrTpcdNm>seed embryo type</seedEmbrTpcdNm><seedMnmmBrdth>seed minimum breadth</seedMnmmBrdth>
      <seedMnmmLngth>seed minimum length</seedMnmmLngth><seedMxmmBrdth>seed maximum breadth</seedMxmmBrdth>
      <seedMxmmLngth>seed maximum length</seedMxmmLngth><seedShpDesc>seed shape description</seedShpDesc>
      <seedShpTpcdNm>seed shape type</seedShpTpcdNm><seedSpecsId>seed species ID</seedSpecsId>
      <seedTpcdNm>seed type</seedTpcdNm>
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

	got, err := client.PlantSeedSearch(context.Background(), application.PlantSeedSearchQuery{
		PageNo:       1,
		NumOfRows:    2,
		ReqSearchWrd: "test-search-word",
		DateFrom:     "test-date-from",
		DateTo:       "test-date-to",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantSeedSearchResult{
		Items: []application.PlantSeedSearchItem{{
			APGFamilyKorNm:   "apg family Korean name",
			APGFamilyNm:      "apg family name",
			BlprdEnmnt:       "bloom period end month",
			BlprdStmnt:       "bloom period start month",
			ClrngMthodCdNm:   "cleaning method",
			FamilyKorNm:      "family Korean name",
			FamilyNm:         "family name",
			FritCdNm:         "fruit type",
			FrssnEnmnt:       "fruit season end month",
			FrssnStmnt:       "fruit season start month",
			LastUpdtDtm:      "last update date time",
			PlantGnrlNm:      "plant general name",
			PlantSpecsScnm:   "plant species scientific name",
			RfrncLtrtrCont:   "reference literature",
			SeedCtsrfcDesc:   "seed surface description",
			SeedCtsrfcTpcdNm: "seed surface type",
			SeedEmbrTpcdNm:   "seed embryo type",
			SeedMnmmBrdth:    "seed minimum breadth",
			SeedMnmmLngth:    "seed minimum length",
			SeedMxmmBrdth:    "seed maximum breadth",
			SeedMxmmLngth:    "seed maximum length",
			SeedShpDesc:      "seed shape description",
			SeedShpTpcdNm:    "seed shape type",
			SeedSpecsID:      "seed species ID",
			SeedTpcdNm:       "seed type",
		}},
		NumOfRows:  2,
		PageNo:     1,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantSeedSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantSeedSearch(context.Background(), application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantSeedSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantSeedSearch(context.Background(), application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 1})
			var apiError *PlantSeedSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantSeedSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantSeedSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantSeedSearch(context.Background(), application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 1})
	var apiError *PlantSeedSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantSeedSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantSeedSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantSeedSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantSeedSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantSeedSearch: response missing resultCode"},
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

			_, err = client.PlantSeedSearch(context.Background(), application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantSeedSearchLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

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
		{name: "exact Korean name", pageNo: 1, numOfRows: 10, reqSearchWrd: "소나무", plantGnrlNm: "소나무"},
		{name: "partial Korean name", pageNo: 1, numOfRows: 10, reqSearchWrd: "소나", plantGnrlNm: "소나무"},
		{name: "uppercase scientific name", pageNo: 1, numOfRows: 10, reqSearchWrd: "Pinus", plantSpecsScnm: "pinus"},
		{name: "lowercase scientific name", pageNo: 1, numOfRows: 10, reqSearchWrd: "pinus", plantSpecsScnm: "pinus"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantSeedSearch(ctx, application.PlantSeedSearchQuery{
				PageNo:       test.pageNo,
				NumOfRows:    test.numOfRows,
				ReqSearchWrd: test.reqSearchWrd,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantSeedSearch returned no items")
			}
			if result.Items[0].PlantGnrlNm == "" || result.Items[0].SeedSpecsID == "" {
				t.Errorf("plantSeedSearch item = %#v, want plantGnrlNm and seedSpecsId", result.Items[0])
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
					t.Errorf("plantSeedSearch items = %#v, want matching result", result.Items)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		first, err := client.PlantSeedSearch(ctx, application.PlantSeedSearchQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantSeedSearch(ctx, application.PlantSeedSearchQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || first.Items[0].SeedSpecsID == second.Items[0].SeedSpecsID {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantSeedSearch(ctx, application.PlantSeedSearchQuery{
			PageNo:       1,
			NumOfRows:    1,
			ReqSearchWrd: "kna-mcp-no-result-20260816",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
