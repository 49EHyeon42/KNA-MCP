package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/insectresource/application/port/outbound"
)

const (
	insectSmplUnitListPath        = insectResourceBasePath + "/insectSmplUnitList"
	insectSmplUnitListSuccessCode = "00"
)

var insectSmplUnitListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.InsectSmplUnitListPort = (*Client)(nil)

type insectSmplUnitListBody struct {
	Items      []insectSmplUnitListItem `xml:"items>item"`
	NumOfRows  int                      `xml:"numOfRows"`
	PageNo     int                      `xml:"pageNo"`
	TotalCount int                      `xml:"totalCount"`
}

type insectSmplUnitListItem struct {
	BspcsInsttNm       string `xml:"bspcsInsttNm"`
	ClarHaslvVal       string `xml:"clarHaslvVal"`
	SmplCllcnDt        string `xml:"smplCllcnDt"`
	GynndTpcd          string `xml:"gynndTpcd"`
	HbttTpcd           string `xml:"hbttTpcd"`
	InsctSmplNo        string `xml:"insctSmplNo"`
	InsctSpecsID       string `xml:"insctSpecsId"`
	InsctSpecsScnm     string `xml:"insctSpecsScnm"`
	LabelUsgCllcnNmplc string `xml:"labelUsgCllcnNmplc"`
	LastUpdtDtm        string `xml:"lastUpdtDtm"`
	PrsrtStcd          string `xml:"prsrtStcd"`
	SlistTpcd          string `xml:"slistTpcd"`
	SmplKindCd         string `xml:"smplKindCd"`
	TorsoLngth         string `xml:"torsoLngth"`
	WingLngth          string `xml:"wingLngth"`
	InsctGnrlNm        string `xml:"insctGnrlNm"`
}

type insectSmplUnitListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body insectSmplUnitListBody `xml:"body"`
}

// InsectSmplUnitListError reports an error returned by insectSmplUnitList.
type InsectSmplUnitListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the insectSmplUnitList error message.
func (e *InsectSmplUnitListError) Error() string {
	return fmt.Sprintf("insectSmplUnitList: API error %s: %s", e.Code, e.Message)
}

// InsectSmplUnitList gets Korea National Arboretum insect specimen details.
func (c *Client) InsectSmplUnitList(ctx context.Context, query application.InsectSmplUnitListQuery) (application.InsectSmplUnitListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, insectSmplUnitListPath)
	if err != nil {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("reqInsctSpecsId", query.ReqInsctSpecsID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: request: %w", err)
	}
	defer response.Body.Close()

	var payload insectSmplUnitListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.InsectSmplUnitListResult{}, &InsectSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.InsectSmplUnitListResult{}, fmt.Errorf("insectSmplUnitList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.InsectSmplUnitListResult{}, errors.New("insectSmplUnitList: response missing resultCode")
	}
	if payload.Header.ResultCode != insectSmplUnitListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = insectSmplUnitListResultMessages[payload.Header.ResultCode]
		}
		return application.InsectSmplUnitListResult{}, &InsectSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body insectSmplUnitListBody) result() application.InsectSmplUnitListResult {
	items := make([]application.InsectSmplUnitListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.InsectSmplUnitListItem{
			BspcsInsttNm:       item.BspcsInsttNm,
			ClarHaslvVal:       item.ClarHaslvVal,
			SmplCllcnDt:        item.SmplCllcnDt,
			GynndTpcd:          item.GynndTpcd,
			HbttTpcd:           item.HbttTpcd,
			InsctSmplNo:        item.InsctSmplNo,
			InsctSpecsID:       item.InsctSpecsID,
			InsctSpecsScnm:     item.InsctSpecsScnm,
			LabelUsgCllcnNmplc: item.LabelUsgCllcnNmplc,
			LastUpdtDtm:        item.LastUpdtDtm,
			PrsrtStcd:          item.PrsrtStcd,
			SlistTpcd:          item.SlistTpcd,
			SmplKindCd:         item.SmplKindCd,
			TorsoLngth:         item.TorsoLngth,
			WingLngth:          item.WingLngth,
			InsctGnrlNm:        item.InsctGnrlNm,
		}
	}

	return application.InsectSmplUnitListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
