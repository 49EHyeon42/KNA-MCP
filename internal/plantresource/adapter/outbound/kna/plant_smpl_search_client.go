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
	plantSmplSearchPath        = plantResourceBasePath + "/plantSmplSearch"
	plantSmplSearchSuccessCode = "00"
)

var plantSmplSearchResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSampleSearchPort = (*Client)(nil)

type plantSmplSearchBody struct {
	Items      []plantSmplSearchItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type plantSmplSearchItem struct {
	Cnt            int    `xml:"cnt"`
	FamilyKorNm    string `xml:"familyKorNm"`
	FamilyNm       string `xml:"familyNm"`
	PlantGnrlNm    string `xml:"plantGnrlNm"`
	PlantSpecsID   string `xml:"plantSpecsId"`
	PlantSpecsScnm string `xml:"plantSpecsScnm"`
}

type plantSmplSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSmplSearchBody `xml:"body"`
}

// PlantSmplSearchError reports an error returned by plantSmplSearch.
type PlantSmplSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSmplSearch error message.
func (e *PlantSmplSearchError) Error() string {
	return fmt.Sprintf("plantSmplSearch: API error %s: %s", e.Code, e.Message)
}

// PlantSampleSearch searches the Korea National Arboretum plant samples.
func (c *Client) PlantSampleSearch(ctx context.Context, query application.PlantSampleSearchQuery) (application.PlantSampleSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSmplSearchPath)
	if err != nil {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNumber))
	values.Set("numOfRows", strconv.Itoa(query.NumberOfRows))
	setQueryValue(values, "reqSearchWrd", query.RequestSearchWord)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSmplSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSampleSearchResult{}, &PlantSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSampleSearchResult{}, fmt.Errorf("plantSmplSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSampleSearchResult{}, errors.New("plantSmplSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSmplSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSmplSearchResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSampleSearchResult{}, &PlantSmplSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSmplSearchBody) result() application.PlantSampleSearchResult {
	items := make([]application.PlantSampleSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSampleSearchItem{
			Count:                      item.Cnt,
			FamilyKoreanName:           item.FamilyKorNm,
			FamilyName:                 item.FamilyNm,
			PlantGeneralName:           item.PlantGnrlNm,
			PlantSpeciesID:             item.PlantSpecsID,
			PlantSpeciesScientificName: item.PlantSpecsScnm,
		}
	}

	return application.PlantSampleSearchResult{
		Items:        items,
		NumberOfRows: body.NumOfRows,
		PageNumber:   body.PageNo,
		TotalCount:   body.TotalCount,
	}
}
