package kna

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application"
	"github.com/49EHyeon42/KNA-MCP/internal/fungiresource/application/port/outbound"
)

const (
	fngsSmplUnitListPath        = fungiResourceBasePath + "/fngsSmplUnitList"
	fngsSmplUnitListSuccessCode = "00"
)

var fngsSmplUnitListResultMessages = map[string]string{
	"00": "NORMAL_SERVICE",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.FngsSmplUnitListPort = (*Client)(nil)

type fngsSmplUnitListBody struct {
	Items      []fngsSmplUnitListItem `xml:"items>item"`
	NumOfRows  int                    `xml:"numOfRows"`
	PageNo     int                    `xml:"pageNo"`
	TotalCount int                    `xml:"totalCount"`
}

type fngsSmplUnitListItem struct {
	ClarDtlDscrt     string `xml:"clarDtlDscrt"`
	ClarHaslvVal     string `xml:"clarHaslvVal"`
	CllcrNm          string `xml:"cllcrNm"`
	FamilyKorNm      string `xml:"familyKorNm"`
	FamilyNm         string `xml:"familyNm"`
	FngsEclgTpcdNm   string `xml:"fngsEclgTpcdNm"`
	FngsGnrlNm       string `xml:"fngsGnrlNm"`
	FngsID           string `xml:"fngsId"`
	FngsScnm         string `xml:"fngsScnm"`
	FngsSmplKindCdNm string `xml:"fngsSmplKindCdNm"`
	FngsSmplNo       string `xml:"fngsSmplNo"`
	GenusKorNm       string `xml:"genusKorNm"`
	GenusNm          string `xml:"genusNm"`
	HbttChrcrCont    string `xml:"hbttChrcrCont"`
	HstCont          string `xml:"hstCont"`
	LastUpdtDtm      string `xml:"lastUpdtDtm"`
	SmplCllcnDt      string `xml:"smplCllcnDt"`
}

type fngsSmplUnitListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body fngsSmplUnitListBody `xml:"body"`
}

// FngsSmplUnitListError reports an error returned by fngsSmplUnitList.
type FngsSmplUnitListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the fngsSmplUnitList error message.
func (e *FngsSmplUnitListError) Error() string {
	return fmt.Sprintf("fngsSmplUnitList: API error %s: %s", e.Code, e.Message)
}

// FngsSmplUnitList gets Korea National Arboretum fungi specimen details.
func (c *Client) FngsSmplUnitList(ctx context.Context, query application.FngsSmplUnitListQuery) (application.FngsSmplUnitListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, fngsSmplUnitListPath)
	if err != nil {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("reqFngsId", query.ReqFngsID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: request: %w", err)
	}
	defer response.Body.Close()

	var payload fngsSmplUnitListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.FngsSmplUnitListResult{}, &FngsSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.FngsSmplUnitListResult{}, fmt.Errorf("fngsSmplUnitList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.FngsSmplUnitListResult{}, errors.New("fngsSmplUnitList: response missing resultCode")
	}
	if payload.Header.ResultCode != fngsSmplUnitListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = fngsSmplUnitListResultMessages[payload.Header.ResultCode]
		}
		return application.FngsSmplUnitListResult{}, &FngsSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body fngsSmplUnitListBody) result() application.FngsSmplUnitListResult {
	items := make([]application.FngsSmplUnitListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.FngsSmplUnitListItem{
			ClarDtlDscrt:     item.ClarDtlDscrt,
			ClarHaslvVal:     item.ClarHaslvVal,
			CllcrNm:          item.CllcrNm,
			FamilyKorNm:      item.FamilyKorNm,
			FamilyNm:         item.FamilyNm,
			FngsEclgTpcdNm:   item.FngsEclgTpcdNm,
			FngsGnrlNm:       item.FngsGnrlNm,
			FngsID:           item.FngsID,
			FngsScnm:         item.FngsScnm,
			FngsSmplKindCdNm: item.FngsSmplKindCdNm,
			FngsSmplNo:       item.FngsSmplNo,
			GenusKorNm:       item.GenusKorNm,
			GenusNm:          item.GenusNm,
			HbttChrcrCont:    item.HbttChrcrCont,
			HstCont:          item.HstCont,
			LastUpdtDtm:      item.LastUpdtDtm,
			SmplCllcnDt:      item.SmplCllcnDt,
		}
	}

	return application.FngsSmplUnitListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
