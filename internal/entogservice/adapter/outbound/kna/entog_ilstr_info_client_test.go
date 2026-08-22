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

func TestEntogIlstrInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/1400119/EntogService/entogIlstrInfo" {
			t.Errorf("path = %q, want %q", request.URL.Path, "/1400119/EntogService/entogIlstrInfo")
		}
		wantQuery := map[string]string{
			"serviceKey": "test+/=",
			"q1":         "test-entognath-pictorial-book-number",
		}
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
<btnc>test-btnc</btnc><cont1>test-cont1</cont1><cont2>test-cont2</cont2><cont3>test-cont3</cont3>
<cont4>test-cont4</cont4><cont5>test-cont5</cont5><cont6>test-cont6</cont6><cont7>test-cont7</cont7>
<cont8>test-cont8</cont8><cont9>test-cont9</cont9><cont10>test-cont10</cont10><cont11>test-cont11</cont11>
<cprtCtnt>test-copyright</cprtCtnt><emrgcCnt>test-emergence-count</emrgcCnt><emrgcEraDscrt>test-emergence-era-description</emrgcEraDscrt>
<entogAthrNm>test-author-name</entogAthrNm><entogEngNm>test-English-name</entogEngNm><entogOfnmKrlngNm>test-Korean-name</entogOfnmKrlngNm>
<entogPilbkNo>test-entognath-pictorial-book-number</entogPilbkNo><entogSpecsNm>test-species-name</entogSpecsNm>
<familyKorNm>test-family-Korean-name</familyKorNm><familyNm>test-family-name</familyNm><genusKorNm>test-genus-Korean-name</genusKorNm>
<genusNm>test-genus-name</genusNm><imgUrl>test-image-URL</imgUrl><mnmmOccrrCnt>test-minimum-occurrence-count</mnmmOccrrCnt>
<mxmmOccrrCnt>test-maximum-occurrence-count</mxmmOccrrCnt><ordKorNm>test-order-Korean-name</ordKorNm><ordNm>test-order-name</ordNm>
</item></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test%2B%2F%3D")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	got, err := client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}
	want := application.EntogIlstrInfoResult{Item: &application.EntogIlstrInfoItem{
		Btnc: "test-btnc", Cont1: "test-cont1", Cont2: "test-cont2", Cont3: "test-cont3", Cont4: "test-cont4",
		Cont5: "test-cont5", Cont6: "test-cont6", Cont7: "test-cont7", Cont8: "test-cont8", Cont9: "test-cont9",
		Cont10: "test-cont10", Cont11: "test-cont11", CprtCtnt: "test-copyright", EmrgcCnt: "test-emergence-count",
		EmrgcEraDscrt: "test-emergence-era-description", EntogAthrNm: "test-author-name", EntogEngNm: "test-English-name",
		EntogOfnmKrlngNm: "test-Korean-name", EntogPilbkNo: "test-entognath-pictorial-book-number", EntogSpecsNm: "test-species-name",
		FamilyKorNm: "test-family-Korean-name", FamilyNm: "test-family-name", GenusKorNm: "test-genus-Korean-name",
		GenusNm: "test-genus-name", ImgURL: "test-image-URL", MnmmOccrrCnt: "test-minimum-occurrence-count",
		MxmmOccrrCnt: "test-maximum-occurrence-count", OrdKorNm: "test-order-Korean-name", OrdNm: "test-order-name",
	}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("result = %#v, want %#v", got, want)
	}
}

func TestEntogIlstrInfoPreservesObservedValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode></header><body><item><cont1>description&#xD;
</cont1><cont2> </cont2><emrgcCnt>      0</emrgcCnt><imgUrl>NONE</imgUrl></item></body></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item == nil || result.Item.Cont1 != "description\r\n" || result.Item.Cont2 != " " || result.Item.EmrgcCnt != "      0" || result.Item.ImgURL != "NONE" {
		t.Errorf("item = %#v, want preserved whitespace and NONE", result.Item)
	}
}

func TestEntogIlstrInfoReturnsNilItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(response, `<response><header><resultCode>00</resultCode><resultMsg>NORMAL SERVICE.</resultMsg></header><body/></response>`)
	}))
	defer server.Close()

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	client.baseURL = server.URL

	result, err := client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-missing-entognath-pictorial-book-number"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Item != nil {
		t.Errorf("item = %#v, want nil", result.Item)
	}
}

func TestEntogIlstrInfoReturnsDocumentedAPIErrors(t *testing.T) {
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

			_, err = client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"})
			var apiError *EntogIlstrInfoError
			if !errors.As(err, &apiError) {
				t.Fatalf("error = %v, want *EntogIlstrInfoError", err)
			}
			if apiError.HTTPStatus != http.StatusOK || apiError.Code != test.code || apiError.Message != test.message {
				t.Errorf("error = %#v", apiError)
			}
		})
	}
}

func TestEntogIlstrInfoReturnsGatewayError(t *testing.T) {
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

	_, err = client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"})
	var apiError *EntogIlstrInfoError
	if !errors.As(err, &apiError) {
		t.Fatalf("error = %v, want *EntogIlstrInfoError", err)
	}
	if apiError.HTTPStatus != http.StatusUnauthorized || apiError.Code != "20" || apiError.Message != "SERVICE_KEY_IS_NULL: 서비스 접근거부" {
		t.Errorf("error = %#v", apiError)
	}
}

func TestEntogIlstrInfoReturnsResponseErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{name: "unexpected HTTP status", statusCode: http.StatusBadGateway, body: `<response/>`, wantError: "entogIlstrInfo: unexpected HTTP status 502 Bad Gateway"},
		{name: "invalid XML", statusCode: http.StatusOK, body: `<response>`, wantError: "entogIlstrInfo: decode response"},
		{name: "missing result code", statusCode: http.StatusOK, body: `<response/>`, wantError: "entogIlstrInfo: response missing resultCode"},
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

			_, err = client.EntogIlstrInfo(context.Background(), application.EntogIlstrInfoQuery{Q1: "test-entognath-pictorial-book-number"})
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestEntogIlstrInfoLive(t *testing.T) {
	serviceKey := requireLiveServiceKey(t)
	client, err := NewClient(serviceKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name     string
		q1       string
		wantBtnc string
	}{
		{name: "search result identifier one", q1: "ZABF0025", wantBtnc: "Coecobrya dubiosa (Yosii, 1956)"},
		{name: "search result identifier two", q1: "ZABK0053", wantBtnc: "Metanura cassagnaui Deharveng & Weiner, 1984"},
		{name: "missing result", q1: "KNA-MCP-NOT-FOUND"},
		{name: "lowercase identifier", q1: "zabf0025"},
		{name: "scientific name identifier", q1: "ZABF002501"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			result, err := client.EntogIlstrInfo(ctx, application.EntogIlstrInfoQuery{Q1: test.q1})
			if err != nil {
				t.Fatal(err)
			}
			if test.wantBtnc == "" {
				if result.Item != nil {
					t.Errorf("item = %#v, want nil", result.Item)
				}
				return
			}
			if result.Item == nil {
				t.Fatal("item is nil")
			}
			if result.Item.EntogPilbkNo != test.q1 || result.Item.Btnc != test.wantBtnc {
				t.Errorf("item = %#v, want q1 %q and btnc %q", result.Item, test.q1, test.wantBtnc)
			}
			if test.q1 == "ZABF0025" && (result.Item.Cont2 != " " || result.Item.EmrgcCnt != "      0" || result.Item.ImgURL != "NONE") {
				t.Errorf("item = %#v, want preserved observed values", result.Item)
			}
		})
	}
}
