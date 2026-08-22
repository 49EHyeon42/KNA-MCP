package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/entogservice/application/port/outbound"
)

const (
	entogSpcmInfoPath        = entogServiceBasePath + "/entogSpcmInfo"
	entogSpcmInfoSuccessCode = "00"
)

var entogSpcmInfoResultMessages = map[string]string{
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

var _ outbound.EntogSpcmInfoPort = (*Client)(nil)

type entogSpcmInfoBody struct {
	Item *entogSpcmInfoItem `xml:"item"`
}

type entogSpcmInfoItem struct {
	Btnc               string `xml:"btnc"`
	ChnNm              string `xml:"chnNm"`
	ClarHaslvVal       string `xml:"clarHaslvVal"`
	ClctDyDesc         string `xml:"clctDyDesc"`
	CprtCtnt           string `xml:"cprtCtnt"`
	EngNm              string `xml:"engNm"`
	EntogGnrlNm        string `xml:"entogGnrlNm"`
	EntogPilbkNo       string `xml:"entogPilbkNo"`
	EntogSmplNo        string `xml:"entogSmplNo"`
	FamilyKorNm        string `xml:"familyKorNm"`
	FamilyNm           string `xml:"familyNm"`
	FrstRgstnDtm       string `xml:"frstRgstnDtm"`
	GenusKorNm         string `xml:"genusKorNm"`
	GenusNm            string `xml:"genusNm"`
	ImgURL             string `xml:"imgUrl"`
	JapNm              string `xml:"japNm"`
	LabelUsgCllcnNmplc string `xml:"labelUsgCllcnNmplc"`
	LastUpdtDtm        string `xml:"lastUpdtDtm"`
	OrdKorNm           string `xml:"ordKorNm"`
	OrdNm              string `xml:"ordNm"`
	PrkNm              string `xml:"prkNm"`
	TorsoLngth         string `xml:"torsoLngth"`
	WingLngth          string `xml:"wingLngth"`
}

type entogSpcmInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body entogSpcmInfoBody `xml:"body"`
}

// EntogSpcmInfoError reports an error returned by entogSpcmInfo.
type EntogSpcmInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the entogSpcmInfo error message.
func (e *EntogSpcmInfoError) Error() string {
	return fmt.Sprintf("entogSpcmInfo: API error %s: %s", e.Code, e.Message)
}

// EntogSpcmInfo gets Korea National Arboretum entognath specimen detail information.
func (c *Client) EntogSpcmInfo(ctx context.Context, query application.EntogSpcmInfoQuery) (application.EntogSpcmInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, entogSpcmInfoPath)
	if err != nil {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("q1", query.Q1)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload entogSpcmInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.EntogSpcmInfoResult{}, &EntogSpcmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.EntogSpcmInfoResult{}, fmt.Errorf("entogSpcmInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.EntogSpcmInfoResult{}, errors.New("entogSpcmInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != entogSpcmInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = entogSpcmInfoResultMessages[payload.Header.ResultCode]
		}
		return application.EntogSpcmInfoResult{}, &EntogSpcmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body entogSpcmInfoBody) result() application.EntogSpcmInfoResult {
	if body.Item == nil {
		return application.EntogSpcmInfoResult{}
	}

	return application.EntogSpcmInfoResult{Item: &application.EntogSpcmInfoItem{
		Btnc:               body.Item.Btnc,
		ChnNm:              body.Item.ChnNm,
		ClarHaslvVal:       body.Item.ClarHaslvVal,
		ClctDyDesc:         body.Item.ClctDyDesc,
		CprtCtnt:           body.Item.CprtCtnt,
		EngNm:              body.Item.EngNm,
		EntogGnrlNm:        body.Item.EntogGnrlNm,
		EntogPilbkNo:       body.Item.EntogPilbkNo,
		EntogSmplNo:        body.Item.EntogSmplNo,
		FamilyKorNm:        body.Item.FamilyKorNm,
		FamilyNm:           body.Item.FamilyNm,
		FrstRgstnDtm:       body.Item.FrstRgstnDtm,
		GenusKorNm:         body.Item.GenusKorNm,
		GenusNm:            body.Item.GenusNm,
		ImgURL:             body.Item.ImgURL,
		JapNm:              body.Item.JapNm,
		LabelUsgCllcnNmplc: body.Item.LabelUsgCllcnNmplc,
		LastUpdtDtm:        body.Item.LastUpdtDtm,
		OrdKorNm:           body.Item.OrdKorNm,
		OrdNm:              body.Item.OrdNm,
		PrkNm:              body.Item.PrkNm,
		TorsoLngth:         body.Item.TorsoLngth,
		WingLngth:          body.Item.WingLngth,
	}}
}
