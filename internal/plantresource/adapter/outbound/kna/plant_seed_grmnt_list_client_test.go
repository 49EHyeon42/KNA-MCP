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

func TestPlantSeedGrmntList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/PlantResource/plantSeedGrmntList" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/PlantResource/plantSeedGrmntList")
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
      <avrgGrmntDcnt>average germination day count</avrgGrmntDcnt>
      <grmntBfrPrcesCont>germination before processing content</grmntBfrPrcesCont>
      <grmntClmdmCont>germination culture medium content</grmntClmdmCont>
      <grmntDscrt>germination description</grmntDscrt>
      <grmntExprmNo>germination experiment number</grmntExprmNo>
      <grmntExprmSeq>germination experiment sequence</grmntExprmSeq>
      <grmntLightCndtn>germination light condition</grmntLightCndtn>
      <grmntRt>germination rate</grmntRt>
      <grmntTmpCndtn>germination temperature condition</grmntTmpCndtn>
      <plantGnrlNm>plant general name</plantGnrlNm>
      <seedNo>seed number</seedNo>
      <seedSpecsId>seed species ID</seedSpecsId>
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

	got, err := client.PlantSeedGrmntList(context.Background(), application.PlantSeedGrmntListQuery{
		PageNo:         1,
		NumOfRows:      2,
		ReqSeedSpecsID: "test-seed-species-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantSeedGrmntListResult{
		Items: []application.PlantSeedGrmntListItem{{
			AvrgGrmntDcnt:     "average germination day count",
			GrmntBfrPrcesCont: "germination before processing content",
			GrmntClmdmCont:    "germination culture medium content",
			GrmntDscrt:        "germination description",
			GrmntExprmNo:      "germination experiment number",
			GrmntExprmSeq:     "germination experiment sequence",
			GrmntLightCndtn:   "germination light condition",
			GrmntRt:           "germination rate",
			GrmntTmpCndtn:     "germination temperature condition",
			PlantGnrlNm:       "plant general name",
			SeedNo:            "seed number",
			SeedSpecsID:       "seed species ID",
		}},
		NumOfRows:  2,
		PageNo:     1,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantSeedGrmntListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantSeedGrmntList(context.Background(), application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 2, ReqSeedSpecsID: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantSeedGrmntListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantSeedGrmntList(context.Background(), application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
			var apiError *PlantSeedGrmntListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantSeedGrmntListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantSeedGrmntListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantSeedGrmntList(context.Background(), application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
	var apiError *PlantSeedGrmntListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantSeedGrmntListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantSeedGrmntListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantSeedGrmntList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantSeedGrmntList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantSeedGrmntList: response missing resultCode"},
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

			_, err = client.PlantSeedGrmntList(context.Background(), application.PlantSeedGrmntListQuery{PageNo: 1, NumOfRows: 1, ReqSeedSpecsID: "test-seed-species-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantSeedGrmntListLive(t *testing.T) {
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
		{name: "second seed species and changed number of rows", pageNo: 1, numOfRows: 2, reqSeedSpecsID: "SS0003338"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.PlantSeedGrmntList(ctx, application.PlantSeedGrmntListQuery{
				PageNo:         test.pageNo,
				NumOfRows:      test.numOfRows,
				ReqSeedSpecsID: test.reqSeedSpecsID,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantSeedGrmntList returned no items")
			}
			if result.Items[0].SeedSpecsID != test.reqSeedSpecsID {
				t.Errorf("seedSpecsId = %q, want %q", result.Items[0].SeedSpecsID, test.reqSeedSpecsID)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.PlantSeedGrmntList(ctx, application.PlantSeedGrmntListQuery{
			PageNo:         2,
			NumOfRows:      1,
			ReqSeedSpecsID: "SS0003338",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 1 || result.TotalCount < 2 || result.PageNo != 2 || result.Items[0].SeedSpecsID != "SS0003338" {
			t.Errorf("result = %#v, want one matching item on the second page", result)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.PlantSeedGrmntList(ctx, application.PlantSeedGrmntListQuery{
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
