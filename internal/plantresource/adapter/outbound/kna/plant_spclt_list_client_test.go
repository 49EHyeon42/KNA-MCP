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

func TestPlantSpcltList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantSpcltListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantSpcltListPath)
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
      <agpFamilyKorNm>agp family Korean name</agpFamilyKorNm><agpFamilyNm>agp family name</agpFamilyNm>
      <extrmCrssScls1Yn>endangered class one yes or no</extrmCrssScls1Yn>
      <extrmCrssScls2Yn>endangered class two yes or no</extrmCrssScls2Yn>
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <plantBrdgFomTpcdNm>plant breeding form type code name</plantBrdgFomTpcdNm>
      <plantGnrlNm>plant general name</plantGnrlNm><plantSpecsScnm>plant species scientific name</plantSpecsScnm>
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

	got, err := client.PlantSpcltList(context.Background(), application.PlantSpcltListQuery{
		PageNo:       2,
		NumOfRows:    1,
		ReqSearchWrd: "test-search-word",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantSpcltListResult{
		Items: []application.PlantSpcltListItem{{
			AgpFamilyKorNm:     "agp family Korean name",
			AgpFamilyNm:        "agp family name",
			ExtrmCrssScls1Yn:   "endangered class one yes or no",
			ExtrmCrssScls2Yn:   "endangered class two yes or no",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			PlantBrdgFomTpcdNm: "plant breeding form type code name",
			PlantGnrlNm:        "plant general name",
			PlantSpecsScnm:     "plant species scientific name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantSpcltListReturnsEmptyItems(t *testing.T) {
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

	result, err := client.PlantSpcltList(context.Background(), application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantSpcltListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantSpcltList(context.Background(), application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1})
			var apiError *PlantSpcltListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantSpcltListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantSpcltListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantSpcltList(context.Background(), application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1})
	var apiError *PlantSpcltListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantSpcltListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantSpcltListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantSpcltList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantSpcltList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantSpcltList: response missing resultCode"},
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

			_, err = client.PlantSpcltList(context.Background(), application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantSpcltListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name             string
		reqSearchWrd     string
		plantGnrlNm      string
		plantSpecsScnm   string
		extrmCrssScls1Yn string
	}{
		{name: "without search word"},
		{name: "exact Korean name", reqSearchWrd: "가거개별꽃", plantGnrlNm: "가거개별꽃"},
		{name: "partial Korean name", reqSearchWrd: "개별꽃", plantGnrlNm: "가거개별꽃"},
		{name: "uppercase scientific name", reqSearchWrd: "Pseudostellaria", plantSpecsScnm: "pseudostellaria"},
		{name: "lowercase scientific name", reqSearchWrd: "pseudostellaria", plantSpecsScnm: "pseudostellaria"},
		{name: "endangered class one", reqSearchWrd: "제주고사리삼", plantGnrlNm: "제주고사리삼", extrmCrssScls1Yn: "Y"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.PlantSpcltList(ctx, application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 10, ReqSearchWrd: test.reqSearchWrd})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantSpcltList returned no items")
			}
			matched := test.plantGnrlNm == "" && test.plantSpecsScnm == ""
			for _, item := range result.Items {
				if test.plantGnrlNm != "" && item.PlantGnrlNm == test.plantGnrlNm && (test.extrmCrssScls1Yn == "" || item.ExtrmCrssScls1Yn == test.extrmCrssScls1Yn) {
					matched = true
				}
				if test.plantSpecsScnm != "" && strings.Contains(strings.ToLower(item.PlantSpecsScnm), test.plantSpecsScnm) {
					matched = true
				}
			}
			if !matched {
				t.Errorf("plantSpcltList items = %#v, want matching result", result.Items)
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.PlantSpcltList(ctx, application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantSpcltList(ctx, application.PlantSpcltListQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || first.Items[0].PlantGnrlNm == second.Items[0].PlantGnrlNm {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.PlantSpcltList(ctx, application.PlantSpcltListQuery{PageNo: 1, NumOfRows: 1, ReqSearchWrd: "KNA-MCP-NO-RESULT"})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
