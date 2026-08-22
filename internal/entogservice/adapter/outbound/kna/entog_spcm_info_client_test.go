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

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
)

func TestEntogSpcmInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/EntogService/entogSpcmInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/EntogService/entogSpcmInfo")
		}
		query := request.URL.Query()
		wantQuery := map[string]string{"serviceKey": "test+/=", "q1": "test-specimen-number"}
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
    <btnc>test-scientific-name</btnc><chnNm>test-Chinese-name</chnNm>
    <clarHaslvVal>test-collection-altitude</clarHaslvVal><clctDyDesc>test-collection-date</clctDyDesc>
    <cprtCtnt>test-copyright</cprtCtnt><engNm>test-English-name</engNm>
    <entogGnrlNm>test-Korean-name</entogGnrlNm><entogPilbkNo>test-pictorial-book-number</entogPilbkNo>
    <entogSmplNo>test-specimen-number</entogSmplNo><familyKorNm>test-family-Korean-name</familyKorNm>
    <familyNm>test-family-name</familyNm><frstRgstnDtm>test-first-registration-date</frstRgstnDtm>
    <genusKorNm>test-genus-Korean-name</genusKorNm><genusNm>test-genus-name</genusNm>
    <imgUrl>test-image-URL</imgUrl><japNm>test-Japanese-name</japNm>
    <labelUsgCllcnNmplc>test-label-collection-place</labelUsgCllcnNmplc><lastUpdtDtm>test-last-update-date</lastUpdtDtm>
    <ordKorNm>test-order-Korean-name</ordKorNm><ordNm>test-order-name</ordNm>
    <prkNm>test-North-Korean-name</prkNm><torsoLngth>test-torso-length</torsoLngth><wingLngth>test-wing-length</wingLngth>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "test-specimen-number"})
	if err != nil {
		t.Fatal(err)
	}
	want := application.EntogSpcmInfoResult{Item: &application.EntogSpcmInfoItem{
		Btnc: "test-scientific-name", ChnNm: "test-Chinese-name", ClarHaslvVal: "test-collection-altitude",
		ClctDyDesc: "test-collection-date", CprtCtnt: "test-copyright", EngNm: "test-English-name",
		EntogGnrlNm: "test-Korean-name", EntogPilbkNo: "test-pictorial-book-number", EntogSmplNo: "test-specimen-number",
		FamilyKorNm: "test-family-Korean-name", FamilyNm: "test-family-name", FrstRgstnDtm: "test-first-registration-date",
		GenusKorNm: "test-genus-Korean-name", GenusNm: "test-genus-name", ImgURL: "test-image-URL",
		JapNm: "test-Japanese-name", LabelUsgCllcnNmplc: "test-label-collection-place", LastUpdtDtm: "test-last-update-date",
		OrdKorNm: "test-order-Korean-name", OrdNm: "test-order-name", PrkNm: "test-North-Korean-name",
		TorsoLngth: "test-torso-length", WingLngth: "test-wing-length",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestEntogSpcmInfoPreservesObservedValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><item><btnc>Test &amp; Test</btnc><chnNm> </chnNm><clarHaslvVal> </clarHaslvVal><entogPilbkNo> </entogPilbkNo><imgUrl>NONE</imgUrl><labelUsgCllcnNmplc>test place   </labelUsgCllcnNmplc><lastUpdtDtm> </lastUpdtDtm><torsoLngth> </torsoLngth><wingLngth> </wingLngth></item></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "test-specimen-number"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item == nil {
		t.Fatal("item is nil")
	}
	if result.Item.Btnc != "Test & Test" || result.Item.ChnNm != " " || result.Item.ClarHaslvVal != " " || result.Item.EntogPilbkNo != " " || result.Item.ImgURL != "NONE" || result.Item.LabelUsgCllcnNmplc != "test place   " || result.Item.LastUpdtDtm != " " || result.Item.TorsoLngth != " " || result.Item.WingLngth != " " {
		t.Errorf("item = %#v, want preserved response values", result.Item)
	}
}

func TestEntogSpcmInfoReturnsNilItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "not-found"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestEntogSpcmInfoReturnsDocumentedAPIErrors(t *testing.T) {
	tests := []struct {
		code    string
		message string
	}{
		{code: "01", message: "APPLICATION_ERROR"},
		{code: "02", message: "DB_ERROR"},
		{code: "03", message: "NODATA_ERROR"},
		{code: "04", message: "HTTP_ERROR"},
		{code: "05", message: "SERVICETIME_OUT"},
		{code: "10", message: "INVALID_REQUEST_PARAMETER_ERROR"},
		{code: "11", message: "NO_MANDATORY_REQUEST_PARAMETERS_ERROR"},
		{code: "12", message: "NO_OPENAPI_SERVICE_ERROR"},
		{code: "20", message: "SERVICE_ACCESS_DENIED_ERROR"},
		{code: "21", message: "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR"},
		{code: "22", message: "LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR"},
		{code: "30", message: "SERVICE_KEY_IS_NOT_REGISTERED_ERROR"},
		{code: "31", message: "DEADLINE_HAS_EXPIRED_ERROR"},
		{code: "32", message: "UNREGISTERED_IP_ERROR"},
		{code: "33", message: "UNSIGNED_CALL_ERROR"},
		{code: "99", message: "UNKNOWN_ERROR"},
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

			_, err = client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "test-specimen-number"})
			var apiError *EntogSpcmInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *EntogSpcmInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestEntogSpcmInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "test-specimen-number"})
	var apiError *EntogSpcmInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *EntogSpcmInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestEntogSpcmInfoReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "entogSpcmInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "entogSpcmInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "entogSpcmInfo: response missing resultCode"},
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

			_, err = client.EntogSpcmInfo(context.Background(), application.EntogSpcmInfoQuery{Q1: "test-specimen-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestEntogSpcmInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		q1   string
	}{
		{name: "first specimen", q1: "CNUE0911181002"},
		{name: "second specimen", q1: "CNUE0911171006"},
		{name: "different scientific name", q1: "CNUE0911261029"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.EntogSpcmInfo(ctx, application.EntogSpcmInfoQuery{Q1: test.q1})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil || result.Item.EntogSmplNo != test.q1 {
				t.Errorf("item = %#v, want entogSmplNo %q", result.Item, test.q1)
			}
		})
	}

	for _, test := range []struct {
		name string
		q1   string
	}{
		{name: "not found", q1: "KNA-MCP-NOT-FOUND-20260822"},
		{name: "lowercase specimen number", q1: "cnue0911181002"},
		{name: "pictorial book number", q1: "ZABF0014"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.EntogSpcmInfo(ctx, application.EntogSpcmInfoQuery{Q1: test.q1})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item != nil {
				t.Errorf("item = %#v, want nil", result.Item)
			}
		})
	}
}
