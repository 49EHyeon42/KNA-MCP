package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/plantresource/application/port/outbound"
)

const (
	plantWordListPath        = plantResourceBasePath + "/plantWordList"
	plantWordListSuccessCode = "00"
)

var plantWordListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantWordListPort = (*Client)(nil)

type plantWordListBody struct {
	Items      []plantWordListItem `xml:"items>item"`
	NumOfRows  int                 `xml:"numOfRows"`
	PageNo     int                 `xml:"pageNo"`
	TotalCount int                 `xml:"totalCount"`
}

type plantWordListItem struct {
	EnglsWrdNm string `xml:"englsWrdNm"`
	KrnWrdNm   string `xml:"krnWrdNm"`
	PrfcnWrdNm string `xml:"prfcnWrdNm"`
	Wrddscrt   string `xml:"wrddscrt"`
}

type plantWordListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantWordListBody `xml:"body"`
}

// PlantWordListError reports an error returned by plantWordList.
type PlantWordListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantWordList error message.
func (e *PlantWordListError) Error() string {
	return fmt.Sprintf("plantWordList: API error %s: %s", e.Code, e.Message)
}

// PlantWordList gets Korea National Arboretum plant word information.
func (c *Client) PlantWordList(ctx context.Context, query application.PlantWordListQuery) (application.PlantWordListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantWordListPath)
	if err != nil {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantWordListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantWordListResult{}, &PlantWordListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantWordListResult{}, fmt.Errorf("plantWordList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantWordListResult{}, errors.New("plantWordList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantWordListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantWordListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantWordListResult{}, &PlantWordListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantWordListBody) result() application.PlantWordListResult {
	items := make([]application.PlantWordListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantWordListItem{
			EnglsWrdNm: item.EnglsWrdNm,
			KrnWrdNm:   item.KrnWrdNm,
			PrfcnWrdNm: item.PrfcnWrdNm,
			Wrddscrt:   item.Wrddscrt,
		}
	}

	return application.PlantWordListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
