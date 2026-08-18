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
	alchnSpcmInfoPath        = lchnServiceBasePath + "/alchnSpcmInfo"
	alchnSpcmInfoSuccessCode = "00"
)

var alchnSpcmInfoResultMessages = map[string]string{
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

var _ outbound.AlchnSpcmInfoPort = (*Client)(nil)

type alchnSpcmInfoBody struct {
	Item *alchnSpcmInfoItem `xml:"item"`
}

type alchnSpcmInfoItem struct {
	Btnc          string `xml:"btnc"`
	ClarDtlDscrt  string `xml:"clarDtlDscrt"`
	CllcrNm       string `xml:"cllcrNm"`
	CltrNm        string `xml:"cltrNm"`
	CprtCtnt      string `xml:"cprtCtnt"`
	EngNm         string `xml:"engNm"`
	ExmneNm       string `xml:"exmneNm"`
	FamilyKorNm   string `xml:"familyKorNm"`
	FamilyNm      string `xml:"familyNm"`
	FrstRgstnDtm  string `xml:"frstRgstnDtm"`
	GenusKorNm    string `xml:"genusKorNm"`
	GenusNm       string `xml:"genusNm"`
	Grdnt         string `xml:"grdnt"`
	HaslvVal      string `xml:"haslvVal"`
	HbttChrcrCont string `xml:"hbttChrcrCont"`
	ImgURL        string `xml:"imgUrl"`
	InsttSmplNo   string `xml:"insttSmplNo"`
	JapNm         string `xml:"japNm"`
	LastUpdtDtm   string `xml:"lastUpdtDtm"`
	LchnGnrlNm    string `xml:"lchnGnrlNm"`
	LchnScnmID    string `xml:"lchnScnmId"`
	LchnSmplNo    string `xml:"lchnSmplNo"`
	OrbrnCd       string `xml:"orbrnCd"`
	PrkNm         string `xml:"prkNm"`
	SmplCllcnDt   string `xml:"smplCllcnDt"`
}

type alchnSpcmInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body alchnSpcmInfoBody `xml:"body"`
}

// AlchnSpcmInfoError reports an error returned by alchnSpcmInfo.
type AlchnSpcmInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the alchnSpcmInfo error message.
func (e *AlchnSpcmInfoError) Error() string {
	return fmt.Sprintf("alchnSpcmInfo: API error %s: %s", e.Code, e.Message)
}

// AlchnSpcmInfo gets Korea National Arboretum lichen specimen detail information.
func (c *Client) AlchnSpcmInfo(ctx context.Context, query application.AlchnSpcmInfoQuery) (application.AlchnSpcmInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, alchnSpcmInfoPath)
	if err != nil {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("q1", query.Q1)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload alchnSpcmInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.AlchnSpcmInfoResult{}, &AlchnSpcmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.AlchnSpcmInfoResult{}, fmt.Errorf("alchnSpcmInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.AlchnSpcmInfoResult{}, errors.New("alchnSpcmInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != alchnSpcmInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = alchnSpcmInfoResultMessages[payload.Header.ResultCode]
		}
		return application.AlchnSpcmInfoResult{}, &AlchnSpcmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body alchnSpcmInfoBody) result() application.AlchnSpcmInfoResult {
	if body.Item == nil {
		return application.AlchnSpcmInfoResult{}
	}

	return application.AlchnSpcmInfoResult{Item: &application.AlchnSpcmInfoItem{
		Btnc:          body.Item.Btnc,
		ClarDtlDscrt:  body.Item.ClarDtlDscrt,
		CllcrNm:       body.Item.CllcrNm,
		CltrNm:        body.Item.CltrNm,
		CprtCtnt:      body.Item.CprtCtnt,
		EngNm:         body.Item.EngNm,
		ExmneNm:       body.Item.ExmneNm,
		FamilyKorNm:   body.Item.FamilyKorNm,
		FamilyNm:      body.Item.FamilyNm,
		FrstRgstnDtm:  body.Item.FrstRgstnDtm,
		GenusKorNm:    body.Item.GenusKorNm,
		GenusNm:       body.Item.GenusNm,
		Grdnt:         body.Item.Grdnt,
		HaslvVal:      body.Item.HaslvVal,
		HbttChrcrCont: body.Item.HbttChrcrCont,
		ImgURL:        body.Item.ImgURL,
		InsttSmplNo:   body.Item.InsttSmplNo,
		JapNm:         body.Item.JapNm,
		LastUpdtDtm:   body.Item.LastUpdtDtm,
		LchnGnrlNm:    body.Item.LchnGnrlNm,
		LchnScnmID:    body.Item.LchnScnmID,
		LchnSmplNo:    body.Item.LchnSmplNo,
		OrbrnCd:       body.Item.OrbrnCd,
		PrkNm:         body.Item.PrkNm,
		SmplCllcnDt:   body.Item.SmplCllcnDt,
	}}
}
