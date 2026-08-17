package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

const (
	insectPrtctListPath        = insectResourceBasePath + "/insectPrtctList"
	insectPrtctListSuccessCode = "00"
)

var insectPrtctListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.InsectPrtctListPort = (*Client)(nil)

type insectPrtctListBody struct {
	Items      []insectPrtctListItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type insectPrtctListItem struct {
	FamilyKorNm    string `xml:"familyKorNm"`
	FamilyNm       string `xml:"familyNm"`
	InsctGnrlNm    string `xml:"insctGnrlNm"`
	InsctPcmtt     string `xml:"insctPcmtt"`
	InsctPilbkNo   string `xml:"insctPilbkNo"`
	InsctSpecsScnm string `xml:"insctSpecsScnm"`
}

type insectPrtctListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body insectPrtctListBody `xml:"body"`
}

// InsectPrtctListError reports an error returned by insectPrtctList.
type InsectPrtctListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the insectPrtctList error message.
func (e *InsectPrtctListError) Error() string {
	return fmt.Sprintf("insectPrtctList: API error %s: %s", e.Code, e.Message)
}

// InsectPrtctList returns endangered insects from the Korea National Arboretum.
func (c *Client) InsectPrtctList(ctx context.Context, query application.InsectPrtctListQuery) (application.InsectPrtctListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, insectPrtctListPath)
	if err != nil {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: request: %w", err)
	}
	defer response.Body.Close()

	var payload insectPrtctListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.InsectPrtctListResult{}, &InsectPrtctListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.InsectPrtctListResult{}, fmt.Errorf("insectPrtctList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.InsectPrtctListResult{}, errors.New("insectPrtctList: response missing resultCode")
	}
	if payload.Header.ResultCode != insectPrtctListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = insectPrtctListResultMessages[payload.Header.ResultCode]
		}
		return application.InsectPrtctListResult{}, &InsectPrtctListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body insectPrtctListBody) result() application.InsectPrtctListResult {
	items := make([]application.InsectPrtctListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.InsectPrtctListItem{
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			InsctGnrlNm:    item.InsctGnrlNm,
			InsctPcmtt:     item.InsctPcmtt,
			InsctPilbkNo:   item.InsctPilbkNo,
			InsctSpecsScnm: item.InsctSpecsScnm,
		}
	}

	return application.InsectPrtctListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
