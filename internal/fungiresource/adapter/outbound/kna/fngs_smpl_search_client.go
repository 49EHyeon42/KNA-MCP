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
	fngsSmplSearchPath        = fungiResourceBasePath + "/fngsSmplSearch"
	fngsSmplSearchSuccessCode = "00"
)

var fngsSmplSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.FngsSmplSearchPort = (*Client)(nil)

type fngsSmplSearchBody struct {
	Items      []fngsSmplSearchItem `xml:"items>item"`
	NumOfRows  int                  `xml:"numOfRows"`
	PageNo     int                  `xml:"pageNo"`
	TotalCount int                  `xml:"totalCount"`
}

type fngsSmplSearchItem struct {
	Cnt         string `xml:"cnt"`
	FamilyKorNm string `xml:"familyKorNm"`
	FamilyNm    string `xml:"familyNm"`
	FngsGnrlNm  string `xml:"fngsGnrlNm"`
	FngsID      string `xml:"fngsId"`
	FngsScnm    string `xml:"fngsScnm"`
	GenusKorNm  string `xml:"genusKorNm"`
	GenusNm     string `xml:"genusNm"`
}

type fngsSmplSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body fngsSmplSearchBody `xml:"body"`
}

// FngsSmplSearchError reports an error returned by fngsSmplSearch.
type FngsSmplSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the fngsSmplSearch error message.
func (e *FngsSmplSearchError) Error() string {
	return fmt.Sprintf("fngsSmplSearch: API error %s: %s", e.Code, e.Message)
}

// FngsSmplSearch searches Korea National Arboretum fungi samples.
func (c *Client) FngsSmplSearch(ctx context.Context, query application.FngsSmplSearchQuery) (application.FngsSmplSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, fngsSmplSearchPath)
	if err != nil {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload fngsSmplSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.FngsSmplSearchResult{}, &FngsSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.FngsSmplSearchResult{}, fmt.Errorf("fngsSmplSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.FngsSmplSearchResult{}, errors.New("fngsSmplSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != fngsSmplSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = fngsSmplSearchResultMessages[payload.Header.ResultCode]
		}
		return application.FngsSmplSearchResult{}, &FngsSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body fngsSmplSearchBody) result() application.FngsSmplSearchResult {
	items := make([]application.FngsSmplSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.FngsSmplSearchItem{
			Cnt:         item.Cnt,
			FamilyKorNm: item.FamilyKorNm,
			FamilyNm:    item.FamilyNm,
			FngsGnrlNm:  item.FngsGnrlNm,
			FngsID:      item.FngsID,
			FngsScnm:    item.FngsScnm,
			GenusKorNm:  item.GenusKorNm,
			GenusNm:     item.GenusNm,
		}
	}

	return application.FngsSmplSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
