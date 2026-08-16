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
			"reqSearchWrd": "test-search-word",
			"dateFrom":     "test-date-from",
			"dateTo":       "test-date-to",
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
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <genusKorNm>genus Korean name</genusKorNm><genusNm>genus name</genusNm>
      <lastUpdtDtm>last update date time</lastUpdtDtm><notRcmmGnrlNm>not recommended general name</notRcmmGnrlNm>
      <plantGnrlNm>plant general name</plantGnrlNm><plantPilbkNo>plant pictorial book number</plantPilbkNo>
      <plantSpecsScnm>plant species scientific name</plantSpecsScnm>
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
		RequestSearchWord: "test-search-word",
		DateFrom:          "test-date-from",
		DateTo:            "test-date-to",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantPictorialBookSearchResult{
		Items: []application.PlantPictorialBookSearchItem{{
			APGFamilyKoreanName:        "apg family Korean name",
			APGFamilyName:              "apg family name",
			FamilyKoreanName:           "family Korean name",
			FamilyName:                 "family name",
			GenusKoreanName:            "genus Korean name",
			GenusName:                  "genus name",
			LastUpdateDateTime:         "last update date time",
			NotRecommendedGeneralName:  "not recommended general name",
			PlantGeneralName:           "plant general name",
			PlantPictorialBookNumber:   "plant pictorial book number",
			PlantSpeciesScientificName: "plant species scientific name",
		}},
		NumberOfRows: 2,
		PageNumber:   1,
		TotalCount:   7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantPictorialBookSearchReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantPictorialBookSearch(context.Background(), application.PlantPictorialBookSearchQuery{PageNumber: 1, NumberOfRows: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
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

func TestPlantPictorialBookSearchReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantPilbkSearch: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantPilbkSearch: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantPilbkSearch: response missing resultCode"},
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

			_, err = client.PlantPictorialBookSearch(context.Background(), application.PlantPictorialBookSearchQuery{PageNumber: 1, NumberOfRows: 1})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
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

	for _, test := range []struct {
		name              string
		pageNumber        int
		numberOfRows      int
		requestSearchWord string
	}{
		{name: "without search word", pageNumber: 1, numberOfRows: 1},
		{name: "with search word and changed pagination", pageNumber: 2, numberOfRows: 2, requestSearchWord: "소나무"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantPictorialBookSearch(ctx, application.PlantPictorialBookSearchQuery{
				PageNumber:        test.pageNumber,
				NumberOfRows:      test.numberOfRows,
				RequestSearchWord: test.requestSearchWord,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantPilbkSearch returned no items")
			}
		})
	}

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantPictorialBookSearch(ctx, application.PlantPictorialBookSearchQuery{
			PageNumber:        1,
			NumberOfRows:      1,
			RequestSearchWord: "kna-mcp-no-result-20260816",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
