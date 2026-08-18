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
	plantSpcltListPath        = plantResourceBasePath + "/plantSpcltList"
	plantSpcltListSuccessCode = "00"
)

var plantSpcltListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSpcltListPort = (*Client)(nil)

type plantSpcltListBody struct {
	Items      []plantSpcltListItem `xml:"items>item"`
	NumOfRows  int                  `xml:"numOfRows"`
	PageNo     int                  `xml:"pageNo"`
	TotalCount int                  `xml:"totalCount"`
}

type plantSpcltListItem struct {
	AgpFamilyKorNm     string `xml:"agpFamilyKorNm"`
	AgpFamilyNm        string `xml:"agpFamilyNm"`
	ExtrmCrssScls1Yn   string `xml:"extrmCrssScls1Yn"`
	ExtrmCrssScls2Yn   string `xml:"extrmCrssScls2Yn"`
	FamilyKorNm        string `xml:"familyKorNm"`
	FamilyNm           string `xml:"familyNm"`
	PlantBrdgFomTpcdNm string `xml:"plantBrdgFomTpcdNm"`
	PlantGnrlNm        string `xml:"plantGnrlNm"`
	PlantSpecsScnm     string `xml:"plantSpecsScnm"`
}

type plantSpcltListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSpcltListBody `xml:"body"`
}

// PlantSpcltListError reports an error returned by plantSpcltList.
type PlantSpcltListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSpcltList error message.
func (e *PlantSpcltListError) Error() string {
	return fmt.Sprintf("plantSpcltList: API error %s: %s", e.Code, e.Message)
}

// PlantSpcltList gets Korea National Arboretum endemic plant information.
func (c *Client) PlantSpcltList(ctx context.Context, query application.PlantSpcltListQuery) (application.PlantSpcltListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSpcltListPath)
	if err != nil {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSpcltListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSpcltListResult{}, &PlantSpcltListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSpcltListResult{}, fmt.Errorf("plantSpcltList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSpcltListResult{}, errors.New("plantSpcltList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSpcltListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSpcltListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSpcltListResult{}, &PlantSpcltListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSpcltListBody) result() application.PlantSpcltListResult {
	items := make([]application.PlantSpcltListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSpcltListItem{
			AgpFamilyKorNm:     item.AgpFamilyKorNm,
			AgpFamilyNm:        item.AgpFamilyNm,
			ExtrmCrssScls1Yn:   item.ExtrmCrssScls1Yn,
			ExtrmCrssScls2Yn:   item.ExtrmCrssScls2Yn,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantSpecsScnm:     item.PlantSpecsScnm,
		}
	}

	return application.PlantSpcltListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
