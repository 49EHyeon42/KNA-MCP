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

func TestAlchnSpcmInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/LchnService/alchnSpcmInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/LchnService/alchnSpcmInfo")
		}
		wantQuery := map[string]string{"serviceKey": "test+/=", "q1": "TEST-SAMPLE-001"}
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
<btnc>Testus lichenii Author</btnc><clarDtlDscrt>test collection site</clarDtlDscrt>
<cllcrNm> </cllcrNm><cltrNm>test collector group</cltrNm><cprtCtnt>test copyright</cprtCtnt><engNm> </engNm><exmneNm>test examiner</exmneNm>
<familyKorNm>테스트과</familyKorNm><familyNm>Testaceae</familyNm><frstRgstnDtm>2020-01-02 03:04:05</frstRgstnDtm>
<genusKorNm>테스트속</genusKorNm><genusNm>Testus</genusNm><grdnt> </grdnt><haslvVal>123</haslvVal><hbttChrcrCont>test substrate</hbttChrcrCont>
<imgUrl>https://example.com/specimens/&#x9;image.jpg</imgUrl><insttSmplNo>TEST-INST-001</insttSmplNo><japNm> </japNm>
<lastUpdtDtm>2020-02-03 04:05:06</lastUpdtDtm><lchnGnrlNm>테스트지의</lchnGnrlNm><lchnScnmId>TEST-SCNM-ID</lchnScnmId>
<lchnSmplNo>TEST-SAMPLE-001</lchnSmplNo><orbrnCd> </orbrnCd><prkNm> </prkNm><smplCllcnDt>20200102</smplCllcnDt>
</item></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: "TEST-SAMPLE-001"})
	if err != nil {
		t.Fatal(err)
	}
	want := application.AlchnSpcmInfoResult{Item: &application.AlchnSpcmInfoItem{
		Btnc:          "Testus lichenii Author",
		ClarDtlDscrt:  "test collection site",
		CllcrNm:       " ",
		CltrNm:        "test collector group",
		CprtCtnt:      "test copyright",
		EngNm:         " ",
		ExmneNm:       "test examiner",
		FamilyKorNm:   "테스트과",
		FamilyNm:      "Testaceae",
		FrstRgstnDtm:  "2020-01-02 03:04:05",
		GenusKorNm:    "테스트속",
		GenusNm:       "Testus",
		Grdnt:         " ",
		HaslvVal:      "123",
		HbttChrcrCont: "test substrate",
		ImgURL:        "https://example.com/specimens/\timage.jpg",
		InsttSmplNo:   "TEST-INST-001",
		JapNm:         " ",
		LastUpdtDtm:   "2020-02-03 04:05:06",
		LchnGnrlNm:    "테스트지의",
		LchnScnmID:    "TEST-SCNM-ID",
		LchnSmplNo:    "TEST-SAMPLE-001",
		OrbrnCd:       " ",
		PrkNm:         " ",
		SmplCllcnDt:   "20200102",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestAlchnSpcmInfoReturnsNilItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: "KNA-MCP-NO-RESULT"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestAlchnSpcmInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: "TEST-SAMPLE-001"})
			var apiError *AlchnSpcmInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *AlchnSpcmInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestAlchnSpcmInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: "TEST-SAMPLE-001"})
	var apiError *AlchnSpcmInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *AlchnSpcmInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusForbidden || apiError.Code != "30" || apiError.Message != "SERVICE_KEY_IS_NOT_REGISTERED_ERROR: 등록되지 않은 서비스키" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestAlchnSpcmInfoReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "alchnSpcmInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "alchnSpcmInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "alchnSpcmInfo: response missing resultCode"},
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

			_, err = client.AlchnSpcmInfo(context.Background(), application.AlchnSpcmInfoQuery{Q1: "TEST-SAMPLE-001"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestAlchnSpcmInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name       string
		q1         string
		wantItem   bool
		wantSmplNo string
	}{
		{name: "search result identifier", q1: "KNKL201702169020", wantItem: true, wantSmplNo: "KNKL201702169020"},
		{name: "second search result identifier", q1: "KNKL201702169022", wantItem: true, wantSmplNo: "KNKL201702169022"},
		{name: "missing result", q1: "KNA-MCP-NO-RESULT"},
		{name: "lowercase identifier", q1: "knkl201702169020"},
		{name: "scientific name identifier", q1: "LS2017000190"},
		{name: "institution specimen number", q1: "KHL0023100"},
		{name: "missing q1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.AlchnSpcmInfo(ctx, application.AlchnSpcmInfoQuery{Q1: test.q1})
			if err != nil {
				t.Fatal(err)
			}
			if (result.Item != nil) != test.wantItem {
				t.Errorf("item = %#v, want present %t", result.Item, test.wantItem)
			}
			if result.Item != nil && result.Item.LchnSmplNo != test.wantSmplNo {
				t.Errorf("lchnSmplNo = %q, want %q", result.Item.LchnSmplNo, test.wantSmplNo)
			}
		})
	}

	t.Run("preserves actual values", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := client.AlchnSpcmInfo(ctx, application.AlchnSpcmInfoQuery{Q1: "KNKL201702169020"})
		if err != nil {
			t.Fatal(err)
		}
		if result.Item == nil {
			t.Fatal("item is nil")
		}
		if result.Item.LchnScnmID != "LS2017000190" || result.Item.InsttSmplNo != "KHL0023100" {
			t.Errorf("item = %#v, want actual identifiers", result.Item)
		}
		if !strings.Contains(result.Item.ImgURL, "\t") {
			t.Errorf("imgUrl = %q, want preserved tab", result.Item.ImgURL)
		}
	})
}
