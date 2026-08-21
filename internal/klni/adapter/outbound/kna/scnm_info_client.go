package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/49EHyeon42/KNA-MCP/internal/klni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/klni/application/port/outbound"
)

const (
	scnmInfoPath        = klniBasePath + "/scnmInfo"
	scnmInfoSuccessCode = "00"
)

var scnmInfoResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.ScnmInfoPort = (*Client)(nil)

type scnmInfoBody struct {
	Item *scnmInfoItem `xml:"item"`
}

type scnmInfoItem struct {
	StpltScnmRltnCdNm string `xml:"stpltScnmRltnCdNm"`
	FalmNm            string `xml:"falmNm"`
	FalnKorNm         string `xml:"falnKorNm"`
	GenusKorNm        string `xml:"genusKorNm"`
	GenusNm           string `xml:"genusNm"`
	LastUpdtDtm       string `xml:"lastUpdtDtm"`
	LchnGnrlNm        string `xml:"lchnGnrlNm"`
	LchnGnrlNm2       string `xml:"lchnGnrlNm2"`
	LchnScnm          string `xml:"lchnScnm"`
	LchnScnmID        string `xml:"lchnScnmId"`
	OrdscLtrtrNm      string `xml:"ordscLtrtrNm"`
	Rmrk              string `xml:"rmrk"`
}

type scnmInfoResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body scnmInfoBody `xml:"body"`
}

// ScnmInfoError reports an error returned by scnmInfo.
type ScnmInfoError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the scnmInfo error message.
func (e *ScnmInfoError) Error() string {
	return fmt.Sprintf("scnmInfo: API error %s: %s", e.Code, e.Message)
}

// ScnmInfo gets Korea National Arboretum lichen scientific name detail information.
func (c *Client) ScnmInfo(ctx context.Context, query application.ScnmInfoQuery) (application.ScnmInfoResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, scnmInfoPath)
	if err != nil {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("reqLchnScnmId", query.ReqLchnScnmID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: request: %w", err)
	}
	defer response.Body.Close()

	var payload scnmInfoResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.ScnmInfoResult{}, &ScnmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.ScnmInfoResult{}, fmt.Errorf("scnmInfo: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.ScnmInfoResult{}, errors.New("scnmInfo: response missing resultCode")
	}
	if payload.Header.ResultCode != scnmInfoSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = scnmInfoResultMessages[payload.Header.ResultCode]
		}
		return application.ScnmInfoResult{}, &ScnmInfoError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body scnmInfoBody) result() application.ScnmInfoResult {
	if body.Item == nil {
		return application.ScnmInfoResult{}
	}

	return application.ScnmInfoResult{Item: &application.ScnmInfoItem{
		StpltScnmRltnCdNm: body.Item.StpltScnmRltnCdNm,
		FalmNm:            body.Item.FalmNm,
		FalnKorNm:         body.Item.FalnKorNm,
		GenusKorNm:        body.Item.GenusKorNm,
		GenusNm:           body.Item.GenusNm,
		LastUpdtDtm:       body.Item.LastUpdtDtm,
		LchnGnrlNm:        body.Item.LchnGnrlNm,
		LchnGnrlNm2:       body.Item.LchnGnrlNm2,
		LchnScnm:          body.Item.LchnScnm,
		LchnScnmID:        body.Item.LchnScnmID,
		OrdscLtrtrNm:      body.Item.OrdscLtrtrNm,
		Rmrk:              body.Item.Rmrk,
	}}
}
