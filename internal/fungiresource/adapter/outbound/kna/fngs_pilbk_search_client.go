package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

const (
	fngsPilbkSearchPath        = fungiResourceBasePath + "/fngsPilbkSearch"
	fngsPilbkSearchSuccessCode = "00"
)

var fngsPilbkSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.FngsPilbkSearchPort = (*Client)(nil)

type fngsPilbkSearchBody struct {
	Items      []fngsPilbkSearchItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type fngsPilbkSearchItem struct {
	FamilyKorNm string `xml:"familyKorNm"`
	FamilyNm    string `xml:"familyNm"`
	FngsGnrlNm  string `xml:"fngsGnrlNm"`
	FngsPilbkNo string `xml:"fngsPilbkNo"`
	FngsScnm    string `xml:"fngsScnm"`
	GenusKorNm  string `xml:"genusKorNm"`
	GenusNm     string `xml:"genusNm"`
	LastUpdtDtm string `xml:"lastUpdtDtm"`
}

type fngsPilbkSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body fngsPilbkSearchBody `xml:"body"`
}

// FngsPilbkSearchError reports an error returned by fngsPilbkSearch.
type FngsPilbkSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the fngsPilbkSearch error message.
func (e *FngsPilbkSearchError) Error() string {
	return fmt.Sprintf("fngsPilbkSearch: API error %s: %s", e.Code, e.Message)
}

// FngsPilbkSearch searches the Korea National Arboretum fungi pictorial book.
func (c *Client) FngsPilbkSearch(ctx context.Context, query application.FngsPilbkSearchQuery) (application.FngsPilbkSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, fngsPilbkSearchPath)
	if err != nil {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	setQueryValue(values, "dateFrom", query.DateFrom)
	setQueryValue(values, "dateTo", query.DateTo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload fngsPilbkSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.FngsPilbkSearchResult{}, &FngsPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.FngsPilbkSearchResult{}, fmt.Errorf("fngsPilbkSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.FngsPilbkSearchResult{}, errors.New("fngsPilbkSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != fngsPilbkSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = fngsPilbkSearchResultMessages[payload.Header.ResultCode]
		}
		return application.FngsPilbkSearchResult{}, &FngsPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body fngsPilbkSearchBody) result() application.FngsPilbkSearchResult {
	items := make([]application.FngsPilbkSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.FngsPilbkSearchItem{
			FamilyKorNm: item.FamilyKorNm,
			FamilyNm:    item.FamilyNm,
			FngsGnrlNm:  item.FngsGnrlNm,
			FngsPilbkNo: item.FngsPilbkNo,
			FngsScnm:    item.FngsScnm,
			GenusKorNm:  item.GenusKorNm,
			GenusNm:     item.GenusNm,
			LastUpdtDtm: item.LastUpdtDtm,
		}
	}

	return application.FngsPilbkSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
