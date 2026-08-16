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
	plantSeedUnitListPath        = plantResourceBasePath + "/plantSeedUnitList"
	plantSeedUnitListSuccessCode = "00"
)

var plantSeedUnitListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSeedUnitListPort = (*Client)(nil)

type plantSeedUnitListBody struct {
	Items      []plantSeedUnitListItem `xml:"items>item"`
	NumOfRows  int                     `xml:"numOfRows"`
	PageNo     int                     `xml:"pageNo"`
	TotalCount int                     `xml:"totalCount"`
}

type plantSeedUnitListItem struct {
	CllcnDate        string `xml:"cllcnDate"`
	PlantGnrlNm      string `xml:"plantGnrlNm"`
	QualtFllnsRt     string `xml:"qualtFllnsRt"`
	SdwghWeght       string `xml:"sdwghWeght"`
	SeedAdmcn        string `xml:"seedAdmcn"`
	SeedCllctPlace   string `xml:"seedCllctPlace"`
	SeedHoldGrainCnt string `xml:"seedHoldGrainCnt"`
	SeedHoldQntt     string `xml:"seedHoldQntt"`
	SeedNo           string `xml:"seedNo"`
	SeedSpecsID      string `xml:"seedSpecsId"`
	StoreChrcrTpcdNm string `xml:"storeChrcrTpcdNm"`
	Vtlfct           string `xml:"vtlfct"`
	VtlfctTestYr     string `xml:"vtlfctTestYr"`
}

type plantSeedUnitListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSeedUnitListBody `xml:"body"`
}

// PlantSeedUnitListError reports an error returned by plantSeedUnitList.
type PlantSeedUnitListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSeedUnitList error message.
func (e *PlantSeedUnitListError) Error() string {
	return fmt.Sprintf("plantSeedUnitList: API error %s: %s", e.Code, e.Message)
}

// PlantSeedUnitList gets Korea National Arboretum plant seed unit information.
func (c *Client) PlantSeedUnitList(ctx context.Context, query application.PlantSeedUnitListQuery) (application.PlantSeedUnitListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSeedUnitListPath)
	if err != nil {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("reqSeedSpecsId", query.ReqSeedSpecsID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSeedUnitListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSeedUnitListResult{}, &PlantSeedUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSeedUnitListResult{}, fmt.Errorf("plantSeedUnitList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSeedUnitListResult{}, errors.New("plantSeedUnitList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSeedUnitListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSeedUnitListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSeedUnitListResult{}, &PlantSeedUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSeedUnitListBody) result() application.PlantSeedUnitListResult {
	items := make([]application.PlantSeedUnitListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSeedUnitListItem{
			CllcnDate:        item.CllcnDate,
			PlantGnrlNm:      item.PlantGnrlNm,
			QualtFllnsRt:     item.QualtFllnsRt,
			SdwghWeght:       item.SdwghWeght,
			SeedAdmcn:        item.SeedAdmcn,
			SeedCllctPlace:   item.SeedCllctPlace,
			SeedHoldGrainCnt: item.SeedHoldGrainCnt,
			SeedHoldQntt:     item.SeedHoldQntt,
			SeedNo:           item.SeedNo,
			SeedSpecsID:      item.SeedSpecsID,
			StoreChrcrTpcdNm: item.StoreChrcrTpcdNm,
			Vtlfct:           item.Vtlfct,
			VtlfctTestYr:     item.VtlfctTestYr,
		}
	}

	return application.PlantSeedUnitListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
