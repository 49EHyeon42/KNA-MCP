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
	plantSmplUnitListPath        = plantResourceBasePath + "/plantSmplUnitList"
	plantSmplUnitListSuccessCode = "00"
)

var plantSmplUnitListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantSmplUnitListPort = (*Client)(nil)

type plantSmplUnitListBody struct {
	Items      []plantSmplUnitListItem `xml:"items>item"`
	NumOfRows  int                     `xml:"numOfRows"`
	PageNo     int                     `xml:"pageNo"`
	TotalCount int                     `xml:"totalCount"`
}

type plantSmplUnitListItem struct {
	AgpFamilyKorNm     string `xml:"agpFamilyKorNm"`
	AgpFamilyNm        string `xml:"agpFamilyNm"`
	BspcsInsttNm       string `xml:"bspcsInsttNm"`
	ClarHaslvVal       string `xml:"clarHaslvVal"`
	ClarNm             string `xml:"clarNm"`
	CllcrNm            string `xml:"cllcrNm"`
	FamilyKorNm        string `xml:"familyKorNm"`
	FamilyNm           string `xml:"familyNm"`
	HbttChrcrCont      string `xml:"hbttChrcrCont"`
	HbttTpcdNm         string `xml:"hbttTpcdNm"`
	PlantBrdgFomTpcdNm string `xml:"plantBrdgFomTpcdNm"`
	PlantGnrlNm        string `xml:"plantGnrlNm"`
	PlantPilbkNo       string `xml:"plantPilbkNo"`
	PlantSmplNo        string `xml:"plantSmplNo"`
	PlantSpecsID       string `xml:"plantSpecsId"`
	PlantSpecsScnm     string `xml:"plantSpecsScnm"`
	SmplCllcnDt        string `xml:"smplCllcnDt"`
	SmplClnyNm         string `xml:"smplClnyNm"`
	SmplKindCdNm       string `xml:"smplKindCdNm"`
	SmplWrdt           string `xml:"smplWrdt"`
	VgttnTpeCdNm       string `xml:"vgttnTpeCdNm"`
}

type plantSmplUnitListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantSmplUnitListBody `xml:"body"`
}

// PlantSmplUnitListError reports an error returned by plantSmplUnitList.
type PlantSmplUnitListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantSmplUnitList error message.
func (e *PlantSmplUnitListError) Error() string {
	return fmt.Sprintf("plantSmplUnitList: API error %s: %s", e.Code, e.Message)
}

// PlantSmplUnitList gets Korea National Arboretum plant specimen details.
func (c *Client) PlantSmplUnitList(ctx context.Context, query application.PlantSmplUnitListQuery) (application.PlantSmplUnitListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantSmplUnitListPath)
	if err != nil {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	values.Set("reqPlantSpecsId", query.ReqPlantSpecsID)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantSmplUnitListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantSmplUnitListResult{}, &PlantSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantSmplUnitListResult{}, fmt.Errorf("plantSmplUnitList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantSmplUnitListResult{}, errors.New("plantSmplUnitList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantSmplUnitListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantSmplUnitListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantSmplUnitListResult{}, &PlantSmplUnitListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantSmplUnitListBody) result() application.PlantSmplUnitListResult {
	items := make([]application.PlantSmplUnitListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantSmplUnitListItem{
			AgpFamilyKorNm:     item.AgpFamilyKorNm,
			AgpFamilyNm:        item.AgpFamilyNm,
			BspcsInsttNm:       item.BspcsInsttNm,
			ClarHaslvVal:       item.ClarHaslvVal,
			ClarNm:             item.ClarNm,
			CllcrNm:            item.CllcrNm,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			HbttChrcrCont:      item.HbttChrcrCont,
			HbttTpcdNm:         item.HbttTpcdNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantPilbkNo:       item.PlantPilbkNo,
			PlantSmplNo:        item.PlantSmplNo,
			PlantSpecsID:       item.PlantSpecsID,
			PlantSpecsScnm:     item.PlantSpecsScnm,
			SmplCllcnDt:        item.SmplCllcnDt,
			SmplClnyNm:         item.SmplClnyNm,
			SmplKindCdNm:       item.SmplKindCdNm,
			SmplWrdt:           item.SmplWrdt,
			VgttnTpeCdNm:       item.VgttnTpeCdNm,
		}
	}

	return application.PlantSmplUnitListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
