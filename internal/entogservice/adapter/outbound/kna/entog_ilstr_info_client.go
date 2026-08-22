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
	entogIlstrInfoPath        = entogServiceBasePath + "/entogIlstrInfo"
	entogIlstrInfoSuccessCode = "00"
)

var entogIlstrInfoResultMessages = map[string]string{
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

var _ outbound.EntogIlstrInfoPort = (*Client)(nil)

type entogIlstrInfoBody struct {
	Item *entogIlstrInfoItem `xml:"item"`
}

type entogIlstrInfoItem struct {
	Btnc             string `xml:"btnc"`
	Cont1            string `xml:"cont1"`
	Cont2            string `xml:"cont2"`
	Cont3            string `xml:"cont3"`
	Cont4            string `xml:"cont4"`
	Cont5            string `xml:"cont5"`
	Cont6            string `xml:"cont6"`
	Cont7            string `xml:"cont7"`
	Cont8            string `xml:"cont8"`
	Cont9            string `xml:"cont9"`
	Cont10           string `xml:"cont10"`
	Cont11           string `xml:"cont11"`
	CprtCtnt         string `xml:"cprtCtnt"`
	EmrgcCnt         string `xml:"emrgcCnt"`
	EmrgcEraDscrt    string `xml:"emrgcEraDscrt"`
	EntogAthrNm      string `xml:"entogAthrNm"`
	EntogEngNm       string `xml:"entogEngNm"`
	EntogOfnmKrlngNm string `xml:"entogOfnmKrlngNm"`
	EntogPilbkNo     string `xml:"entogPilbkNo"`
	EntogSpecsNm     string `xml:"entogSpecsNm"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	ImgURL           string `xml:"imgUrl"`
	MnmmOccrrCnt     string `xml:"mnmmOccrrCnt"`
	MxmmOccrrCnt     string `xml:"mxmmOccrrCnt"`
	OrdKorNm         string `xml:"ordKorNm"`
	OrdNm            string `xml:"ordNm"`
}

type entogIlstrInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body entogIlstrInfoBody `xml:"body"`
}

// EntogIlstrInfoError reports an error returned by entogIlstrInfo.
type EntogIlstrInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the entogIlstrInfo error message.
func (e *EntogIlstrInfoError) Error() string {
	return fmt.Sprintf("entogIlstrInfo: API error %s: %s", e.Code, e.Message)
}

// EntogIlstrInfo gets Korea National Arboretum entognath pictorial book detail information.
func (c *Client) EntogIlstrInfo(ctx context.Context, query application.EntogIlstrInfoQuery) (application.EntogIlstrInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, entogIlstrInfoPath)
	if err != nil {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("q1", query.Q1)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload entogIlstrInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.EntogIlstrInfoResult{}, &EntogIlstrInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.EntogIlstrInfoResult{}, fmt.Errorf("entogIlstrInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.EntogIlstrInfoResult{}, errors.New("entogIlstrInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != entogIlstrInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = entogIlstrInfoResultMessages[payload.Header.ResultCode]
		}
		return application.EntogIlstrInfoResult{}, &EntogIlstrInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body entogIlstrInfoBody) result() application.EntogIlstrInfoResult {
	if body.Item == nil {
		return application.EntogIlstrInfoResult{}
	}

	return application.EntogIlstrInfoResult{Item: &application.EntogIlstrInfoItem{
		Btnc:             body.Item.Btnc,
		Cont1:            body.Item.Cont1,
		Cont2:            body.Item.Cont2,
		Cont3:            body.Item.Cont3,
		Cont4:            body.Item.Cont4,
		Cont5:            body.Item.Cont5,
		Cont6:            body.Item.Cont6,
		Cont7:            body.Item.Cont7,
		Cont8:            body.Item.Cont8,
		Cont9:            body.Item.Cont9,
		Cont10:           body.Item.Cont10,
		Cont11:           body.Item.Cont11,
		CprtCtnt:         body.Item.CprtCtnt,
		EmrgcCnt:         body.Item.EmrgcCnt,
		EmrgcEraDscrt:    body.Item.EmrgcEraDscrt,
		EntogAthrNm:      body.Item.EntogAthrNm,
		EntogEngNm:       body.Item.EntogEngNm,
		EntogOfnmKrlngNm: body.Item.EntogOfnmKrlngNm,
		EntogPilbkNo:     body.Item.EntogPilbkNo,
		EntogSpecsNm:     body.Item.EntogSpecsNm,
		FamilyKorNm:      body.Item.FamilyKorNm,
		FamilyNm:         body.Item.FamilyNm,
		GenusKorNm:       body.Item.GenusKorNm,
		GenusNm:          body.Item.GenusNm,
		ImgURL:           body.Item.ImgURL,
		MnmmOccrrCnt:     body.Item.MnmmOccrrCnt,
		MxmmOccrrCnt:     body.Item.MxmmOccrrCnt,
		OrdKorNm:         body.Item.OrdKorNm,
		OrdNm:            body.Item.OrdNm,
	}}
}
