package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/childservice/application/port/outbound"
)

const (
	childPilbkSearchPath        = childServiceBasePath + "/childPilbkSearch"
	childPilbkSearchSuccessCode = "00"
)

var childPilbkSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.ChildPilbkSearchPort = (*Client)(nil)

type childPilbkSearchBody struct {
	Items      []childPilbkSearchItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type childPilbkSearchItem struct {
	BiogyNm           string `xml:"biogyNm"`
	ChildLvbngPilbkNo string `xml:"childLvbngPilbkNo"`
	FamilyKorNm       string `xml:"familyKorNm"`
	FamilyNm          string `xml:"familyNm"`
	LvbngTpcdNm       string `xml:"lvbngTpcdNm"`
	LvngKrlngNm       string `xml:"lvngKrlngNm"`
}

type childPilbkSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body childPilbkSearchBody `xml:"body"`
}

// ChildPilbkSearchError reports an error returned by childPilbkSearch.
type ChildPilbkSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the childPilbkSearch error message.
func (e *ChildPilbkSearchError) Error() string {
	return fmt.Sprintf("childPilbkSearch: API error %s: %s", e.Code, e.Message)
}

// ChildPilbkSearch searches the Korea National Arboretum child pictorial book.
func (c *Client) ChildPilbkSearch(ctx context.Context, query application.ChildPilbkSearchQuery) (application.ChildPilbkSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, childPilbkSearchPath)
	if err != nil {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload childPilbkSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.ChildPilbkSearchResult{}, &ChildPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.ChildPilbkSearchResult{}, fmt.Errorf("childPilbkSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.ChildPilbkSearchResult{}, errors.New("childPilbkSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != childPilbkSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = childPilbkSearchResultMessages[payload.Header.ResultCode]
		}
		return application.ChildPilbkSearchResult{}, &ChildPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body childPilbkSearchBody) result() application.ChildPilbkSearchResult {
	items := make([]application.ChildPilbkSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.ChildPilbkSearchItem{
			BiogyNm:           item.BiogyNm,
			ChildLvbngPilbkNo: item.ChildLvbngPilbkNo,
			FamilyKorNm:       item.FamilyKorNm,
			FamilyNm:          item.FamilyNm,
			LvbngTpcdNm:       item.LvbngTpcdNm,
			LvngKrlngNm:       item.LvngKrlngNm,
		}
	}

	return application.ChildPilbkSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
