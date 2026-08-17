package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application"
	"github.com/49EHyeon42/KNA-MCP/internal/kpni/application/port/outbound"
)

const (
	gnrlNmLtrtrSearchPath        = kpniBasePath + "/gnrlNmLtrtrSearch"
	gnrlNmLtrtrSearchSuccessCode = "00"
)

var gnrlNmLtrtrSearchResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.GnrlNmLtrtrSearchPort = (*Client)(nil)

type gnrlNmLtrtrSearchBody struct {
	Items      []gnrlNmLtrtrSearchItem `xml:"items>item"`
	NumOfRows  int                     `xml:"numOfRows"`
	PageNo     int                     `xml:"pageNo"`
	TotalCount int                     `xml:"totalCount"`
}

type gnrlNmLtrtrSearchItem struct {
	RcmmnTpcdNm      string `xml:"rcmmnTpcdNm"`
	LtrtrInfrmNm     string `xml:"ltrtrInfrmNm"`
	LvbngFrlngTpcdNm string `xml:"lvbngFrlngTpcdNm"`
	PlantGnrlNm      string `xml:"plantGnrlNm"`
	PlantSpecsScnm   string `xml:"plantSpecsScnm"`
}

type gnrlNmLtrtrSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body gnrlNmLtrtrSearchBody `xml:"body"`
}

// GnrlNmLtrtrSearchError reports an error returned by gnrlNmLtrtrSearch.
type GnrlNmLtrtrSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the gnrlNmLtrtrSearch error message.
func (e *GnrlNmLtrtrSearchError) Error() string {
	return fmt.Sprintf("gnrlNmLtrtrSearch: API error %s: %s", e.Code, e.Message)
}

// GnrlNmLtrtrSearch gets Korea National Arboretum plant general name literature information.
func (c *Client) GnrlNmLtrtrSearch(ctx context.Context, query application.GnrlNmLtrtrSearchQuery) (application.GnrlNmLtrtrSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, gnrlNmLtrtrSearchPath)
	if err != nil {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	if query.ReqPlantGnrlNm != "" {
		values.Set("reqPlantGnrlNm", query.ReqPlantGnrlNm)
	}
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload gnrlNmLtrtrSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.GnrlNmLtrtrSearchResult{}, &GnrlNmLtrtrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.GnrlNmLtrtrSearchResult{}, fmt.Errorf("gnrlNmLtrtrSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.GnrlNmLtrtrSearchResult{}, errors.New("gnrlNmLtrtrSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != gnrlNmLtrtrSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = gnrlNmLtrtrSearchResultMessages[payload.Header.ResultCode]
		}
		return application.GnrlNmLtrtrSearchResult{}, &GnrlNmLtrtrSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body gnrlNmLtrtrSearchBody) result() application.GnrlNmLtrtrSearchResult {
	items := make([]application.GnrlNmLtrtrSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.GnrlNmLtrtrSearchItem{
			RcmmnTpcdNm:      item.RcmmnTpcdNm,
			LtrtrInfrmNm:     item.LtrtrInfrmNm,
			LvbngFrlngTpcdNm: item.LvbngFrlngTpcdNm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
		}
	}

	return application.GnrlNmLtrtrSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
