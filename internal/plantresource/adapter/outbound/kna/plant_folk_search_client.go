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
	plantFolkSearchPath        = plantResourceBasePath + "/plantFolkSearch"
	plantFolkSearchSuccessCode = "00"
)

var plantFolkSearchResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantFolkSearchPort = (*Client)(nil)

type plantFolkSearchBody struct {
	Items      []plantFolkSearchItem `xml:"items>item"`
	NumOfRows  int                   `xml:"numOfRows"`
	PageNo     int                   `xml:"pageNo"`
	TotalCount int                   `xml:"totalCount"`
}

type plantFolkSearchItem struct {
	FlcstPlantIdntfDscrt string `xml:"flcstPlantIdntfDscrt"`
	FlpltID              string `xml:"flpltId"`
	PlantBrdgFomTpcdNm   string `xml:"plantBrdgFomTpcdNm"`
	PlantGnrlNm          string `xml:"plantGnrlNm"`
	PlantSpecsScnm       string `xml:"plantSpecsScnm"`
	Ptnt                 string `xml:"ptnt"`
}

type plantFolkSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantFolkSearchBody `xml:"body"`
}

// PlantFolkSearchError reports an error returned by plantFolkSearch.
type PlantFolkSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantFolkSearch error message.
func (e *PlantFolkSearchError) Error() string {
	return fmt.Sprintf("plantFolkSearch: API error %s: %s", e.Code, e.Message)
}

// PlantFolkSearch searches the Korea National Arboretum folk plants.
func (c *Client) PlantFolkSearch(ctx context.Context, query application.PlantFolkSearchQuery) (application.PlantFolkSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantFolkSearchPath)
	if err != nil {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantFolkSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantFolkSearchResult{}, &PlantFolkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantFolkSearchResult{}, fmt.Errorf("plantFolkSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantFolkSearchResult{}, errors.New("plantFolkSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != plantFolkSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantFolkSearchResultMessages[payload.Header.ResultCode]
		}
		return application.PlantFolkSearchResult{}, &PlantFolkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantFolkSearchBody) result() application.PlantFolkSearchResult {
	items := make([]application.PlantFolkSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantFolkSearchItem{
			FlcstPlantIdntfDscrt: item.FlcstPlantIdntfDscrt,
			FlpltID:              item.FlpltID,
			PlantBrdgFomTpcdNm:   item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:          item.PlantGnrlNm,
			PlantSpecsScnm:       item.PlantSpecsScnm,
			Ptnt:                 item.Ptnt,
		}
	}

	return application.PlantFolkSearchResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
