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
	plantFolkAreaListPath        = plantResourceBasePath + "/plantFolkAreaList"
	plantFolkAreaListSuccessCode = "00"
)

var plantFolkAreaListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantFolkAreaListPort = (*Client)(nil)

type plantFolkAreaListBody struct {
	Items      []plantFolkAreaListItem `xml:"items>item"`
	NumOfRows  int                     `xml:"numOfRows"`
	PageNo     int                     `xml:"pageNo"`
	TotalCount int                     `xml:"totalCount"`
}

type plantFolkAreaListItem struct {
	FlcstPlantExmnnAraTpcdNm string `xml:"flcstPlantExmnnAraTpcdNm"`
	FlcstPlantLcltDscrt      string `xml:"flcstPlantLcltDscrt"`
	FlcstPlantPrpseDscrt     string `xml:"flcstPlantPrpseDscrt"`
	FlpltID                  string `xml:"flpltId"`
	PlantBrdgFomTpcdNm       string `xml:"plantBrdgFomTpcdNm"`
	PlantGnrlNm              string `xml:"plantGnrlNm"`
	PlantSpecsScnm           string `xml:"plantSpecsScnm"`
}

type plantFolkAreaListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantFolkAreaListBody `xml:"body"`
}

// PlantFolkAreaListError reports an error returned by plantFolkAreaList.
type PlantFolkAreaListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantFolkAreaList error message.
func (e *PlantFolkAreaListError) Error() string {
	return fmt.Sprintf("plantFolkAreaList: API error %s: %s", e.Code, e.Message)
}

// PlantFolkAreaList gets Korea National Arboretum folk plant area information.
func (c *Client) PlantFolkAreaList(ctx context.Context, query application.PlantFolkAreaListQuery) (application.PlantFolkAreaListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantFolkAreaListPath)
	if err != nil {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("flpltId", query.FlpltID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantFolkAreaListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantFolkAreaListResult{}, &PlantFolkAreaListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantFolkAreaListResult{}, fmt.Errorf("plantFolkAreaList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantFolkAreaListResult{}, errors.New("plantFolkAreaList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantFolkAreaListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantFolkAreaListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantFolkAreaListResult{}, &PlantFolkAreaListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantFolkAreaListBody) result() application.PlantFolkAreaListResult {
	items := make([]application.PlantFolkAreaListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantFolkAreaListItem{
			FlcstPlantExmnnAraTpcdNm: item.FlcstPlantExmnnAraTpcdNm,
			FlcstPlantLcltDscrt:      item.FlcstPlantLcltDscrt,
			FlcstPlantPrpseDscrt:     item.FlcstPlantPrpseDscrt,
			FlpltID:                  item.FlpltID,
			PlantBrdgFomTpcdNm:       item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:              item.PlantGnrlNm,
			PlantSpecsScnm:           item.PlantSpecsScnm,
		}
	}

	return application.PlantFolkAreaListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
