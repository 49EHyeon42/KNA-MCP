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

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
)

func TestScnmInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != scnmInfoPath {
			t.Errorf("path = %q, want %q", request.URL.Path, scnmInfoPath)
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":     "test+/=",
			"reqPlantScnmId": "1004701",
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
  <body><item>
    <apgFalmKorNm>APG family Korean name</apgFalmKorNm>
    <apgFalmNm>APG family name</apgFalmNm>
    <biogyNmTpcdNm>biology name type code name</biogyNmTpcdNm>
    <cltvaYn>cultivation yes or no</cltvaYn>
    <eclgDstrbYn>ecological disturbance yes or no</eclgDstrbYn>
    <extcCncrnsYn>exotic concern yes or no</extcCncrnsYn>
    <extcPlantCdNm>exotic plant code name</extcPlantCdNm>
    <extcPlantYn>exotic plant yes or no</extcPlantYn>
    <falmKorNm>family Korean name</falmKorNm>
    <falmNm>family name</falmNm>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <ltrtrInfrmNm>literature information name</ltrtrInfrmNm>
    <plantBrdgFomTpcdNm>plant breeding form type code name</plantBrdgFomTpcdNm>
    <plantChnNm>plant Chinese name</plantChnNm>
    <plantEngNm>plant English name</plantEngNm>
    <plantGnrlNm>plant general name</plantGnrlNm>
    <plantGnrlNm2>plant general name 2</plantGnrlNm2>
    <plantJpnNm>plant Japanese name</plantJpnNm>
    <plantScnmId>1004701</plantScnmId>
    <plantSpecsScnm>plant species scientific name</plantSpecsScnm>
    <rareTpcdNm>rare type code name</rareTpcdNm>
    <relPlantSpecsScnm>related plant species scientific name</relPlantSpecsScnm>
    <relScnmTpcdNm>related scientific name type code name</relScnmTpcdNm>
    <rmrk>remark</rmrk>
    <rrnssPlantYn>rareness plant yes or no</rrnssPlantYn>
    <spcltPlantCdNm>specialty plant code name</spcltPlantCdNm>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqPlantScnmID: "1004701"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ScnmInfoResult{Item: &application.ScnmInfoItem{
		APGFalmKorNm:       "APG family Korean name",
		APGFalmNm:          "APG family name",
		BiogyNmTpcdNm:      "biology name type code name",
		CltvaYn:            "cultivation yes or no",
		EclgDstrbYn:        "ecological disturbance yes or no",
		ExtcCncrnsYn:       "exotic concern yes or no",
		ExtcPlantCdNm:      "exotic plant code name",
		ExtcPlantYn:        "exotic plant yes or no",
		FalmKorNm:          "family Korean name",
		FalmNm:             "family name",
		GenusKorNm:         "genus Korean name",
		GenusNm:            "genus name",
		LtrtrInfrmNm:       "literature information name",
		PlantBrdgFomTpcdNm: "plant breeding form type code name",
		PlantChnNm:         "plant Chinese name",
		PlantEngNm:         "plant English name",
		PlantGnrlNm:        "plant general name",
		PlantGnrlNm2:       "plant general name 2",
		PlantJpnNm:         "plant Japanese name",
		PlantScnmID:        "1004701",
		PlantSpecsScnm:     "plant species scientific name",
		RareTpcdNm:         "rare type code name",
		RelPlantSpecsScnm:  "related plant species scientific name",
		RelScnmTpcdNm:      "related scientific name type code name",
		Rmrk:               "remark",
		RrnssPlantYn:       "rareness plant yes or no",
		SpcltPlantCdNm:     "specialty plant code name",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestScnmInfoReturnsEmptyItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqPlantScnmID: "not-found"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestScnmInfoReturnsDocumentedAPIErrors(t *testing.T) {
	for _, test := range []struct {
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
	} {
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

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqPlantScnmID: "1004701"})
			var apiError *ScnmInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ScnmInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestScnmInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqPlantScnmID: "1004701"})
	var apiError *ScnmInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ScnmInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestScnmInfoReturnsResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "scnmInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "scnmInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "scnmInfo: response missing resultCode"},
	} {
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

			_, err = client.ScnmInfo(context.Background(), application.ScnmInfoQuery{ReqPlantScnmID: "1004701"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestScnmInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, plantScnmID := range []string{"1004701", "1002511"} {
		t.Run(plantScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqPlantScnmID: plantScnmID})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.PlantScnmID != plantScnmID || result.Item.PlantGnrlNm == "" || result.Item.PlantSpecsScnm == "" {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	for _, plantScnmID := range []string{"999999999999", "not-an-id"} {
		t.Run("empty "+plantScnmID, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqPlantScnmID: plantScnmID})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item != nil {
				t.Errorf("item = %#v, want nil", result.Item)
			}
		})
	}

	t.Run("invalid service key", func(t *testing.T) {
		client, err := NewClient("invalid-service-key")
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		_, err = client.ScnmInfo(ctx, application.ScnmInfoQuery{ReqPlantScnmID: "1004701"})
		var apiError *ScnmInfoError
		if !errors.As(err, &apiError) {
			t.Fatalf("error = %v, want *ScnmInfoError", err)
		}
		if apiError.Code != "30" {
			t.Errorf("error = %#v, want code 30", apiError)
		}
	})
}
