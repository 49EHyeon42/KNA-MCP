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
	"testing"
	"time"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
)

func TestPlantPictorialBookSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantPilbkSearchPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantPilbkSearchPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":   "test+/=",
			"pageNo":       "1",
			"numOfRows":    "2",
			"reqSearchWrd": "소나무",
			"dateFrom":     "20250101",
			"dateTo":       "20251231",
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
      <apgFamilyKorNm>소나무과</apgFamilyKorNm><apgFamilyNm>Pinaceae</apgFamilyNm>
      <familyKorNm>소나무과</familyKorNm><familyNm>Pinaceae</familyNm>
      <genusKorNm>소나무속</genusKorNm><genusNm>Pinus</genusNm>
      <lastUpdtDtm>20040323</lastUpdtDtm><notRcmmGnrlNm>세잎소나무,삼엽송</notRcmmGnrlNm>
      <plantGnrlNm>리기다소나무</plantGnrlNm><plantPilbkNo>31665</plantPilbkNo>
      <plantSpecsScnm>Pinus rigida Mill.</plantSpecsScnm>
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

	got, err := client.PlantPictorialBookSearch(context.Background(), application.PlantPictorialBookSearchQuery{
		PageNumber:        1,
		NumberOfRows:      2,
		RequestSearchWord: "소나무",
		DateFrom:          "20250101",
		DateTo:            "20251231",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantPictorialBookSearchResult{
		Items: []application.PlantPictorialBookSearchItem{{
			APGFamilyKoreanName:        "소나무과",
			APGFamilyName:              "Pinaceae",
			FamilyKoreanName:           "소나무과",
			FamilyName:                 "Pinaceae",
			GenusKoreanName:            "소나무속",
			GenusName:                  "Pinus",
			LastUpdateDateTime:         "20040323",
			NotRecommendedGeneralName:  "세잎소나무,삼엽송",
			PlantGeneralName:           "리기다소나무",
			PlantPictorialBookNumber:   "31665",
			PlantSpeciesScientificName: "Pinus rigida Mill.",
		}},
		NumberOfRows: 2,
		PageNumber:   1,
		TotalCount:   7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantPictorialBookSearchReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantPictorialBookSearch(context.Background(), application.PlantPictorialBookSearchQuery{PageNumber: 1, NumberOfRows: 1})
			var apiError *PlantPilbkSearchError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantPilbkSearchError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantPictorialBookSearchReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantPictorialBookSearch(context.Background(), application.PlantPictorialBookSearchQuery{PageNumber: 1, NumberOfRows: 1})
	var apiError *PlantPilbkSearchError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantPilbkSearchError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantPictorialBookSearchLive(t *testing.T) {
	serviceKey := os.Getenv("DATA_GO_KR_SERVICE_KEY")
	if serviceKey == "" {
		t.Skip("DATA_GO_KR_SERVICE_KEY is not set")
	}

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.PlantPictorialBookSearch(ctx, application.PlantPictorialBookSearchQuery{
		PageNumber:        1,
		NumberOfRows:      1,
		RequestSearchWord: "소나무",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) == 0 {
		t.Fatal("plantPilbkSearch returned no items")
	}
}
