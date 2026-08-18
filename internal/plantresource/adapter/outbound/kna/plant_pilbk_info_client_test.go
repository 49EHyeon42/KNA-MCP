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

func TestPlantPilbkInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/PlantResource/plantPilbkInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/PlantResource/plantPilbkInfo")
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":      "test+/=",
			"reqPlantPilbkNo": "test-book-number",
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
  <body><item>
    <apgFamilyKorNm>apg family Korean name</apgFamilyKorNm><apgFamilyNm>apg family name</apgFamilyNm>
    <bfofMthod>pest control method</bfofMthod><brdMthdDesc>breeding method description</brdMthdDesc>
    <bugInfo>bug information</bugInfo><dstrb>distribution</dstrb><engNm>English name</engNm>
    <familyKorNm>family Korean name</familyKorNm><familyNm>family name</familyNm>
    <farmSpftDesc>farm feature description</farmSpftDesc><genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm><grwEvrntDesc>growth environment description</grwEvrntDesc>
    <inductionDesc>induction description</inductionDesc><lastUpdtDtm>last update date time</lastUpdtDtm>
    <notRcmmGnrlNm>not recommended general name</notRcmmGnrlNm><note>note</note>
    <orplcNm>origin place name</orplcNm><osDstrb>overseas distribution</osDstrb>
    <plantGnrlNm>plant general name</plantGnrlNm><plantPilbkNo>plant pictorial book number</plantPilbkNo>
    <plantSpecsScnm>plant species scientific name</plantSpecsScnm><prtcPlnDesc>protection plan description</prtcPlnDesc>
    <rrngGubun>growth classification</rrngGubun><rrngType>growth type</rrngType><shpe>shape</shpe>
    <smlrPlntDesc>similar plant description</smlrPlntDesc><spft>feature</spft>
    <useMthdDesc>use method description</useMthdDesc><woodDesc>wood description</woodDesc>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.PlantPilbkInfo(context.Background(), application.PlantPilbkInfoQuery{
		ReqPlantPilbkNo: "test-book-number",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.PlantPilbkInfoResult{
		APGFamilyKorNm: "apg family Korean name",
		APGFamilyNm:    "apg family name",
		BfofMthod:      "pest control method",
		BrdMthdDesc:    "breeding method description",
		BugInfo:        "bug information",
		Dstrb:          "distribution",
		EngNm:          "English name",
		FamilyKorNm:    "family Korean name",
		FamilyNm:       "family name",
		FarmSpftDesc:   "farm feature description",
		GenusKorNm:     "genus Korean name",
		GenusNm:        "genus name",
		GrwEvrntDesc:   "growth environment description",
		InductionDesc:  "induction description",
		LastUpdtDtm:    "last update date time",
		NotRcmmGnrlNm:  "not recommended general name",
		Note:           "note",
		OrplcNm:        "origin place name",
		OsDstrb:        "overseas distribution",
		PlantGnrlNm:    "plant general name",
		PlantPilbkNo:   "plant pictorial book number",
		PlantSpecsScnm: "plant species scientific name",
		PrtcPlnDesc:    "protection plan description",
		RrngGubun:      "growth classification",
		RrngType:       "growth type",
		Shpe:           "shape",
		SmlrPlntDesc:   "similar plant description",
		Spft:           "feature",
		UseMthdDesc:    "use method description",
		WoodDesc:       "wood description",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestPlantPilbkInfoReturnsEmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.PlantPilbkInfo(context.Background(), application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "unknown"})
	if err != nil {
		t.Fatal(err)
	}
	if result != (application.PlantPilbkInfoResult{}) {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestPlantPilbkInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.PlantPilbkInfo(context.Background(), application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "test-book-number"})
			var apiError *PlantPilbkInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *PlantPilbkInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestPlantPilbkInfoReturnsGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(response, `<OpenAPI_ServiceResponse><cmmMsgHeader><errMsg>SERVICE_KEY_IS_NOT_REGISTERED_ERROR</errMsg><returnAuthMsg>unregistered service key</returnAuthMsg><returnReasonCode>30</returnReasonCode></cmmMsgHeader></OpenAPI_ServiceResponse>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.PlantPilbkInfo(context.Background(), application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "test-book-number"})
	var apiError *PlantPilbkInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *PlantPilbkInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: unregistered service key" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestPlantPilbkInfoReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "plantPilbkInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "plantPilbkInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "plantPilbkInfo: response missing resultCode"},
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

			_, err = client.PlantPilbkInfo(context.Background(), application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "test-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestPlantPilbkInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, number := range []string{"31662", "31665"} {
		t.Run(number, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.PlantPilbkInfo(ctx, application.PlantPilbkInfoQuery{ReqPlantPilbkNo: number})
			if err != nil {
				t.Fatal(err)
			}
			if result.PlantPilbkNo != number || result.PlantGnrlNm == "" {
				t.Errorf("result = %#v", result)
			}
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := client.PlantPilbkInfo(ctx, application.PlantPilbkInfoQuery{ReqPlantPilbkNo: "999999999"})
	if err != nil {
		t.Fatal(err)
	}
	if result != (application.PlantPilbkInfoResult{}) {
		t.Errorf("unknown result = %#v, want empty result", result)
	}
}
