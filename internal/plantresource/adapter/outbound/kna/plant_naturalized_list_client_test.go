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

func TestPlantNaturalizedList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantNaturalizedListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantNaturalizedListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "2",
			"numOfRows":    "1",
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
      <agpFamilyNm>agp family name</agpFamilyNm><apgFamilyKorNm>apg family Korean name</apgFamilyKorNm>
      <blprdEnmnt>bloom period end month</blprdEnmnt><blprdStmnt>bloom period start month</blprdStmnt>
      <distrAraDscrt>distribution area description</distrAraDscrt><eclgDstrbYn>ecological disturbance yes or no</eclgDstrbYn>
      <extcPlantCdNm>exotic plant code name</extcPlantCdNm><familyKorNm>family Korean name</familyKorNm>
      <familyNm>family name</familyNm><frtTpcdNm>fruit type code name</frtTpcdNm>
      <lastUpdtDtm>last update date time</lastUpdtDtm><ntldgTpcdNm>naturalization degree type name</ntldgTpcdNm>
      <ntrlzEraTpcdNm>naturalization era type name</ntrlzEraTpcdNm><orplcNm>original place name</orplcNm>
      <plantBrdgFomTpcdNm>plant breeding form type name</plantBrdgFomTpcdNm>
      <plantDistrGrcd>plant distribution grade code</plantDistrGrcd><plantDistrQntt>plant distribution quantity</plantDistrQntt>
      <plantDistrQnttGrcd>plant distribution quantity grade code</plantDistrQnttGrcd>
      <plantEngNm>plant English name</plantEngNm><plantGnrlNm>plant general name</plantGnrlNm>
      <plantJpnNm>plant Japanese name</plantJpnNm><plantLfcclTpcdNm>plant life cycle type name</plantLfcclTpcdNm>
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

	got, err := client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{
		PageNo:       2,
		NumOfRows:    1,
		ReqSearchWrd: "test-search-word",
		DateFrom:     "test-date-from",
		DateTo:       "test-date-to",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantNaturalizedListResult{
		Items: []application.PlantNaturalizedListItem{{
			AgpFamilyNm:        "agp family name",
			APGFamilyKorNm:     "apg family Korean name",
			BlprdEnmnt:         "bloom period end month",
			BlprdStmnt:         "bloom period start month",
			DistrAraDscrt:      "distribution area description",
			EclgDstrbYn:        "ecological disturbance yes or no",
			ExtcPlantCdNm:      "exotic plant code name",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			FrtTpcdNm:          "fruit type code name",
			LastUpdtDtm:        "last update date time",
			NtldgTpcdNm:        "naturalization degree type name",
			NtrlzEraTpcdNm:     "naturalization era type name",
			OrplcNm:            "original place name",
			PlantBrdgFomTpcdNm: "plant breeding form type name",
			PlantDistrGrcd:     "plant distribution grade code",
			PlantDistrQntt:     "plant distribution quantity",
			PlantDistrQnttGrcd: "plant distribution quantity grade code",
			PlantEngNm:         "plant English name",
			PlantGnrlNm:        "plant general name",
			PlantJpnNm:         "plant Japanese name",
			PlantLfcclTpcdNm:   "plant life cycle type name",
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

func TestPlantNaturalizedListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		for _, key := range []string{"reqSearchWrd", "dateFrom", "dateTo"} {
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

	result, err := client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantNaturalizedListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1})
			var apiError *PlantNaturalizedListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantNaturalizedListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantNaturalizedListReturnsObservedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>99</resultCode><resultMsg>ORA-00908: missing NULL keyword</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1, DateFrom: "20230221"})
	var apiError *PlantNaturalizedListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantNaturalizedListError", err)
	}
	if apiError.HTTPStatus != http.StatusOK || apiError.Code != "99" || apiError.Message != "ORA-00908: missing NULL keyword" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantNaturalizedListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1})
	var apiError *PlantNaturalizedListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantNaturalizedListError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantNaturalizedListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantNaturalizedList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantNaturalizedList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantNaturalizedList: response missing resultCode"},
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

			_, err = client.PlantNaturalizedList(context.Background(), application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantNaturalizedListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name           string
		reqSearchWrd   string
		plantGnrlNm    string
		plantSpecsScnm string
	}{
		{name: "without search word"},
		{name: "exact Korean name", reqSearchWrd: "가는끈끈이장구채", plantGnrlNm: "가는끈끈이장구채"},
		{name: "partial Korean name", reqSearchWrd: "끈끈이", plantGnrlNm: "가는끈끈이장구채"},
		{name: "uppercase scientific name", reqSearchWrd: "Silene", plantSpecsScnm: "silene"},
		{name: "lowercase scientific name", reqSearchWrd: "silene", plantSpecsScnm: "silene"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantNaturalizedList(ctx, application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 10, ReqSearchWrd: test.reqSearchWrd})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantNaturalizedList returned no items")
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
					t.Errorf("plantNaturalizedList items = %#v, want matching result", result.Items)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		first, err := client.PlantNaturalizedList(ctx, application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.PlantNaturalizedList(ctx, application.PlantNaturalizedListQuery{PageNo: 2, NumOfRows: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 1 || len(second.Items) != 1 || first.TotalCount != second.TotalCount || first.TotalCount < 2 || first.Items[0].PlantGnrlNm == second.Items[0].PlantGnrlNm {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantNaturalizedList(ctx, application.PlantNaturalizedListQuery{PageNo: 1, NumOfRows: 1, ReqSearchWrd: "KNA-MCP-NO-RESULT"})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
