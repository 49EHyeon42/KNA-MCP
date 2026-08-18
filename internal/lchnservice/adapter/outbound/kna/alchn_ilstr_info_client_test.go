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

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
)

func TestAlchnIlstrInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/LchnService/alchnIlstrInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/LchnService/alchnIlstrInfo")
		}
		wantQuery := map[string]string{"serviceKey": "test+/=", "q1": "test-lichen-pictorial-book-number"}
		query := request.URL.Query()
		if len(query) != len(wantQuery) {
			t.Errorf("query key count = %d, want %d", len(query), len(wantQuery))
		}
		for key, want := range wantQuery {
			if query.Get(key) != want {
				t.Errorf("query %s = %q, want %q", key, query.Get(key), want)
			}
		}

		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header><body><item>
<btnc>lichen scientific name</btnc><cont1> </cont1><cont2>cont2</cont2><cont3>cont3</cont3><cont4>cont4</cont4>
<cont5>cont5</cont5><cont6>cont6</cont6><cont7>cont7</cont7><cont8>cont8</cont8><cont9>cont9</cont9>
<cont10>cont10</cont10><cont11>cont11</cont11><cont12>cont12</cont12><cprtCtnt>copyright</cprtCtnt>
<engNm> </engNm><familyKorNm> </familyKorNm><familyNm>family name</familyNm><frstRgstnDtm>first registration date time</frstRgstnDtm>
<genusKorNm>genus Korean name</genusKorNm><genusNm>genus name</genusNm><imgUrl>http://example.com/lichen.jpg</imgUrl><japNm> </japNm>
<lastUpdtDtm>last update date time</lastUpdtDtm><lchnGnrlNm>lichen general name</lchnGnrlNm><lchnInfrpNm> </lchnInfrpNm>
<lchnPilbkNo>test-lichen-pictorial-book-number</lchnPilbkNo><lchnScnmId>lichen scientific name ID</lchnScnmId><lchnTtnm>lichen species epithet</lchnTtnm><prkNm> </prkNm>
</item></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: "test-lichen-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}
	want := application.AlchnIlstrInfoResult{Item: &application.AlchnIlstrInfoItem{
		Btnc: "lichen scientific name", Cont1: " ", Cont2: "cont2", Cont3: "cont3", Cont4: "cont4",
		Cont5: "cont5", Cont6: "cont6", Cont7: "cont7", Cont8: "cont8", Cont9: "cont9", Cont10: "cont10",
		Cont11: "cont11", Cont12: "cont12", CprtCtnt: "copyright", EngNm: " ", FamilyKorNm: " ", FamilyNm: "family name",
		FrstRgstnDtm: "first registration date time", GenusKorNm: "genus Korean name", GenusNm: "genus name",
		ImgURL: "http://example.com/lichen.jpg", JapNm: " ", LastUpdtDtm: "last update date time",
		LchnGnrlNm: "lichen general name", LchnInfrpNm: " ", LchnPilbkNo: "test-lichen-pictorial-book-number", LchnScnmID: "lichen scientific name ID",
		LchnTtnm: "lichen species epithet", PrkNm: " ",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestAlchnIlstrInfoReturnsNilItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: "test-missing-lichen-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestAlchnIlstrInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: "test-lichen-pictorial-book-number"})
			var apiError *AlchnIlstrInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *AlchnIlstrInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestAlchnIlstrInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: "test-lichen-pictorial-book-number"})
	var apiError *AlchnIlstrInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *AlchnIlstrInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestAlchnIlstrInfoReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "alchnIlstrInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "alchnIlstrInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "alchnIlstrInfo: response missing resultCode"},
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

			_, err = client.AlchnIlstrInfo(context.Background(), application.AlchnIlstrInfoQuery{Q1: "test-lichen-pictorial-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestAlchnIlstrInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name     string
		q1       string
		wantItem bool
	}{
		{name: "search result identifier", q1: "LC10000061", wantItem: true},
		{name: "documented identifier", q1: "LC10000042", wantItem: true},
		{name: "missing result", q1: "LC99999999"},
		{name: "lowercase identifier", q1: "lc10000061"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.AlchnIlstrInfo(ctx, application.AlchnIlstrInfoQuery{Q1: test.q1})
			if err != nil {
				t.Fatal(err)
			}
			if (result.Item != nil) != test.wantItem {
				t.Errorf("item = %#v, want present %t", result.Item, test.wantItem)
			}
			if result.Item != nil && result.Item.LchnPilbkNo != test.q1 {
				t.Errorf("lchnPilbkNo = %q, want %q", result.Item.LchnPilbkNo, test.q1)
			}
		})
	}
}
