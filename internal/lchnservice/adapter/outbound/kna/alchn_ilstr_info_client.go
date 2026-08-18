package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lchnservice/application/port/outbound"
)

const (
	alchnIlstrInfoPath        = lchnServiceBasePath + "/alchnIlstrInfo"
	alchnIlstrInfoSuccessCode = "00"
)

var alchnIlstrInfoResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"01": "APPLICATION_ERROR",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"04": "HTTP_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"12": "NO_OPENAPI_SERVICE_ERROR",
	"20": "SERVICE_ACCESS_DENIED_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"22": "LIMITED_NUMBER_OF_SERVICE_REQUESTS_EXCEEDS_ERROR",
	"30": "SERVICE_KEY_IS_NOT_REGISTERED_ERROR",
	"31": "DEADLINE_HAS_EXPIRED_ERROR",
	"32": "UNREGISTERED_IP_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
	"99": "UNKNOWN_ERROR",
}

var _ outbound.AlchnIlstrInfoPort = (*Client)(nil)

type alchnIlstrInfoBody struct {
	Item *alchnIlstrInfoItem `xml:"item"`
}

type alchnIlstrInfoItem struct {
	Btnc         string `xml:"btnc"`
	Cont1        string `xml:"cont1"`
	Cont2        string `xml:"cont2"`
	Cont3        string `xml:"cont3"`
	Cont4        string `xml:"cont4"`
	Cont5        string `xml:"cont5"`
	Cont6        string `xml:"cont6"`
	Cont7        string `xml:"cont7"`
	Cont8        string `xml:"cont8"`
	Cont9        string `xml:"cont9"`
	Cont10       string `xml:"cont10"`
	Cont11       string `xml:"cont11"`
	Cont12       string `xml:"cont12"`
	CprtCtnt     string `xml:"cprtCtnt"`
	EngNm        string `xml:"engNm"`
	FamilyKorNm  string `xml:"familyKorNm"`
	FamilyNm     string `xml:"familyNm"`
	FrstRgstnDtm string `xml:"frstRgstnDtm"`
	GenusKorNm   string `xml:"genusKorNm"`
	GenusNm      string `xml:"genusNm"`
	ImgURL       string `xml:"imgUrl"`
	JapNm        string `xml:"japNm"`
	LastUpdtDtm  string `xml:"lastUpdtDtm"`
	LchnGnrlNm   string `xml:"lchnGnrlNm"`
	LchnInfrpNm  string `xml:"lchnInfrpNm"`
	LchnPilbkNo  string `xml:"lchnPilbkNo"`
	LchnScnmID   string `xml:"lchnScnmId"`
	LchnTtnm     string `xml:"lchnTtnm"`
	PrkNm        string `xml:"prkNm"`
}

type alchnIlstrInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body alchnIlstrInfoBody `xml:"body"`
}

// AlchnIlstrInfoError reports an error returned by alchnIlstrInfo.
type AlchnIlstrInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the alchnIlstrInfo error message.
func (e *AlchnIlstrInfoError) Error() string {
	return fmt.Sprintf("alchnIlstrInfo: API error %s: %s", e.Code, e.Message)
}

// AlchnIlstrInfo gets Korea National Arboretum lichen pictorial book detail information.
func (c *Client) AlchnIlstrInfo(ctx context.Context, query application.AlchnIlstrInfoQuery) (application.AlchnIlstrInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, alchnIlstrInfoPath)
	if err != nil {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("q1", query.Q1)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload alchnIlstrInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.AlchnIlstrInfoResult{}, &AlchnIlstrInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.AlchnIlstrInfoResult{}, fmt.Errorf("alchnIlstrInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.AlchnIlstrInfoResult{}, errors.New("alchnIlstrInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != alchnIlstrInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = alchnIlstrInfoResultMessages[payload.Header.ResultCode]
		}
		return application.AlchnIlstrInfoResult{}, &AlchnIlstrInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body alchnIlstrInfoBody) result() application.AlchnIlstrInfoResult {
	if body.Item == nil {
		return application.AlchnIlstrInfoResult{}
	}

	return application.AlchnIlstrInfoResult{Item: &application.AlchnIlstrInfoItem{
		Btnc:         body.Item.Btnc,
		Cont1:        body.Item.Cont1,
		Cont2:        body.Item.Cont2,
		Cont3:        body.Item.Cont3,
		Cont4:        body.Item.Cont4,
		Cont5:        body.Item.Cont5,
		Cont6:        body.Item.Cont6,
		Cont7:        body.Item.Cont7,
		Cont8:        body.Item.Cont8,
		Cont9:        body.Item.Cont9,
		Cont10:       body.Item.Cont10,
		Cont11:       body.Item.Cont11,
		Cont12:       body.Item.Cont12,
		CprtCtnt:     body.Item.CprtCtnt,
		EngNm:        body.Item.EngNm,
		FamilyKorNm:  body.Item.FamilyKorNm,
		FamilyNm:     body.Item.FamilyNm,
		FrstRgstnDtm: body.Item.FrstRgstnDtm,
		GenusKorNm:   body.Item.GenusKorNm,
		GenusNm:      body.Item.GenusNm,
		ImgURL:       body.Item.ImgURL,
		JapNm:        body.Item.JapNm,
		LastUpdtDtm:  body.Item.LastUpdtDtm,
		LchnGnrlNm:   body.Item.LchnGnrlNm,
		LchnInfrpNm:  body.Item.LchnInfrpNm,
		LchnPilbkNo:  body.Item.LchnPilbkNo,
		LchnScnmID:   body.Item.LchnScnmID,
		LchnTtnm:     body.Item.LchnTtnm,
		PrkNm:        body.Item.PrkNm,
	}}
}
