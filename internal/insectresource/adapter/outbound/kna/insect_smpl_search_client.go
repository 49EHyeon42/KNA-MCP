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
	insectSmplSearchPath        = insectResourceBasePath + "/insectSmplSearch"
	insectSmplSearchSuccessCode = "00"
)

var insectSmplSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.InsectSmplSearchPort = (*Client)(nil)

type insectSmplSearchBody struct {
	Items      []insectSmplSearchItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type insectSmplSearchItem struct {
	Cnt              string `xml:"cnt"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	InsctGnrlNm      string `xml:"insctGnrlNm"`
	InsctSpecsID     string `xml:"insctSpecsId"`
	InsctSpecsScnm   string `xml:"insctSpecsScnm"`
	SubFamilyKorNm   string `xml:"subFamilyKorNm"`
	SubFamilyNm      string `xml:"subFamilyNm"`
	SuperFamilyKorNm string `xml:"superFamilyKorNm"`
	SuperFamilyNm    string `xml:"superFamilyNm"`
}

type insectSmplSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body insectSmplSearchBody `xml:"body"`
}

// InsectSmplSearchError reports an error returned by insectSmplSearch.
type InsectSmplSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the insectSmplSearch error message.
func (e *InsectSmplSearchError) Error() string {
	return fmt.Sprintf("insectSmplSearch: API error %s: %s", e.Code, e.Message)
}

// InsectSmplSearch searches the Korea National Arboretum insect samples.
func (c *Client) InsectSmplSearch(ctx context.Context, query application.InsectSmplSearchQuery) (application.InsectSmplSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, insectSmplSearchPath)
	if err != nil {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload insectSmplSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.InsectSmplSearchResult{}, &InsectSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.InsectSmplSearchResult{}, fmt.Errorf("insectSmplSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.InsectSmplSearchResult{}, errors.New("insectSmplSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != insectSmplSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = insectSmplSearchResultMessages[payload.Header.ResultCode]
		}
		return application.InsectSmplSearchResult{}, &InsectSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body insectSmplSearchBody) result() application.InsectSmplSearchResult {
	items := make([]application.InsectSmplSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.InsectSmplSearchItem{
			Cnt:              item.Cnt,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			InsctGnrlNm:      item.InsctGnrlNm,
			InsctSpecsID:     item.InsctSpecsID,
			InsctSpecsScnm:   item.InsctSpecsScnm,
			SubFamilyKorNm:   item.SubFamilyKorNm,
			SubFamilyNm:      item.SubFamilyNm,
			SuperFamilyKorNm: item.SuperFamilyKorNm,
			SuperFamilyNm:    item.SuperFamilyNm,
		}
	}

	return application.InsectSmplSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
