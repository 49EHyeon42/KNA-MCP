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

func TestInsectSmplUnitList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != insectSmplUnitListPath {
			t.Errorf("path = %q, want %q", request.URL.Path, insectSmplUnitListPath)
		}

		query := request.URL.Query()
		wantQuery := map[string]string{
			"serviceKey":      "test+/=",
			"pageNo":          "2",
			"numOfRows":       "1",
			"reqInsctSpecsId": "test-insect-species-id",
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
      <bspcsInsttNm>specimen holding institution</bspcsInsttNm>
      <clarHaslvVal>collection site elevation</clarHaslvVal><smplCllcnDt>specimen collection date</smplCllcnDt>
      <gynndTpcd>sex type</gynndTpcd><hbttTpcd>habitat type</hbttTpcd>
      <insctSmplNo>insect specimen number</insctSmplNo><insctSpecsId>insect species ID</insctSpecsId>
      <insctSpecsScnm>insect species scientific name</insctSpecsScnm>
      <labelUsgCllcnNmplc>label collection place name</labelUsgCllcnNmplc><lastUpdtDtm> </lastUpdtDtm>
      <prsrtStcd>preservation status</prsrtStcd><slistTpcd>minute insect type</slistTpcd>
      <smplKindCd>specimen type</smplKindCd><torsoLngth>torso length</torsoLngth>
      <wingLngth>wing length</wingLngth><insctGnrlNm>insect general name</insctGnrlNm>
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

	got, err := client.InsectSmplUnitList(context.Background(), application.InsectSmplUnitListQuery{
		PageNo:          2,
		NumOfRows:       1,
		ReqInsctSpecsID: "test-insect-species-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := application.InsectSmplUnitListResult{
		Items: []application.InsectSmplUnitListItem{{
			BspcsInsttNm:       "specimen holding institution",
			ClarHaslvVal:       "collection site elevation",
			SmplCllcnDt:        "specimen collection date",
			GynndTpcd:          "sex type",
			HbttTpcd:           "habitat type",
			InsctSmplNo:        "insect specimen number",
			InsctSpecsID:       "insect species ID",
			InsctSpecsScnm:     "insect species scientific name",
			LabelUsgCllcnNmplc: "label collection place name",
			LastUpdtDtm:        " ",
			PrsrtStcd:          "preservation status",
			SlistTpcd:          "minute insect type",
			SmplKindCd:         "specimen type",
			TorsoLngth:         "torso length",
			WingLngth:          "wing length",
			InsctGnrlNm:        "insect general name",
		}},
		NumOfRows:  1,
		PageNo:     2,
		TotalCount: 7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestInsectSmplUnitListReturnsEmptyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><items/><numOfRows>1</numOfRows><pageNo>1</pageNo><totalCount>0</totalCount></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.InsectSmplUnitList(context.Background(), application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "unknown-id"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 || result.TotalCount != 0 {
		t.Errorf("result = %#v, want empty result", result)
	}
}

func TestInsectSmplUnitListReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.InsectSmplUnitList(context.Background(), application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "test-id"})
			var apiError *InsectSmplUnitListError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *InsectSmplUnitListError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestInsectSmplUnitListReturnsGatewayError(t *testing.T) {
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

	_, err = client.InsectSmplUnitList(context.Background(), application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "test-id"})
	var apiError *InsectSmplUnitListError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *InsectSmplUnitListError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestInsectSmplUnitListReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "insectSmplUnitList: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "insectSmplUnitList: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "insectSmplUnitList: response missing resultCode"},
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

			_, err = client.InsectSmplUnitList(context.Background(), application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 1, ReqInsctSpecsID: "test-id"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestInsectSmplUnitListLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)

	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name            string
		reqInsctSpecsID string
	}{
		{name: "first insect species ID", reqInsctSpecsID: "I000008533"},
		{name: "second insect species ID", reqInsctSpecsID: "I000019060"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.InsectSmplUnitList(ctx, application.InsectSmplUnitListQuery{
				PageNo:          1,
				NumOfRows:       2,
				ReqInsctSpecsID: test.reqInsctSpecsID,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Items) != 2 || result.TotalCount < len(result.Items) {
				t.Fatalf("result = %#v, want 2 items and a valid totalCount", result)
			}
			for _, item := range result.Items {
				if item.InsctSpecsID != test.reqInsctSpecsID || item.InsctSmplNo == "" {
					t.Errorf("item = %#v, want species ID %q and specimen number", item, test.reqInsctSpecsID)
				}
			}
		})
	}

	t.Run("changed page", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		first, err := client.InsectSmplUnitList(ctx, application.InsectSmplUnitListQuery{PageNo: 1, NumOfRows: 2, ReqInsctSpecsID: "I000008533"})
		if err != nil {
			t.Fatal(err)
		}
		second, err := client.InsectSmplUnitList(ctx, application.InsectSmplUnitListQuery{PageNo: 2, NumOfRows: 2, ReqInsctSpecsID: "I000008533"})
		if err != nil {
			t.Fatal(err)
		}
		if len(first.Items) != 2 || len(second.Items) != 2 || first.TotalCount != second.TotalCount || first.Items[0].InsctSmplNo == second.Items[0].InsctSmplNo {
			t.Errorf("first = %#v, second = %#v, want distinct pages", first, second)
		}
	})

	t.Run("without result", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.InsectSmplUnitList(ctx, application.InsectSmplUnitListQuery{
			PageNo:          1,
			NumOfRows:       1,
			ReqInsctSpecsID: "I999999999",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Items) != 0 || result.TotalCount != 0 {
			t.Errorf("result = %#v, want empty result", result)
		}
	})
}
