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
	plantNaturalizedListPath        = plantResourceBasePath + "/plantNaturalizedList"
	plantNaturalizedListSuccessCode = "00"
)

var plantNaturalizedListResultMessages = map[string]string{
	"00": "NORMAL SERVICE.",
	"02": "DB_ERROR",
	"03": "NODATA_ERROR",
	"05": "SERVICETIME_OUT",
	"10": "INVALID_REQUEST_PARAMETER_ERROR",
	"11": "NO_MANDATORY_REQUEST_PARAMETERS_ERROR",
	"21": "TEMPORARILY_DISABLE_THE_SERVICEKEY_ERROR",
	"33": "UNSIGNED_CALL_ERROR",
}

var _ outbound.PlantNaturalizedListPort = (*Client)(nil)

type plantNaturalizedListBody struct {
	Items      []plantNaturalizedListItem `xml:"items>item"`
	NumOfRows  int                        `xml:"numOfRows"`
	PageNo     int                        `xml:"pageNo"`
	TotalCount int                        `xml:"totalCount"`
}

type plantNaturalizedListItem struct {
	AgpFamilyNm        string `xml:"agpFamilyNm"`
	APGFamilyKorNm     string `xml:"apgFamilyKorNm"`
	BlprdEnmnt         string `xml:"blprdEnmnt"`
	BlprdStmnt         string `xml:"blprdStmnt"`
	DistrAraDscrt      string `xml:"distrAraDscrt"`
	EclgDstrbYn        string `xml:"eclgDstrbYn"`
	ExtcPlantCdNm      string `xml:"extcPlantCdNm"`
	FamilyKorNm        string `xml:"familyKorNm"`
	FamilyNm           string `xml:"familyNm"`
	FrtTpcdNm          string `xml:"frtTpcdNm"`
	LastUpdtDtm        string `xml:"lastUpdtDtm"`
	NtldgTpcdNm        string `xml:"ntldgTpcdNm"`
	NtrlzEraTpcdNm     string `xml:"ntrlzEraTpcdNm"`
	OrplcNm            string `xml:"orplcNm"`
	PlantBrdgFomTpcdNm string `xml:"plantBrdgFomTpcdNm"`
	PlantDistrGrcd     string `xml:"plantDistrGrcd"`
	PlantDistrQntt     string `xml:"plantDistrQntt"`
	PlantDistrQnttGrcd string `xml:"plantDistrQnttGrcd"`
	PlantEngNm         string `xml:"plantEngNm"`
	PlantGnrlNm        string `xml:"plantGnrlNm"`
	PlantJpnNm         string `xml:"plantJpnNm"`
	PlantLfcclTpcdNm   string `xml:"plantLfcclTpcdNm"`
	PlantSpecsScnm     string `xml:"plantSpecsScnm"`
}

