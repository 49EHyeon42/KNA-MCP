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

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
)

func TestChildPilbkInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/ChildService2/childPilbkInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/ChildService2/childPilbkInfo")
		}
		if got := request.Header.Get("Accept"); got != "application/xml" {
			t.Errorf("Accept = %q, want application/xml", got)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":           "test+/=",
			"reqChildLvbngPilbkNo": "test-child-pictorial-book-number",
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
    <biogyNm>biology name</biogyNm>
    <childLvbngPilbkNo>test-child-pictorial-book-number</childLvbngPilbkNo>
    <extrmCrss>extinction crisis</extrmCrss>
    <familyKorNm>family Korean name</familyKorNm>
    <familyNm>family name</familyNm>
    <genusKorNm>genus Korean name</genusKorNm>
    <genusNm>genus name</genusNm>
    <hbttFieldYn>field habitat flag</hbttFieldYn>
    <hbttFrestYn>forest habitat flag</hbttFrestYn>
    <hbttRiverYn>river habitat flag</hbttRiverYn>
    <lvbngDscrt>living thing description &lt;br/&gt; next</lvbngDscrt>
    <lvbngTpcdNm>living thing type code name</lvbngTpcdNm>
    <lvngKrlngNm>living thing Korean name</lvngKrlngNm>
    <prtctSpecsTpcdNm> </prtctSpecsTpcdNm>
  </item></body>
</response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}

	want := application.ChildPilbkInfoResult{Item: &application.ChildPilbkInfoItem{
		BiogyNm:           "biology name",
		ChildLvbngPilbkNo: "test-child-pictorial-book-number",
		ExtrmCrss:         "extinction crisis",
		FamilyKorNm:       "family Korean name",
		FamilyNm:          "family name",
		GenusKorNm:        "genus Korean name",
		GenusNm:           "genus name",
		HbttFieldYn:       "field habitat flag",
		HbttFrestYn:       "forest habitat flag",
		HbttRiverYn:       "river habitat flag",
		LvbngDscrt:        "living thing description <br/> next",
		LvbngTpcdNm:       "living thing type code name",
		LvngKrlngNm:       "living thing Korean name",
		PrtctSpecsTpcdNm:  " ",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestChildPilbkInfoReturnsEmptyItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "999999999999"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestChildPilbkInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"})
			var apiError *ChildPilbkInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *ChildPilbkInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestChildPilbkInfoReturnsObservedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>99</resultCode><resultMsg>ORA-01722: invalid number
</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	_, err = client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "invalid-id"})
	var apiError *ChildPilbkInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ChildPilbkInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusOK || apiError.Code != "99" || apiError.Message != "ORA-01722: invalid number\n" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestChildPilbkInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"})
	var apiError *ChildPilbkInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *ChildPilbkInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestChildPilbkInfoReturnsResponseErrors(t *testing.T) {
	for _, test := range []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "childPilbkInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "childPilbkInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "childPilbkInfo: response missing resultCode"},
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

			_, err = client.ChildPilbkInfo(context.Background(), application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "test-child-pictorial-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestChildPilbkInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	searchResult, err := client.ChildPilbkSearch(ctx, application.ChildPilbkSearchQuery{PageNo: 1, NumOfRows: 2})
	cancel()
	if err != nil {
		t.Fatal(err)
	}
	if len(searchResult.Items) != 2 {
		t.Fatalf("childPilbkSearch items = %d, want 2", len(searchResult.Items))
	}

	for index, searchItem := range searchResult.Items {
		t.Run(fmt.Sprintf("search result %d", index+1), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.ChildPilbkInfo(ctx, application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: searchItem.ChildLvbngPilbkNo})
			if err != nil {
				t.Fatal(err)
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.ChildLvbngPilbkNo != searchItem.ChildLvbngPilbkNo || result.Item.LvngKrlngNm == "" || result.Item.BiogyNm == "" {
				t.Errorf("item = %#v", result.Item)
			}
		})
	}

	t.Run("empty not-found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.ChildPilbkInfo(ctx, application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "999999999999"})
		if err != nil {
			t.Fatal(err)
		}
		if result.Item != nil {
			t.Errorf("item = %#v, want nil", result.Item)
		}
	})

	t.Run("invalid number", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		_, err := client.ChildPilbkInfo(ctx, application.ChildPilbkInfoQuery{ReqChildLvbngPilbkNo: "invalid-id"})
		var apiError *ChildPilbkInfoError
		if !errors.As(err, &apiError) || apiError.Code != "99" || !strings.Contains(apiError.Message, "ORA-01722") {
			t.Errorf("error = %#v, want observed resultCode 99", err)
		}
	})
}
