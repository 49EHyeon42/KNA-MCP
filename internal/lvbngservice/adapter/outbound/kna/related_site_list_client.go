package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application"
	"github.com/49EHyeon42/KNA-MCP/internal/lvbngservice/application/port/outbound"
)

const (
	relatedSiteListPath        = "/1400119/LvbngService2/relatedSiteList"
	relatedSiteListSuccessCode = "00"
)

var relatedSiteListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.RelatedSiteListPort = (*Client)(nil)

type relatedSiteListBody struct {
	Items      []relatedSiteListItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type relatedSiteListItem struct {
	LvbngTpcdNm string `xml:"lvbngTpcdNm"`
	SiteCtgryNm string `xml:"siteCtgryNm"`
	SiteNm      string `xml:"siteNm"`
	SiteURL     string `xml:"siteUrl"`
}

type relatedSiteListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body relatedSiteListBody `xml:"body"`
}

// RelatedSiteListError reports an error returned by relatedSiteList.
type RelatedSiteListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the relatedSiteList error message.
func (e *RelatedSiteListError) Error() string {
	return fmt.Sprintf("relatedSiteList: API error %s: %s", e.Code, e.Message)
}

// RelatedSiteList gets Korea National Arboretum biological related site information.
func (c *Client) RelatedSiteList(ctx context.Context, query application.RelatedSiteListQuery) (application.RelatedSiteListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, relatedSiteListPath)
	if err != nil {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: request: %w", err)
	}
	defer response.Body.Close()

	var payload relatedSiteListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.RelatedSiteListResult{}, &RelatedSiteListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.RelatedSiteListResult{}, fmt.Errorf("relatedSiteList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.RelatedSiteListResult{}, errors.New("relatedSiteList: response missing resultCode")
	}
	if payload.Header.ResultCode != relatedSiteListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = relatedSiteListResultMessages[payload.Header.ResultCode]
		}
		return application.RelatedSiteListResult{}, &RelatedSiteListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body relatedSiteListBody) result() application.RelatedSiteListResult {
	items := make([]application.RelatedSiteListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.RelatedSiteListItem{
			LvbngTpcdNm: item.LvbngTpcdNm,
			SiteCtgryNm: item.SiteCtgryNm,
			SiteNm:      item.SiteNm,
			SiteURL:     item.SiteURL,
		}
	}

	return application.RelatedSiteListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
