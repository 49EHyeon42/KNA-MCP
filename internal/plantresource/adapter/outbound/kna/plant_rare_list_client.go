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
	plantRareListPath        = plantResourceBasePath + "/plantRareList"
	plantRareListSuccessCode = "00"
)

var plantRareListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantRareListPort = (*Client)(nil)

type plantRareListBody struct {
	Items      []plantRareListItem `xml:"items>item"`
	NumOfRows  int                 `xml:"numOfRows"`
	PageNo     int                 `xml:"pageNo"`
	TotalCount int                 `xml:"totalCount"`
}

type plantRareListItem struct {
	AgpFamilyNm      string `xml:"agpFamilyNm"`
	APGFamilyKorNm   string `xml:"apgFamilyKorNm"`
	ExtrmCrssScls1Yn string `xml:"extrmCrssScls1Yn"`
	ExtrmCrssScls2Yn string `xml:"extrmCrssScls2Yn"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	PlantGnrlNm      string `xml:"plantGnrlNm"`
	PlantSpecsScnm   string `xml:"plantSpecsScnm"`
	RareTpcdNm       string `xml:"rareTpcdNm"`
}

type plantRareListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantRareListBody `xml:"body"`
}

// PlantRareListError reports an error returned by plantRareList.
type PlantRareListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantRareList error message.
func (e *PlantRareListError) Error() string {
	return fmt.Sprintf("plantRareList: API error %s: %s", e.Code, e.Message)
}

// PlantRareList gets Korea National Arboretum rare plant information.
func (c *Client) PlantRareList(ctx context.Context, query application.PlantRareListQuery) (application.PlantRareListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantRareListPath)
	if err != nil {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantRareListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantRareListResult{}, &PlantRareListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantRareListResult{}, fmt.Errorf("plantRareList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantRareListResult{}, errors.New("plantRareList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantRareListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantRareListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantRareListResult{}, &PlantRareListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantRareListBody) result() application.PlantRareListResult {
	items := make([]application.PlantRareListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantRareListItem{
			AgpFamilyNm:      item.AgpFamilyNm,
			APGFamilyKorNm:   item.APGFamilyKorNm,
			ExtrmCrssScls1Yn: item.ExtrmCrssScls1Yn,
			ExtrmCrssScls2Yn: item.ExtrmCrssScls2Yn,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			PlantGnrlNm:      item.PlantGnrlNm,
			PlantSpecsScnm:   item.PlantSpecsScnm,
			RareTpcdNm:       item.RareTpcdNm,
		}
	}

	return application.PlantRareListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
