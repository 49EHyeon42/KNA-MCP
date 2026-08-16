package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/application"
	"github.com/49EHyeon42/KNA-MCP/internal/application/port/outbound"
)

const (
	plantPilbkSearchPath        = plantResourceBasePath + "/plantPilbkSearch"
	plantPilbkSearchSuccessCode = "00"
)

var plantPilbkSearchResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantPictorialBookSearchPort = (*Client)(nil)

type plantPilbkSearchBody struct {
	Items      []plantPilbkSearchItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type plantPilbkSearchItem struct {
	APGFamilyKorNm string `xml:"apgFamilyKorNm"`
	APGFamilyNm    string `xml:"apgFamilyNm"`
	FamilyKorNm    string `xml:"familyKorNm"`
	FamilyNm       string `xml:"familyNm"`
	GenusKorNm     string `xml:"genusKorNm"`
	GenusNm        string `xml:"genusNm"`
	LastUpdtDtm    string `xml:"lastUpdtDtm"`
	NotRcmmGnrlNm  string `xml:"notRcmmGnrlNm"`
	PlantGnrlNm    string `xml:"plantGnrlNm"`
	PlantPilbkNo   string `xml:"plantPilbkNo"`
	PlantSpecsScnm string `xml:"plantSpecsScnm"`
}

type plantPilbkSearchResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantPilbkSearchBody `xml:"body"`
}

// PlantPilbkSearchError reports an error returned by plantPilbkSearch.
type PlantPilbkSearchError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantPilbkSearch error message.
func (e *PlantPilbkSearchError) Error() string {
	return fmt.Sprintf("plantPilbkSearch: API error %s: %s", e.Code, e.Message)
}

// PlantPictorialBookSearch searches the Korea National Arboretum plant pictorial book.
func (c *Client) PlantPictorialBookSearch(ctx context.Context, query application.PlantPictorialBookSearchQuery) (application.PlantPictorialBookSearchResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantPilbkSearchPath)
	if err != nil {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNumber))
	values.Set("numOfRows", strconv.Itoa(query.NumberOfRows))
	setQueryValue(values, "reqSearchWrd", query.RequestSearchWord)
	setQueryValue(values, "dateFrom", query.DateFrom)
	setQueryValue(values, "dateTo", query.DateTo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantPilbkSearchResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantPictorialBookSearchResult{}, &PlantPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantPictorialBookSearchResult{}, fmt.Errorf("plantPilbkSearch: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantPictorialBookSearchResult{}, errors.New("plantPilbkSearch: response missing resultCode")
	}
	if payload.Header.ResultCode != plantPilbkSearchSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantPilbkSearchResultMessages[payload.Header.ResultCode]
		}
		return application.PlantPictorialBookSearchResult{}, &PlantPilbkSearchError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantPilbkSearchBody) result() application.PlantPictorialBookSearchResult {
	items := make([]application.PlantPictorialBookSearchItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantPictorialBookSearchItem{
			APGFamilyKoreanName:        item.APGFamilyKorNm,
			APGFamilyName:              item.APGFamilyNm,
			FamilyKoreanName:           item.FamilyKorNm,
			FamilyName:                 item.FamilyNm,
			GenusKoreanName:            item.GenusKorNm,
			GenusName:                  item.GenusNm,
			LastUpdateDateTime:         item.LastUpdtDtm,
			NotRecommendedGeneralName:  item.NotRcmmGnrlNm,
			PlantGeneralName:           item.PlantGnrlNm,
			PlantPictorialBookNumber:   item.PlantPilbkNo,
			PlantSpeciesScientificName: item.PlantSpecsScnm,
		}
	}

	return application.PlantPictorialBookSearchResult{
		Items:        items,
		NumberOfRows: body.NumOfRows,
		PageNumber:   body.PageNo,
		TotalCount:   body.TotalCount,
	}
}