type plantNaturalizedListResponse struct {
	Header struct {
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	GatewayHeader struct {
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
	Body plantNaturalizedListBody `xml:"body"`
}

// PlantNaturalizedListError reports an error returned by plantNaturalizedList.
type PlantNaturalizedListError struct {
	HTTPStatus int
	Code       string
	Message    string
}

// Error returns the plantNaturalizedList error message.
func (e *PlantNaturalizedListError) Error() string {
	return fmt.Sprintf("plantNaturalizedList: API error %s: %s", e.Code, e.Message)
}

// PlantNaturalizedList gets Korea National Arboretum naturalized plant information.
func (c *Client) PlantNaturalizedList(ctx context.Context, query application.PlantNaturalizedListQuery) (application.PlantNaturalizedListResult, error) {
	endpoint, err := url.JoinPath(c.baseURL, plantNaturalizedListPath)
	if err != nil {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: build URL: %w", err)
	}

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: parse URL: %w", err)
	}

	values := requestURL.Query()
	values.Set("serviceKey", c.serviceKey)
	values.Set("pageNo", strconv.Itoa(query.PageNo))
	values.Set("numOfRows", strconv.Itoa(query.NumOfRows))
	setQueryValue(values, "reqSearchWrd", query.ReqSearchWrd)
	setQueryValue(values, "dateFrom", query.DateFrom)
	setQueryValue(values, "dateTo", query.DateTo)
	requestURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: create request: %w", err)
	}
	request.Header.Set("Accept", "application/xml")

	response, err := c.do(request)
	if err != nil {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: request: %w", err)
	}
	defer response.Body.Close()

	var payload plantNaturalizedListResponse
	decodeErr := xml.NewDecoder(response.Body).Decode(&payload)
	if payload.GatewayHeader.ReturnReasonCode != "" {
		message := payload.GatewayHeader.ErrMsg
		if payload.GatewayHeader.ReturnAuthMsg != "" {
			message += ": " + payload.GatewayHeader.ReturnAuthMsg
		}
		return application.PlantNaturalizedListResult{}, &PlantNaturalizedListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.GatewayHeader.ReturnReasonCode,
			Message:    message,
		}
	}
	if response.StatusCode != http.StatusOK {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: unexpected HTTP status %s", response.Status)
	}
	if decodeErr != nil {
		return application.PlantNaturalizedListResult{}, fmt.Errorf("plantNaturalizedList: decode response: %w", decodeErr)
	}
	if payload.Header.ResultCode == "" {
		return application.PlantNaturalizedListResult{}, errors.New("plantNaturalizedList: response missing resultCode")
	}
	if payload.Header.ResultCode != plantNaturalizedListSuccessCode {
		message := payload.Header.ResultMsg
		if message == "" {
			message = plantNaturalizedListResultMessages[payload.Header.ResultCode]
		}
		return application.PlantNaturalizedListResult{}, &PlantNaturalizedListError{
			HTTPStatus: response.StatusCode,
			Code:       payload.Header.ResultCode,
			Message:    message,
		}
	}

	return payload.Body.result(), nil
}

func (body plantNaturalizedListBody) result() application.PlantNaturalizedListResult {
	items := make([]application.PlantNaturalizedListItem, len(body.Items))
	for i, item := range body.Items {
		items[i] = application.PlantNaturalizedListItem{
			AgpFamilyNm:        item.AgpFamilyNm,
			APGFamilyKorNm:     item.APGFamilyKorNm,
			BlprdEnmnt:         item.BlprdEnmnt,
			BlprdStmnt:         item.BlprdStmnt,
			DistrAraDscrt:      item.DistrAraDscrt,
			EclgDstrbYn:        item.EclgDstrbYn,
			ExtcPlantCdNm:      item.ExtcPlantCdNm,
			FamilyKorNm:        item.FamilyKorNm,
			FamilyNm:           item.FamilyNm,
			FrtTpcdNm:          item.FrtTpcdNm,
			LastUpdtDtm:        item.LastUpdtDtm,
			NtldgTpcdNm:        item.NtldgTpcdNm,
			NtrlzEraTpcdNm:     item.NtrlzEraTpcdNm,
			OrplcNm:            item.OrplcNm,
			PlantBrdgFomTpcdNm: item.PlantBrdgFomTpcdNm,
			PlantDistrGrcd:     item.PlantDistrGrcd,
			PlantDistrQntt:     item.PlantDistrQntt,
			PlantDistrQnttGrcd: item.PlantDistrQnttGrcd,
			PlantEngNm:         item.PlantEngNm,
			PlantGnrlNm:        item.PlantGnrlNm,
			PlantJpnNm:         item.PlantJpnNm,
			PlantLfcclTpcdNm:   item.PlantLfcclTpcdNm,
			PlantSpecsScnm:     item.PlantSpecsScnm,
		}
	}

	return application.PlantNaturalizedListResult{
		Items:      items,
		NumOfRows:  body.NumOfRows,
		PageNo:     body.PageNo,
		TotalCount: body.TotalCount,
	}
}
