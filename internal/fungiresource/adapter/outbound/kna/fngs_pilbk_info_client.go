package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

const (
	fngsPilbkInfoPath        = fungiResourceBasePath + "/fngsPilbkInfo"
	fngsPilbkInfoSuccessCode = "00"
)

var fngsPilbkInfoResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.FngsPilbkInfoPort = (*Client)(nil)

type fngsPilbkInfoBody struct {
	Item *fngsPilbkInfoItem `xml:"item"`
}

type fngsPilbkInfoItem struct {
	MshrmColorCdNm      string `xml:"mshrmColorCdNm"`
	CrpphFomTpcdNm      string `xml:"crpphFomTpcdNm"`
	FamilyKorNm         string `xml:"familyKorNm"`
	FamilyNm            string `xml:"familyNm"`
	FngsEclgTpcdNm      string `xml:"fngsEclgTpcdNm"`
	FngsGnrlNm          string `xml:"fngsGnrlNm"`
	FngsPilbkNo         string `xml:"fngsPilbkNo"`
	FngsPrpseTpcdNm     string `xml:"fngsPrpseTpcdNm"`
	FngsScnm            string `xml:"fngsScnm"`
	GenusKorNm          string `xml:"genusKorNm"`
	GenusNm             string `xml:"genusNm"`
	GrwEvrntDesc        string `xml:"grwEvrntDesc"`
	LastUpdtDtm         string `xml:"lastUpdtDtm"`
	MicroShpe           string `xml:"microShpe"`
	MshrmTpcdNm         string `xml:"mshrmTpcdNm"`
	OccrrSsnNm          string `xml:"occrrSsnNm"`
	RsrcActoClsscTpcdNm string `xml:"rsrcActoClsscTpcdNm"`
	Shpe                string `xml:"shpe"`
}

type fngsPilbkInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body fngsPilbkInfoBody `xml:"body"`
}

// FngsPilbkInfoError reports an error returned by fngsPilbkInfo.
type FngsPilbkInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the fngsPilbkInfo error message.
func (e *FngsPilbkInfoError) Error() string {
	return fmt.Sprintf("fngsPilbkInfo: API error %s: %s", e.Code, e.Message)
}

// FngsPilbkInfo gets Korea National Arboretum fungi pictorial book detail information.
func (c *Client) FngsPilbkInfo(ctx context.Context, query application.FngsPilbkInfoQuery) (application.FngsPilbkInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, fngsPilbkInfoPath)
	if err != nil {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqFngsPilbkNo", query.ReqFngsPilbkNo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload fngsPilbkInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.FngsPilbkInfoResult{}, &FngsPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.FngsPilbkInfoResult{}, fmt.Errorf("fngsPilbkInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.FngsPilbkInfoResult{}, errors.New("fngsPilbkInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != fngsPilbkInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = fngsPilbkInfoResultMessages[payload.Header.ResultCode]
		}
		return application.FngsPilbkInfoResult{}, &FngsPilbkInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body fngsPilbkInfoBody) result() application.FngsPilbkInfoResult {
	if body.Item == nil {
		return application.FngsPilbkInfoResult{}
	}

	return application.FngsPilbkInfoResult{Item: &application.FngsPilbkInfoItem{
		MshrmColorCdNm:      body.Item.MshrmColorCdNm,
		CrpphFomTpcdNm:      body.Item.CrpphFomTpcdNm,
		FamilyKorNm:         body.Item.FamilyKorNm,
		FamilyNm:            body.Item.FamilyNm,
		FngsEclgTpcdNm:      body.Item.FngsEclgTpcdNm,
		FngsGnrlNm:          body.Item.FngsGnrlNm,
		FngsPilbkNo:         body.Item.FngsPilbkNo,
		FngsPrpseTpcdNm:     body.Item.FngsPrpseTpcdNm,
		FngsScnm:            body.Item.FngsScnm,
		GenusKorNm:          body.Item.GenusKorNm,
		GenusNm:             body.Item.GenusNm,
		GrwEvrntDesc:        body.Item.GrwEvrntDesc,
		LastUpdtDtm:         body.Item.LastUpdtDtm,
		MicroShpe:           body.Item.MicroShpe,
		MshrmTpcdNm:         body.Item.MshrmTpcdNm,
		OccrrSsnNm:          body.Item.OccrrSsnNm,
		RsrcActoClsscTpcdNm: body.Item.RsrcActoClsscTpcdNm,
		Shpe:                body.Item.Shpe,
	}}
}
