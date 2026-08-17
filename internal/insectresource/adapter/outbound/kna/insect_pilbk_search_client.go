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
	insectPilbkSearchPath        = insectResourceBasePath + "/insectPilbkSearch"
	insectPilbkSearchSuccessCode = "00"
)

var insectPilbkSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.InsectPilbkSearchPort = (*Client)(nil)

type insectPilbkSearchBody struct {
	Items      []insectPilbkSearchItem `xml:"items>item"`
	NumOfRows  int                     `xml:"numOfRows"`
	PageNo     int                     `xml:"pageNo"`
	TotalCount int                     `xml:"totalCount"`
}

type insectPilbkSearchItem struct {
	FamilyKorNm    string `xml:"familyKorNm"`
	FamilyNm       string `xml:"familyNm"`
	GenusKorNm     string `xml:"genusKorNm"`
	GenusNm        string `xml:"genusNm"`
	InsctGnrlNm    string `xml:"insctGnrlNm"`
	InsctPilbkNo   string `xml:"insctPilbkNo"`
	InsctSpecsScnm string `xml:"insctSpecsScnm"`
	LastUpdtDtm    string `xml:"lastUpdtDtm"`
}

type insectPilbkSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body insectPilbkSearchBody `xml:"body"`
}

// InsectPilbkSearchError reports an error returned by insectPilbkSearch.
type InsectPilbkSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the insectPilbkSearch error message.
func (e *InsectPilbkSearchError) Error() string {
	return fmt.Sprintf("insectPilbkSearch: API error %s: %s", e.Code, e.Message)
}

// InsectPilbkSearch searches the Korea National Arboretum insect pictorial book.
func (c *Client) InsectPilbkSearch(ctx context.Context, query application.InsectPilbkSearchQuery) (application.InsectPilbkSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, insectPilbkSearchPath)
	if err != nil {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	if query.ReqSearchWrd != "" {
		values.Set("reqSearchWrd", query.ReqSearchWrd)
	}
	if query.DateFrom != "" {
		values.Set("dateFrom", query.DateFrom)
	}
	if query.DateTo != "" {
		values.Set("dateTo", query.DateTo)
	}
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload insectPilbkSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.InsectPilbkSearchResult{}, &InsectPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.InsectPilbkSearchResult{}, fmt.Errorf("insectPilbkSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.InsectPilbkSearchResult{}, errors.New("insectPilbkSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != insectPilbkSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = insectPilbkSearchResultMessages[payload.Header.ResultCode]
		}
		return application.InsectPilbkSearchResult{}, &InsectPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body insectPilbkSearchBody) result() application.InsectPilbkSearchResult {
	items := make([]application.InsectPilbkSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.InsectPilbkSearchItem{
			FamilyKorNm:    item.FamilyKorNm,
			FamilyNm:       item.FamilyNm,
			GenusKorNm:     item.GenusKorNm,
			GenusNm:        item.GenusNm,
			InsctGnrlNm:    item.InsctGnrlNm,
			InsctPilbkNo:   item.InsctPilbkNo,
			InsctSpecsScnm: item.InsctSpecsScnm,
			LastUpdtDtm:    item.LastUpdtDtm,
		}
	}

	return application.InsectPilbkSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
