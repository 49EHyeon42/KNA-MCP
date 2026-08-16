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
	plantSeedGrmntListPath        = plantResourceBasePath + "/plantSeedGrmntList"
	plantSeedGrmntListSuccessCode = "00"
)

var plantSeedGrmntListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSeedGrmntListPort = (*Client)(nil)

type plantSeedGrmntListBody struct {
	Items      []plantSeedGrmntListItem `xml:"items>item"`
	NumOfRows  int                      `xml:"numOfRows"`
	PageNo     int                      `xml:"pageNo"`
	TotalCount int                      `xml:"totalCount"`
}

type plantSeedGrmntListItem struct {
	AvrgGrmntDcnt     string `xml:"avrgGrmntDcnt"`
	GrmntBfrPrcesCont string `xml:"grmntBfrPrcesCont"`
	GrmntClmdmCont    string `xml:"grmntClmdmCont"`
	GrmntDscrt        string `xml:"grmntDscrt"`
	GrmntExprmNo      string `xml:"grmntExprmNo"`
	GrmntExprmSeq     string `xml:"grmntExprmSeq"`
	GrmntLightCndtn   string `xml:"grmntLightCndtn"`
	GrmntRt           string `xml:"grmntRt"`
	GrmntTmpCndtn     string `xml:"grmntTmpCndtn"`
	PlantGnrlNm       string `xml:"plantGnrlNm"`
	SeedNo            string `xml:"seedNo"`
	SeedSpecsID       string `xml:"seedSpecsId"`
}

type plantSeedGrmntListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSeedGrmntListBody `xml:"body"`
}

// PlantSeedGrmntListError reports an error returned by plantSeedGrmntList.
type PlantSeedGrmntListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSeedGrmntList error message.
func (e *PlantSeedGrmntListError) Error() string {
	return fmt.Sprintf("plantSeedGrmntList: API error %s: %s", e.Code, e.Message)
}

// PlantSeedGrmntList gets Korea National Arboretum plant seed germination information.
func (c *Client) PlantSeedGrmntList(ctx context.Context, query application.PlantSeedGrmntListQuery) (application.PlantSeedGrmntListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSeedGrmntListPath)
	if err != nil {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("reqSeedSpecsId", query.ReqSeedSpecsID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSeedGrmntListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSeedGrmntListResult{}, &PlantSeedGrmntListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSeedGrmntListResult{}, fmt.Errorf("plantSeedGrmntList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSeedGrmntListResult{}, errors.New("plantSeedGrmntList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSeedGrmntListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSeedGrmntListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSeedGrmntListResult{}, &PlantSeedGrmntListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSeedGrmntListBody) result() application.PlantSeedGrmntListResult {
	items := make([]application.PlantSeedGrmntListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSeedGrmntListItem{
			AvrgGrmntDcnt:     item.AvrgGrmntDcnt,
			GrmntBfrPrcesCont: item.GrmntBfrPrcesCont,
			GrmntClmdmCont:    item.GrmntClmdmCont,
			GrmntDscrt:        item.GrmntDscrt,
			GrmntExprmNo:      item.GrmntExprmNo,
			GrmntExprmSeq:     item.GrmntExprmSeq,
			GrmntLightCndtn:   item.GrmntLightCndtn,
			GrmntRt:           item.GrmntRt,
			GrmntTmpCndtn:     item.GrmntTmpCndtn,
			PlantGnrlNm:       item.PlantGnrlNm,
			SeedNo:            item.SeedNo,
			SeedSpecsID:       item.SeedSpecsID,
		}
	}

	return application.PlantSeedGrmntListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
