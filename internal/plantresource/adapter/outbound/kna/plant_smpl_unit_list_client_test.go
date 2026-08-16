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

func TestPlantSmplUnitList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != plantSmplUnitListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, plantSmplUnitListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":      "test+/=",
			"pageNo":          "1",
			"numOfRows":       "2",
			"reqPlantSpecsId": "test-plant-species-id",
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
      <bspcsInsttNm>specimen holding institution</bspcsInsttNm><clarHaslvVal>collection site elevation</clarHaslvVal>
      <clarNm>collection site</clarNm><cllcrNm>collector name</cllcrNm>
      <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
      <hbttChrcrCont>habitat characteristics</hbttChrcrCont><hbttTpcdNm>habitat type</hbttTpcdNm>
      <plantBrdgFomTpcdNm>plant reproductive form</plantBrdgFomTpcdNm><plantGnrlNm>plant general name</plantGnrlNm>
      <plantPilbkNo>plant pictorial book number</plantPilbkNo><plantSmplNo>plant specimen number</plantSmplNo>
      <plantSpecsId>plant species ID</plantSpecsId><plantSpecsScnm>plant species scientific name</plantSpecsScnm>
      <smplCllcnDt>specimen collection date</smplCllcnDt><smplClnyNm>specimen community name</smplClnyNm>
      <smplKindCdNm>specimen type</smplKindCdNm><smplWrdt>specimen preparation date</smplWrdt>
      <vgttnTpeCdNm>vegetation type</vgttnTpeCdNm>
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

	got, err := client.PlantSmplUnitList(context.Background(), application.PlantSmplUnitListQuery{
		PageNo:          1,
		NumOfRows:       2,
		ReqPlantSpecsID: "test-plant-species-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantSmplUnitListResult{
		Items: []application.PlantSmplUnitListItem{{
			AgpFamilyKorNm:     "agp family Korean name",
			AgpFamilyNm:        "agp family name",
			BspcsInsttNm:       "specimen holding institution",
			ClarHaslvVal:       "collection site elevation",
			ClarNm:             "collection site",
			CllcrNm:            "collector name",
			FamilyKorNm:        "family Korean name",
			FamilyNm:           "family name",
			HbttChrcrCont:      "habitat characteristics",
			HbttTpcdNm:         "habitat type",
			PlantBrdgFomTpcdNm: "plant reproductive form",
			PlantGnrlNm:        "plant general name",
			PlantPilbkNo:       "plant pictorial book number",
			PlantSmplNo:        "plant specimen number",
			PlantSpecsID:       "plant species ID",
			PlantSpecsScnm:     "plant species scientific name",
			SmplCllcnDt:        "specimen collection date",
			SmplClnyNm:         "specimen community name",
			SmplKindCdNm:       "specimen type",
			SmplWrdt:           "specimen preparation date",
			VgttnTpeCdNm:       "vegetation type",
		}},
		NumOfRows:  2,
		PageNo:     1,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantSmplUnitListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>2</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantSmplUnitList(context.Background(), application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 2, ReqPlantSpecsID: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantSmplUnitListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantSmplUnitList(context.Background(), application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqPlantSpecsID: "test-plant-species-id"})
			var apiError *PlantSmplUnitListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantSmplUnitListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantSmplUnitListReturnsGatewayError(t *testing.T) {
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

	_, err = client.PlantSmplUnitList(context.Background(), application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqPlantSpecsID: "test-plant-species-id"})
	var apiError *PlantSmplUnitListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantSmplUnitListError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantSmplUnitListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantSmplUnitList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantSmplUnitList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantSmplUnitList: response missing resultCode"},
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

			_, err = client.PlantSmplUnitList(context.Background(), application.PlantSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqPlantSpecsID: "test-plant-species-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantSmplUnitListLive(t *testing.T) {
	serviceKey := os.Getenv("DATA_GO_KR_SERVICE_KEY")
	if serviceKey == "" {
		t.Skip("DATA_GO_KR_SERVICE_KEY is not set")
	}

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name            string
		pageNo          int
		numOfRows       int
		reqPlantSpecsID string
	}{
		{name: "first plant species", pageNo: 1, numOfRows: 1, reqPlantSpecsID: "P000004958"},
		{name: "second plant species and changed pagination", pageNo: 2, numOfRows: 2, reqPlantSpecsID: "P000004954"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := client.PlantSmplUnitList(ctx, application.PlantSmplUnitListQuery{
				PageNo:          test.pageNo,
				NumOfRows:       test.numOfRows,
				ReqPlantSpecsID: test.reqPlantSpecsID,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) == 0 {
				t.Fatal("plantSmplUnitList returned no items")
			}
		})
	}

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := client.PlantSmplUnitList(ctx, application.PlantSmplUnitListQuery{
			PageNo:          1,
			NumOfRows:       1,
			ReqPlantSpecsID: "P999999999",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
