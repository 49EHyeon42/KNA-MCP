package kna

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
)

func TestPlantFolkAreaList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantFolkAreaListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantFolkAreaListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"pageNo":     "2",
			"numOfRows":  "1",
			"flpltId":    "test-folk-plant-id",
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
      <flcstPlantExmnnAraTpcdNm>folk plant examination area type name</flcstPlantExmnnAraTpcdNm>
      <flcstPlantLcltDscrt>folk plant locality description</flcstPlantLcltDscrt>
      <flcstPlantPrpseDscrt>folk plant purpose description</flcstPlantPrpseDscrt>
      <flpltId>folk plant ID</flpltId>
      <plantBrdgFomTpcdNm>plant breeding form type name</plantBrdgFomTpcdNm>
      <plantGnrlNm>plant general name</plantGnrlNm>
      <plantSpecsScnm>plant species scientific name</plantSpecsScnm>
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

	got, err := client.PlantFolkAreaList(context.Background(), application.PlantFolkAreaListQuery{
		PageNo:    2,
		NumOfRows: 1,
		FlpltID:   "test-folk-plant-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantFolkAreaListResult{
		Items: []application.PlantFolkAreaListItem{{
			FlcstPlantExmnnAraTpcdNm: "folk plant examination area type name",
			FlcstPlantLcltDscrt:      "folk plant locality description",
			FlcstPlantPrpseDscrt:     "folk plant purpose description",
			FlpltID:                  "folk plant ID",
			PlantBrdgFomTpcdNm:       "plant breeding form type name",
			PlantGnrlNm:              "plant general name",
			PlantSpecsScnm:           "plant species scientific name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantFolkAreaListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantFolkAreaList(context.Background(), application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantFolkAreaListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantFolkAreaList(context.Background(), application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "test-folk-plant-id"})
			var apiError *PlantFolkAreaListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantFolkAreaListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantFolkAreaListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantFolkAreaList(context.Background(), application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "test-folk-plant-id"})
	var apiError *PlantFolkAreaListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantFolkAreaListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantFolkAreaListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantFolkAreaList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantFolkAreaList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantFolkAreaList: response missing resultCode"},
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

			_, err = client.PlantFolkAreaList(context.Background(), application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "test-folk-plant-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantFolkAreaListLive(t *testing.T) {
	serviceKey := os.Getenv("DATA_GO_KR_SERVICE_KEY")
	if serviceKey == "" {
		t.Skip("DATA_GO_KR_SERVICE_KEY is not set")
	}

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name        string
		flpltID     string
		plantGnrlNm string
	}{
		{name: "first folk plant", flpltID: "2014000403", plantGnrlNm: "가는금불초"},
		{name: "second folk plant", flpltID: "2014000404", plantGnrlNm: "금불초"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantFolkAreaList(ctx, application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 10, FlpltID: test.flpltID})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) != 1 || result.TotalCount != 1 || result.Items[0].FlpltID != test.flpltID || result.Items[0].PlantGnrlNm != test.plantGnrlNm {
				t.Errorf("result = %#v, want one %s record for %s", result, test.plantGnrlNm, test.flpltID)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		first, err := client.PlantFolkAreaList(ctx, application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "2014000365"})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantFolkAreaList(ctx, application.PlantFolkAreaListQuery{PageNo: 2, NumOfRows: 1, FlpltID: "2014000365"})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != 5 || second.TotalCount != 5 || first.Items[0].FlcstPlantExmnnAraTpcdNm == second.Items[0].FlcstPlantExmnnAraTpcdNm {
			t.Errorf("first = %#v, second = %#v, want distinct pages with totalCount 5", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantFolkAreaList(ctx, application.PlantFolkAreaListQuery{PageNo: 1, NumOfRows: 1, FlpltID: "KNA-MCP-NO-RESULT"})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
