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

func TestPlantSeedUnitList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantSeedUnitListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantSeedUnitListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":     "test+/=",
			"pageNo":         "1",
			"numOfRows":      "2",
			"reqSeedSpecsId": "test-seed-species-id",
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
      <cllcnDate>collection date</cllcnDate><plantGnrlNm>plant general name</plantGnrlNm>
      <qualtFllnsRt>quality fullness rate</qualtFllnsRt><sdwghWeght>thousand seed weight</sdwghWeght>
      <seedAdmcn>seed air-dry moisture content</seedAdmcn><seedCllctPlace>seed collection place</seedCllctPlace>
      <seedHoldGrainCnt>seed holding grain count</seedHoldGrainCnt><seedHoldQntt>seed holding quantity</seedHoldQntt>
      <seedNo>seed number</seedNo><seedSpecsId>seed species ID</seedSpecsId>
      <storeChrcrTpcdNm>storage characteristic type</storeChrcrTpcdNm><vtlfct>vitality rate</vtlfct>
      <vtlfctTestYr>vitality test year</vtlfctTestYr>
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

	got, err := client.PlantSeedUnitList(context.Background(), application.PlantSeedUnitListQuery{
		PageNo:         1,
		NumOfRows:      2,
		ReqSeedSpecsID: "test-seed-species-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantSeedUnitListResult{
		Items: []application.PlantSeedUnitListItem{{
			CllcnDate:        "collection date",
			PlantGnrlNm:      "plant general name",
			QualtFllnsRt:     "quality fullness rate",
			SdwghWeght:       "thousand seed weight",
			SeedAdmcn:        "seed air-dry moisture content",
			SeedCllctPlace:   "seed collection place",
			SeedHoldGrainCnt: "seed holding grain count",
			SeedHoldQntt:     "seed holding quantity",
			SeedNo:           "seed number",
			SeedSpecsID:      "seed species ID",
			StoreChrcrTpcdNm: "storage characteristic type",
			Vtlfct:           "vitality rate",
			VtlfctTestYr:     "vitality test year",
		}},
		NumOfRows:  2,
		PageNo:     1,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantSeedUnitListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantSeedUnitList(context.Background(), application.PlantSeedUnitListQuery{PageNo: 1, NumOfRows: 2, ReqSeedSpecsID: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantSeedUnitListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantSeedUnitList(context.Background(), application.PlantSeedUnitListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
			var apiError *PlantSeedUnitListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantSeedUnitListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantSeedUnitListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantSeedUnitList(context.Background(), application.PlantSeedUnitListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
	var apiError *PlantSeedUnitListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantSeedUnitListError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantSeedUnitListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantSeedUnitList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantSeedUnitList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantSeedUnitList: response missing resultCode"},
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

			_, err = client.PlantSeedUnitList(context.Background(), application.PlantSeedUnitListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantSeedUnitListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name           string
		pageNo         int
		numOfRows      int
		reqSeedSpecsID string
	}{
		{name: "first seed species", pageNo: 1, numOfRows: 1, reqSeedSpecsID: "SS0002847"},
		{name: "second seed species and changed number of rows", pageNo: 1, numOfRows: 2, reqSeedSpecsID: "SS0003223"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantSeedUnitList(ctx, application.PlantSeedUnitListQuery{
				PageNo:         test.pageNo,
				NumOfRows:      test.numOfRows,
				ReqSeedSpecsID: test.reqSeedSpecsID,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantSeedUnitList returned no items")
			}
			if result.Items[0].SeedSpecsID != test.reqSeedSpecsID {
				t.Errorf("seedSpecsId = %q, want %q", result.Items[0].SeedSpecsID, test.reqSeedSpecsID)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantSeedUnitList(ctx, application.PlantSeedUnitListQuery{
			PageNo:         2,
			NumOfRows:      1,
			ReqSeedSpecsID: "SS0002847",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 1 {
			t.Errorf("result = %#v, want empty second page with totalCount 1", result)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantSeedUnitList(ctx, application.PlantSeedUnitListQuery{
			PageNo:         1,
			NumOfRows:      1,
			ReqSeedSpecsID: "KNA-MCP-NO-RESULT",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
