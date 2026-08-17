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

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
)

func TestInsectPilbkInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != insectPilbkInfoPath {
			t.Errorf("path = %q, want %q", request.URL.Path, insectPilbkInfoPath)
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":      "test+/=",
			"reqInsctPilbkNo": "test-insect-pictorial-book-number",
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
    <ecoDsrct>ecology description</ecoDsrct>
    <eggDsrct>egg description</eggDsrct>
    <emrgcCnt>emergence count</emrgcCnt>
    <emrgcEraDscrt>emergence era description</emrgcEraDscrt>
    <familyKorNm>family Korean name</familyKorNm>
    <familyNm>family name</familyNm>
    <femaleDsrct>female description</femaleDsrct>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <gnrlDsrct>general description</gnrlDsrct>
    <habitDsrct>habit description</habitDsrct>
    <insctEngNm>insect English name</insctEngNm>
    <insctGnrlNm>insect general name</insctGnrlNm>
    <insctPilbkNo>test-insect-pictorial-book-number</insctPilbkNo>
    <insctSpecsScnm>insect species scientific name</insctSpecsScnm>
    <larvaDsrct>larva description</larvaDsrct>
    <lastUpdtDtm> </lastUpdtDtm>
    <maleDsrct>male description</maleDsrct>
    <mnmmOccrrCnt>minimum occurrence count</mnmmOccrrCnt>
    <mxmmOccrrCnt>maximum occurrence count</mxmmOccrrCnt>
    <ordKorNm>order Korean name</ordKorNm>
    <ordNm>order name</ordNm>
    <pestDsrct>pest control description</pestDsrct>
    <pupaDsrct>pupa description</pupaDsrct>
    <referDsrct>reference description</referDsrct>
    <subFamilyKorNm>subfamily Korean name</subFamilyKorNm>
    <subFamilyNm>subfamily name</subFamilyNm>
    <superFamilyKorNm>superfamily Korean name</superFamilyKorNm>
    <superFamilyNm>superfamily name</superFamilyNm>
    <winterDsrct>winter description</winterDsrct>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.InsectPilbkInfo(context.Background(), application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.InsectPilbkInfoResult{Item: &application.InsectPilbkInfoItem{
		EcoDsrct:         "ecology description",
		EggDsrct:         "egg description",
		EmrgcCnt:         "emergence count",
		EmrgcEraDscrt:    "emergence era description",
		FamilyKorNm:      "family Korean name",
		FamilyNm:         "family name",
		FemaleDsrct:      "female description",
		GenusKorNm:       "genus Korean name",
		GenusNm:          "genus name",
		GnrlDsrct:        "general description",
		HabitDsrct:       "habit description",
		InsctEngNm:       "insect English name",
		InsctGnrlNm:      "insect general name",
		InsctPilbkNo:     "test-insect-pictorial-book-number",
		InsctSpecsScnm:   "insect species scientific name",
		LarvaDsrct:       "larva description",
		LastUpdtDtm:      " ",
		MaleDsrct:        "male description",
		MnmmOccrrCnt:     "minimum occurrence count",
		MxmmOccrrCnt:     "maximum occurrence count",
		OrdKorNm:         "order Korean name",
		OrdNm:            "order name",
		PestDsrct:        "pest control description",
		PupaDsrct:        "pupa description",
		ReferDsrct:       "reference description",
		SubFamilyKorNm:   "subfamily Korean name",
		SubFamilyNm:      "subfamily name",
		SuperFamilyKorNm: "superfamily Korean name",
		SuperFamilyNm:    "superfamily name",
		WinterDsrct:      "winter description",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestInsectPilbkInfoReturnsEmptyItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.InsectPilbkInfo(context.Background(), application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "not-found"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestInsectPilbkInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.InsectPilbkInfo(context.Background(), application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"})
			var apiError *InsectPilbkInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *InsectPilbkInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestInsectPilbkInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.InsectPilbkInfo(context.Background(), application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"})
	var apiError *InsectPilbkInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *InsectPilbkInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestInsectPilbkInfoReturnsResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "insectPilbkInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "insectPilbkInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "insectPilbkInfo: response missing resultCode"},
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

			_, err = client.InsectPilbkInfo(context.Background(), application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "test-insect-pictorial-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestInsectPilbkInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, insctPilbkNo := range []string{"ZREP0001", "ZRED0001"} {
		t.Run(insctPilbkNo, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.InsectPilbkInfo(ctx, application.InsectPilbkInfoQuery{ReqInsctPilbkNo: insctPilbkNo})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.InsctPilbkNo != insctPilbkNo || result.Item.InsctGnrlNm == "" || result.Item.InsctSpecsScnm == "" {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	t.Run("empty not-found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.InsectPilbkInfo(ctx, application.InsectPilbkInfoQuery{ReqInsctPilbkNo: "kna-mcp-no-result-20260817"})
		if err != nil {
			t.Fatal(err)
		}
		if result.Item != nil {
			t.Errorf("item = %#v, want nil", result.Item)
		}
	})
}
